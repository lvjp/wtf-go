package token

import (
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInMemoryStore(t *testing.T) {
	t.Run("Create", func(t *testing.T) {
		t.Run("empty token ID", func(t *testing.T) {
			var memory memoryStore
			err := memory.Create(t.Context(), StoredToken{})
			require.ErrorIs(t, err, ErrEmptyID)
		})

		t.Run("sequence", func(t *testing.T) {
			var memory memoryStore

			tokenA := newRandomToken()
			tokenB := newRandomToken()

			require.NoError(t, memory.Create(t.Context(), tokenA))
			require.NoError(t, memory.Create(t.Context(), tokenB))

			require.Equal(t, 2, syncMapCount(&memory.tokens))
			requireSyncMapEqual(t, &memory.tokens, tokenA.ID, tokenA)
			requireSyncMapEqual(t, &memory.tokens, tokenB.ID, tokenB)
		})

		t.Run("collision", func(t *testing.T) {
			var memory memoryStore

			token := newRandomToken()

			firstErr := memory.Create(t.Context(), token)
			require.NoError(t, firstErr)

			secondErr := memory.Create(t.Context(), token)
			require.ErrorIs(t, secondErr, ErrAlreadyExists)
		})
	})

	t.Run("Read", func(t *testing.T) {
		t.Run("empty store", func(t *testing.T) {
			var memory memoryStore

			token, err := memory.Read(t.Context(), uuid.NewString())
			require.ErrorIs(t, err, ErrNotFound)
			require.Nil(t, token)
			require.Equal(t, 0, syncMapCount(&memory.tokens))
		})

		t.Run("different tokenID", func(t *testing.T) {
			var memory memoryStore

			tokenA := newRandomToken()
			tokenB := newRandomToken()
			require.NotEqual(t, tokenA.ID, tokenB.ID)

			storeErr := memory.Create(t.Context(), tokenA)
			require.NoError(t, storeErr)

			loadedToken, loadErr := memory.Read(t.Context(), tokenB.ID)
			require.ErrorIs(t, loadErr, ErrNotFound)
			require.Nil(t, loadedToken)
		})

		t.Run("existing tokenID", func(t *testing.T) {
			var memory memoryStore

			token := newRandomToken()

			storeErr := memory.Create(t.Context(), token)
			require.NoError(t, storeErr)

			loadedToken, loadErr := memory.Read(t.Context(), token.ID)
			require.NoError(t, loadErr)
			require.Equal(t, token, *loadedToken)
		})
	})

	t.Run("Delete", func(t *testing.T) {
		t.Run("empty store", func(t *testing.T) {
			var memory memoryStore

			err := memory.Delete(t.Context(), uuid.NewString())
			require.NoError(t, err)
			require.Equal(t, 0, syncMapCount(&memory.tokens))
		})

		t.Run("different tokenID", func(t *testing.T) {
			var memory memoryStore

			tokenA := newRandomToken()
			tokenB := newRandomToken()
			require.NotEqual(t, tokenA.ID, tokenB.ID)

			storeErr := memory.Create(t.Context(), tokenA)
			require.NoError(t, storeErr)

			deleteErr := memory.Delete(t.Context(), tokenB.ID)
			require.NoError(t, deleteErr)
		})

		t.Run("existing tokenID", func(t *testing.T) {
			var memory memoryStore

			token := newRandomToken()

			storeErr := memory.Create(t.Context(), token)
			require.NoError(t, storeErr)

			deleteErr := memory.Delete(t.Context(), token.ID)
			require.NoError(t, deleteErr)

			loadedToken, loadErr := memory.Read(t.Context(), token.ID)
			require.ErrorIs(t, loadErr, ErrNotFound)
			require.Nil(t, loadedToken)
		})
	})

	t.Run("concurrency", func(t *testing.T) {
		var wg sync.WaitGroup
		var memory memoryStore

		for range 10 {
			wg.Go(func() {
				for range 50 {
					token := newRandomToken()

					_, loadErr := memory.Read(t.Context(), token.ID)
					if !assert.ErrorIs(t, loadErr, ErrNotFound) {
						return
					}

					storeErr := memory.Create(t.Context(), token)
					if !assert.NoError(t, storeErr) {
						return
					}

					deleteErr := memory.Delete(t.Context(), token.ID)
					if !assert.NoError(t, deleteErr) {
						return
					}
					_, loadErr = memory.Read(t.Context(), token.ID)
					assert.ErrorIs(t, loadErr, ErrNotFound)
				}
			})
		}

		wg.Wait()
	})
}

func newRandomToken() StoredToken {
	data := uuid.New()
	return StoredToken{
		ID:       uuid.NewString(),
		Data:     data[:],
		NotAfter: time.Now().Add(time.Hour),
	}
}

func syncMapCount(m *sync.Map) int {
	count := 0
	m.Range(func(_, _ any) bool {
		count++
		return true
	})
	return count
}

func requireSyncMapEqual(t *testing.T, m *sync.Map, key string, expected StoredToken) {
	raw, exists := m.Load(key)
	require.True(t, exists, "expected key %q to exist in sync.Map", key)
	require.NotNil(t, raw, "expected value for key %q to be non-nil", key)

	actual, ok := raw.(*StoredToken)
	require.True(t, ok, "expected value for key %q to be of type *token.Token, not %T", key, actual)
	require.Equal(t, expected, *actual)
}
