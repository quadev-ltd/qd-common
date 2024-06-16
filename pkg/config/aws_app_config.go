package config

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/appconfigdata"
	"github.com/spf13/viper"

	awsConstants "github.com/quadev-ltd/qd-common/pkg/aws"
)

// Address is the address of a service
type Address struct {
	Host string
	Port string
}

// Config is the configuration of the application
type Config struct {
	TLSEnabled                bool    `mapstructure:"tls_enabled"`
	EmailVerificationEndpoint string  `mapstructure:"email_verification_endpoint"`
	GatewayService            Address `mapstructure:"gateway_service"`
	EmailService              Address `mapstructure:"email_service"`
	AuthenticationService     Address `mapstructure:"authentication_service"`
	VisualizationService      Address `mapstructure:"visualization_service"`
}

// Parameters is the parameters needed to load the configuration
type Parameters struct {
	Region                 string
	AWSKey                 string
	AWSSecret              string
	ApplicationID          string
	EnvironmentID          string
	ConfigurationProfileID string
}

// Load loads the configuration from AWS AppConfig
func (config *Config) Load(env, awsKey, awsSecret string) error {
	var awsDetails awsConstants.AppConfig
	switch env {
	case LocalEnvironment:
		awsDetails = awsConstants.LocalConfig
	case DevelopmentEnvironment:
		awsDetails = awsConstants.DevConfig
	case ProductionEnvironment:
		awsDetails = awsConstants.ProdConfig
	default:
		awsDetails = awsConstants.DevConfig
	}

	configParameters := Parameters{
		Region:                 awsDetails.Region,
		AWSKey:                 awsKey,
		AWSSecret:              awsSecret,
		ApplicationID:          awsDetails.ApplicationID,
		EnvironmentID:          awsDetails.EnvironmentID,
		ConfigurationProfileID: awsDetails.ConfigurationProfileID,
	}
	sess, err := session.NewSession(
		&aws.Config{
			Region: aws.String(configParameters.Region),
			Credentials: credentials.NewStaticCredentials(
				configParameters.AWSKey,
				configParameters.AWSSecret,
				"",
			),
		},
	)
	if err != nil {
		return err
	}

	svc := appconfigdata.New(sess)

	startSessionOutput, err := svc.StartConfigurationSession(
		&appconfigdata.StartConfigurationSessionInput{
			ApplicationIdentifier:          aws.String(configParameters.ApplicationID),
			EnvironmentIdentifier:          aws.String(configParameters.EnvironmentID),
			ConfigurationProfileIdentifier: aws.String(configParameters.ConfigurationProfileID),
		},
	)
	if err != nil {
		return err
	}

	latestConfigOutput, err := svc.GetLatestConfiguration(&appconfigdata.GetLatestConfigurationInput{
		ConfigurationToken: startSessionOutput.InitialConfigurationToken,
	})

	if err != nil {
		return err
	}

	// Process configuration data (e.g., set environment variables)
	configData := latestConfigOutput.Configuration
	viper.SetConfigType("yaml")
	if err := viper.ReadConfig(strings.NewReader(string(configData))); err != nil {
		return fmt.Errorf("Error reading YAML content into viper: %v", err)
	}
	if err := viper.Unmarshal(&config); err != nil {
		return fmt.Errorf("Error unmarshaling YAML content: %v", err)
	}
	return nil
}
