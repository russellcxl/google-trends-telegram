package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

var config *Config

type Config struct {
	GoogleClient GoogleClient `yaml:"google_client"`
	Redis        Redis        `yaml:"redis"`
}

type GoogleClient struct {
	DefaultParams [][]string `yaml:"default_params"`
	Daily         Daily      `yaml:"daily"`
}

type Daily struct {
	ListCount              int `yaml:"list_count"`
	RefreshIntervalMinutes int `yaml:"refresh_interval_mins"`
}

type Redis struct {
	RetryCount int `yaml:"retry_count"`
}

func GetConfig() *Config {
	if config != nil {
		return config
	}
	configPath, found := os.LookupEnv("CONFIG_PATH")
	if !found {
		log.Fatalf("failed to find config path in env")
	}
	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("failed to read config file: %v", err)
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("error unmarshaling YAML: %v", err)
	}
	return config
}
