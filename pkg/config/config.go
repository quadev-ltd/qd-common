package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// AppConfig represents the common structure of application configurations.
type AppConfig struct {
	Environment string
	Verbose     bool
}

// GetEnvironment returns the environment of the application.
func GetEnvironment() string {
	env := os.Getenv(AppEnvironmentKey)
	if env == "" {
		env = LocalEnvironment
	}
	return env
}

// GetVerbose returns the verbose mode of the application.
func GetVerbose() bool {
	return os.Getenv(VerboseKey) == "true"
}

// SetupConfig loads the configuration from the given path and environment.
func SetupConfig(path, env string) (*viper.Viper, error) {
	v := viper.New()

	// Load base configuration file first
	v.SetConfigType("yml")
	v.AddConfigPath(path)
	v.SetConfigName("config.template")
	err := v.MergeInConfig()
	if err != nil {
		log.Warn().Msgf("Error loading base configuration file: %v", err)
	}

	// Load environment-specific configuration file
	v.SetConfigName(fmt.Sprintf("config.%s", env))
	err = v.MergeInConfig()
	if err != nil {
		log.Warn().Msgf("Error loading environment-specific configuration file: %v", err)
	}

	// Bind environment variables
	prefix := fmt.Sprintf("%s_ENV", strings.ToUpper(env))
	v.AutomaticEnv()
	v.SetEnvPrefix(prefix)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	return v, nil
}
