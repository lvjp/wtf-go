// Package chain implements the chain of responsibility pattern.
//
// A value passes through a sequence of middlewares, each able to transform it,
// short-circuit the chain, or delegate to the next handler:
//
//	value → [middleware 1] → [middleware 2] → [handler]
//
// Each handler returns only an error — there is no separate return value. This
// is a deliberate design choice to respect Go conventions: returning a value
// alongside an error implies the value is meaningful even on failure, which Go
// discourages. Embedding the response inside the error would be equally wrong.
//
// For request/reply pipelines, pass a pointer to a mutable struct as Value and
// let the handler populate a response field — the caller owns the convention:
//
//	type WorkingContext struct {
//	    Request  int
//	    Response string
//	}
//
// See https://en.wikipedia.org/wiki/Chain-of-responsibility_pattern
package chain

import (
	"context"
)

type Handler[Value any] interface {
	Handle(ctx context.Context, value Value) error
}

type HandlerFunc[Value any] func(ctx context.Context, value Value) error

func (fn HandlerFunc[Value]) Handle(ctx context.Context, value Value) error {
	return fn(ctx, value)
}

type Middleware[Value any] interface {
	Do(ctx context.Context, value Value, next Handler[Value]) error
}

type MiddlewareFunc[Value any] func(ctx context.Context, value Value, next Handler[Value]) error

func (fn MiddlewareFunc[Value]) Do(ctx context.Context, value Value, next Handler[Value]) error {
	return fn(ctx, value, next)
}

type chainLink[Value any] struct {
	next Handler[Value]
	with Middleware[Value]
}

func (cl chainLink[Value]) Handle(ctx context.Context, value Value) error {
	return cl.with.Do(ctx, value, cl.next)
}

func NewChain[Value any](h Handler[Value], with ...Middleware[Value]) Handler[Value] {
	for i := len(with) - 1; i >= 0; i-- {
		h = chainLink[Value]{
			next: h,
			with: with[i],
		}
	}

	return h
}
