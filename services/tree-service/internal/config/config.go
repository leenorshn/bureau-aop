package config

import (
	"os"
)

type Config struct {
	MongoURI   string
	MongoDBName string
	RedisURL   string
	Port       string
}

func Load() *Config {
	return &Config{
		MongoURI:   getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDBName: getEnv("MONGO_DB_NAME", "bureau"),
		RedisURL:   getEnv("REDIS_URL", ""),
		Port:       getEnv("TREE_SERVICE_PORT", "8082"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}


