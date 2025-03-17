package config

import (
	"log"

	"github.com/alexey-dobry/goodwords/internal/models"
	cfg "github.com/spf13/viper"
)

type Config struct {
	BadWords        []string
	ListOfEndpoints []models.EndpointData
}

func ReadConfig() *Config {

	var configData Config

	cfg.SetConfigName("config")
	cfg.SetConfigType("toml")
	cfg.AddConfigPath("../")

	if err := cfg.ReadInConfig(); err != nil {
		log.Fatalf("failed to read config: %s", err)
	}

	if err := cfg.UnmarshalKey("prohibited_words", configData.BadWords); err != nil {
		log.Fatalf("Failed to write prohibied_words data from config: %s", err)
	}

	if err := cfg.UnmarshalKey("endpoints", configData.ListOfEndpoints); err != nil {
		log.Fatalf("Failed to write endpoints data from config: %s", err)
	}

	return &configData
}
