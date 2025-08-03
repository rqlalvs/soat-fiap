package handlers

import (
	"encoding/json"
	"net/http"
	"soat-fiap/internal/core/domain"
	"soat-fiap/internal/core/ports"

	"github.com/gorilla/mux"
)

type PagamentoHandler struct {
	pagamentoService ports.PagamentoService
}

func NovoPagamentoHandler(pagamentoService ports.PagamentoService) *PagamentoHandler {
	return &PagamentoHandler{
		pagamentoService: pagamentoService,
	}
}

type CriarPagamentoRequest struct {
	PedidoID string                `json:"pedido_id"`
	Valor    float64               `json:"valor"`
	Tipo     domain.TipoPagamento  `json:"tipo"`
}

type WebhookPagamentoRequest struct {
	ExternalID string                `json:"external_id"`
	Status     domain.StatusPagamento `json:"status"`
	Action     string                 `json:"action,omitempty"`
}

// CriarPagamento cria um novo pagamento para um pedido.
// @Summary Criar pagamento
// @Tags pagamento
// @Accept json
// @Produce json
// @Param pagamento body CriarPagamentoRequest true "Dados do pagamento"
// @Success 201 {object} domain.Pagamento
// @Failure 400 {string} string "Erro ao criar pagamento"
// @Router /pagamentos [post]
func (h *PagamentoHandler) CriarPagamento(w http.ResponseWriter, r *http.Request) {
	var req CriarPagamentoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Erro ao decodificar requisição: "+err.Error(), http.StatusBadRequest)
		return
	}

	pagamento, err := h.pagamentoService.CriarPagamento(r.Context(), req.PedidoID, req.Valor, req.Tipo)
	if err != nil {
		http.Error(w, "Erro ao criar pagamento: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(pagamento)
}

// BuscarPagamentoPorID busca um pagamento pelo ID.
// @Summary Buscar pagamento por ID
// @Tags pagamento
// @Produce json
// @Param id path string true "ID do pagamento"
// @Success 200 {object} domain.Pagamento
// @Failure 404 {string} string "Pagamento não encontrado"
// @Router /pagamentos/{id} [get]
func (h *PagamentoHandler) BuscarPagamentoPorID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	pagamento, err := h.pagamentoService.BuscarPagamentoPorID(r.Context(), id)
	if err != nil {
		http.Error(w, "Erro ao buscar pagamento: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if pagamento == nil {
		http.Error(w, "Pagamento não encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pagamento)
}

// ConsultarStatusPagamentoPorPedido consulta o status de pagamento de um pedido.
// @Summary Consultar status de pagamento por pedido
// @Tags pagamento
// @Produce json
// @Param pedido_id path string true "ID do pedido"
// @Success 200 {object} domain.Pagamento
// @Failure 404 {string} string "Pagamento não encontrado"
// @Router /pedidos/{pedido_id}/pagamento [get]
func (h *PagamentoHandler) ConsultarStatusPagamentoPorPedido(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pedidoID := vars["pedido_id"]

	pagamento, err := h.pagamentoService.ConsultarStatusPagamento(r.Context(), pedidoID)
	if err != nil {
		http.Error(w, "Erro ao consultar status do pagamento: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if pagamento == nil {
		http.Error(w, "Pagamento não encontrado para o pedido", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pagamento)
}

// WebhookPagamento recebe notificações de pagamento do Mercado Pago ou sistema externo.
// @Summary Webhook de pagamento
// @Tags pagamento
// @Accept json
// @Produce json
// @Param webhook body WebhookPagamentoRequest true "Dados do webhook"
// @Success 200 {object} map[string]string
// @Failure 400 {string} string "Erro ao processar webhook"
// @Router /webhook/pagamento [post]
func (h *PagamentoHandler) WebhookPagamento(w http.ResponseWriter, r *http.Request) {
	var req WebhookPagamentoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Erro ao decodificar webhook: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validar dados do webhook
	if req.ExternalID == "" {
		http.Error(w, "external_id é obrigatório", http.StatusBadRequest)
		return
	}

	if req.Status == "" {
		http.Error(w, "status é obrigatório", http.StatusBadRequest)
		return
	}

	err := h.pagamentoService.ProcessarWebhookPagamento(r.Context(), req.ExternalID, req.Status)
	if err != nil {
		http.Error(w, "Erro ao processar webhook: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Webhook processado com sucesso",
	})
}

// CancelarPagamento cancela um pagamento.
// @Summary Cancelar pagamento
// @Tags pagamento
// @Produce json
// @Param id path string true "ID do pagamento"
// @Success 200 {object} map[string]string
// @Failure 400 {string} string "Erro ao cancelar pagamento"
// @Router /pagamentos/{id}/cancelar [patch]
func (h *PagamentoHandler) CancelarPagamento(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := h.pagamentoService.CancelarPagamento(r.Context(), id)
	if err != nil {
		http.Error(w, "Erro ao cancelar pagamento: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Pagamento cancelado com sucesso",
	})
}

// SimularPagamento simula aprovação/recusa de pagamento para testes (mock).
// @Summary Simular pagamento (mock)
// @Tags pagamento
// @Accept json
// @Produce json
// @Param external_id path string true "ID externo do pagamento"
// @Param status query string true "Status para simular (APROVADO/RECUSADO)"
// @Success 200 {object} map[string]string
// @Failure 400 {string} string "Erro ao simular pagamento"
// @Router /mock/pagamento/{external_id}/simular [post]
func (h *PagamentoHandler) SimularPagamento(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	externalID := vars["external_id"]
	
	statusStr := r.URL.Query().Get("status")
	if statusStr == "" {
		http.Error(w, "Parâmetro 'status' é obrigatório", http.StatusBadRequest)
		return
	}

	var status domain.StatusPagamento
	switch statusStr {
	case "APROVADO":
		status = domain.StatusPagamentoAprovado
	case "RECUSADO":
		status = domain.StatusPagamentoRecusado
	case "CANCELADO":
		status = domain.StatusPagamentoCancelado
	default:
		http.Error(w, "Status inválido. Use: APROVADO, RECUSADO ou CANCELADO", http.StatusBadRequest)
		return
	}

	err := h.pagamentoService.ProcessarWebhookPagamento(r.Context(), externalID, status)
	if err != nil {
		http.Error(w, "Erro ao simular pagamento: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":     "success",
		"message":    "Pagamento simulado com sucesso",
		"external_id": externalID,
		"new_status": string(status),
	})
}
