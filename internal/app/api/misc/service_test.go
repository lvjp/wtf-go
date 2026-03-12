package misc

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHealth(t *testing.T) {
	svc := NewService()

	resp, err := svc.Health(t.Context())
	require.NoError(t, err)
	require.Equal(t, "Status: OK", resp)
}

func TestVersion(t *testing.T) {
	svc := NewService()

	// Revision and Time are expected to be set during build time, but we can't
	// guarantee that in tests, so we just check that they are not empty.
	resp, err := svc.Version(t.Context())
	require.NoError(t, err)
	require.NotNil(t, resp, "resp should not be nil")
	require.NotEmpty(t, resp.Go, "resp.Go")
	require.NotEmpty(t, resp.Platform, "resp.Platform")
}
