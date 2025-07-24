package handlers

import (
	"encoding/json"
	"net/http"
	"regexp"

	"soat-fiap/internal/core/ports"

	"github.com/gorilla/mux"
)

type ClienteHandler struct {
	clienteService ports.ClienteService
}

func NovoClienteHandler(clienteService ports.ClienteService) *ClienteHandler {
	return &ClienteHandler{
		clienteService: clienteService,
	}
}

type CriarClienteRequest struct {
	Nome     string `json:"nome"`
	CPF      string `json:"cpf"`
	Email    string `json:"email"`
	Telefone string `json:"telefone"`
}

// CriarCliente cria um novo cliente.
// @Summary Criar cliente
// @Tags clientes
// @Accept json
// @Produce json
// @Param cliente body CriarClienteRequest true "Dados do cliente"
// @Success 201 {object} domain.Cliente
// @Failure 400 {string} string "Erro ao criar cliente"
// @Router /clientes [post]
func (h *ClienteHandler) CriarCliente(w http.ResponseWriter, r *http.Request) {
	var req CriarClienteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Erro ao decodificar requisição: "+err.Error(), http.StatusBadRequest)
		return
	}

	cliente, err := h.clienteService.CriarCliente(r.Context(), req.Nome, req.CPF, req.Email, req.Telefone)
	if err != nil {
		http.Error(w, "Erro ao criar cliente: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(cliente)
}

// BuscarClientePorID retorna um cliente pelo ID.
// @Summary Buscar cliente por ID
// @Tags clientes
// @Produce json
// @Param id path string true "ID do cliente"
// @Success 200 {object} domain.Cliente
// @Failure 404 {string} string "Cliente não encontrado"
// @Router /clientes/{id} [get]
func (h *ClienteHandler) BuscarClientePorID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	cliente, err := h.clienteService.BuscarClientePorID(r.Context(), id)
	if err != nil {
		http.Error(w, "Erro ao buscar cliente: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if cliente == nil {
		http.Error(w, "Cliente não encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cliente)
}

// BuscarClientePorCPF retorna um cliente pelo CPF.
// @Summary Buscar cliente por CPF
// @Tags clientes
// @Produce json
// @Param cpf path string true "CPF do cliente"
// @Success 200 {object} domain.Cliente
// @Failure 404 {string} string "Cliente não encontrado"
// @Router /clientes/cpf/{cpf} [get]
func (h *ClienteHandler) BuscarClientePorCPF(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cpf := vars["cpf"]

	re := regexp.MustCompile(`[^0-9]`)
	cpf = re.ReplaceAllString(cpf, "")

	cliente, err := h.clienteService.BuscarClientePorCPF(r.Context(), cpf)
	if err != nil {
		http.Error(w, "Erro ao buscar cliente: "+err.Error(), http.StatusBadRequest)
		return
	}

	if cliente == nil {
		http.Error(w, "Cliente não encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cliente)
}

// ListarClientes retorna todos os clientes.
// @Summary Listar clientes
// @Tags clientes
// @Produce json
// @Success 200 {array} domain.Cliente
// @Failure 500 {string} string "Erro ao listar clientes"
// @Router /clientes [get]
func (h *ClienteHandler) ListarClientes(w http.ResponseWriter, r *http.Request) {
	clientes, err := h.clienteService.ListarClientes(r.Context())
	if err != nil {
		http.Error(w, "Erro ao listar clientes: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(clientes)
}

type AtualizarClienteRequest struct {
	Nome     string `json:"nome"`
	CPF      string `json:"cpf"`
	Email    string `json:"email"`
	Telefone string `json:"telefone"`
}

// AtualizarCliente atualiza um cliente existente.
// @Summary Atualizar cliente
// @Tags clientes
// @Accept json
// @Produce json
// @Param id path string true "ID do cliente"
// @Param cliente body AtualizarClienteRequest true "Dados do cliente"
// @Success 200 {object} domain.Cliente
// @Failure 400 {string} string "Erro ao atualizar cliente"
// @Router /clientes/{id} [put]
func (h *ClienteHandler) AtualizarCliente(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req AtualizarClienteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Erro ao decodificar requisição: "+err.Error(), http.StatusBadRequest)
		return
	}

	clienteExistente, err := h.clienteService.BuscarClientePorID(r.Context(), id)
	if err != nil {
		http.Error(w, "Erro ao buscar cliente: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if clienteExistente == nil {
		http.Error(w, "Cliente não encontrado", http.StatusNotFound)
		return
	}
	clienteExistente.Nome = req.Nome
	clienteExistente.CPF = req.CPF
	clienteExistente.Email = req.Email
	clienteExistente.Telefone = req.Telefone

	err = h.clienteService.AtualizarCliente(r.Context(), clienteExistente)
	if err != nil {
		http.Error(w, "Erro ao atualizar cliente: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(clienteExistente)
}

// DeletarCliente remove um cliente pelo ID.
// @Summary Deletar cliente
// @Tags clientes
// @Produce json
// @Param id path string true "ID do cliente"
// @Success 204 {string} string "Cliente deletado"
// @Failure 500 {string} string "Erro ao deletar cliente"
// @Router /clientes/{id} [delete]
func (h *ClienteHandler) DeletarCliente(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := h.clienteService.DeletarCliente(r.Context(), id)
	if err != nil {
		http.Error(w, "Erro ao deletar cliente: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
