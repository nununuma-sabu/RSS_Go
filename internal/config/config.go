package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Feeds []Feed `yaml:"feeds"`
}

type Feed struct {
	URL      string `yaml:"url"`
	Category string `yaml:"category"`
}

func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
