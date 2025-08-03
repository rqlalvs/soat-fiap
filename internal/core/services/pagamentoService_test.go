package services

import (
	"context"
	"testing"
	"time"

	"soat-fiap/internal/core/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock do PagamentoRepository
type MockPagamentoRepository struct {
	mock.Mock
}

func (m *MockPagamentoRepository) Criar(ctx context.Context, pagamento *domain.Pagamento) error {
	args := m.Called(ctx, pagamento)
	return args.Error(0)
}

func (m *MockPagamentoRepository) BuscarPorID(ctx context.Context, id string) (*domain.Pagamento, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.Pagamento), args.Error(1)
}

func (m *MockPagamentoRepository) BuscarPorPedidoID(ctx context.Context, pedidoID string) (*domain.Pagamento, error) {
	args := m.Called(ctx, pedidoID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Pagamento), args.Error(1)
}

func (m *MockPagamentoRepository) BuscarPorExternalID(ctx context.Context, externalID string) (*domain.Pagamento, error) {
	args := m.Called(ctx, externalID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Pagamento), args.Error(1)
}

func (m *MockPagamentoRepository) Atualizar(ctx context.Context, pagamento *domain.Pagamento) error {
	args := m.Called(ctx, pagamento)
	return args.Error(0)
}

func (m *MockPagamentoRepository) ListarPorStatus(ctx context.Context, status domain.StatusPagamento) ([]*domain.Pagamento, error) {
	args := m.Called(ctx, status)
	return args.Get(0).([]*domain.Pagamento), args.Error(1)
}

func (m *MockPagamentoRepository) ListarExpirados(ctx context.Context) ([]*domain.Pagamento, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*domain.Pagamento), args.Error(1)
}

type MockPedidoService struct {
	mock.Mock
}

func (m *MockPedidoService) CriarPedido(ctx context.Context, clienteID *string, itens []domain.ItemPedido) (*domain.Pedido, error) {
	args := m.Called(ctx, clienteID, itens)
	return args.Get(0).(*domain.Pedido), args.Error(1)
}

func (m *MockPedidoService) BuscarPedidoPorID(ctx context.Context, id string) (*domain.Pedido, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Pedido), args.Error(1)
}

func (m *MockPedidoService) ListarPedidos(ctx context.Context) ([]*domain.Pedido, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*domain.Pedido), args.Error(1)
}

func (m *MockPedidoService) ListarPedidosPorStatus(ctx context.Context, status domain.StatusPedido) ([]*domain.Pedido, error) {
	args := m.Called(ctx, status)
	return args.Get(0).([]*domain.Pedido), args.Error(1)
}

func (m *MockPedidoService) ListarPedidosPorCliente(ctx context.Context, clienteID string) ([]*domain.Pedido, error) {
	args := m.Called(ctx, clienteID)
	return args.Get(0).([]*domain.Pedido), args.Error(1)
}

func (m *MockPedidoService) AtualizarStatusPedido(ctx context.Context, id string, status domain.StatusPedido) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func TestPagamentoService_CriarPagamento(t *testing.T) {
	ctx := context.Background()
	mockPagamentoRepo := new(MockPagamentoRepository)
	mockPedidoService := new(MockPedidoService)

	service := NovoPagamentoService(mockPagamentoRepo, mockPedidoService)

	pedido := &domain.Pedido{
		ID:         "pedido-123",
		ValorTotal: 25.50,
		Status:     domain.StatusRecebido,
		CreatedAt:  time.Now(),
	}

	mockPedidoService.On("BuscarPedidoPorID", ctx, "pedido-123").Return(pedido, nil)
	mockPagamentoRepo.On("BuscarPorPedidoID", ctx, "pedido-123").Return(nil, nil) // Não existe pagamento
	mockPagamentoRepo.On("Criar", ctx, mock.AnythingOfType("*domain.Pagamento")).Return(nil)

	pagamento, err := service.CriarPagamento(ctx, "pedido-123", 25.50, domain.TipoPagamentoPix)

	assert.NoError(t, err)
	assert.NotNil(t, pagamento)
	assert.Equal(t, "pedido-123", pagamento.PedidoID)
	assert.Equal(t, 25.50, pagamento.Valor)
	assert.Equal(t, domain.TipoPagamentoPix, pagamento.Tipo)
	assert.Equal(t, domain.StatusPagamentoPendente, pagamento.Status)
	assert.NotEmpty(t, pagamento.QRCodeData)
	assert.NotEmpty(t, pagamento.ExternalID)
	assert.NotNil(t, pagamento.ExpiresAt)

	mockPedidoService.AssertExpectations(t)
	mockPagamentoRepo.AssertExpectations(t)
}

func TestPagamentoService_ProcessarWebhookPagamento_Aprovado(t *testing.T) {
	ctx := context.Background()
	mockPagamentoRepo := new(MockPagamentoRepository)
	mockPedidoService := new(MockPedidoService)

	service := NovoPagamentoService(mockPagamentoRepo, mockPedidoService)

	pagamento := &domain.Pagamento{
		ID:         "pagamento-123",
		PedidoID:   "pedido-123",
		Valor:      25.50,
		Status:     domain.StatusPagamentoPendente,
		ExternalID: "MP_12345678",
		CreatedAt:  time.Now(),
	}

	mockPagamentoRepo.On("BuscarPorExternalID", ctx, "MP_12345678").Return(pagamento, nil)
	mockPagamentoRepo.On("Atualizar", ctx, mock.AnythingOfType("*domain.Pagamento")).Return(nil)
	mockPedidoService.On("AtualizarStatusPedido", ctx, "pedido-123", domain.StatusEmPreparacao).Return(nil)

	err := service.ProcessarWebhookPagamento(ctx, "MP_12345678", domain.StatusPagamentoAprovado)

	assert.NoError(t, err)
	assert.Equal(t, domain.StatusPagamentoAprovado, pagamento.Status)
	assert.NotNil(t, pagamento.ProcessedAt)

	mockPagamentoRepo.AssertExpectations(t)
	mockPedidoService.AssertExpectations(t)
}

func TestPagamentoService_ConsultarStatusPagamento_Expirado(t *testing.T) {
	ctx := context.Background()
	mockPagamentoRepo := new(MockPagamentoRepository)
	mockPedidoService := new(MockPedidoService)

	service := NovoPagamentoService(mockPagamentoRepo, mockPedidoService)

	expiresAt := time.Now().Add(-1 * time.Hour) 
	pagamento := &domain.Pagamento{
		ID:         "pagamento-123",
		PedidoID:   "pedido-123",
		Valor:      25.50,
		Status:     domain.StatusPagamentoPendente,
		Tipo:       domain.TipoPagamentoPix,
		ExpiresAt:  &expiresAt,
		CreatedAt:  time.Now().Add(-2 * time.Hour),
	}

	mockPagamentoRepo.On("BuscarPorPedidoID", ctx, "pedido-123").Return(pagamento, nil)
	mockPagamentoRepo.On("Atualizar", ctx, mock.AnythingOfType("*domain.Pagamento")).Return(nil)

	resultado, err := service.ConsultarStatusPagamento(ctx, "pedido-123")

	assert.NoError(t, err)
	assert.NotNil(t, resultado)
	assert.Equal(t, domain.StatusPagamentoCancelado, resultado.Status)

	mockPagamentoRepo.AssertExpectations(t)
}
