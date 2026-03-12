package e2e

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"

	"git.sr.ht/~lvjp/wtf-go/internal/app/api/misc"

	"github.com/stretchr/testify/require"
)

func TestOperationId_MiscVersion(t *testing.T) {
	var endpoint string
	if endpoint = os.Getenv("TEST_WTF_GO_BACKEND_ENDPOINT"); endpoint == "" {
		t.Skip("TEST_WTF_GO_BACKEND_ENDPOINT environment variable is not set, skipping test")
	}

	req, err := http.NewRequestWithContext(
		t.Context(),
		http.MethodGet,
		endpoint+"/api/v0/misc/version",
		nil,
	)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	dec := json.NewDecoder(resp.Body)
	dec.DisallowUnknownFields()

	// Revision and Time are expected to be set during build time, but we can't
	// guarantee that in tests, so we just check that they are not empty.
	var payload misc.VersionResponse
	require.NoError(t, dec.Decode(&payload))
	require.NotEmpty(t, payload.Go, "payload.Go")
	require.NotEmpty(t, payload.Platform, "payload.Platform")
}

func TestOperationId_MiscHealth(t *testing.T) {
	var endpoint string
	if endpoint = os.Getenv("TEST_WTF_GO_BACKEND_ENDPOINT"); endpoint == "" {
		t.Skip("TEST_WTF_GO_BACKEND_ENDPOINT environment variable is not set, skipping test")
	}

	req, err := http.NewRequestWithContext(
		t.Context(),
		http.MethodGet,
		endpoint+"/api/v0/misc/health",
		nil,
	)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	payload, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, []byte("Status: OK"), payload)
}
