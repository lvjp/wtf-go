package token

import (
	"encoding"
	"testing"
	"testing/synctest"
	"time"

	"github.com/stretchr/testify/require"
)

func TestToken_IsExpired(t *testing.T) {
	testCases := []struct {
		name      string
		add       time.Duration
		isExpired bool
	}{
		{
			name:      "past",
			add:       -time.Nanosecond,
			isExpired: true,
		},
		{
			name:      "now",
			add:       0,
			isExpired: false,
		},
		{
			name:      "future",
			add:       time.Nanosecond,
			isExpired: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var token Token[interface {
				encoding.BinaryMarshaler
				encoding.BinaryUnmarshaler
			}]

			synctest.Test(t, func(t *testing.T) {
				token.NotAfter = time.Now().Add(tc.add)
				require.Equal(t, tc.isExpired, token.IsExpired())
			})
		})
	}
}
