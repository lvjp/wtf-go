package token

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"testing/synctest"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestNewManager(t *testing.T) {
	t.Run("nil nil", func(t *testing.T) {
		m, err := NewManager[*TestData](nil, nil)
		require.ErrorIs(t, err, errNilGenerator)
		require.Nil(t, m)
	})

	t.Run("nil generator", func(t *testing.T) {
		m, err := NewManager[*TestData](nil, NewMemoryStore())
		require.ErrorIs(t, err, errNilGenerator)
		require.Nil(t, m)
	})

	t.Run("nil store", func(t *testing.T) {
		m, err := NewManager[*TestData](UUIDGenerator, nil)
		require.ErrorIs(t, err, errNilStore)
		require.Nil(t, m)
	})

	t.Run("normal", func(t *testing.T) {
		m, err := NewManager[*TestData](UUIDGenerator, NewMemoryStore())
		require.NoError(t, err)
		require.NotNil(t, m)
	})
}

func TestManager_Create(t *testing.T) {
	t.Run("invalid ttl", func(t *testing.T) {
		m, err := NewManager[*TestData](UUIDGenerator, NewMemoryStore())
		require.NoError(t, err)

		token, err := m.Create(t.Context(), 0, newTestData())
		require.ErrorIs(t, err, errNonPositiveTTL)
		require.Nil(t, token)
	})

	t.Run("ID generator error", func(t *testing.T) {
		m, err := NewManager[*TestData](generateTokenIDError, NewMemoryStore())
		require.NoError(t, err)
		require.NotNil(t, m)

		token, err := m.Create(t.Context(), time.Second, newTestData())
		require.ErrorContains(t, err, "token.Manager: ID generator error: generateTokenIDError")
		require.Nil(t, token)
	})

	t.Run("empty ID", func(t *testing.T) {
		m, err := NewManager[*TestData](generateTokenIDEmpty, NewMemoryStore())
		require.NoError(t, err)

		token, err := m.Create(t.Context(), time.Second, newTestData())
		require.ErrorIs(t, err, ErrEmptyID)
		require.Nil(t, token)
	})

	t.Run("marshall error", func(t *testing.T) {
		m, err := NewManager[*TestData](UUIDGenerator, NewMemoryStore())
		require.NoError(t, err)
		require.NotNil(t, m)

		data := &TestData{
			Key:   "invalid:key",
			Value: "value",
		}
		token, err := m.Create(t.Context(), time.Second, data)
		require.ErrorIs(t, err, ErrMarshalData)
		require.ErrorContains(t, err, "invalid character ':' in key")
		require.Nil(t, token)
	})

	t.Run("normal", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			m, err := NewManager[*TestData](UUIDGenerator, NewMemoryStore())
			require.NoError(t, err)
			require.NotNil(t, m)

			data := newTestData()
			const ttl = time.Second
			token, err := m.Create(t.Context(), ttl, data)
			require.NoError(t, err)
			require.NotNil(t, token)
			require.NotEmpty(t, token.ID)
			require.Equal(t, data, token.Data)
			require.Equal(t, time.Now().Add(ttl), token.NotAfter)
		})
	})
}

func TestManager_Lookup(t *testing.T) {
	t.Run("not found", func(t *testing.T) {
		m, err := NewManager[*TestData](UUIDGenerator, NewMemoryStore())
		require.NoError(t, err)
		require.NotNil(t, m)

		token, err := m.Lookup(t.Context(), uuid.NewString())
		require.ErrorIs(t, err, ErrNotFound)
		require.Nil(t, token)
	})

	t.Run("expiration", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			m, err := NewManager[*TestData](UUIDGenerator, NewMemoryStore())
			require.NoError(t, err)
			require.NotNil(t, m)

			const ttl = 42 * time.Second
			expirationTime := time.Now().Add(ttl)
			token, err := m.Create(t.Context(), ttl, newTestData())
			require.NoError(t, err)
			require.NotNil(t, token)
			require.Equal(t, expirationTime, token.NotAfter)

			// Check that token can be read before NotAfter time
			time.Sleep(time.Until(expirationTime) - time.Nanosecond)
			require.True(t, time.Now().Before(expirationTime))
			beforeToken, err := m.Lookup(t.Context(), token.ID)
			require.NoError(t, err)
			require.Equal(t, token, beforeToken)

			// Check that token can still be read at exactly NotAfter time
			time.Sleep(time.Nanosecond)
			require.Equal(t, time.Now(), expirationTime)
			atToken, err := m.Lookup(t.Context(), token.ID)
			require.NoError(t, err)
			require.Equal(t, token, atToken)

			// Check that token cannot be read after NotAfter time
			time.Sleep(time.Nanosecond)
			require.True(t, time.Now().After(expirationTime))
			afterToken, err := m.Lookup(t.Context(), token.ID)
			require.ErrorIs(t, err, ErrExpiredToken)
			require.Nil(t, afterToken)
		})
	})

	t.Run("unmarshal error", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			raw := StoredToken{
				ID:       uuid.NewString(),
				Data:     []byte("random data that cannot be unmarshalled"),
				NotAfter: time.Now().Add(time.Second),
			}

			store := NewMemoryStore()
			err := store.Create(t.Context(), raw)
			require.NoError(t, err)

			m, err := NewManager[*TestData](UUIDGenerator, store)
			require.NoError(t, err)
			require.NotNil(t, m)
			token, err := m.Lookup(t.Context(), raw.ID)
			require.ErrorIs(t, err, ErrUnmarshalData)
			require.ErrorContains(t, err, "invalid data format")
			require.Nil(t, token)
		})
	})
}

func TestManager_Revoke(t *testing.T) {
	m, err := NewManager[*TestData](UUIDGenerator, NewMemoryStore())
	require.NoError(t, err)
	require.NotNil(t, m)

	t.Run("empty", func(t *testing.T) {
		err := m.Revoke(t.Context(), "")
		require.NoError(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		err := m.Revoke(t.Context(), uuid.NewString())
		require.NoError(t, err)
	})

	t.Run("normal", func(t *testing.T) {
		data := newTestData()
		token, err := m.Create(t.Context(), time.Second, data)
		require.NoError(t, err)
		require.NotNil(t, token)

		err = m.Revoke(t.Context(), token.ID)
		require.NoError(t, err)

		revokedToken, err := m.Lookup(t.Context(), token.ID)
		require.ErrorIs(t, err, ErrNotFound)
		require.Nil(t, revokedToken)
	})
}

func generateTokenIDEmpty() (string, error) {
	return "", nil
}

func generateTokenIDError() (string, error) {
	return "", errors.New("generateTokenIDError")
}

type TestData struct {
	Key   string
	Value string
}

func (td *TestData) MarshalBinary() ([]byte, error) {
	if strings.ContainsRune(td.Key, ':') {
		return nil, fmt.Errorf("invalid character ':' in key")
	}

	if strings.ContainsRune(td.Value, ':') {
		return nil, fmt.Errorf("invalid character ':' in value")
	}

	return []byte(td.Key + ":" + td.Value), nil
}

func (td *TestData) UnmarshalBinary(data []byte) error {
	parts := strings.Split(string(data), ":")
	if len(parts) != 2 {
		return fmt.Errorf("invalid data format")
	}

	td.Key = parts[0]
	td.Value = parts[1]

	return nil
}

func newTestData() *TestData {
	return &TestData{
		Key:   uuid.NewString(),
		Value: uuid.NewString(),
	}
}
