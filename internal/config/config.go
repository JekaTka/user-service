package config

import "github.com/spf13/viper"

type Environment string

const (
	ProdEnvironment Environment = "production"
)

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variable.
type Config struct {
	Environment       Environment `mapstructure:"ENVIRONMENT"`
	HTTPServerAddress string      `mapstructure:"HTTP_SERVER_ADDRESS"`
	GRPCServerAddress string      `mapstructure:"GRPC_SERVER_ADDRESS"`
	DBSource          string      `mapstructure:"DB_SOURCE"`
	MigrationURL      string      `mapstructure:"MIGRATION_URL"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (*Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var config Config
	err = viper.Unmarshal(&config)

	return &config, nil
}
