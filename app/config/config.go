package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost        string
	DBUser        string
	DBPassword    string
	DBName        string
	DBPort        string
	SSLMode       string
	JWTSecret     string
	Port          string
	R2Endpoint    string
	R2AccessKeyID string
	R2SecretKey   string
	R2BucketName  string
	R2PublicURL   string
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	return &Config{
		DBHost:        getEnv("DB_HOST", "localhost"),
		DBUser:        getEnv("DB_USER", "root"),
		DBPassword:    getEnv("DB_PASSWORD", "password"),
		DBName:        getEnv("DB_NAME", "webuye_sportif"),
		DBPort:        getEnv("DB_PORT", "3306"),
		SSLMode:       getEnv("DB_SSLMODE", "disable"),
		JWTSecret:     getEnv("JWT_SECRET", "super-secret-key"),
		Port:          getEnv("PORT", "8002"),
		R2Endpoint:    getEnv("R2_ENDPOINT", ""),
		R2AccessKeyID: getEnv("R2_ACCESS_KEY_ID", ""),
		R2SecretKey:   getEnv("R2_SECRET_ACCESS_KEY", ""),
		R2BucketName:  getEnv("R2_BUCKET_NAME", ""),
		R2PublicURL:   getEnv("R2_PUBLIC_URL", ""),
		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       getEnv("REDIS_DB", "0"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
