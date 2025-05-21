package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	config "soat-fiap/configs"
	"soat-fiap/internal/adapters/primary/handlers"
	"soat-fiap/internal/adapters/secondary/repositories"
	"soat-fiap/internal/core/services"
	"soat-fiap/internal/routes"
	mysql "soat-fiap/pkg/database"

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

	db, err := mysql.ConectarMySQL(
		config.DBHost,
		config.DBPort,
		config.DBUser,
		config.DBPassword,
		config.DBName,
	)
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}
	defer db.Close()

	clienteRepository := repositories.NovoClienteRepository(db)
	produtoRepository := repositories.NovoProdutoRepository(db)
	pedidoRepository := repositories.NovoPedidoRepository(db)

	clienteService := services.NovoClienteService(clienteRepository)
	produtoService := services.NovoProdutoService(produtoRepository)
	pedidoService := services.NovoPedidoService(pedidoRepository, produtoRepository)

	clienteHandler := handlers.NovoClienteHandler(clienteService)
	produtoHandler := handlers.NovoProdutoHandler(produtoService)
	pedidoHandler := handlers.NovoPedidoHandler(pedidoService)
	healthHandler := handlers.NovoHealthHandler(AppVersion)

	router := mux.NewRouter()
	routes.ConfigurarRotas(router, clienteHandler, produtoHandler, pedidoHandler, healthHandler)

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
