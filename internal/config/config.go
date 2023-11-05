package config

import (
	"os"

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

func NewBaseConfigFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	toml.NewEncoder(file).Encode(AppConfig{
		LoggerConfig: logger.LoggerConfig{
			Level: string(logger.LogLevelRelease),
		},
		ServerConfig: ServerConfig{
			Ip:   "127.0.0.1",
			Port: 8000,
		},
	})
	return nil
}
