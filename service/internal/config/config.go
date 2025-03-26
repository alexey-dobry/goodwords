package config

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	cfg "github.com/spf13/viper"
)

type Config struct {
	BadWords        []string             `validate:"required"`
	ListOfEndpoints []ConfigEndpointData `validate:"required,dive"`
}

type ConfigEndpointData struct {
	URL        string `validate:"required" mapstructure:"url"`
	MaxTime    int    `validate:"required" mapstructure:"max_time"`
	MaxRetries int    `validate:"required" mapstructure:"max_retries"`
	ReturnData string `validate:"required" mapstructure:"return_data"`
}

func ReadConfig() (*Config, error) {

	var configData Config

	var V = validator.New()

	cfg.SetConfigName("config")
	cfg.SetConfigType("toml")
	cfg.AddConfigPath("./config")

	if err := cfg.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("Failed to read config: %s", err)
	}

	if err := cfg.UnmarshalKey("bad_words", &configData.BadWords); err != nil {
		return nil, fmt.Errorf("Failed to write prohibied_words data from config: %s", err)
	}

	if err := cfg.UnmarshalKey("list_of_endpoints", &configData.ListOfEndpoints); err != nil {
		return nil, fmt.Errorf("Failed to write endpoints data from config: %s", err)
	}

	if err := V.Struct(configData); err != nil {
		return nil, fmt.Errorf("Config data validatiopn failed: %s", err)
	}

	return &configData, nil
}
