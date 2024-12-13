package config

import (
	"os"
	"path/filepath"

	"github.com/go-playground/validator"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

const (
	configPathEnvName     = "SPEC_FILE"
	configFileNameDefault = "./.config.yml"
)

type (
	specWithMetaConfig struct {
		Spec Config `yaml:"spec"`
	}

	// Config ...
	Config struct {
		// General ...
		General General `yaml:"general" validate:"required"`
		LLM     LLM     `yaml:"llm" validate:"required"`
	}

	// General config
	General struct {
		// APIKey ...
		APIKey string `yaml:"api_key" validate:"required"`
		// ShutdownWaitSec is the number of secs the server will wait
		// before shutting down after it receives an exit signal
		ShutdownWaitSec int `yaml:"graceful_shutdown_wait_time_sec" validate:"required"`
		// LogLevel ...
		LogLevel string `yaml:"log_level" validate:"required"`
	}

	// LLM config
	LLM struct {
		Provider string `yaml:"provider" validate:"required"`
		Keywords string `yaml:"keywords"`
		Gemini   Gemini `yaml:"gemini" validate:"required"`
	}

	// Gemini config under llm
	Gemini struct {
		Model string `yaml:"model" validate:"required"`
		// MaxRequestsPerMinute ...
		MaxRequestsPerMinute int `yaml:"max_requests_per_minute" validate:"required"`
	}
)

// Load loads all configurations in to a new Config struct
func Load() (*Config, error) {
	configFilePath := os.Getenv(configPathEnvName)
	if configFilePath == "" {
		workingDir, err := os.Getwd()
		if err != nil {
			return nil, errors.Wrap(err, "unable to get current working directory")
		}
		configFilePath = filepath.Join(workingDir, configFileNameDefault)
	}
	// reading app file config
	configFile, err := os.Open(configFilePath)
	if err != nil {
		return nil, errors.Wrap(err, "cannot open config file")
	}

	var spec specWithMetaConfig
	err = yaml.NewDecoder(configFile).Decode(&spec)
	if err != nil {
		return nil, errors.Wrap(err, "cannot unmarshal config data")
	}

	config := &spec.Spec

	// validating app file configs
	v := validator.New()
	err = v.Struct(config)
	if err != nil {
		return nil, errors.Wrap(err, "config file is not valid")
	}
	return config, nil
}
