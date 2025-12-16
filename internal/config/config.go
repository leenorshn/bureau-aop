package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoURI             string
	MongoDBName          string
	JWTSecret            string
	JWTRefreshSecret     string
	JWTAccessExp         time.Duration
	JWTRefreshExp        time.Duration
	AdminSeedEmail       string
	AdminSeedPassword    string
	AppPort              string
	AppEnv               string
	BinaryThreshold      float64
	BinaryCommissionRate float64
	DefaultProductPrice  float64
	// Nouveaux paramètres pour l'algorithme binaire amélioré
	BinaryCycleValue      float64
	BinaryDailyCycleLimit int
	BinaryWeeklyCycleLimit int
	BinaryMinVolumePerLeg float64
}

func Load() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		MongoURI:             getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDBName:          getEnv("MONGO_DB_NAME", "mlm_db"),
		JWTSecret:            getEnv("JWT_SECRET", "your-secret-key"),
		JWTRefreshSecret:     getEnv("JWT_REFRESH_SECRET", "your-refresh-secret-key"),
		JWTAccessExp:         getDurationEnv("JWT_ACCESS_EXP", 15*time.Minute),
		JWTRefreshExp:        getDurationEnv("JWT_REFRESH_EXP", 7*24*time.Hour),
		AdminSeedEmail:       getEnv("ADMIN_SEED_EMAIL", "admin@example.com"),
		AdminSeedPassword:    getEnv("ADMIN_SEED_PASSWORD", "admin123"),
		AppPort:              getEnv("APP_PORT", "4000"),
		AppEnv:               getEnv("APP_ENV", "development"),
		BinaryThreshold:      getFloatEnv("BINARY_THRESHOLD", 100.0),
		BinaryCommissionRate: getFloatEnv("BINARY_COMMISSION_RATE", 0.1),
		DefaultProductPrice:  getFloatEnv("DEFAULT_PRODUCT_PRICE", 50.0),
		// Nouveaux paramètres pour l'algorithme binaire amélioré
		BinaryCycleValue:      getFloatEnv("BINARY_CYCLE_VALUE", 20.0),
		BinaryDailyCycleLimit: getIntEnv("BINARY_DAILY_CYCLE_LIMIT", 4),
		BinaryWeeklyCycleLimit: getIntEnv("BINARY_WEEKLY_CYCLE_LIMIT", 0),
		BinaryMinVolumePerLeg: getFloatEnv("BINARY_MIN_VOLUME_PER_LEG", 1.0),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getFloatEnv(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
