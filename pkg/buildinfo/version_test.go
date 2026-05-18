package buildinfo

import (
	"runtime/debug"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuildInfo(t *testing.T) {
	testCases := []struct {
		name         string
		buildInfo    BuildInfo
		rawBuildInfo debug.BuildInfo
		line         string
	}{
		{
			name: "normal",
			buildInfo: BuildInfo{
				Revision:     "5a1b17fe93cbb941c954108336b33149d17cf398",
				RevisionTime: "2026-02-19T23:14:00Z",
				Modified:     false,
				GoVersion:    "go1.26.0",
				GoOS:         "linux",
				GoArch:       "386",
			},
			rawBuildInfo: debug.BuildInfo{
				GoVersion: "go1.26.0",
				Settings: []debug.BuildSetting{
					{Key: keyRevision, Value: "5a1b17fe93cbb941c954108336b33149d17cf398"},
					{Key: keyRevisionTime, Value: "2026-02-19T23:14:00Z"},
					{Key: keyModified, Value: "false"},
					{Key: keyGoOS, Value: "linux"},
					{Key: keyGoArch, Value: "386"},
				},
			},
			line: "5a1b17fe93cbb941c954108336b33149d17cf398 2026-02-19T23:14:00Z go1.26.0 linux/386",
		},
		{
			name: "modified",
			buildInfo: BuildInfo{
				Revision:     "c3c77a198d92135318aafabc5b710c5860f916eb",
				RevisionTime: "2026-05-17T20:44:02Z",
				Modified:     true,
				GoVersion:    "go1.26.3",
				GoOS:         "js",
				GoArch:       "wasm",
			},
			rawBuildInfo: debug.BuildInfo{
				GoVersion: "go1.26.3",
				Settings: []debug.BuildSetting{
					{Key: keyRevision, Value: "c3c77a198d92135318aafabc5b710c5860f916eb"},
					{Key: keyRevisionTime, Value: "2026-05-17T20:44:02Z"},
					{Key: keyModified, Value: "true"},
					{Key: keyGoOS, Value: "js"},
					{Key: keyGoArch, Value: "wasm"},
				},
			},
			line: "c3c77a198d92135318aafabc5b710c5860f916eb 2026-05-17T20:44:02Z go1.26.3 js/wasm (modified)",
		},
	}

	t.Run("String", func(t *testing.T) {
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				require.Equal(t, tc.line, tc.buildInfo.String())
			})
		}
	})

	t.Run("Load", func(t *testing.T) {
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				var actual BuildInfo
				actual.Load(&tc.rawBuildInfo)
				require.Equal(t, tc.buildInfo, actual)
			})
		}
	})

	t.Run("defaults", func(t *testing.T) {
		bi := BuildInfo{}
		bi.setDefaults()

		require.Equal(t, "- - - -/-", bi.String())
	})
}

func TestGet(t *testing.T) {
	require.Equal(t, buildInfo, Get())
}
