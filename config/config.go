package config

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"strconv"
)

type Configuration struct {
	Port        int
	PathProxy   string
	LoggerLevel string
}

func Load(path string) (*Configuration, error) {
	config := &Configuration{}

	if path != "" {
		configFile, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}

		if err := yaml.Unmarshal(configFile, config); err != nil {
			return nil, fmt.Errorf("error parsing YAML: %w", err)
		}
	}

	loadFromEnv(config)

	if err := config.Validate(); err != nil {
		return nil, err
	}
	return config, nil
}

func loadFromEnv(cfg *Configuration) {
	if cfg.Port == 0 {
		if envPort := os.Getenv("PORT"); envPort != "" {
			if port, err := strconv.Atoi(envPort); err == nil {
				cfg.Port = port
			}
		}
	}

	if cfg.PathProxy == "" {
		if envProxy := os.Getenv("PATH_PROXY"); envProxy != "" {
			cfg.PathProxy = envProxy
		} else {
			cfg.PathProxy = "https://jsonplaceholder.typicode.com/posts/"
		}
	}

	if cfg.LoggerLevel == "" {
		if envLogLevel := os.Getenv("LOGGER_LEVEL"); envLogLevel != "" {
			cfg.LoggerLevel = envLogLevel
		}
	}
}

func (cfg *Configuration) Validate() error {
	if cfg.Port <= 0 {
		return errors.New("port must be positive")
	}
	if cfg.PathProxy == "" {
		return errors.New("proxy path is required")
	}
	if cfg.LoggerLevel != "INFO" && cfg.LoggerLevel != "DEBUG" && cfg.LoggerLevel != "WARN" && cfg.LoggerLevel != "ERROR" {
		return errors.New("Invalid Level")
	}
	return nil
}
