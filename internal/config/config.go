package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	Port       string
}

func LoadConfig() *Config {
	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return &Config{
		DBHost:     getEnv("DB_HOST", DBHostDef),
		DBPort:     getEnv("DB_PORT", DBPortDef),
		DBUser:     getEnv("DB_USER", DBUserDef),
		DBPassword: getEnv("DB_PASSWORD", DBPasswordDef),
		DBName:     getEnv("DB_NAME", DBNameDef),
		Port:       getEnv("PORT", PortDef),
	}
}
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
