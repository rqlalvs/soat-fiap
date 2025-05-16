package handlers

import (
	"encoding/json"
	"net/http"
	"soat-fiap/internal/core/domain"
	"soat-fiap/internal/core/ports"

	"github.com/gorilla/mux"
)

type ProdutoHandler struct {
	produtoService ports.ProdutoService
}

func NovoProdutoHandler(produtoService ports.ProdutoService) *ProdutoHandler {
	return &ProdutoHandler{
		produtoService: produtoService,
	}
}

type CriarProdutoRequest struct {
	Nome      string           `json:"nome"`
	Descricao string           `json:"descricao"`
	Preco     float64          `json:"preco"`
	Categoria domain.Categoria `json:"categoria"`
}

func (h *ProdutoHandler) CriarProduto(w http.ResponseWriter, r *http.Request) {
	var req CriarProdutoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Erro ao decodificar requisição: "+err.Error(), http.StatusBadRequest)
		return
	}

	produto, err := h.produtoService.CriarProduto(r.Context(), req.Nome, req.Descricao, req.Preco, req.Categoria)
	if err != nil {
		http.Error(w, "Erro ao criar produto: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(produto)
}

func (h *ProdutoHandler) BuscarProdutoPorID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	produto, err := h.produtoService.BuscarProdutoPorID(r.Context(), id)
	if err != nil {
		http.Error(w, "Erro ao buscar produto: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if produto == nil {
		http.Error(w, "Produto não encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(produto)
}

func (h *ProdutoHandler) ListarProdutos(w http.ResponseWriter, r *http.Request) {
	categoria := r.URL.Query().Get("categoria")
	var produtos []*domain.Produto
	var err error

	if categoria != "" {
		produtos, err = h.produtoService.ListarProdutosPorCategoria(r.Context(), domain.Categoria(categoria))
	} else {
		produtos, err = h.produtoService.ListarProdutos(r.Context())
	}

	if err != nil {
		http.Error(w, "Erro ao listar produtos: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(produtos)
}

type AtualizarProdutoRequest struct {
	Nome       string           `json:"nome"`
	Descricao  string           `json:"descricao"`
	Preco      float64          `json:"preco"`
	Categoria  domain.Categoria `json:"categoria"`
	Disponivel bool             `json:"disponivel"`
}

func (h *ProdutoHandler) AtualizarProduto(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req AtualizarProdutoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Erro ao decodificar requisição: "+err.Error(), http.StatusBadRequest)
		return
	}

	produtoExistente, err := h.produtoService.BuscarProdutoPorID(r.Context(), id)
	if err != nil {
		http.Error(w, "Erro ao buscar produto: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if produtoExistente == nil {
		http.Error(w, "Produto não encontrado", http.StatusNotFound)
		return
	}

	produtoExistente.Nome = req.Nome
	produtoExistente.Descricao = req.Descricao
	produtoExistente.Preco = req.Preco
	produtoExistente.Categoria = req.Categoria
	produtoExistente.Disponivel = req.Disponivel

	err = h.produtoService.AtualizarProduto(r.Context(), produtoExistente)
	if err != nil {
		http.Error(w, "Erro ao atualizar produto: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(produtoExistente)
}

func (h *ProdutoHandler) DeletarProduto(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := h.produtoService.DeletarProduto(r.Context(), id)
	if err != nil {
		http.Error(w, "Erro ao deletar produto: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
