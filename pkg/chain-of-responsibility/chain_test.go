package chain

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewChain(t *testing.T) {
	handlerValue := "handler-value"
	handlerError := errors.New("handler-error")

	handler := HandlerFunc[string](func(_ context.Context, value string) error {
		require.Equal(t, handlerValue, value)
		return handlerError
	})

	chain := NewChain(handler)

	actualError := chain.Handle(context.Background(), handlerValue)
	require.Same(t, handlerError, actualError)
}

func TestNewChainMiddlewareExecutionOrder(t *testing.T) {
	const input = "input"
	var order []string

	makeMiddleware := func(name string) Middleware[string] {
		return MiddlewareFunc[string](func(ctx context.Context, value string, next Handler[string]) error {
			order = append(order, name+":before")
			err := next.Handle(ctx, value)
			order = append(order, name+":after")
			return err
		})
	}

	handler := HandlerFunc[string](func(_ context.Context, value string) error {
		require.Equal(t, input, value)
		order = append(order, "handler")
		return nil
	})

	chain := NewChain(handler, makeMiddleware("m1"), makeMiddleware("m2"), makeMiddleware("m3"))
	err := chain.Handle(context.Background(), input)
	require.NoError(t, err)

	require.Equal(t, []string{
		"m1:before", "m2:before", "m3:before",
		"handler",
		"m3:after", "m2:after", "m1:after",
	}, order)
}

func TestNewChain_middlewareShortCircuit(t *testing.T) {
	handlerCalled := false
	shortCircuitErr := errors.New("short-circuit")

	handler := HandlerFunc[string](func(context.Context, string) error {
		handlerCalled = true
		return errors.New("short-circuited-output")
	})

	shortCircuit := MiddlewareFunc[string](func(context.Context, string, Handler[string]) error {
		return shortCircuitErr
	})

	passThrough := MiddlewareFunc[string](func(ctx context.Context, input string, next Handler[string]) error {
		return next.Handle(ctx, input)
	})

	chain := NewChain(handler, passThrough, shortCircuit, passThrough)
	err := chain.Handle(context.Background(), "input")

	require.False(t, handlerCalled)
	require.ErrorIs(t, err, shortCircuitErr)
}

func TestNewChain_middlewareInputTransformation(t *testing.T) {
	handler := HandlerFunc[string](func(_ context.Context, input string) error {
		require.Equal(t, "transformed-input", input)
		return nil
	})

	transform := MiddlewareFunc[string](func(ctx context.Context, input string, next Handler[string]) error {
		return next.Handle(ctx, "transformed-"+input)
	})

	chain := NewChain(handler, transform)
	err := chain.Handle(context.Background(), "input")

	require.NoError(t, err)
}

func TestNewChain_contextPropagation(t *testing.T) {
	type key string
	const k key = "k"

	handler := HandlerFunc[string](func(ctx context.Context, input string) error {
		require.Equal(t, "set-by-middleware", ctx.Value(k))
		return nil
	})

	setCtx := MiddlewareFunc[string](func(ctx context.Context, input string, next Handler[string]) error {
		return next.Handle(context.WithValue(ctx, k, "set-by-middleware"), input)
	})

	chain := NewChain(handler, setCtx)
	err := chain.Handle(context.Background(), "input")
	require.NoError(t, err)
}

func ExampleNewChain_notification() {
	addOne := MiddlewareFunc[int](func(ctx context.Context, value int, next Handler[int]) error {
		return next.Handle(ctx, value+1)
	})
	multiplyByTwo := MiddlewareFunc[int](func(ctx context.Context, value int, next Handler[int]) error {
		return next.Handle(ctx, value*2)
	})

	h := HandlerFunc[int](func(_ context.Context, value int) error {
		fmt.Println("handler received value:", value)
		return nil
	})

	chain := NewChain(h, multiplyByTwo, addOne)
	err := chain.Handle(context.Background(), 3)
	if err != nil {
		panic("unexpected error: " + err.Error())
	}

	// Output:
	// handler received value: 7
}

func ExampleNewChain_requestReply() {
	type WorkingContext struct {
		Request  int
		Response string
	}

	addOne := MiddlewareFunc[*WorkingContext](func(ctx context.Context, wc *WorkingContext, next Handler[*WorkingContext]) error {
		wc.Request += 1
		return next.Handle(ctx, wc)
	})
	multiplyByTwo := MiddlewareFunc[*WorkingContext](func(ctx context.Context, wc *WorkingContext, next Handler[*WorkingContext]) error {
		wc.Request *= 2
		return next.Handle(ctx, wc)
	})

	h := HandlerFunc[(*WorkingContext)](func(_ context.Context, wc *WorkingContext) error {
		wc.Response = strconv.Itoa(wc.Request)
		return nil
	})

	chain := NewChain(h, multiplyByTwo, addOne)
	wc := &WorkingContext{Request: 3}
	err := chain.Handle(context.Background(), wc)
	if err != nil {
		panic("unexpected error: " + err.Error())
	}
	fmt.Println("computed value:", wc.Response)

	// Output:
	// computed value: 7
}
