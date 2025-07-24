package services

import (
	"context"
	"testing"

	"soat-fiap/internal/core/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockProdutoRepository struct {
	mock.Mock
}

func (m *MockProdutoRepository) Criar(ctx context.Context, produto *domain.Produto) error {
	args := m.Called(ctx, produto)
	return args.Error(0)
}
func (m *MockProdutoRepository) BuscarPorID(ctx context.Context, id string) (*domain.Produto, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.Produto), args.Error(1)
}
func (m *MockProdutoRepository) Listar(ctx context.Context) ([]*domain.Produto, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*domain.Produto), args.Error(1)
}
func (m *MockProdutoRepository) ListarPorCategoria(ctx context.Context, categoria domain.Categoria) ([]*domain.Produto, error) {
	args := m.Called(ctx, categoria)
	return args.Get(0).([]*domain.Produto), args.Error(1)
}
func (m *MockProdutoRepository) Atualizar(ctx context.Context, produto *domain.Produto) error {
	args := m.Called(ctx, produto)
	return args.Error(0)
}
func (m *MockProdutoRepository) Deletar(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestCriarProduto_Sucesso(t *testing.T) {
	repo := new(MockProdutoRepository)
	service := NovoProdutoService(repo)
	ctx := context.Background()

	repo.On("Criar", ctx, mock.AnythingOfType("*domain.Produto")).Return(nil)

	produto, err := service.CriarProduto(ctx, "Coca-Cola", "Refrigerante", 5.0, domain.CategoriaBebida)
	assert.NoError(t, err)
	assert.NotNil(t, produto)
	assert.Equal(t, "Coca-Cola", produto.Nome)
}

func TestCriarProduto_CategoriaInvalida(t *testing.T) {
	repo := new(MockProdutoRepository)
	service := NovoProdutoService(repo)
	ctx := context.Background()

	produto, err := service.CriarProduto(ctx, "Coca-Cola", "Refrigerante", 5.0, "INVALIDA")
	assert.Error(t, err)
	assert.Nil(t, produto)
}

func TestBuscarProdutoPorID_Sucesso(t *testing.T) {
	repo := new(MockProdutoRepository)
	service := NovoProdutoService(repo)
	ctx := context.Background()

	repo.On("BuscarPorID", ctx, "1").Return(&domain.Produto{ID: "1", Nome: "Coca-Cola"}, nil)

	produto, err := service.BuscarProdutoPorID(ctx, "1")
	assert.NoError(t, err)
	assert.NotNil(t, produto)
	assert.Equal(t, "Coca-Cola", produto.Nome)
}

func TestAtualizarProduto_ProdutoNaoEncontrado(t *testing.T) {
	repo := new(MockProdutoRepository)
	service := NovoProdutoService(repo)
	ctx := context.Background()

	repo.On("BuscarPorID", ctx, "1").Return((*domain.Produto)(nil), nil)

	err := service.AtualizarProduto(ctx, &domain.Produto{ID: "1", Nome: "Novo"})
	assert.Error(t, err)
}

func TestDeletarProduto_Sucesso(t *testing.T) {
	repo := new(MockProdutoRepository)
	service := NovoProdutoService(repo)
	ctx := context.Background()

	repo.On("Deletar", ctx, "1").Return(nil)

	err := service.DeletarProduto(ctx, "1")
	assert.NoError(t, err)
}
