package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type App struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

type Http struct {
	Port string `yaml:"port"`
}

type DbConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Database string `yaml:"database"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Schema   string `yaml:"schema"`
}

type Log struct {
	Level string `yaml:"log_level"`
}

type Config struct {
	App  `yaml:"app"`
	Http `yaml:"http"`
	Db   DbConfig `yaml:"db"`
	Log  `yaml:"logger"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	err := cleanenv.ReadConfig("./config/config.yml", cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}
	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, err
}
