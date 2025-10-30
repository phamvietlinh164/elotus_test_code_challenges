package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

var Cfg *Config

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

type Jwt struct {
	Secret string `yaml:"secret"`
	Ttl    int    `yaml:"ttl"`
}

type Config struct {
	App      App      `yaml:"app"`
	Database Database `yaml:"database"`
	Jwt      Jwt      `yaml:"jwt"`
}

func Load() error {
	data, err := os.ReadFile("internal/config/config.yaml")
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(data, &Cfg); err != nil {
		return err
	}

	if Cfg.App.Port == "" {
		Cfg.App.Port = ":8080"
	}

	fmt.Printf("Loaded config: %+v\n", Cfg)
	return nil
}
