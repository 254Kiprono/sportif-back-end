package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost           string
	DBUser           string
	DBPassword       string
	DBName           string
	DBPort           string
	SSLMode          string
	JWTSecret        string
	Port             string
	B2Endpoint       string
	B2Region         string
	B2BucketName     string
	B2KeyID          string
	B2ApplicationKey string
	RedisHost        string
	RedisPort        string
	RedisPassword    string
	RedisDB          string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	return &Config{
		DBHost:           getEnv("DB_HOST", "localhost"),
		DBUser:           getEnv("DB_USER", "root"),
		DBPassword:       getEnv("DB_PASSWORD", "password"),
		DBName:           getEnv("DB_NAME", "webuye_sportif"),
		DBPort:           getEnv("DB_PORT", "3306"),
		SSLMode:          getEnv("DB_SSLMODE", "disable"),
		JWTSecret:        getEnv("JWT_SECRET", "super-secret-key"),
		Port:             getEnv("PORT", "8002"),
		B2Endpoint:       getEnv("B2_ENDPOINT", ""),
		B2Region:         getEnv("B2_REGION", ""),
		B2BucketName:     getEnv("B2_BUCKET_NAME", ""),
		B2KeyID:          getEnv("B2_KEY_ID", ""),
		B2ApplicationKey: getEnv("B2_APPLICATION_KEY", ""),
		RedisHost:        getEnv("REDIS_HOST", "localhost"),
		RedisPort:        getEnv("REDIS_PORT", "6379"),
		RedisPassword:    getEnv("REDIS_PASSWORD", ""),
		RedisDB:          getEnv("REDIS_DB", "0"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
