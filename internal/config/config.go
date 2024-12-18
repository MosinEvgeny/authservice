package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
	DatabaseURL   string `mapstructure:"DATABASE_URL"`
	JWTSecret     string `mapstructure:"JWT_SECRET"`
	Email         string `mapstructure:"EMAIL"`
	EmailPassword string `mapstructure:"EMAIL_PASSWORD"`
}

func LoadConfig() (*Config, error) {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}
