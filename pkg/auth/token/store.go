package token

import (
	"context"
	"errors"
	"time"
)

var ErrNotFound = errors.New("token.Store: token not found")
var ErrAlreadyExists = errors.New("token.Store: ID already exists")
var ErrEmptyID = errors.New("token.Store: ID cannot be empty")
var ErrMarshalData = errors.New("token.Store: failed to marshal associated data")
var ErrUnmarshalData = errors.New("token.Store: failed to unmarshal associated data")

// Store defines the interface for token storage.
// Store is only responsible for storing, retrieving, refreshing and deleting tokens.
// Not for data validation.
// It abstracts away the underlying storage mechanism, allowing for different implementations.
type Store interface {
	// Create a new token into the store.
	// If the token ID is empty, ErrEmptyID is returned.
	// If a token with the same ID already exists, ErrAlreadyExists is returned.
	Create(ctx context.Context, token StoredToken) error

	// Read retrieves the token with the given ID from the store.
	// If the token ID is empty, ErrEmptyID is returned.
	// If the token is not found, ErrNotFound is returned.
	Read(ctx context.Context, id string) (*StoredToken, error)

	// Delete removes the token with the given ID from the store.
	// No error is returned if the token does not exist or if the ID is empty.
	Delete(ctx context.Context, id string) error
}

// StoredToken represents the token data that is stored in the Store.
type StoredToken struct {
	ID       string
	Data     []byte
	NotAfter time.Time
}
