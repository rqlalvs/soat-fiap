package routes

import (
	"net/http"

	"soat-fiap/internal/adapters/primary/handlers"

	"github.com/gorilla/mux"
)

func ConfigurarRotas(r *mux.Router, clienteHandler *handlers.ClienteHandler, produtoHandler *handlers.ProdutoHandler, pedidoHandler *handlers.PedidoHandler, healthHandler *handlers.HealthHandler) {
	api := r.PathPrefix("/api/v1").Subrouter()

	api.HandleFunc("/health", healthHandler.HealthCheck).Methods(http.MethodGet)

	api.HandleFunc("/clientes", clienteHandler.CriarCliente).Methods(http.MethodPost)
	api.HandleFunc("/clientes", clienteHandler.ListarClientes).Methods(http.MethodGet)
	api.HandleFunc("/clientes/cpf/{cpf}", clienteHandler.BuscarClientePorCPF).Methods(http.MethodGet)
	api.HandleFunc("/clientes/{id}", clienteHandler.BuscarClientePorID).Methods(http.MethodGet)
	api.HandleFunc("/clientes/{id}", clienteHandler.AtualizarCliente).Methods(http.MethodPut)
	api.HandleFunc("/clientes/{id}", clienteHandler.DeletarCliente).Methods(http.MethodDelete)

	api.HandleFunc("/produtos", produtoHandler.CriarProduto).Methods(http.MethodPost)
	api.HandleFunc("/produtos", produtoHandler.ListarProdutos).Methods(http.MethodGet)
	api.HandleFunc("/produtos/{id}", produtoHandler.BuscarProdutoPorID).Methods(http.MethodGet)
	api.HandleFunc("/produtos/{id}", produtoHandler.AtualizarProduto).Methods(http.MethodPut)
	api.HandleFunc("/produtos/{id}", produtoHandler.DeletarProduto).Methods(http.MethodDelete)

	api.HandleFunc("/checkout", pedidoHandler.FakeCheckout).Methods(http.MethodPost)
	api.HandleFunc("/pedidos", pedidoHandler.FakeCheckout).Methods(http.MethodPost)
	api.HandleFunc("/pedidos", pedidoHandler.ListarPedidos).Methods(http.MethodGet)
	api.HandleFunc("/pedidos/{id}", pedidoHandler.BuscarPedidoPorID).Methods(http.MethodGet)
	api.HandleFunc("/pedidos/{id}/status", pedidoHandler.AtualizarStatusPedido).Methods(http.MethodPatch)
}
