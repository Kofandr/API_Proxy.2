package main

import (
	"context"
	"github.com/Kofandr/API_Proxy.2/config"
	"github.com/Kofandr/API_Proxy.2/internal/logger"
	"github.com/Kofandr/API_Proxy.2/internal/server"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading Config.env file: %s", err)
	}

	pathConfig := os.Getenv("CONFIG_PATH")
	cfg, err := config.Load(pathConfig)
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	logger := logger.New(cfg.LoggerLevel)

	mainServer := server.New(logger, cfg)

	go func() {
		if err := mainServer.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server crash")

		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	<-signalChan

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	logger.Info("Shutting down...")
	if err := mainServer.Shutdown(ctx); err != nil {
		logger.Error("Shutdown failed", "error", err)
	} else {
		logger.Info("Server stopped")
	}

}
