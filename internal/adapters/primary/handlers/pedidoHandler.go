package handlers

import (
	"encoding/json"
	"net/http"
	"soat-fiap/internal/adapters/primary/dtos"
	"soat-fiap/internal/adapters/primary/presenters"
	"soat-fiap/internal/core/domain"
	"soat-fiap/internal/core/ports"
	"sort"
	"time"

	"github.com/gorilla/mux"
)

type PedidoHandler struct {
	pedidoService    ports.PedidoService
	pagamentoService ports.PagamentoService
	presenter        *presenters.PedidoPresenter
}

func NovoPedidoHandler(
	pedidoService ports.PedidoService, 
	pagamentoService ports.PagamentoService,
) *PedidoHandler {
	return &PedidoHandler{
		pedidoService:    pedidoService,
		pagamentoService: pagamentoService,
		presenter:        presenters.NewPedidoPresenter(),
	}
}


type AtualizarStatusRequest struct {
	Status domain.StatusPedido `json:"status"`
}

// CheckoutPedido cria um novo pedido e gera pagamento PIX.
// @Summary Checkout do pedido
// @Tags pedidos
// @Accept json
// @Produce json
// @Param pedido body dtos.PedidoRequest true "Dados do pedido"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {string} string "Erro ao criar pedido"
// @Router /checkout [post]
func (h *PedidoHandler) CheckoutPedido(w http.ResponseWriter, r *http.Request) {
	var req dtos.PedidoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		status, response := h.presenter.PresentError(err)
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(response)
		return
	}

	pedido, err := req.ToDomain()
	if err != nil {
		status, response := h.presenter.PresentError(err)
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(response)
		return
	}

	createdPedido, err := h.pedidoService.CriarPedido(r.Context(), pedido.ClienteID, pedido.Itens)
	if err != nil {
		status, response := h.presenter.PresentError(err)
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(response)
		return
	}

	pagamento, err := h.pagamentoService.CriarPagamento(
		r.Context(), 
		createdPedido.ID, 
		createdPedido.ValorTotal, 
		domain.TipoPagamentoPix,
	)
	if err != nil {
		status, response := h.presenter.PresentError(err)
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := map[string]interface{}{
		"pedido_id":    createdPedido.ID,
		"status":       string(createdPedido.Status),
		"valor":        createdPedido.ValorTotal,
		"pagamento_id": pagamento.ID,
		"qr_code":      pagamento.QRCode,
		"external_id":  pagamento.ExternalID,
		"expires_at":   pagamento.ExpiresAt,
		"message":      "Pedido criado com sucesso. Realize o pagamento via PIX para prosseguir.",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// ListarPedidos retorna pedidos com ordenação específica: Pronto > Em preparação > Recebido, mais antigos primeiro, não exibe finalizados.
// @Summary Listar pedidos
// @Tags pedidos
// @Produce json
// @Param status query string false "Status do pedido"
// @Param cliente_id query string false "ID do cliente"
// @Success 200 {array} domain.Pedido
// @Failure 500 {string} string "Erro ao listar pedidos"
// @Router /pedidos [get]
func (h *PedidoHandler) ListarPedidos(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	clienteID := r.URL.Query().Get("cliente_id")

	var pedidos []*domain.Pedido
	var err error

	switch {
	case status != "":
		pedidos, err = h.pedidoService.ListarPedidosPorStatus(r.Context(), domain.StatusPedido(status))
	case clienteID != "":
		pedidos, err = h.pedidoService.ListarPedidosPorCliente(r.Context(), clienteID)
	default:
		pedidos, err = h.pedidoService.ListarPedidos(r.Context())
	}

	if err != nil {
		status, response := h.presenter.PresentError(err)
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(response)
		return
	}

	var filteredPedidos []*domain.Pedido
	for _, p := range pedidos {
		if p.Status != domain.StatusFinalizado {
			filteredPedidos = append(filteredPedidos, p)
		}
	}


	sort.Slice(filteredPedidos, func(i, j int) bool {
		if filteredPedidos[i].Status == filteredPedidos[j].Status {
			return filteredPedidos[i].CreatedAt.Before(filteredPedidos[j].CreatedAt)
		}
		return getStatusPriority(filteredPedidos[i].Status) > getStatusPriority(filteredPedidos[j].Status)
	})

	response := h.presenter.PresentList(filteredPedidos)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// BuscarPedidoPorID retorna um pedido pelo ID.
// @Summary Buscar pedido por ID
// @Tags pedidos
// @Produce json
// @Param id path string true "ID do pedido"
// @Success 200 {object} domain.Pedido
// @Failure 404 {string} string "Pedido não encontrado"
// @Router /pedidos/{id} [get]
func (h *PedidoHandler) BuscarPedidoPorID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	pedido, err := h.pedidoService.BuscarPedidoPorID(r.Context(), id)
	if err != nil {
		status, response := h.presenter.PresentError(err)
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(response)
		return
	}

	if pedido == nil {
		http.Error(w, "Pedido não encontrado", http.StatusNotFound)
		return
	}

	response := h.presenter.Present(pedido)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// AtualizarStatusPedido atualiza o status de um pedido.
// @Summary Atualizar status do pedido
// @Tags pedidos
// @Accept json
// @Produce json
// @Param id path string true "ID do pedido"
// @Param status body AtualizarStatusRequest true "Novo status"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {string} string "Erro ao atualizar status do pedido"
// @Router /pedidos/{id}/status [patch]
func (h *PedidoHandler) AtualizarStatusPedido(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req AtualizarStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		status, response := h.presenter.PresentError(err)
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(response)
		return
	}

	err := h.pedidoService.AtualizarStatusPedido(r.Context(), id, req.Status)
	if err != nil {
		status, response := h.presenter.PresentError(err)
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(response)
		return
	}

	status, response := h.presenter.PresentSuccess("Status do pedido atualizado com sucesso", nil)
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

// ConsultarStatusPagamento retorna o status de pagamento do pedido.
// @Summary Consultar status de pagamento do pedido
// @Tags pagamentos
// @Produce json
// @Param id path string true "ID do pedido"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {string} string "Pedido não encontrado"
// @Router /pedidos/{id}/pagamento/status [get]
func (h *PedidoHandler) ConsultarStatusPagamento(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	pagamento, err := h.pagamentoService.ConsultarStatusPagamento(r.Context(), id)
	if err != nil {
		http.Error(w, "Erro ao consultar status do pagamento: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if pagamento == nil {
		http.Error(w, "Pagamento não encontrado para o pedido", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"pedido_id":    pagamento.PedidoID,
		"pagamento_id": pagamento.ID,
		"status":       string(pagamento.Status),
		"valor":        pagamento.Valor,
		"tipo":         string(pagamento.Tipo),
		"qr_code":      pagamento.QRCode,
		"external_id":  pagamento.ExternalID,
		"expires_at":   pagamento.ExpiresAt,
		"processed_at": pagamento.ProcessedAt,
	})
}

func getStatusPriority(status domain.StatusPedido) int {
	switch status {
	case domain.StatusPronto:
		return 3
	case domain.StatusEmPreparacao:
		return 2
	case domain.StatusRecebido:
		return 1
	default:
		return 0
	}
}
