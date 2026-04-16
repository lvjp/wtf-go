package token

import (
	"context"
	"fmt"
	"slices"
	"sync"
)

// NewMemoryStore return an in-memory implementation of the Store interface.
// This is not suitable for production use, but can be used for testing and development purposes.
func NewMemoryStore() Store {
	return &memoryStore{}
}

type memoryStore struct {
	tokens sync.Map
}

func (ms *memoryStore) Create(_ context.Context, token StoredToken) error {
	if token.ID == "" {
		return ErrEmptyID
	}

	// Prevent external modification
	token.Data = slices.Clone(token.Data)

	_, loaded := ms.tokens.LoadOrStore(token.ID, &token)
	if loaded {
		return ErrAlreadyExists
	}

	return nil
}

func (ms *memoryStore) Read(_ context.Context, id string) (*StoredToken, error) {
	if id == "" {
		return nil, ErrEmptyID
	}

	raw, ok := ms.tokens.Load(id)
	if !ok {
		return nil, ErrNotFound
	}

	v, ok := raw.(*StoredToken)
	if !ok {
		panic(fmt.Sprintf("token.memoryStore: unexpected stored value type: %T", raw))
	}

	// Prevent external modification
	v.Data = slices.Clone(v.Data)

	return v, nil
}

func (ms *memoryStore) Delete(_ context.Context, id string) error {
	ms.tokens.Delete(id)

	return nil
}
