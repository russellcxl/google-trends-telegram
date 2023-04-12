package utils

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

var config *Config

type Config struct {
	GoogleClient GoogleClient `yaml:"google_client"`
}

type GoogleClient struct {
	DefaultParams [][]string `yaml:"default_params"`
	Daily         Daily      `yaml:"daily"`
}

type Daily struct {
	ListCount int `yaml:"list_count"`
}

func GetConfig(path string) *Config {
	if config != nil {
		return config
	}
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("failed to read config file: %v", err)
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("error unmarshaling YAML: %v", err)
	}
	return config
}
