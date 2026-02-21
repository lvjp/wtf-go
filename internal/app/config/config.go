package config

import (
	"bytes"
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-yaml"
)

type Config struct {
	Server Server `yaml:"server"`
	Log    Log    `yaml:"log"`
}

type Server struct {
	ListenAddress *string `yaml:"listen_address" validate:"omitempty,hostname_port"`
}

type Log struct {
	Level  string `validate:"omitempty,oneofci=debug info warn error fatal panic"`
	Format string `validate:"omitempty,oneof=console json"`
}

func LoadFromFile(path string) (*Config, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config file reading: %v", err)
	}

	dec := yaml.NewDecoder(
		bytes.NewReader(raw),
		yaml.Strict(),
		yaml.DisallowUnknownField(),
	)

	var ret Config
	ret.SetDefaults()

	if err := dec.Decode(&ret); err != nil {
		return nil, fmt.Errorf("config file parsing: %v", err)
	}

	validate := validator.New()
	if err := validate.Struct(&ret); err != nil {
		return nil, fmt.Errorf("config file validation: %v", err)
	}

	return &ret, nil
}

func (cfg *Config) SetDefaults() {
	if cfg.Server.ListenAddress == nil {
		defaultAddr := ":8080"
		cfg.Server.ListenAddress = &defaultAddr
	}
	if cfg.Log.Level == "" {
		cfg.Log.Level = "info"
	}
	if cfg.Log.Format == "" {
		cfg.Log.Format = "json"
	}
}
