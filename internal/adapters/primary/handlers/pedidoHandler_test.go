package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"soat-fiap/internal/adapters/primary/dtos"
	"soat-fiap/internal/core/domain"
	"testing"
	"time"

	"github.com/gorilla/mux"
)

type MockPedidoService struct {
	CriarPedidoFunc             func(ctx context.Context, clienteID *string, itens []domain.ItemPedido) (*domain.Pedido, error)
	BuscarPedidoPorIDFunc       func(ctx context.Context, id string) (*domain.Pedido, error)
	ListarPedidosFunc           func(ctx context.Context) ([]*domain.Pedido, error)
	ListarPedidosPorStatusFunc  func(ctx context.Context, status domain.StatusPedido) ([]*domain.Pedido, error)
	ListarPedidosPorClienteFunc func(ctx context.Context, clienteID string) ([]*domain.Pedido, error)
	AtualizarStatusPedidoFunc   func(ctx context.Context, id string, status domain.StatusPedido) error
}

func (m *MockPedidoService) CriarPedido(ctx context.Context, clienteID *string, itens []domain.ItemPedido) (*domain.Pedido, error) {
	if m.CriarPedidoFunc != nil {
		return m.CriarPedidoFunc(ctx, clienteID, itens)
	}
	return nil, nil
}

func (m *MockPedidoService) BuscarPedidoPorID(ctx context.Context, id string) (*domain.Pedido, error) {
	if m.BuscarPedidoPorIDFunc != nil {
		return m.BuscarPedidoPorIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockPedidoService) ListarPedidos(ctx context.Context) ([]*domain.Pedido, error) {
	if m.ListarPedidosFunc != nil {
		return m.ListarPedidosFunc(ctx)
	}
	return nil, nil
}

func (m *MockPedidoService) ListarPedidosPorStatus(ctx context.Context, status domain.StatusPedido) ([]*domain.Pedido, error) {
	if m.ListarPedidosPorStatusFunc != nil {
		return m.ListarPedidosPorStatusFunc(ctx, status)
	}
	return nil, nil
}

func (m *MockPedidoService) ListarPedidosPorCliente(ctx context.Context, clienteID string) ([]*domain.Pedido, error) {
	if m.ListarPedidosPorClienteFunc != nil {
		return m.ListarPedidosPorClienteFunc(ctx, clienteID)
	}
	return nil, nil
}

func (m *MockPedidoService) AtualizarStatusPedido(ctx context.Context, id string, status domain.StatusPedido) error {
	if m.AtualizarStatusPedidoFunc != nil {
		return m.AtualizarStatusPedidoFunc(ctx, id, status)
	}
	return nil
}

func TestNovoPedidoHandler(t *testing.T) {
	mockPedidoService := &MockPedidoService{}
	mockPagamentoService := &MockPagamentoService{}
	handler := NovoPedidoHandler(mockPedidoService, mockPagamentoService)

	if handler == nil {
		t.Fatal("NovoPedidoHandler() retornou nil")
	}

	if handler.pedidoService == nil {
		t.Error("pedidoService não foi inicializado")
	}

	if handler.pagamentoService == nil {
		t.Error("pagamentoService não foi inicializado")
	}
}

func TestPedidoHandler_CheckoutPedido(t *testing.T) {
	clienteID := "cliente-123"
	
	tests := []struct {
		name                  string
		requestBody           interface{}
		mockCriarPedidoFunc   func(ctx context.Context, clienteID *string, itens []domain.ItemPedido) (*domain.Pedido, error)
		mockCriarPagamentoFunc func(ctx context.Context, pedidoID string, valor float64, tipo domain.TipoPagamento) (*domain.Pagamento, error)
		expectedStatus        int
	}{
		{
			name: "Sucesso ao fazer checkout",
			requestBody: map[string]interface{}{
				"cliente_id": clienteID,
				"itens": []map[string]interface{}{
					{
						"produto_id": "prod-1",
						"quantidade": 2,
						"observacao": "Sem cebola",
					},
				},
			},
			mockCriarPedidoFunc: func(ctx context.Context, clienteID *string, itens []domain.ItemPedido) (*domain.Pedido, error) {
				return &domain.Pedido{
					ID:         "pedido-123",
					ClienteID:  clienteID,
					ValorTotal: 50.00,
					Status:     domain.StatusRecebido,
					Itens:      itens,
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
				}, nil
			},
			mockCriarPagamentoFunc: func(ctx context.Context, pedidoID string, valor float64, tipo domain.TipoPagamento) (*domain.Pagamento, error) {
				expiresAt := time.Now().Add(30 * time.Minute)
				return &domain.Pagamento{
					ID:         "pag-123",
					PedidoID:   pedidoID,
					Valor:      valor,
					Status:     domain.StatusPagamentoPendente,
					Tipo:       tipo,
					QRCodeData: "qr-code",
					ExternalID: "ext-123",
					ExpiresAt:  &expiresAt,
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
				}, nil
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Erro ao decodificar requisição",
			requestBody:    "invalid json",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Erro ao criar pedido",
			requestBody: map[string]interface{}{
				"cliente_id": clienteID,
				"itens": []map[string]interface{}{
					{
						"produto_id": "prod-1",
						"quantidade": 2,
					},
				},
			},
			mockCriarPedidoFunc: func(ctx context.Context, clienteID *string, itens []domain.ItemPedido) (*domain.Pedido, error) {
				return nil, errors.New("erro ao criar pedido")
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Erro ao criar pagamento",
			requestBody: map[string]interface{}{
				"cliente_id": clienteID,
				"itens": []map[string]interface{}{
					{
						"produto_id": "prod-1",
						"quantidade": 2,
					},
				},
			},
			mockCriarPedidoFunc: func(ctx context.Context, clienteID *string, itens []domain.ItemPedido) (*domain.Pedido, error) {
				return &domain.Pedido{
					ID:         "pedido-123",
					ClienteID:  clienteID,
					ValorTotal: 50.00,
					Status:     domain.StatusRecebido,
					Itens:      itens,
				}, nil
			},
			mockCriarPagamentoFunc: func(ctx context.Context, pedidoID string, valor float64, tipo domain.TipoPagamento) (*domain.Pagamento, error) {
				return nil, errors.New("erro ao criar pagamento")
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPedidoService := &MockPedidoService{
				CriarPedidoFunc: tt.mockCriarPedidoFunc,
			}
			mockPagamentoService := &MockPagamentoService{
				CriarPagamentoFunc: tt.mockCriarPagamentoFunc,
			}
			handler := NovoPedidoHandler(mockPedidoService, mockPagamentoService)

			var body []byte
			var err error
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.requestBody)
				if err != nil {
					t.Fatalf("Erro ao serializar request: %v", err)
				}
			}

			req := httptest.NewRequest(http.MethodPost, "/checkout", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			handler.CheckoutPedido(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Status = %d, esperado %d", w.Code, tt.expectedStatus)
			}
		})
	}
}

func TestPedidoHandler_ListarPedidos(t *testing.T) {
	tests := []struct {
		name                   string
		queryParams            string
		mockListarFunc         func(ctx context.Context) ([]*domain.Pedido, error)
		mockListarStatusFunc   func(ctx context.Context, status domain.StatusPedido) ([]*domain.Pedido, error)
		mockListarClienteFunc  func(ctx context.Context, clienteID string) ([]*domain.Pedido, error)
		expectedStatus         int
	}{
		{
			name:        "Sucesso ao listar todos os pedidos",
			queryParams: "",
			mockListarFunc: func(ctx context.Context) ([]*domain.Pedido, error) {
				return []*domain.Pedido{
					{ID: "1", Status: domain.StatusRecebido, CreatedAt: time.Now()},
					{ID: "2", Status: domain.StatusEmPreparacao, CreatedAt: time.Now()},
				}, nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "Sucesso ao listar por status",
			queryParams: "?status=RECEBIDO",
			mockListarStatusFunc: func(ctx context.Context, status domain.StatusPedido) ([]*domain.Pedido, error) {
				return []*domain.Pedido{
					{ID: "1", Status: status, CreatedAt: time.Now()},
				}, nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "Sucesso ao listar por cliente",
			queryParams: "?cliente_id=cliente-123",
			mockListarClienteFunc: func(ctx context.Context, clienteID string) ([]*domain.Pedido, error) {
				return []*domain.Pedido{
					{ID: "1", ClienteID: &clienteID, CreatedAt: time.Now()},
				}, nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "Erro ao listar pedidos",
			queryParams: "",
			mockListarFunc: func(ctx context.Context) ([]*domain.Pedido, error) {
				return nil, errors.New("erro no banco")
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPedidoService := &MockPedidoService{
				ListarPedidosFunc:           tt.mockListarFunc,
				ListarPedidosPorStatusFunc:  tt.mockListarStatusFunc,
				ListarPedidosPorClienteFunc: tt.mockListarClienteFunc,
			}
			mockPagamentoService := &MockPagamentoService{}
			handler := NovoPedidoHandler(mockPedidoService, mockPagamentoService)

			req := httptest.NewRequest(http.MethodGet, "/pedidos"+tt.queryParams, nil)
			w := httptest.NewRecorder()

			handler.ListarPedidos(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Status = %d, esperado %d", w.Code, tt.expectedStatus)
			}
		})
	}
}

func TestPedidoHandler_BuscarPedidoPorID(t *testing.T) {
	tests := []struct {
		name           string
		pedidoID       string
		mockFunc       func(ctx context.Context, id string) (*domain.Pedido, error)
		expectedStatus int
	}{
		{
			name:     "Sucesso ao buscar pedido",
			pedidoID: "pedido-123",
			mockFunc: func(ctx context.Context, id string) (*domain.Pedido, error) {
				clienteID := "cliente-123"
				return &domain.Pedido{
					ID:         id,
					ClienteID:  &clienteID,
					ValorTotal: 100.00,
					Status:     domain.StatusRecebido,
					Itens:      []domain.ItemPedido{},
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
				}, nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:     "Pedido não encontrado",
			pedidoID: "pedido-999",
			mockFunc: func(ctx context.Context, id string) (*domain.Pedido, error) {
				return nil, nil
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:     "Erro ao buscar pedido",
			pedidoID: "pedido-123",
			mockFunc: func(ctx context.Context, id string) (*domain.Pedido, error) {
				return nil, errors.New("erro no banco")
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPedidoService := &MockPedidoService{
				BuscarPedidoPorIDFunc: tt.mockFunc,
			}
			mockPagamentoService := &MockPagamentoService{}
			handler := NovoPedidoHandler(mockPedidoService, mockPagamentoService)

			req := httptest.NewRequest(http.MethodGet, "/pedidos/"+tt.pedidoID, nil)
			req = mux.SetURLVars(req, map[string]string{"id": tt.pedidoID})
			w := httptest.NewRecorder()

			handler.BuscarPedidoPorID(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Status = %d, esperado %d", w.Code, tt.expectedStatus)
			}
		})
	}
}

func TestPedidoHandler_AtualizarStatusPedido(t *testing.T) {
	tests := []struct {
		name           string
		pedidoID       string
		requestBody    interface{}
		mockFunc       func(ctx context.Context, id string, status domain.StatusPedido) error
		expectedStatus int
	}{
		{
			name:     "Sucesso ao atualizar status",
			pedidoID: "pedido-123",
			requestBody: AtualizarStatusRequest{
				Status: domain.StatusEmPreparacao,
			},
			mockFunc: func(ctx context.Context, id string, status domain.StatusPedido) error {
				return nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:     "Erro ao atualizar status",
			pedidoID: "pedido-123",
			requestBody: AtualizarStatusRequest{
				Status: domain.StatusEmPreparacao,
			},
			mockFunc: func(ctx context.Context, id string, status domain.StatusPedido) error {
				return errors.New("erro ao atualizar")
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Erro ao decodificar requisição",
			pedidoID:       "pedido-123",
			requestBody:    "invalid json",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPedidoService := &MockPedidoService{
				AtualizarStatusPedidoFunc: tt.mockFunc,
			}
			mockPagamentoService := &MockPagamentoService{}
			handler := NovoPedidoHandler(mockPedidoService, mockPagamentoService)

			var body []byte
			var err error
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.requestBody)
				if err != nil {
					t.Fatalf("Erro ao serializar request: %v", err)
				}
			}

			req := httptest.NewRequest(http.MethodPatch, "/pedidos/"+tt.pedidoID+"/status", bytes.NewBuffer(body))
			req = mux.SetURLVars(req, map[string]string{"id": tt.pedidoID})
			w := httptest.NewRecorder()

			handler.AtualizarStatusPedido(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Status = %d, esperado %d", w.Code, tt.expectedStatus)
			}
		})
	}
}

func TestPedidoHandler_ConsultarStatusPagamento(t *testing.T) {
	tests := []struct {
		name           string
		pedidoID       string
		mockFunc       func(ctx context.Context, pedidoID string) (*domain.Pagamento, error)
		expectedStatus int
	}{
		{
			name:     "Sucesso ao consultar status de pagamento",
			pedidoID: "pedido-123",
			mockFunc: func(ctx context.Context, pedidoID string) (*domain.Pagamento, error) {
				return &domain.Pagamento{
					ID:         "pag-123",
					PedidoID:   pedidoID,
					Status:     domain.StatusPagamentoAprovado,
					Valor:      100.00,
					Tipo:       domain.TipoPagamentoPix,
					QRCodeData: "qr-code",
					ExternalID: "ext-123",
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
				}, nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:     "Pagamento não encontrado",
			pedidoID: "pedido-999",
			mockFunc: func(ctx context.Context, pedidoID string) (*domain.Pagamento, error) {
				return nil, nil
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:     "Erro ao consultar status",
			pedidoID: "pedido-123",
			mockFunc: func(ctx context.Context, pedidoID string) (*domain.Pagamento, error) {
				return nil, errors.New("erro no banco")
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPedidoService := &MockPedidoService{}
			mockPagamentoService := &MockPagamentoService{
				ConsultarStatusPagamentoFunc: tt.mockFunc,
			}
			handler := NovoPedidoHandler(mockPedidoService, mockPagamentoService)

			req := httptest.NewRequest(http.MethodGet, "/pedidos/"+tt.pedidoID+"/pagamento/status", nil)
			req = mux.SetURLVars(req, map[string]string{"id": tt.pedidoID})
			w := httptest.NewRecorder()

			handler.ConsultarStatusPagamento(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Status = %d, esperado %d", w.Code, tt.expectedStatus)
			}
		})
	}
}

func TestGetStatusPrioridade(t *testing.T) {
	tests := []struct {
		name     string
		status   domain.StatusPedido
		expected int
	}{
		{"Pronto", domain.StatusPronto, 3},
		{"Em Preparação", domain.StatusEmPreparacao, 2},
		{"Recebido", domain.StatusRecebido, 1},
		{"Finalizado", domain.StatusFinalizado, 0},
		{"Cancelado", domain.StatusCancelado, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getStatusPrioridade(tt.status)
			if result != tt.expected {
				t.Errorf("getStatusPrioridade(%s) = %d, esperado %d", tt.status, result, tt.expected)
			}
		})
	}
}
