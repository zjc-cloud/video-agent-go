package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	API      APIConfig
	Storage  StorageConfig
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

type ServerConfig struct {
	Host string
	Port int
}

type APIConfig struct {
	OpenAIKey string
}

type StorageConfig struct {
	Type   string
	Bucket string
	Region string
}

var AppConfig *Config

func Init() {
	// Load .env file if exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	port, _ := strconv.Atoi(getEnv("SERVER_PORT", "8080"))
	dbPort, _ := strconv.Atoi(getEnv("DB_PORT", "3306"))

	AppConfig = &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     dbPort,
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", "video_agent"),
		},
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
			Port: port,
		},
		API: APIConfig{
			OpenAIKey: getEnv("OPENAI_API_KEY", ""),
		},
		Storage: StorageConfig{
			Type:   getEnv("STORAGE_TYPE", "local"),
			Bucket: getEnv("CLOUD_BUCKET", ""),
			Region: getEnv("CLOUD_REGION", "us-west-2"),
		},
	}

	// Validate required config
	if AppConfig.API.OpenAIKey == "" {
		log.Fatal("OPENAI_API_KEY is required")
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
