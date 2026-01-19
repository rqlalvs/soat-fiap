package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	ServerPort          string
	DBHost              string
	DBPort              string
	DBUser              string
	DBPassword          string
	DBName              string
	LogLevel            string
	SwaggerEnable       bool
	ClienteServiceURL   string
	ProdutoServiceURL   string
	PedidoServiceURL    string
	PagamentoServiceURL string
}

func LoadConfig() *Config {
	serverPort := getEnvOrPanic("SERVER_PORT", "8080")
	dbHost := getEnv("DB_HOST", "")
	dbPort := getEnv("DB_PORT", "")
	dbUser := getEnv("DB_USER", "")
	dbPassword := getEnv("DB_PASSWORD", "")
	dbName := getEnv("DB_NAME", "")
	logLevel := getEnv("LOG_LEVEL", "info")
	swaggerEnable := getEnvAsBool("SWAGGER_ENABLE", true)
	clienteServiceURL := getEnv("CLIENTE_SERVICE_URL", "http://localhost:8081")
	produtoServiceURL := getEnv("PRODUTO_SERVICE_URL", "http://localhost:8082")
	pedidoServiceURL := getEnv("PEDIDO_SERVICE_URL", "http://localhost:8083")
	pagamentoServiceURL := getEnv("PAGAMENTO_SERVICE_URL", "http://localhost:8084")

	return &Config{
		ServerPort:          serverPort,
		DBHost:              dbHost,
		DBPort:              dbPort,
		DBUser:              dbUser,
		DBPassword:          dbPassword,
		DBName:              dbName,
		LogLevel:            logLevel,
		SwaggerEnable:       swaggerEnable,
		ClienteServiceURL:   clienteServiceURL,
		ProdutoServiceURL:   produtoServiceURL,
		PedidoServiceURL:    pedidoServiceURL,
		PagamentoServiceURL: pagamentoServiceURL,
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
