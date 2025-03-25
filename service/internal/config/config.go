package config

import (
	"errors"
	"fmt"

	"github.com/alexey-dobry/goodwords/internal/models"
	_ "github.com/alexey-dobry/goodwords/internal/validator"
	cfg "github.com/spf13/viper"
)

type Config struct {
	BadWords        []string              `validate:"required"`
	ListOfEndpoints []models.EndpointData `validate:"required"`
}

func ReadConfig() (*Config, error) {

	var configData Config

	cfg.SetConfigName("config")
	cfg.SetConfigType("toml")
	cfg.AddConfigPath("../")

	if err := cfg.ReadInConfig(); err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to read config: %s", err))
	}

	if err := cfg.UnmarshalKey("prohibited_words", &configData.BadWords); err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to write prohibied_words data from config: %s", err))
	}

	if err := cfg.UnmarshalKey("endpoints", &configData.ListOfEndpoints); err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to write endpoints data from config: %s", err))
	}

	return &configData, nil
}
