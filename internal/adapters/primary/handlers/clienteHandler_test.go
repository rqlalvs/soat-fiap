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

type MockClienteService struct {
	CriarClienteFunc         func(ctx context.Context, nome, cpf, email, telefone string) (*domain.Cliente, error)
	BuscarClientePorIDFunc   func(ctx context.Context, id string) (*domain.Cliente, error)
	BuscarClientePorCPFFunc  func(ctx context.Context, cpf string) (*domain.Cliente, error)
	ListarClientesFunc       func(ctx context.Context) ([]*domain.Cliente, error)
	AtualizarClienteFunc     func(ctx context.Context, cliente *domain.Cliente) error
	DeletarClienteFunc       func(ctx context.Context, id string) error
}

func (m *MockClienteService) CriarCliente(ctx context.Context, nome, cpf, email, telefone string) (*domain.Cliente, error) {
	if m.CriarClienteFunc != nil {
		return m.CriarClienteFunc(ctx, nome, cpf, email, telefone)
	}
	return nil, nil
}

func (m *MockClienteService) BuscarClientePorID(ctx context.Context, id string) (*domain.Cliente, error) {
	if m.BuscarClientePorIDFunc != nil {
		return m.BuscarClientePorIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockClienteService) BuscarClientePorCPF(ctx context.Context, cpf string) (*domain.Cliente, error) {
	if m.BuscarClientePorCPFFunc != nil {
		return m.BuscarClientePorCPFFunc(ctx, cpf)
	}
	return nil, nil
}

func (m *MockClienteService) ListarClientes(ctx context.Context) ([]*domain.Cliente, error) {
	if m.ListarClientesFunc != nil {
		return m.ListarClientesFunc(ctx)
	}
	return nil, nil
}

func (m *MockClienteService) AtualizarCliente(ctx context.Context, cliente *domain.Cliente) error {
	if m.AtualizarClienteFunc != nil {
		return m.AtualizarClienteFunc(ctx, cliente)
	}
	return nil
}

func (m *MockClienteService) DeletarCliente(ctx context.Context, id string) error {
	if m.DeletarClienteFunc != nil {
		return m.DeletarClienteFunc(ctx, id)
	}
	return nil
}

func TestNovoClienteHandler(t *testing.T) {
	mockService := &MockClienteService{}
	handler := NovoClienteHandler(mockService)

	if handler == nil {
		t.Fatal("NovoClienteHandler() retornou nil")
	}

	if handler.clienteService == nil {
		t.Error("clienteService não foi inicializado")
	}
}

func TestClienteHandler_CriarCliente(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockFunc       func(ctx context.Context, nome, cpf, email, telefone string) (*domain.Cliente, error)
		expectedStatus int
	}{
		{
			name: "Sucesso ao criar cliente",
			requestBody: CriarClienteRequest{
				Nome:     "João Silva",
				CPF:      "12345678909",
				Email:    "joao@example.com",
				Telefone: "11999999999",
			},
			mockFunc: func(ctx context.Context, nome, cpf, email, telefone string) (*domain.Cliente, error) {
				return &domain.Cliente{
					ID:        "123",
					Nome:      nome,
					CPF:       cpf,
					Email:     email,
					Telefone:  telefone,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Erro ao criar cliente - serviço retorna erro",
			requestBody: CriarClienteRequest{
				Nome:     "João Silva",
				CPF:      "12345678909",
				Email:    "joao@example.com",
				Telefone: "11999999999",
			},
			mockFunc: func(ctx context.Context, nome, cpf, email, telefone string) (*domain.Cliente, error) {
				return nil, errors.New("CPF inválido")
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
			mockService := &MockClienteService{
				CriarClienteFunc: tt.mockFunc,
			}
			handler := NovoClienteHandler(mockService)

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

			req := httptest.NewRequest(http.MethodPost, "/clientes", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			handler.CriarCliente(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Status = %d, esperado %d", w.Code, tt.expectedStatus)
			}
		})
	}
}

func TestClienteHandler_BuscarClientePorID(t *testing.T) {
	tests := []struct {
		name           string
		clienteID      string
		mockFunc       func(ctx context.Context, id string) (*domain.Cliente, error)
		expectedStatus int
	}{
		{
			name:      "Sucesso ao buscar cliente",
			clienteID: "123",
			mockFunc: func(ctx context.Context, id string) (*domain.Cliente, error) {
				return &domain.Cliente{
					ID:    id,
					Nome:  "João Silva",
					CPF:   "12345678909",
					Email: "joao@example.com",
				}, nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:      "Cliente não encontrado",
			clienteID: "999",
			mockFunc: func(ctx context.Context, id string) (*domain.Cliente, error) {
				return nil, nil
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:      "Erro ao buscar cliente",
			clienteID: "123",
			mockFunc: func(ctx context.Context, id string) (*domain.Cliente, error) {
				return nil, errors.New("erro no banco")
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockClienteService{
				BuscarClientePorIDFunc: tt.mockFunc,
			}
			handler := NovoClienteHandler(mockService)

			req := httptest.NewRequest(http.MethodGet, "/clientes/"+tt.clienteID, nil)
			req = mux.SetURLVars(req, map[string]string{"id": tt.clienteID})
			w := httptest.NewRecorder()

			handler.BuscarClientePorID(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Status = %d, esperado %d", w.Code, tt.expectedStatus)
			}
		})
	}
}

func TestClienteHandler_BuscarClientePorCPF(t *testing.T) {
	tests := []struct {
		name           string
		cpf            string
		mockFunc       func(ctx context.Context, cpf string) (*domain.Cliente, error)
		expectedStatus int
	}{
		{
			name: "Sucesso ao buscar por CPF",
			cpf:  "123.456.789-09",
			mockFunc: func(ctx context.Context, cpf string) (*domain.Cliente, error) {
				return &domain.Cliente{
					ID:   "123",
					Nome: "João Silva",
					CPF:  cpf,
				}, nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "CPF não encontrado",
			cpf:  "999.999.999-99",
			mockFunc: func(ctx context.Context, cpf string) (*domain.Cliente, error) {
				return nil, nil
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "Erro ao buscar por CPF",
			cpf:  "123.456.789-09",
			mockFunc: func(ctx context.Context, cpf string) (*domain.Cliente, error) {
				return nil, errors.New("erro no banco")
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockClienteService{
				BuscarClientePorCPFFunc: tt.mockFunc,
			}
			handler := NovoClienteHandler(mockService)

			req := httptest.NewRequest(http.MethodGet, "/clientes/cpf/"+tt.cpf, nil)
			req = mux.SetURLVars(req, map[string]string{"cpf": tt.cpf})
			w := httptest.NewRecorder()

			handler.BuscarClientePorCPF(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Status = %d, esperado %d", w.Code, tt.expectedStatus)
			}
		})
	}
}

func TestClienteHandler_ListarClientes(t *testing.T) {
	tests := []struct {
		name           string
		mockFunc       func(ctx context.Context) ([]*domain.Cliente, error)
		expectedStatus int
	}{
		{
			name: "Sucesso ao listar clientes",
			mockFunc: func(ctx context.Context) ([]*domain.Cliente, error) {
				return []*domain.Cliente{
					{ID: "1", Nome: "Cliente 1"},
					{ID: "2", Nome: "Cliente 2"},
				}, nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Lista vazia",
			mockFunc: func(ctx context.Context) ([]*domain.Cliente, error) {
				return []*domain.Cliente{}, nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Erro ao listar",
			mockFunc: func(ctx context.Context) ([]*domain.Cliente, error) {
				return nil, errors.New("erro no banco")
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockClienteService{
				ListarClientesFunc: tt.mockFunc,
			}
			handler := NovoClienteHandler(mockService)

			req := httptest.NewRequest(http.MethodGet, "/clientes", nil)
			w := httptest.NewRecorder()

			handler.ListarClientes(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Status = %d, esperado %d", w.Code, tt.expectedStatus)
			}
		})
	}
}

func TestClienteHandler_AtualizarCliente(t *testing.T) {
	tests := []struct {
		name            string
		clienteID       string
		requestBody     interface{}
		mockBuscarFunc  func(ctx context.Context, id string) (*domain.Cliente, error)
		mockAtualizarFunc func(ctx context.Context, cliente *domain.Cliente) error
		expectedStatus  int
	}{
		{
			name:      "Sucesso ao atualizar cliente",
			clienteID: "123",
			requestBody: AtualizarClienteRequest{
				Nome:     "João Silva Atualizado",
				CPF:      "12345678909",
				Email:    "joao.novo@example.com",
				Telefone: "11988888888",
			},
			mockBuscarFunc: func(ctx context.Context, id string) (*domain.Cliente, error) {
				return &domain.Cliente{
					ID:   id,
					Nome: "João Silva",
					CPF:  "12345678909",
				}, nil
			},
			mockAtualizarFunc: func(ctx context.Context, cliente *domain.Cliente) error {
				return nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:      "Cliente não encontrado",
			clienteID: "999",
			requestBody: AtualizarClienteRequest{
				Nome: "João Silva",
			},
			mockBuscarFunc: func(ctx context.Context, id string) (*domain.Cliente, error) {
				return nil, nil
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Erro ao decodificar requisição",
			clienteID:      "123",
			requestBody:    "invalid json",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:      "Erro ao buscar cliente",
			clienteID: "123",
			requestBody: AtualizarClienteRequest{
				Nome: "João Silva",
			},
			mockBuscarFunc: func(ctx context.Context, id string) (*domain.Cliente, error) {
				return nil, errors.New("erro no banco")
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:      "Erro ao atualizar cliente",
			clienteID: "123",
			requestBody: AtualizarClienteRequest{
				Nome: "João Silva",
			},
			mockBuscarFunc: func(ctx context.Context, id string) (*domain.Cliente, error) {
				return &domain.Cliente{ID: id}, nil
			},
			mockAtualizarFunc: func(ctx context.Context, cliente *domain.Cliente) error {
				return errors.New("erro ao atualizar")
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockClienteService{
				BuscarClientePorIDFunc: tt.mockBuscarFunc,
				AtualizarClienteFunc:   tt.mockAtualizarFunc,
			}
			handler := NovoClienteHandler(mockService)

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

			req := httptest.NewRequest(http.MethodPut, "/clientes/"+tt.clienteID, bytes.NewBuffer(body))
			req = mux.SetURLVars(req, map[string]string{"id": tt.clienteID})
			w := httptest.NewRecorder()

			handler.AtualizarCliente(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Status = %d, esperado %d", w.Code, tt.expectedStatus)
			}
		})
	}
}

func TestClienteHandler_DeletarCliente(t *testing.T) {
	tests := []struct {
		name           string
		clienteID      string
		mockFunc       func(ctx context.Context, id string) error
		expectedStatus int
	}{
		{
			name:      "Sucesso ao deletar cliente",
			clienteID: "123",
			mockFunc: func(ctx context.Context, id string) error {
				return nil
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:      "Erro ao deletar cliente",
			clienteID: "123",
			mockFunc: func(ctx context.Context, id string) error {
				return errors.New("erro no banco")
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockClienteService{
				DeletarClienteFunc: tt.mockFunc,
			}
			handler := NovoClienteHandler(mockService)

			req := httptest.NewRequest(http.MethodDelete, "/clientes/"+tt.clienteID, nil)
			req = mux.SetURLVars(req, map[string]string{"id": tt.clienteID})
			w := httptest.NewRecorder()

			handler.DeletarCliente(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Status = %d, esperado %d", w.Code, tt.expectedStatus)
			}
		})
	}
}
