package models

import _ "github.com/alexey-dobry/goodwords/internal/validator"

type EndpointData struct {
	URL        string `mapstructure:"url" validate:"required"`
	MaxTime    int    `mapstructure:"max_time" validate:"required"`
	MaxRetries int    `mapstructure:"max_retries" validate:"required"`
	ReturnData string `mapstructure:"return_data" validate:"required"`
}
