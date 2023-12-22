package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"

	"github.com/serg-pe/signals/internal/config"
	"github.com/serg-pe/signals/internal/server"
	"github.com/serg-pe/signals/pkg/logger"
	"go.uber.org/zap"
)

const (
	cfgFileName = "config.toml"
)

func readConfigFile() config.AppConfig {
	curDir, err := os.Getwd()
	if err != nil {
		panic(fmt.Errorf("failed to get current working directory path: %s", err.Error()))
	}

	absPath := filepath.Join(curDir, cfgFileName)

	cfg, err := config.NewFromFile(absPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			createConfigFileAndExit(absPath)
		} else {
			panic(fmt.Errorf("failed to read config file status: %s", err.Error()))
		}
	}

	return cfg
}

func createConfigFileAndExit(path string) {
	fmt.Println("config file not found, creating config.toml...")
	err := config.NewBaseConfigFile(path)
	if err != nil {
		panic(fmt.Errorf("failed to init config file: %s", err))
	}
	fmt.Println("edit config.toml and restart server")
	fmt.Println("press any key to continue")
	os.Stdin.Read(nil)
	os.Exit(0)
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	defer func() {
		fmt.Println("shutdown signal")
		cancel()
	}()

	cfg := readConfigFile()

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

		err := server.Run(ctx)
		if err == http.ErrServerClosed {
			logger.Info("server stoped")
		} else {
			logger.Error("server stoped with error", zap.Error(err))
		}
	}()
	logger.Info("server successfully started", zap.String("ip", cfg.Ip), zap.Uint16("port", cfg.Port))

	go func() {
		<-sig
		logger.Info("stop server")
		server.Stop(ctx)
	}()

	wg.Wait()
}
