package token

import (
	"errors"
	"testing"
	"testing/iotest"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestUUIDGenerator(t *testing.T) {
	t.Run("MultipleCalls", func(t *testing.T) {
		const callCount = 10
		results := make(map[string]bool, callCount)

		for i := range callCount {
			uuid, err := UUIDGenerator()
			require.NoError(t, err)
			require.NotEmpty(t, uuid)
			require.NotContains(t, results, uuid, "UUID collision with %q on round #%d", uuid, i+1)
		}
	})

	t.Run("error", func(t *testing.T) {
		readErr := errors.New("expected read error")
		expected := "UUID generation error: " + readErr.Error()

		uuid.SetRand(iotest.ErrReader(readErr))
		t.Cleanup(func() {
			uuid.SetRand(nil)
		})

		uuid, err := UUIDGenerator()
		require.ErrorContains(t, err, expected)
		require.Empty(t, uuid)
	})
}
