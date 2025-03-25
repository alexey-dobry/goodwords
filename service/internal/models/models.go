package models

type EndpointData struct {
	URL        string `validate:"required" mapstructure:"url"`
	MaxTime    int    `validate:"required" mapstructure:"max_time"`
	MaxRetries int    `validate:"required" mapstructure:"max_retries"`
	ReturnData string `validate:"required" mapstructure:"return_data"`
}
