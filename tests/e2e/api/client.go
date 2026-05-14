package e2e

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/lvjp/wtf-go/pkg/api"
)

func NewClient(t *testing.T) *api.Client {
	var endpoint string
	if endpoint = os.Getenv("TEST_WTF_GO_BACKEND_ENDPOINT"); endpoint == "" {
		t.Skip("TEST_WTF_GO_BACKEND_ENDPOINT environment variable is not set, skipping test")
	}

	c, err := api.NewClient(endpoint + "/api/v0")
	require.NoError(t, err)

	return c
}
