package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/goccy/go-yaml"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestConfigBuilder_empty(t *testing.T) {
	cfg, usedPath, err := New()
	require.NoError(t, err)
	require.Empty(t, usedPath)
	require.NotNil(t, cfg)
	require.NotZero(t, cfg)
}

func TestConfigBuilder_WithDefaults(t *testing.T) {
	cfg, usedPath, buildErr := New(WithDefaults())
	require.NoError(t, buildErr)
	require.Empty(t, usedPath)
	require.NotNil(t, cfg)

	validationErr := cfg.Validate()
	require.NoError(t, validationErr)

	var zero Config
	require.NotEqual(t, &zero, cfg)
}

func TestConfigBuilder_WithLogLevel(t *testing.T) {
	const logLevel = "my-log-level"

	cfg, usedPath, buildErr := New(WithLogLevel(logLevel))
	require.NoError(t, buildErr)
	require.Empty(t, usedPath)
	require.NotNil(t, cfg)
	require.Equal(t, logLevel, cfg.Log.Level)
}

func TestWithEnvVars(t *testing.T) {
	const envKey = "WTF_GO_LOG.LEVEL"
	const envValue = "my-env-log-level"

	t.Setenv(envKey, envValue)

	cfg, usedPath, buildErr := New(WithEnvVars())
	require.NoError(t, buildErr)
	require.Empty(t, usedPath)
	require.NotNil(t, cfg)
	require.Equal(t, envValue, cfg.Log.Level)
}

func TestConfigBuilder_WithConfigFile(t *testing.T) {
	expectedUsedPath, expectedConfig := generateDefaultYaml(t)

	actualConfig, actualUsedPath, err := New(WithConfigFile(expectedUsedPath))
	require.NoError(t, err)
	require.Equal(t, expectedUsedPath, actualUsedPath)
	require.Equal(t, expectedConfig, actualConfig)
}

func generateDefaultYaml(t *testing.T) (string, *Config) {
	path := filepath.Join(t.TempDir(), uuid.NewString()+".yaml")

	// TODO: generate this instead of hardcoding it here. The dummy usage of testing/quick produces
	// invalid yaml string with control charactect, so we need to find a way to generate valid yaml
	// content.
	config := Config{
		Server: Server{
			ListenAddress: "generated:server.listen_address",
		},
		Log: Log{
			Level:  "generated:log.level",
			Format: "generated:log.format",
		},
	}

	raw, err := yaml.Marshal(config)
	require.NoError(t, err)

	err = os.WriteFile(path, raw, 0600)
	require.NoError(t, err)

	return path, &config
}
