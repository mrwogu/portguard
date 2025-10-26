package main

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	defaultConfigPath = "/etc/portguard/config.yaml"
	defaultListenPort = "8888"
)

func loadConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	if cfg.Server.Port == "" {
		cfg.Server.Port = defaultListenPort
	}
	if cfg.Server.Timeout == 0 {
		cfg.Server.Timeout = 2 * time.Second
	}

	return &cfg, nil
}
