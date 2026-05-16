package e2e

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/lvjp/wtf-go/pkg/api"
)

func TestOperationId_AuthTokenCreate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		c := NewClient(t)

		resp, err := c.AuthTokenCreate(t.Context(), "admin@localhost")
		require.NoError(t, err)
		require.NotEmpty(t, resp.ID)
		require.NotZero(t, resp.NotAfter)
	})

	t.Run("invalid subject", func(t *testing.T) {
		c := NewClient(t)

		resp, err := c.AuthTokenCreate(t.Context(), "unknown@localhost")

		var httpErr *api.HTTPRequestError
		require.ErrorAs(t, err, &httpErr)
		require.Equal(t, 403, httpErr.HTTPResponse.StatusCode)
		require.Nil(t, resp)
	})
}

func TestOperationId_AuthTokenRevoke(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		c := NewClient(t)

		_, err := c.AuthTokenCreate(t.Context(), "admin@localhost")
		require.NoError(t, err)

		err = c.AuthTokenRevoke(t.Context())
		require.NoError(t, err)
	})

	t.Run("twice", func(t *testing.T) {
		c := NewClient(t)

		resp, err := c.AuthTokenCreate(t.Context(), "admin@localhost")
		require.NoError(t, err)

		err = c.AuthTokenRevoke(t.Context())
		require.NoError(t, err)

		err = c.AuthTokenRevoke(t.Context(), api.WithAuthToken(resp.ID))
		var httpErr *api.HTTPRequestError
		require.ErrorAs(t, err, &httpErr)
		require.Equal(t, 403, httpErr.HTTPResponse.StatusCode)
	})

	t.Run("without token", func(t *testing.T) {
		c := NewClient(t)

		err := c.AuthTokenRevoke(t.Context())
		var httpErr *api.HTTPRequestError
		require.ErrorAs(t, err, &httpErr)
		require.Equal(t, 401, httpErr.HTTPResponse.StatusCode)
	})

	t.Run("invalid token", func(t *testing.T) {
		c := NewClient(t, api.WithAuthToken("invalid-token"))

		err := c.AuthTokenRevoke(t.Context())
		var httpErr *api.HTTPRequestError
		require.ErrorAs(t, err, &httpErr)
		require.Equal(t, 403, httpErr.HTTPResponse.StatusCode)
	})
}
