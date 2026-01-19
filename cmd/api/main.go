package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "soat-fiap/docs" // swagger docs

	config "soat-fiap/configs"
	"soat-fiap/internal/adapters/primary/handlers"
	"soat-fiap/internal/core/services"
	"soat-fiap/internal/routes"
	"soat-fiap/pkg/clients"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
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
	cfg := config.LoadConfig()

	// Inicializar clients HTTP para os serviços externos
	clienteClient := clients.NovoClienteClient(cfg.ClienteServiceURL)
	produtoClient := clients.NovoProdutoClient(cfg.ProdutoServiceURL)
	pedidoClient := clients.NovoPedidoClient(cfg.PedidoServiceURL)
	pagamentoClient := clients.NovoPagamentoClient(cfg.PagamentoServiceURL)

	// Inicializar services com os clients
	clienteService := services.NovoClienteService(clienteClient)
	produtoService := services.NovoProdutoService(produtoClient)
	pedidoService := services.NovoPedidoService(pedidoClient, produtoClient, clienteClient, pagamentoClient)
	pagamentoService := services.NovoPagamentoService(pagamentoClient, pedidoClient)

	clienteHandler := handlers.NovoClienteHandler(clienteService)
	produtoHandler := handlers.NovoProdutoHandler(produtoService)
	pedidoHandler := handlers.NovoPedidoHandler(pedidoService, pagamentoService)
	pagamentoHandler := handlers.NovoPagamentoHandler(pagamentoService)
	healthHandler := handlers.NovoHealthHandler(AppVersion)

	router := mux.NewRouter()
	routes.ConfigurarRotas(router, clienteHandler, produtoHandler, pedidoHandler, pagamentoHandler, healthHandler)

	// Swagger
	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("doc.json"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
	))

	server := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Servidor iniciado na porta %s", cfg.ServerPort)
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
