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

// CriarProduto cria um novo produto.
// @Summary Criar produto
// @Tags produtos
// @Accept json
// @Produce json
// @Param produto body CriarProdutoRequest true "Dados do produto"
// @Success 201 {object} domain.Produto
// @Failure 400 {string} string "Erro ao criar produto"
// @Router /produtos [post]
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

// ListarProdutos retorna todos os produtos, podendo filtrar por categoria.
// @Summary Listar produtos
// @Tags produtos
// @Produce json
// @Param categoria query string false "Categoria do produto"
// @Success 200 {array} domain.Produto
// @Failure 500 {string} string "Erro ao listar produtos"
// @Router /produtos [get]
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

// AtualizarProduto atualiza um produto existente.
// @Summary Atualizar produto
// @Tags produtos
// @Accept json
// @Produce json
// @Param id path string true "ID do produto"
// @Param produto body AtualizarProdutoRequest true "Dados do produto"
// @Success 200 {object} domain.Produto
// @Failure 400 {string} string "Erro ao atualizar produto"
// @Router /produtos/{id} [put]
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

// DeletarProduto remove um produto pelo ID.
// @Summary Deletar produto
// @Tags produtos
// @Produce json
// @Param id path string true "ID do produto"
// @Success 204 {string} string "Produto deletado"
// @Failure 500 {string} string "Erro ao deletar produto"
// @Router /produtos/{id} [delete]
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
