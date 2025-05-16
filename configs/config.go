package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	ServerPort    string
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	LogLevel      string
	SwaggerEnable bool
}

func LoadConfig() *Config {
	serverPort := getEnvOrPanic("SERVER_PORT", "8080")
	dbHost := getEnvOrPanic("DB_HOST", "")
	dbPort := getEnvOrPanic("DB_PORT", "")
	dbUser := getEnvOrPanic("DB_USER", "")
	dbPassword := getEnvOrPanic("DB_PASSWORD", "")
	dbName := getEnvOrPanic("DB_NAME", "")
	logLevel := getEnv("LOG_LEVEL", "info")
	swaggerEnable := getEnvAsBool("SWAGGER_ENABLE", true)

	return &Config{
		ServerPort:    serverPort,
		DBHost:        dbHost,
		DBPort:        dbPort,
		DBUser:        dbUser,
		DBPassword:    dbPassword,
		DBName:        dbName,
		LogLevel:      logLevel,
		SwaggerEnable: swaggerEnable,
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvOrPanic(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	if defaultValue != "" {
		return defaultValue
	}
	panic(fmt.Sprintf("Variável de ambiente obrigatória não definida: %s", key))
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, strconv.FormatBool(defaultValue))
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}
