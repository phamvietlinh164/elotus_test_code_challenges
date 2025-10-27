package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type App struct {
	Port string `yaml:"port"`
	Mode string `yaml:"mode"`
	Log  string `yaml:"log"`
}

type Database struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Name     string `yaml:"name"`
}

type Config struct {
	App      App      `yaml:"app"`
	Database Database `yaml:"database"`
}

func Load() (*Config, error) {
	data, err := os.ReadFile("internal/config/config.yaml")
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if cfg.App.Port == "" {
		cfg.App.Port = ":8080"
	}

	fmt.Printf("Loaded config: %+v\n", cfg)
	return &cfg, nil
}
