package config

import "os"

type Config struct {
	TreeServiceURL string
	Port           string
}

func Load() *Config {
	return &Config{
		TreeServiceURL: getEnv("TREE_SERVICE_URL", "http://localhost:8082"),
		Port:           getEnv("GATEWAY_PORT", "8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}



