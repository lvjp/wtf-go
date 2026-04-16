package token

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var errNilGenerator = errors.New("token.NewManager: ID generator is nil")
var errNilStore = errors.New("token.NewManager: store is nil")
var ErrExpiredToken = errors.New("token.Manager: expired token")
var errNonPositiveTTL = errors.New("token.manager: TTL should be positive")

type Manager[
	T TokenData,
] interface {
	Create(ctx context.Context, ttl time.Duration, data T) (*Token[T], error)
	Lookup(ctx context.Context, ID string) (*Token[T], error)
	Revoke(ctx context.Context, ID string) error
}

func NewManager[
	T interface {
		TokenData
		*TBase
	},
	TBase any,
](gen IDGenerator, store Store) (Manager[T], error) {
	if gen == nil {
		return nil, errNilGenerator
	}

	if store == nil {
		return nil, errNilStore
	}

	ret := &manager[T, TBase]{
		gen:   gen,
		store: store,
	}

	return ret, nil
}

type manager[
	T interface {
		TokenData
		*TBase
	},
	TBase any,
] struct {
	gen   IDGenerator
	store Store
}

func (m *manager[T, TBase]) Create(ctx context.Context, ttl time.Duration, data T) (*Token[T], error) {
	if err := validateTTL(ttl); err != nil {
		return nil, err
	}

	toStore := StoredToken{
		NotAfter: time.Now().Add(ttl),
	}

	var err error
	toStore.ID, err = m.gen()
	if err != nil {
		return nil, fmt.Errorf("token.Manager: ID generator error: %v", err)
	}

	toStore.Data, err = data.MarshalBinary()
	if err != nil {
		return nil, errors.Join(ErrMarshalData, err)
	}

	if err := m.store.Create(ctx, toStore); err != nil {
		return nil, fmt.Errorf("token.Manager: %w", err)
	}

	ret := &Token[T]{
		ID:       toStore.ID,
		Data:     data,
		NotAfter: toStore.NotAfter,
	}

	return ret, nil
}

func (m *manager[T, TBase]) Lookup(ctx context.Context, id string) (*Token[T], error) {
	stored, err := m.store.Read(ctx, id)
	if err != nil {
		return nil, err
	}

	if time.Now().After(stored.NotAfter) {
		return nil, ErrExpiredToken
	}

	var data T = new(TBase)
	if err := data.UnmarshalBinary(stored.Data); err != nil {
		return nil, errors.Join(ErrUnmarshalData, err)
	}

	ret := &Token[T]{
		ID:       id,
		Data:     data,
		NotAfter: stored.NotAfter,
	}

	return ret, nil
}

func (m *manager[T, TBase]) Revoke(ctx context.Context, id string) error {
	return m.store.Delete(ctx, id)
}

func validateTTL(ttl time.Duration) error {
	if ttl <= 0 {
		return errNonPositiveTTL
	}

	return nil
}
