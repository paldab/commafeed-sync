package config

import (
	"fmt"
	"os"
)

const (
	defaultConfigPath = "/commafeed/feeds.yaml"
)

func getEnvOrDefault(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func GetFeedsConfigPath() string {
	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		configPath = defaultConfigPath
	}

	return configPath
}

func GetCFUrl() (string, error) {
	url := os.Getenv("COMMAFEED_URL")

	if url == "" {
		return "", fmt.Errorf("COMMAFEED_URL is not set")
	}

	return url, nil
}

func GetCredentials() (string, string) {
	const (
		defaultUsername = "admin"
		defaultPassword = "admin"
	)

	user := getEnvOrDefault("COMMAFEED_USERNAME", defaultUsername)
	pass := getEnvOrDefault("COMMAFEED_PASSWORD", defaultPassword)

	return user, pass
}
