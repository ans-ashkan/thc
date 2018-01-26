package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config stores API configurations
type config struct {
	APIKey      string
	APISecret   string
	Token       string
	TokenSecret string
	Owner       string
	OwnerID     string
}

// CheckRequiredConfigs returns error if required settings are missing
func (config *config) CheckRequiredConfigs() error {
	if config.APIKey == "" {
		return fmt.Errorf("APIKey is null or empty")
	}

	if config.APISecret == "" {
		return fmt.Errorf("APISecret is null or empty")
	}

	if config.Token == "" {
		return fmt.Errorf("Token is null or empty")
	}

	if config.TokenSecret == "" {
		return fmt.Errorf("TokenSecret is null or empty")
	}

	if config.Owner == "" {
		return fmt.Errorf("TokenSecret is null or empty")
	}

	if config.OwnerID == "" {
		return fmt.Errorf("TokenSecret is null or empty")
	}

	return nil
}

var configCache *config

// GetConfig read config or returns cached config
func GetConfig() *config {
	if configCache == nil {
		viper.SetConfigName("config")
		viper.AddConfigPath(".")
		err := viper.ReadInConfig()
		if err != nil {
			panic(fmt.Errorf("Fatal error config file: %s", err))
		}

		configCache = &config{}
		viper.Unmarshal(configCache)
		if err := configCache.CheckRequiredConfigs(); err != nil {
			panic(fmt.Errorf("Fatal config error. %s", err))
		}
	}

	return configCache
}
