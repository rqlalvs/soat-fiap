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
