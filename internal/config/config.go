package config

import (
	"github.com/BurntSushi/toml"
	"github.com/serg-pe/signals/pkg/logger"
)

type AppConfig struct {
	logger.LoggerConfig `toml:"logger"`
	ServerConfig        `toml:"server"`
}

type ServerConfig struct {
	Ip   string `toml:"ip"`
	Port uint16 `toml:"port"`
}

func NewFromFile(path string) (AppConfig, error) {
	cfg := AppConfig{}
	_, err := toml.DecodeFile(path, &cfg)
	return cfg, err
}
