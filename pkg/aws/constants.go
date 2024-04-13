package aws

// Config is the configuration of the AWS
type Config struct {
	Key    string `mapstructure:"key"`
	Secret string `mapstructure:"secret"`
}

// AppConfig is the configuration of the AWS AppConfig
type AppConfig struct {
	Region                 string
	ApplicationID          string
	EnvironmentID          string
	ConfigurationProfileID string
}

// LocalConfig is the configuration for the development environment
var LocalConfig = AppConfig{
	Region:                 "eu-west-1",
	ApplicationID:          "wwrvnjk",
	EnvironmentID:          "iv285bh",
	ConfigurationProfileID: "9va7r42",
}

// DevConfig is the configuration for the development environment
var DevConfig = AppConfig{
	Region:                 "eu-west-1",
	ApplicationID:          "wwrvnjk",
	EnvironmentID:          "iv285bh",
	ConfigurationProfileID: "9va7r42",
}
