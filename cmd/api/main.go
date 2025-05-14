package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	config "soat-fiap/configs"
	"soat-fiap/internal/adapters/primary/handlers"
	"soat-fiap/internal/adapters/secondary/repositories"
	"soat-fiap/internal/core/services"
	"soat-fiap/internal/routes"
	sqlite "soat-fiap/pkg/database"

	"github.com/gorilla/mux"
)

const (
	AppVersion = "1.0.0"
)

// @title           API SOAT-FIAP
// @version         1.0
// @description     API de exemplo com arquitetura hexagonal
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.example.com/support
// @contact.email  support@example.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1

func main() {
	config := config.LoadConfig()

	ensureDirExists(filepath.Dir(config.DatabasePath))

	db, err := sqlite.ConectarDB(config.DatabasePath)
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}
	defer db.Close()

	clienteRepository := repositories.NovoSQLiteClienteRepository(db)

	clienteService := services.NovoClienteService(clienteRepository)

	clienteHandler := handlers.NovoClienteHandler(clienteService)
	healthHandler := handlers.NovoHealthHandler(AppVersion)

	router := mux.NewRouter()
	routes.ConfigurarRotas(router, clienteHandler, healthHandler)

	server := &http.Server{
		Addr:         ":" + config.ServerPort,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Servidor iniciado na porta %s", config.ServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Erro ao iniciar servidor: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Desligando servidor...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Erro ao desligar servidor: %v", err)
	}

	log.Println("Servidor encerrado")
}

func ensureDirExists(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			log.Fatalf("Erro ao criar diretório %s: %v", dir, err)
		}
	}
}
