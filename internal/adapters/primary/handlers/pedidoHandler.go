package handlers

import (
	"encoding/json"
	"net/http"
	"soat-fiap/internal/core/domain"
	"soat-fiap/internal/core/ports"

	"github.com/gorilla/mux"
)

type PedidoHandler struct {
	pedidoService ports.PedidoService
}

func NovoPedidoHandler(pedidoService ports.PedidoService) *PedidoHandler {
	return &PedidoHandler{
		pedidoService: pedidoService,
	}
}

type CriarPedidoRequest struct {
	ClienteID *string                  `json:"cliente_id,omitempty"`
	Itens     []CriarItemPedidoRequest `json:"itens"`
}

type CriarItemPedidoRequest struct {
	ProdutoID  string `json:"produto_id"`
	Quantidade int    `json:"quantidade"`
	Observacao string `json:"observacao,omitempty"`
}

type AtualizarStatusRequest struct {
	Status domain.StatusPedido `json:"status"`
}

func (h *PedidoHandler) FakeCheckout(w http.ResponseWriter, r *http.Request) {
	var req CriarPedidoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Erro ao decodificar requisição: "+err.Error(), http.StatusBadRequest)
		return
	}

	var itens []domain.ItemPedido
	for _, item := range req.Itens {
		itens = append(itens, domain.ItemPedido{
			ProdutoID:  item.ProdutoID,
			Quantidade: item.Quantidade,
			Observacao: item.Observacao,
		})
	}

	pedido, err := h.pedidoService.CriarPedido(r.Context(), req.ClienteID, itens)
	if err != nil {
		http.Error(w, "Erro ao criar pedido: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Pedido criado e enviado para fila com sucesso",
		"pedido":  pedido,
	})
}

func (h *PedidoHandler) ListarPedidos(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	clienteID := r.URL.Query().Get("cliente_id")

	var pedidos []*domain.Pedido
	var err error

	if status != "" {
		pedidos, err = h.pedidoService.ListarPedidosPorStatus(r.Context(), domain.StatusPedido(status))
	} else if clienteID != "" {
		pedidos, err = h.pedidoService.ListarPedidosPorCliente(r.Context(), clienteID)
	} else {
		pedidos, err = h.pedidoService.ListarPedidos(r.Context())
	}

	if err != nil {
		http.Error(w, "Erro ao listar pedidos: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pedidos)
}

func (h *PedidoHandler) BuscarPedidoPorID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	pedido, err := h.pedidoService.BuscarPedidoPorID(r.Context(), id)
	if err != nil {
		http.Error(w, "Erro ao buscar pedido: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if pedido == nil {
		http.Error(w, "Pedido não encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pedido)
}

func (h *PedidoHandler) AtualizarStatusPedido(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req AtualizarStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Erro ao decodificar requisição: "+err.Error(), http.StatusBadRequest)
		return
	}

	err := h.pedidoService.AtualizarStatusPedido(r.Context(), id, req.Status)
	if err != nil {
		http.Error(w, "Erro ao atualizar status do pedido: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Status do pedido atualizado com sucesso",
	})
}
