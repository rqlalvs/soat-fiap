package handlers

import (
	"encoding/json"
	"net/http"
	"soat-fiap/internal/core/domain"
	"soat-fiap/internal/core/ports"
	"sort"
	"time"

	"github.com/gorilla/mux"
)

type PedidoHandler struct {
	pedidoService    ports.PedidoService
	pagamentoService ports.PagamentoService
}

func NovoPedidoHandler(pedidoService ports.PedidoService, pagamentoService ports.PagamentoService) *PedidoHandler {
	return &PedidoHandler{
		pedidoService:    pedidoService,
		pagamentoService: pagamentoService,
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

type CheckoutResponse struct {
	PedidoID    string      `json:"pedido_id"`
	Status      string      `json:"status"`
	Valor       float64     `json:"valor"`
	PagamentoID string      `json:"pagamento_id"`
	QRCode      string      `json:"qr_code"`
	ExternalID  string      `json:"external_id"`
	ExpiresAt   *time.Time  `json:"expires_at"`
	Message     string      `json:"message"`
}

type AtualizarStatusResponse struct {
	Message string `json:"message"`
}

type ConsultarStatusPagamentoResponse struct {
	PedidoID    string      `json:"pedido_id"`
	PagamentoID string      `json:"pagamento_id"`
	Status      string      `json:"status"`
	Valor       float64     `json:"valor"`
	Tipo        string      `json:"tipo"`
	QRCode      string      `json:"qr_code"`
	ExternalID  string      `json:"external_id"`
	ExpiresAt   *time.Time  `json:"expires_at"`
	ProcessedAt *time.Time  `json:"processed_at"`
}

// CheckoutPedido cria um novo pedido e gera pagamento PIX.
// @Summary Checkout do pedido
// @Tags pedidos
// @Accept json
// @Produce json
// @Param pedido body CriarPedidoRequest true "Dados do pedido"
// @Success 201 {object} CheckoutResponse
// @Failure 400 {string} string "Erro ao criar pedido"
// @Router /checkout [post]
func (h *PedidoHandler) CheckoutPedido(w http.ResponseWriter, r *http.Request) {
	var req CriarPedidoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Erro ao decodificar requisição: "+err.Error(), http.StatusBadRequest)
		return
	}

	itens := make([]domain.ItemPedido, len(req.Itens))
	for i, item := range req.Itens {
		itens[i] = domain.ItemPedido{
			ProdutoID:  item.ProdutoID,
			Quantidade: item.Quantidade,
			Observacao: item.Observacao,
		}
	}

	pedido, err := h.pedidoService.CriarPedido(r.Context(), req.ClienteID, itens)
	if err != nil {
		http.Error(w, "Erro ao criar pedido: "+err.Error(), http.StatusBadRequest)
		return
	}

	pagamento, err := h.pagamentoService.CriarPagamento(r.Context(), pedido.ID, pedido.ValorTotal, domain.TipoPagamentoPix)
	if err != nil {
		http.Error(w, "Erro ao criar pagamento: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(CheckoutResponse{
		PedidoID:    pedido.ID,
		Status:      string(pedido.Status),
		Valor:       pedido.ValorTotal,
		PagamentoID: pagamento.ID,
		QRCode:      pagamento.QRCodeData,
		ExternalID:  pagamento.ExternalID,
		ExpiresAt:   pagamento.ExpiresAt,
		Message:     "Pedido criado com sucesso. Realize o pagamento via PIX para prosseguir.",
	})
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

	pedidosFiltrados := make([]*domain.Pedido, 0)
	for _, pedido := range pedidos {
		if pedido.Status != domain.StatusFinalizado {
			pedidosFiltrados = append(pedidosFiltrados, pedido)
		}
	}

	sort.Slice(pedidosFiltrados, func(i, j int) bool {
		pedidoA, pedidoB := pedidosFiltrados[i], pedidosFiltrados[j]
		
		prioridadeA := getStatusPrioridade(pedidoA.Status)
		prioridadeB := getStatusPrioridade(pedidoB.Status)
		
		if prioridadeA != prioridadeB {
			return prioridadeA < prioridadeB
		}
		
		return pedidoA.CreatedAt.Before(pedidoB.CreatedAt)
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pedidosFiltrados)
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

// AtualizarStatusPedido atualiza o status de um pedido.
// @Summary Atualizar status do pedido
// @Tags pedidos
// @Accept json
// @Produce json
// @Param id path string true "ID do pedido"
// @Param status body AtualizarStatusRequest true "Novo status"
// @Success 200 {object} AtualizarStatusResponse
// @Failure 400 {string} string "Erro ao atualizar status do pedido"
// @Router /pedidos/{id}/status [patch]
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
	json.NewEncoder(w).Encode(AtualizarStatusResponse{
		Message: "Status do pedido atualizado com sucesso",
	})
}

// ConsultarStatusPagamento retorna o status de pagamento do pedido.
// @Summary Consultar status de pagamento do pedido
// @Tags pagamentos
// @Produce json
// @Param id path string true "ID do pedido"
// @Success 200 {object} ConsultarStatusPagamentoResponse
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
	json.NewEncoder(w).Encode(ConsultarStatusPagamentoResponse{
		PedidoID:    pagamento.PedidoID,
		PagamentoID: pagamento.ID,
		Status:      string(pagamento.Status),
		Valor:       pagamento.Valor,
		Tipo:        string(pagamento.Tipo),
		QRCode:      pagamento.QRCodeData,
		ExternalID:  pagamento.ExternalID,
		ExpiresAt:   pagamento.ExpiresAt,
		ProcessedAt: pagamento.ProcessedAt,
	})
}

func getStatusPrioridade(status domain.StatusPedido) int {
	switch status {
	case domain.StatusPronto:
		return 1 
	case domain.StatusEmPreparacao:
		return 2
	case domain.StatusRecebido:
		return 3
	case domain.StatusCancelado:
		return 4
	default:
		return 5
	}
}
