package models

type EndpointData struct {
	URL        string `mapstructure:"url"`
	MaxTime    int    `mapstructure:"max_time"`
	MaxRetries int    `mapstructure:"max_retries"`
	ReturnData string `mapstructure:"return_data"`
}
