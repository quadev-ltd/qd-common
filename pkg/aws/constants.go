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

// Application ID and Region
const (
	ApplicationID = "wwrvnjk"
	Region        = "eu-west-1"
)

// LocalConfig is the configuration for the local/development environment
var LocalConfig = AppConfig{
	Region:                 Region,
	ApplicationID:          ApplicationID,
	EnvironmentID:          "gl16bm7",
	ConfigurationProfileID: "du1vndf",
}

// DevConfig is the configuration for the development environment
var DevConfig = AppConfig{
	Region:                 Region,
	ApplicationID:          ApplicationID,
	EnvironmentID:          "gl16bm7",
	ConfigurationProfileID: "du1vndf",
}

// ProdConfig is the configuration for the production environment
var ProdConfig = AppConfig{
	Region:                 Region,
	ApplicationID:          ApplicationID,
	EnvironmentID:          "9t9usa0",
	ConfigurationProfileID: "2yy0ge2",
}
