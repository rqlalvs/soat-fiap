package config

import (
	"os"
	"strconv"
)

type Config struct {
	ServerPort    string
	DatabasePath  string
	LogLevel      string
	SwaggerEnable bool
}

func LoadConfig() *Config {
	serverPort := getEnv("SERVER_PORT", "8080")
	databasePath := getEnv("DATABASE_PATH", "./data/db.sqlite")
	logLevel := getEnv("LOG_LEVEL", "info")
	swaggerEnable := getEnvAsBool("SWAGGER_ENABLE", true)

	return &Config{
		ServerPort:    serverPort,
		DatabasePath:  databasePath,
		LogLevel:      logLevel,
		SwaggerEnable: swaggerEnable,
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, strconv.FormatBool(defaultValue))
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}
