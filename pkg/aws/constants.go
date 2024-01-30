package aws

// AWSConfig is the configuration of the AWS
type AWSConfig struct {
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
	Region:                 "us-east-1",
	ApplicationID:          "222hh4s",
	EnvironmentID:          "i0lf5xh",
	ConfigurationProfileID: "mfgfaot",
}

// DevConfig is the configuration for the development environment
var DevConfig = AppConfig{
	Region:                 "us-east-1",
	ApplicationID:          "222hh4s",
	EnvironmentID:          "i0lf5xh",
	ConfigurationProfileID: "mfgfaot",
}
