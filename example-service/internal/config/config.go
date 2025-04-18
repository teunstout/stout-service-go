package config

import "os"

// Config holds the application configuration.
type Config struct {
	AppName string
}

// LoadConfig loads configuration from environment variables.
func LoadConfig() *Config {
	return &Config{
		AppName: os.Getenv("APP_NAME"),
	}
}
