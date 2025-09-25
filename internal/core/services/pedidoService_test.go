package services

import (
	"context"
	"testing"

	"soat-fiap/internal/core/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPedidoRepository struct {
	mock.Mock
}

func (m *MockPedidoRepository) Criar(ctx context.Context, pedido *domain.Pedido) error {
	args := m.Called(ctx, pedido)
	return args.Error(0)
}
func (m *MockPedidoRepository) BuscarPorID(ctx context.Context, id string) (*domain.Pedido, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.Pedido), args.Error(1)
}
func (m *MockPedidoRepository) Listar(ctx context.Context) ([]*domain.Pedido, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*domain.Pedido), args.Error(1)
}
func (m *MockPedidoRepository) ListarPorStatus(ctx context.Context, status domain.StatusPedido) ([]*domain.Pedido, error) {
	args := m.Called(ctx, status)
	return args.Get(0).([]*domain.Pedido), args.Error(1)
}
func (m *MockPedidoRepository) ListarPorCliente(ctx context.Context, clienteID string) ([]*domain.Pedido, error) {
	args := m.Called(ctx, clienteID)
	return args.Get(0).([]*domain.Pedido), args.Error(1)
}
func (m *MockPedidoRepository) Atualizar(ctx context.Context, pedido *domain.Pedido) error {
	args := m.Called(ctx, pedido)
	return args.Error(0)
}

type MockProdutoRepositoryPedido struct {
	mock.Mock
}

func (m *MockProdutoRepositoryPedido) BuscarPorID(ctx context.Context, id string) (*domain.Produto, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.Produto), args.Error(1)
}
func (m *MockProdutoRepositoryPedido) Criar(ctx context.Context, produto *domain.Produto) error {
	return nil
}
func (m *MockProdutoRepositoryPedido) Listar(ctx context.Context) ([]*domain.Produto, error) {
	return nil, nil
}
func (m *MockProdutoRepositoryPedido) ListarPorCategoria(ctx context.Context, categoria domain.Categoria) ([]*domain.Produto, error) {
	return nil, nil
}
func (m *MockProdutoRepositoryPedido) Atualizar(ctx context.Context, produto *domain.Produto) error {
	return nil
}
func (m *MockProdutoRepositoryPedido) Deletar(ctx context.Context, id string) error {
	return nil
}

type MockPagamentoRepositoryPedido struct {
	mock.Mock
}

func (m *MockPagamentoRepositoryPedido) Criar(ctx context.Context, pagamento *domain.Pagamento) error {
	return nil
}
func (m *MockPagamentoRepositoryPedido) BuscarPorID(ctx context.Context, id string) (*domain.Pagamento, error) {
	return nil, nil
}
func (m *MockPagamentoRepositoryPedido) BuscarPorPedidoID(ctx context.Context, pedidoID string) (*domain.Pagamento, error) {
	return nil, nil
}
func (m *MockPagamentoRepositoryPedido) BuscarPorExternalID(ctx context.Context, externalID string) (*domain.Pagamento, error) {
	return nil, nil
}
func (m *MockPagamentoRepositoryPedido) Atualizar(ctx context.Context, pagamento *domain.Pagamento) error {
	return nil
}
func (m *MockPagamentoRepositoryPedido) ListarExpirados(ctx context.Context) ([]*domain.Pagamento, error) {
	return nil, nil
}
func (m *MockPagamentoRepositoryPedido) ListarPorStatus(ctx context.Context, status domain.StatusPagamento) ([]*domain.Pagamento, error) {
	return nil, nil
}

func TestCriarPedido_Sucesso(t *testing.T) {
	pedidoRepo := new(MockPedidoRepository)
	produtoRepo := new(MockProdutoRepositoryPedido)
	pagamentoRepo := new(MockPagamentoRepositoryPedido)
	service := NovoPedidoService(pedidoRepo, produtoRepo, pagamentoRepo)
	ctx := context.Background()

	item := domain.ItemPedido{ProdutoID: "1", Quantidade: 2}
	produto := &domain.Produto{ID: "1", Nome: "Coca-Cola", Preco: 5.0, Disponivel: true}

	produtoRepo.On("BuscarPorID", ctx, "1").Return(produto, nil)
	pedidoRepo.On("Criar", ctx, mock.AnythingOfType("*domain.Pedido")).Return(nil)

	pedido, err := service.CriarPedido(ctx, nil, []domain.ItemPedido{item})
	assert.NoError(t, err)
	assert.NotNil(t, pedido)
	assert.Equal(t, domain.StatusRecebido, pedido.Status)
}

func TestCriarPedido_ProdutoNaoEncontrado(t *testing.T) {
	pedidoRepo := new(MockPedidoRepository)
	produtoRepo := new(MockProdutoRepositoryPedido)
	pagamentoRepo := new(MockPagamentoRepositoryPedido)
	service := NovoPedidoService(pedidoRepo, produtoRepo, pagamentoRepo)
	ctx := context.Background()

	produtoRepo.On("BuscarPorID", ctx, "1").Return((*domain.Produto)(nil), nil)

	pedido, err := service.CriarPedido(ctx, nil, []domain.ItemPedido{{ProdutoID: "1", Quantidade: 1}})
	assert.Error(t, err)
	assert.Nil(t, pedido)
}

func TestCriarPedido_ProdutoIndisponivel(t *testing.T) {
	pedidoRepo := new(MockPedidoRepository)
	produtoRepo := new(MockProdutoRepositoryPedido)
	pagamentoRepo := new(MockPagamentoRepositoryPedido)
	service := NovoPedidoService(pedidoRepo, produtoRepo, pagamentoRepo)
	ctx := context.Background()

	produtoRepo.On("BuscarPorID", ctx, "1").Return(&domain.Produto{ID: "1", Nome: "Coca-Cola", Preco: 5.0, Disponivel: false}, nil)

	pedido, err := service.CriarPedido(ctx, nil, []domain.ItemPedido{{ProdutoID: "1", Quantidade: 1}})
	assert.Error(t, err)
	assert.Nil(t, pedido)
}

func TestBuscarPedidoPorID_Sucesso(t *testing.T) {
	pedidoRepo := new(MockPedidoRepository)
	produtoRepo := new(MockProdutoRepositoryPedido)
	pagamentoRepo := new(MockPagamentoRepositoryPedido)
	service := NovoPedidoService(pedidoRepo, produtoRepo, pagamentoRepo)
	ctx := context.Background()

	pedidoRepo.On("BuscarPorID", ctx, "1").Return(&domain.Pedido{ID: "1"}, nil)

	pedido, err := service.BuscarPedidoPorID(ctx, "1")
	assert.NoError(t, err)
	assert.NotNil(t, pedido)
	assert.Equal(t, "1", pedido.ID)
}

func TestAtualizarStatusPedido_StatusInvalido(t *testing.T) {
	pedidoRepo := new(MockPedidoRepository)
	produtoRepo := new(MockProdutoRepositoryPedido)
	pagamentoRepo := new(MockPagamentoRepositoryPedido)
	service := NovoPedidoService(pedidoRepo, produtoRepo, pagamentoRepo)
	ctx := context.Background()

	err := service.AtualizarStatusPedido(ctx, "1", "INVALIDO")
	assert.Error(t, err)
}

func TestAtualizarStatusPedido_PedidoNaoEncontrado(t *testing.T) {
	pedidoRepo := new(MockPedidoRepository)
	produtoRepo := new(MockProdutoRepositoryPedido)
	pagamentoRepo := new(MockPagamentoRepositoryPedido)
	service := NovoPedidoService(pedidoRepo, produtoRepo, pagamentoRepo)
	ctx := context.Background()

	pedidoRepo.On("BuscarPorID", ctx, "1").Return((*domain.Pedido)(nil), nil)

	err := service.AtualizarStatusPedido(ctx, "1", domain.StatusPronto)
	assert.Error(t, err)
}
