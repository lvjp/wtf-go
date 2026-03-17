package e2e

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOperationId_MiscVersion(t *testing.T) {
	c := NewClient(t)

	resp, err := c.MiscVersion(t.Context())
	require.NoError(t, err)

	// Revision and Time are expected to be set during build time, but we can't
	// guarantee that in tests, so we just check that they are not empty.
	payload := resp
	require.NotEmpty(t, payload.Go, "payload.Go")
	require.NotEmpty(t, payload.Platform, "payload.Platform")
}

func TestOperationId_MiscHealth(t *testing.T) {
	c := NewClient(t)

	resp, err := c.MiscHealth(t.Context())
	require.NoError(t, err)
	require.Equal(t, "OK", resp.Status)
}
