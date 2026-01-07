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

type MockProdutoService struct {
	CriarProdutoFunc               func(ctx context.Context, nome, descricao string, preco float64, categoria domain.Categoria) (*domain.Produto, error)
	BuscarProdutoPorIDFunc         func(ctx context.Context, id string) (*domain.Produto, error)
	ListarProdutosFunc             func(ctx context.Context) ([]*domain.Produto, error)
	ListarProdutosPorCategoriaFunc func(ctx context.Context, categoria domain.Categoria) ([]*domain.Produto, error)
	AtualizarProdutoFunc           func(ctx context.Context, produto *domain.Produto) error
	DeletarProdutoFunc             func(ctx context.Context, id string) error
}

func (m *MockProdutoService) CriarProduto(ctx context.Context, nome, descricao string, preco float64, categoria domain.Categoria) (*domain.Produto, error) {
	if m.CriarProdutoFunc != nil {
		return m.CriarProdutoFunc(ctx, nome, descricao, preco, categoria)
	}
	return nil, nil
}

func (m *MockProdutoService) BuscarProdutoPorID(ctx context.Context, id string) (*domain.Produto, error) {
	if m.BuscarProdutoPorIDFunc != nil {
		return m.BuscarProdutoPorIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockProdutoService) ListarProdutos(ctx context.Context) ([]*domain.Produto, error) {
	if m.ListarProdutosFunc != nil {
		return m.ListarProdutosFunc(ctx)
	}
	return nil, nil
}

func (m *MockProdutoService) ListarProdutosPorCategoria(ctx context.Context, categoria domain.Categoria) ([]*domain.Produto, error) {
	if m.ListarProdutosPorCategoriaFunc != nil {
		return m.ListarProdutosPorCategoriaFunc(ctx, categoria)
	}
	return nil, nil
}

func (m *MockProdutoService) AtualizarProduto(ctx context.Context, produto *domain.Produto) error {
	if m.AtualizarProdutoFunc != nil {
		return m.AtualizarProdutoFunc(ctx, produto)
	}
	return nil
}

func (m *MockProdutoService) DeletarProduto(ctx context.Context, id string) error {
	if m.DeletarProdutoFunc != nil {
		return m.DeletarProdutoFunc(ctx, id)
	}
	return nil
}

func TestNovoProdutoHandler(t *testing.T) {
	mockService := &MockProdutoService{}
	handler := NovoProdutoHandler(mockService)

	if handler == nil {
		t.Fatal("NovoProdutoHandler() retornou nil")
	}

	if handler.produtoService == nil {
		t.Error("produtoService não foi inicializado")
	}
}

func TestProdutoHandler_CriarProduto(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockFunc       func(ctx context.Context, nome, descricao string, preco float64, categoria domain.Categoria) (*domain.Produto, error)
		expectedStatus int
	}{
		{
			name: "Sucesso ao criar produto",
			requestBody: CriarProdutoRequest{
				Nome:      "Hamburguer",
				Descricao: "Hamburguer artesanal",
				Preco:     25.90,
				Categoria: domain.CategoriaLanche,
			},
			mockFunc: func(ctx context.Context, nome, descricao string, preco float64, categoria domain.Categoria) (*domain.Produto, error) {
				return &domain.Produto{
					ID:         "123",
					Nome:       nome,
					Descricao:  descricao,
					Preco:      preco,
					Categoria:  categoria,
					Disponivel: true,
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
				}, nil
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Erro ao criar produto",
			requestBody: CriarProdutoRequest{
				Nome:      "Hamburguer",
				Descricao: "Hamburguer artesanal",
				Preco:     25.90,
				Categoria: domain.CategoriaLanche,
			},
			mockFunc: func(ctx context.Context, nome, descricao string, preco float64, categoria domain.Categoria) (*domain.Produto, error) {
				return nil, errors.New("erro ao salvar produto")
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
			mockService := &MockProdutoService{
				CriarProdutoFunc: tt.mockFunc,
			}
			handler := NovoProdutoHandler(mockService)

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

			req := httptest.NewRequest(http.MethodPost, "/produtos", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			handler.CriarProduto(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Status = %d, esperado %d", w.Code, tt.expectedStatus)
			}
		})
	}
}

func TestProdutoHandler_BuscarProdutoPorID(t *testing.T) {
	tests := []struct {
		name           string
		produtoID      string
		mockFunc       func(ctx context.Context, id string) (*domain.Produto, error)
		expectedStatus int
	}{
		{
			name:      "Sucesso ao buscar produto",
			produtoID: "123",
			mockFunc: func(ctx context.Context, id string) (*domain.Produto, error) {
				return &domain.Produto{
					ID:        id,
					Nome:      "Hamburguer",
					Preco:     25.90,
					Categoria: domain.CategoriaLanche,
				}, nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:      "Produto não encontrado",
			produtoID: "999",
			mockFunc: func(ctx context.Context, id string) (*domain.Produto, error) {
				return nil, nil
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:      "Erro ao buscar produto",
			produtoID: "123",
			mockFunc: func(ctx context.Context, id string) (*domain.Produto, error) {
				return nil, errors.New("erro no banco")
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockProdutoService{
				BuscarProdutoPorIDFunc: tt.mockFunc,
			}
			handler := NovoProdutoHandler(mockService)

			req := httptest.NewRequest(http.MethodGet, "/produtos/"+tt.produtoID, nil)
			req = mux.SetURLVars(req, map[string]string{"id": tt.produtoID})
			w := httptest.NewRecorder()

			handler.BuscarProdutoPorID(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Status = %d, esperado %d", w.Code, tt.expectedStatus)
			}
		})
	}
}

func TestProdutoHandler_ListarProdutos(t *testing.T) {
	tests := []struct {
		name                  string
		categoria             string
		mockListarFunc        func(ctx context.Context) ([]*domain.Produto, error)
		mockListarCatFunc     func(ctx context.Context, categoria domain.Categoria) ([]*domain.Produto, error)
		expectedStatus        int
	}{
		{
			name:      "Sucesso ao listar todos os produtos",
			categoria: "",
			mockListarFunc: func(ctx context.Context) ([]*domain.Produto, error) {
				return []*domain.Produto{
					{ID: "1", Nome: "Produto 1"},
					{ID: "2", Nome: "Produto 2"},
				}, nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:      "Sucesso ao listar por categoria",
			categoria: "LANCHE",
			mockListarCatFunc: func(ctx context.Context, categoria domain.Categoria) ([]*domain.Produto, error) {
				return []*domain.Produto{
					{ID: "1", Nome: "Hamburguer", Categoria: categoria},
				}, nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:      "Erro ao listar produtos",
			categoria: "",
			mockListarFunc: func(ctx context.Context) ([]*domain.Produto, error) {
				return nil, errors.New("erro no banco")
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:      "Erro ao listar por categoria",
			categoria: "LANCHE",
			mockListarCatFunc: func(ctx context.Context, categoria domain.Categoria) ([]*domain.Produto, error) {
				return nil, errors.New("erro no banco")
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockProdutoService{
				ListarProdutosFunc:             tt.mockListarFunc,
				ListarProdutosPorCategoriaFunc: tt.mockListarCatFunc,
			}
			handler := NovoProdutoHandler(mockService)

			url := "/produtos"
			if tt.categoria != "" {
				url += "?categoria=" + tt.categoria
			}

			req := httptest.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()

			handler.ListarProdutos(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Status = %d, esperado %d", w.Code, tt.expectedStatus)
			}
		})
	}
}

func TestProdutoHandler_AtualizarProduto(t *testing.T) {
	tests := []struct {
		name              string
		produtoID         string
		requestBody       interface{}
		mockBuscarFunc    func(ctx context.Context, id string) (*domain.Produto, error)
		mockAtualizarFunc func(ctx context.Context, produto *domain.Produto) error
		expectedStatus    int
	}{
		{
			name:      "Sucesso ao atualizar produto",
			produtoID: "123",
			requestBody: AtualizarProdutoRequest{
				Nome:       "Hamburguer Atualizado",
				Descricao:  "Nova descrição",
				Preco:      30.00,
				Categoria:  domain.CategoriaLanche,
				Disponivel: true,
			},
			mockBuscarFunc: func(ctx context.Context, id string) (*domain.Produto, error) {
				return &domain.Produto{
					ID:   id,
					Nome: "Hamburguer",
				}, nil
			},
			mockAtualizarFunc: func(ctx context.Context, produto *domain.Produto) error {
				return nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:      "Produto não encontrado",
			produtoID: "999",
			requestBody: AtualizarProdutoRequest{
				Nome: "Produto",
			},
			mockBuscarFunc: func(ctx context.Context, id string) (*domain.Produto, error) {
				return nil, nil
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Erro ao decodificar requisição",
			produtoID:      "123",
			requestBody:    "invalid json",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:      "Erro ao buscar produto",
			produtoID: "123",
			requestBody: AtualizarProdutoRequest{
				Nome: "Produto",
			},
			mockBuscarFunc: func(ctx context.Context, id string) (*domain.Produto, error) {
				return nil, errors.New("erro no banco")
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:      "Erro ao atualizar produto",
			produtoID: "123",
			requestBody: AtualizarProdutoRequest{
				Nome: "Produto",
			},
			mockBuscarFunc: func(ctx context.Context, id string) (*domain.Produto, error) {
				return &domain.Produto{ID: id}, nil
			},
			mockAtualizarFunc: func(ctx context.Context, produto *domain.Produto) error {
				return errors.New("erro ao atualizar")
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockProdutoService{
				BuscarProdutoPorIDFunc: tt.mockBuscarFunc,
				AtualizarProdutoFunc:   tt.mockAtualizarFunc,
			}
			handler := NovoProdutoHandler(mockService)

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

			req := httptest.NewRequest(http.MethodPut, "/produtos/"+tt.produtoID, bytes.NewBuffer(body))
			req = mux.SetURLVars(req, map[string]string{"id": tt.produtoID})
			w := httptest.NewRecorder()

			handler.AtualizarProduto(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Status = %d, esperado %d", w.Code, tt.expectedStatus)
			}
		})
	}
}

func TestProdutoHandler_DeletarProduto(t *testing.T) {
	tests := []struct {
		name           string
		produtoID      string
		mockFunc       func(ctx context.Context, id string) error
		expectedStatus int
	}{
		{
			name:      "Sucesso ao deletar produto",
			produtoID: "123",
			mockFunc: func(ctx context.Context, id string) error {
				return nil
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:      "Erro ao deletar produto",
			produtoID: "123",
			mockFunc: func(ctx context.Context, id string) error {
				return errors.New("erro no banco")
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockProdutoService{
				DeletarProdutoFunc: tt.mockFunc,
			}
			handler := NovoProdutoHandler(mockService)

			req := httptest.NewRequest(http.MethodDelete, "/produtos/"+tt.produtoID, nil)
			req = mux.SetURLVars(req, map[string]string{"id": tt.produtoID})
			w := httptest.NewRecorder()

			handler.DeletarProduto(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Status = %d, esperado %d", w.Code, tt.expectedStatus)
			}
		})
	}
}
