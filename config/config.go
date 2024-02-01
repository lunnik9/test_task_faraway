package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	RedisConfig RedisConfig
	MongoConfig MongoConfig
	ServerPort  string
}

type RedisConfig struct {
	Address  string
	Password string
	DB       int
}

type MongoConfig struct {
	URI string
}

func LoadConfig() AppConfig {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, relying on environment variables.")
	}

	port := getEnv("SERVER_PORT", "8080")

	return AppConfig{
		RedisConfig: RedisConfig{
			Address:  getEnv("REDIS_ADDRESS", "localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		MongoConfig: MongoConfig{
			URI: getEnv("MONGO_URI", "localhost:27017"),
		},
		ServerPort: port,
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}
