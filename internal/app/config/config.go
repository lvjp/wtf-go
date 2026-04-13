package config

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
)

type Config struct {
	Server Server
	Log    Log
}

type Server struct {
	ListenAddress string `yaml:"listen_address" validate:"omitempty,hostname_port"`
}

type Log struct {
	Level  string `validate:"omitempty,oneofci=debug info warn error fatal panic"`
	Format string `validate:"omitempty,oneof=console json"`
}

func (cfg *Config) Validate() error {
	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		return fmt.Errorf("config validation: %v", err)
	}

	return nil
}

type NewOption func(*factory)

func WithDefaults() NewOption {
	return func(f *factory) {
		f.viper.SetDefault("server.listen_address", ":8080")
		f.viper.SetDefault("log.level", "info")
		f.viper.SetDefault("log.format", "json")
	}
}

func WithLogLevel(level string) NewOption {
	return func(f *factory) {
		f.viper.Set("log.level", level)
	}
}

func WithEnvVars() NewOption {
	return func(f *factory) {
		f.viper.SetEnvPrefix("wtf_go")
		f.viper.AutomaticEnv()
	}
}

func WithConfigLookup() NewOption {
	return func(f *factory) {
		f.withConfigLookup = true

		f.viper.SetConfigName("config")
		f.viper.SetConfigType("yaml")
		f.viper.AddConfigPath("/etc/wtf-go")
		f.viper.AddConfigPath("$HOME/.config/wtf-go")
	}
}

func WithConfigFile(path string) NewOption {
	return func(f *factory) {
		f.withConfigFile = true

		f.viper.SetConfigFile(path)
	}
}

func New(opts ...NewOption) (loaded *Config, configFileUsed string, err error) {
	f := factory{
		viper: viper.New(),
	}

	for _, opt := range opts {
		opt(&f)
	}

	return f.create()
}

type factory struct {
	viper *viper.Viper

	withConfigLookup bool
	withConfigFile   bool
}

func (f *factory) create() (*Config, string, error) {
	switch {
	case f.withConfigLookup && f.withConfigFile:
		return nil, "", fmt.Errorf("cannot use both config lookup and config file options")

	case f.withConfigFile:
		if err := f.viper.ReadInConfig(); err != nil {
			return nil, f.viper.ConfigFileUsed(), fmt.Errorf("config file loading: %v", err)
		}

	case f.withConfigLookup:
		if err := f.viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				return nil, "", fmt.Errorf("config file lookup: %v", err)
			}
		}
	}

	// Verify that all config keys are mapped to the struct fields,
	// to avoid silent errors due to typos.
	if err := f.viper.UnmarshalExact(&Config{}, useYAMLTagname); err != nil {
		return nil, f.viper.ConfigFileUsed(), fmt.Errorf("config build struct: %v", err)
	}

	var cfg Config

	// UnmarshalExact doesn't work with env vars, so we need to set the fields manually here.
	cfg.Server.ListenAddress = f.viper.GetString("server.listen_address")
	cfg.Log.Level = f.viper.GetString("log.level")
	cfg.Log.Format = f.viper.GetString("log.format")

	return &cfg, f.viper.ConfigFileUsed(), nil
}

// useYAMLTagname evit the repition of the same field name in both yaml tag and mapstructure tag,
// and also avoid the risk of them being different due to typos.
func useYAMLTagname(dc *mapstructure.DecoderConfig) {
	dc.TagName = "yaml"
}
