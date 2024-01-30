package aws

type AWSAppConfig struct {
	Region                 string
	ApplicationID          string
	EnvironmentID          string
	ConfigurationProfileID string
}

var DevConfig = AWSAppConfig{
	Region:                 "us-east-1",
	ApplicationID:          "222hh4s",
	EnvironmentID:          "i0lf5xh",
	ConfigurationProfileID: "mfgfaot",
}
