package e2e

import (
	"os"
	"testing"

	"github.com/lvjp/wtf-go/pkg/api"
)

func NewClient(t *testing.T, optsFn ...api.OptionsFn) *api.Client {
	var endpoint string
	if endpoint = os.Getenv("TEST_WTF_GO_BACKEND_ENDPOINT"); endpoint == "" {
		t.Skip("TEST_WTF_GO_BACKEND_ENDPOINT environment variable is not set, skipping test")
	}

	return api.NewClient(endpoint+"/api/v0", optsFn...)
}
