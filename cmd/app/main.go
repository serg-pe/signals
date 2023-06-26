package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/serg-pe/signals/internal/config"
	"github.com/serg-pe/signals/internal/server"
	"github.com/serg-pe/signals/pkg/logger"
	"go.uber.org/zap"
)

const (
	cfgFileName = "config.toml"
)

func main() {
	curDir, err := os.Getwd()
	if err != nil {
		panic(fmt.Errorf("failed to get working directory path: %s", err.Error()))
	}

	cfg, err := config.NewFromFile(filepath.Join(curDir, cfgFileName))
	if err != nil {
		panic(fmt.Errorf("failed to open config file: %s", err))
	}

	logger, err := logger.New(cfg.LoggerConfig)
	if err != nil {
		panic(fmt.Errorf("failed to init logger: %s", err.Error()))
	}

	server, err := server.New(cfg.ServerConfig, logger)
	if err != nil {
		logger.Fatal("failed to start server", zap.Error(err))
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()

		err := server.Run()
		if err != nil {
			logger.Fatal("server stoped", zap.Error(err))
		}
	}()
	logger.Info("server successfully started", zap.String("ip", cfg.Ip), zap.Uint16("port", cfg.Port))

	wg.Wait()

}
