package config

import "os"

type Config struct {
	AppName string
}

func LoadConfig() *Config {
	return &Config{
		AppName: os.Getenv("APP_NAME"),
	}
}
