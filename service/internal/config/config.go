package config

import (
	//"log"
	//"strings"

	"github.com/alexey-dobry/goodwords/internal/models"
	//"github.com/spf13/viper"
)

type Config struct {
	BadWords        []string
	ListOfEndpoints []models.EndpointData
}

func ReadConfig() *Config {
	// TODO
	return nil
}
