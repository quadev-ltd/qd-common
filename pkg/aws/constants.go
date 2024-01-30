package aws

// AWSAppConfig is the configuration of the AWS AppConfig
type AWSAppConfig struct {
	Region                 string
	ApplicationID          string
	EnvironmentID          string
	ConfigurationProfileID string
}

// DevConfig is the configuration for the development environment
var DevConfig = AWSAppConfig{
	Region:                 "us-east-1",
	ApplicationID:          "222hh4s",
	EnvironmentID:          "i0lf5xh",
	ConfigurationProfileID: "mfgfaot",
}
