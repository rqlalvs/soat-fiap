package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"soat-fiap/internal/core/domain"
	"testing"
	"time"

	"github.com/gorilla/mux"
)

type MockPagamentoService struct {
	CriarPagamentoFunc            func(ctx context.Context, pedidoID string, valor float64, tipo domain.TipoPagamento) (*domain.Pagamento, error)
	BuscarPagamentoPorIDFunc      func(ctx context.Context, id string) (*domain.Pagamento, error)
	ConsultarStatusPagamentoFunc  func(ctx context.Context, pedidoID string) (*domain.Pagamento, error)
	ProcessarWebhookPagamentoFunc func(ctx context.Context, externalID string, status domain.StatusPagamento) error
	CancelarPagamentoFunc         func(ctx context.Context, id string) error
}

func (m *MockPagamentoService) CriarPagamento(ctx context.Context, pedidoID string, valor float64, tipo domain.TipoPagamento) (*domain.Pagamento, error) {
	if m.CriarPagamentoFunc != nil {
		return m.CriarPagamentoFunc(ctx, pedidoID, valor, tipo)
	}
	return nil, nil
}

func (m *MockPagamentoService) BuscarPagamentoPorID(ctx context.Context, id string) (*domain.Pagamento, error) {
	if m.BuscarPagamentoPorIDFunc != nil {
		return m.BuscarPagamentoPorIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockPagamentoService) ConsultarStatusPagamento(ctx context.Context, pedidoID string) (*domain.Pagamento, error) {
	if m.ConsultarStatusPagamentoFunc != nil {
		return m.ConsultarStatusPagamentoFunc(ctx, pedidoID)
	}
	return nil, nil
}

func (m *MockPagamentoService) ProcessarWebhookPagamento(ctx context.Context, externalID string, status domain.StatusPagamento) error {
	if m.ProcessarWebhookPagamentoFunc != nil {
		return m.ProcessarWebhookPagamentoFunc(ctx, externalID, status)
	}
	return nil
}

func (m *MockPagamentoService) CancelarPagamento(ctx context.Context, id string) error {
	if m.CancelarPagamentoFunc != nil {
		return m.CancelarPagamentoFunc(ctx, id)
	}
	return nil
}

func TestNovoPagamentoHandler(t *testing.T) {
	mockService := &MockPagamentoService{}
	handler := NovoPagamentoHandler(mockService)

	if handler == nil {
		t.Fatal("NovoPagamentoHandler() retornou nil")
	}

	if handler.pagamentoService == nil {
		t.Error("pagamentoService não foi inicializado")
	}
}

func TestPagamentoHandler_CriarPagamento(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockFunc       func(ctx context.Context, pedidoID string, valor float64, tipo domain.TipoPagamento) (*domain.Pagamento, error)
		expectedStatus int
	}{
		{
			name: "Sucesso ao criar pagamento",
			requestBody: CriarPagamentoRequest{
				PedidoID: "pedido-123",
				Valor:    100.50,
				Tipo:     domain.TipoPagamentoPix,
			},
			mockFunc: func(ctx context.Context, pedidoID string, valor float64, tipo domain.TipoPagamento) (*domain.Pagamento, error) {
				expiresAt := time.Now().Add(30 * time.Minute)
				return &domain.Pagamento{
					ID:         "pag-123",
					PedidoID:   pedidoID,
					Valor:      valor,
					Status:     domain.StatusPagamentoPendente,
					Tipo:       tipo,
					QRCodeData: "qr-code-data",
					ExternalID: "ext-123",
					ExpiresAt:  &expiresAt,
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
				}, nil
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Erro ao criar pagamento",
			requestBody: CriarPagamentoRequest{
				PedidoID: "pedido-123",
				Valor:    100.50,
				Tipo:     domain.TipoPagamentoPix,
			},
			mockFunc: func(ctx context.Context, pedidoID string, valor float64, tipo domain.TipoPagamento) (*domain.Pagamento, error) {
				return nil, errors.New("erro ao criar pagamento")
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Erro ao decodificar requisição",
			requestBody:    "invalid json",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockPagamentoService{
				CriarPagamentoFunc: tt.mockFunc,
			}
			handler := NovoPagamentoHandler(mockService)

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

			req := httptest.NewRequest(http.MethodPost, "/pagamentos", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			handler.CriarPagamento(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Status = %d, esperado %d", w.Code, tt.expectedStatus)
			}
		})
	}
}

func TestPagamentoHandler_BuscarPagamentoPorID(t *testing.T) {
	tests := []struct {
		name           string
		pagamentoID    string
		mockFunc       func(ctx context.Context, id string) (*domain.Pagamento, error)
		expectedStatus int
	}{
		{
			name:        "Sucesso ao buscar pagamento",
			pagamentoID: "pag-123",
			mockFunc: func(ctx context.Context, id string) (*domain.Pagamento, error) {
				return &domain.Pagamento{
					ID:       id,
					PedidoID: "pedido-123",
					Valor:    100.50,
					Status:   domain.StatusPagamentoPendente,
				}, nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "Pagamento não encontrado",
			pagamentoID: "pag-999",
			mockFunc: func(ctx context.Context, id string) (*domain.Pagamento, error) {
				return nil, nil
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:        "Erro ao buscar pagamento",
			pagamentoID: "pag-123",
			mockFunc: func(ctx context.Context, id string) (*domain.Pagamento, error) {
				return nil, errors.New("erro no banco")
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockPagamentoService{
				BuscarPagamentoPorIDFunc: tt.mockFunc,
			}
			handler := NovoPagamentoHandler(mockService)

			req := httptest.NewRequest(http.MethodGet, "/pagamentos/"+tt.pagamentoID, nil)
			req = mux.SetURLVars(req, map[string]string{"id": tt.pagamentoID})
			w := httptest.NewRecorder()

			handler.BuscarPagamentoPorID(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Status = %d, esperado %d", w.Code, tt.expectedStatus)
			}
		})
	}
}

func TestPagamentoHandler_ConsultarStatusPagamentoPorPedido(t *testing.T) {
	tests := []struct {
		name           string
		pedidoID       string
		mockFunc       func(ctx context.Context, pedidoID string) (*domain.Pagamento, error)
		expectedStatus int
	}{
		{
			name:     "Sucesso ao consultar status",
			pedidoID: "pedido-123",
			mockFunc: func(ctx context.Context, pedidoID string) (*domain.Pagamento, error) {
				return &domain.Pagamento{
					ID:       "pag-123",
					PedidoID: pedidoID,
					Status:   domain.StatusPagamentoAprovado,
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
			mockService := &MockPagamentoService{
				ConsultarStatusPagamentoFunc: tt.mockFunc,
			}
			handler := NovoPagamentoHandler(mockService)

			req := httptest.NewRequest(http.MethodGet, "/pedidos/"+tt.pedidoID+"/pagamento", nil)
			req = mux.SetURLVars(req, map[string]string{"pedido_id": tt.pedidoID})
			w := httptest.NewRecorder()

			handler.ConsultarStatusPagamentoPorPedido(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Status = %d, esperado %d", w.Code, tt.expectedStatus)
			}
		})
	}
}

func TestPagamentoHandler_WebhookPagamento(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockFunc       func(ctx context.Context, externalID string, status domain.StatusPagamento) error
		expectedStatus int
	}{
		{
			name: "Sucesso ao processar webhook",
			requestBody: WebhookPagamentoRequest{
				ExternalID: "ext-123",
				Status:     domain.StatusPagamentoAprovado,
			},
			mockFunc: func(ctx context.Context, externalID string, status domain.StatusPagamento) error {
				return nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Erro ao processar webhook",
			requestBody: WebhookPagamentoRequest{
				ExternalID: "ext-123",
				Status:     domain.StatusPagamentoAprovado,
			},
			mockFunc: func(ctx context.Context, externalID string, status domain.StatusPagamento) error {
				return errors.New("erro ao processar")
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Erro ao decodificar webhook",
			requestBody:    "invalid json",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "ExternalID vazio",
			requestBody: WebhookPagamentoRequest{
				ExternalID: "",
				Status:     domain.StatusPagamentoAprovado,
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Status vazio",
			requestBody: WebhookPagamentoRequest{
				ExternalID: "ext-123",
				Status:     "",
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockPagamentoService{
				ProcessarWebhookPagamentoFunc: tt.mockFunc,
			}
			handler := NovoPagamentoHandler(mockService)

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

			req := httptest.NewRequest(http.MethodPost, "/webhook/pagamento", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			handler.WebhookPagamento(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Status = %d, esperado %d", w.Code, tt.expectedStatus)
			}
		})
	}
}

func TestPagamentoHandler_CancelarPagamento(t *testing.T) {
	tests := []struct {
		name           string
		pagamentoID    string
		mockFunc       func(ctx context.Context, id string) error
		expectedStatus int
	}{
		{
			name:        "Sucesso ao cancelar pagamento",
			pagamentoID: "pag-123",
			mockFunc: func(ctx context.Context, id string) error {
				return nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "Erro ao cancelar pagamento",
			pagamentoID: "pag-123",
			mockFunc: func(ctx context.Context, id string) error {
				return errors.New("erro ao cancelar")
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockPagamentoService{
				CancelarPagamentoFunc: tt.mockFunc,
			}
			handler := NovoPagamentoHandler(mockService)

			req := httptest.NewRequest(http.MethodPatch, "/pagamentos/"+tt.pagamentoID+"/cancelar", nil)
			req = mux.SetURLVars(req, map[string]string{"id": tt.pagamentoID})
			w := httptest.NewRecorder()

			handler.CancelarPagamento(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Status = %d, esperado %d", w.Code, tt.expectedStatus)
			}
		})
	}
}

func TestPagamentoHandler_SimularPagamento(t *testing.T) {
	tests := []struct {
		name           string
		externalID     string
		status         string
		mockFunc       func(ctx context.Context, externalID string, status domain.StatusPagamento) error
		expectedStatus int
	}{
		{
			name:       "Sucesso ao simular aprovação",
			externalID: "ext-123",
			status:     "APROVADO",
			mockFunc: func(ctx context.Context, externalID string, status domain.StatusPagamento) error {
				return nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:       "Sucesso ao simular recusa",
			externalID: "ext-123",
			status:     "RECUSADO",
			mockFunc: func(ctx context.Context, externalID string, status domain.StatusPagamento) error {
				return nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:       "Sucesso ao simular cancelamento",
			externalID: "ext-123",
			status:     "CANCELADO",
			mockFunc: func(ctx context.Context, externalID string, status domain.StatusPagamento) error {
				return nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Status vazio",
			externalID:     "ext-123",
			status:         "",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Status inválido",
			externalID:     "ext-123",
			status:         "INVALIDO",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:       "Erro ao simular pagamento",
			externalID: "ext-123",
			status:     "APROVADO",
			mockFunc: func(ctx context.Context, externalID string, status domain.StatusPagamento) error {
				return errors.New("erro ao simular")
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockPagamentoService{
				ProcessarWebhookPagamentoFunc: tt.mockFunc,
			}
			handler := NovoPagamentoHandler(mockService)

			url := "/mock/pagamento/" + tt.externalID + "/simular"
			if tt.status != "" {
				url += "?status=" + tt.status
			}

			req := httptest.NewRequest(http.MethodPost, url, nil)
			req = mux.SetURLVars(req, map[string]string{"external_id": tt.externalID})
			w := httptest.NewRecorder()

			handler.SimularPagamento(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Status = %d, esperado %d", w.Code, tt.expectedStatus)
			}
		})
	}
}
