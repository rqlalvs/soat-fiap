package routes

import (
	"net/http"

	"soat-fiap/internal/adapters/primary/handlers"

	"github.com/gorilla/mux"
)

func ConfigurarRotas(r *mux.Router, clienteHandler *handlers.ClienteHandler, healthHandler *handlers.HealthHandler) {
	api := r.PathPrefix("/api/v1").Subrouter()

	api.HandleFunc("/health", healthHandler.HealthCheck).Methods(http.MethodGet)

	api.HandleFunc("/clientes", clienteHandler.CriarCliente).Methods(http.MethodPost)
	api.HandleFunc("/clientes", clienteHandler.ListarClientes).Methods(http.MethodGet)
	api.HandleFunc("/clientes/cpf/{cpf}", clienteHandler.BuscarClientePorCPF).Methods(http.MethodGet)
	api.HandleFunc("/clientes/{id}", clienteHandler.BuscarClientePorID).Methods(http.MethodGet)
	api.HandleFunc("/clientes/{id}", clienteHandler.AtualizarCliente).Methods(http.MethodPut)
	api.HandleFunc("/clientes/{id}", clienteHandler.DeletarCliente).Methods(http.MethodDelete)
}
