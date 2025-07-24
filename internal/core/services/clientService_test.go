package services

import (
	"context"
	"testing"

	"soat-fiap/internal/core/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockClienteRepository struct {
	mock.Mock
}

func (m *MockClienteRepository) Criar(ctx context.Context, cliente *domain.Cliente) error {
	args := m.Called(ctx, cliente)
	return args.Error(0)
}
func (m *MockClienteRepository) BuscarPorID(ctx context.Context, id string) (*domain.Cliente, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.Cliente), args.Error(1)
}
func (m *MockClienteRepository) BuscarPorCPF(ctx context.Context, cpf string) (*domain.Cliente, error) {
	args := m.Called(ctx, cpf)
	return args.Get(0).(*domain.Cliente), args.Error(1)
}
func (m *MockClienteRepository) Listar(ctx context.Context) ([]*domain.Cliente, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*domain.Cliente), args.Error(1)
}
func (m *MockClienteRepository) Atualizar(ctx context.Context, cliente *domain.Cliente) error {
	args := m.Called(ctx, cliente)
	return args.Error(0)
}
func (m *MockClienteRepository) Deletar(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestCriarCliente_Sucesso(t *testing.T) {
	repo := new(MockClienteRepository)
	service := NovoClienteService(repo)
	ctx := context.Background()

	cpfValido := "39053344705"
	repo.On("BuscarPorCPF", ctx, cpfValido).Return((*domain.Cliente)(nil), nil)
	repo.On("Criar", ctx, mock.AnythingOfType("*domain.Cliente")).Return(nil)

	cliente, err := service.CriarCliente(ctx, "João", cpfValido, "joao@email.com", "11999999999")
	assert.NoError(t, err)
	assert.NotNil(t, cliente)
	assert.Equal(t, "João", cliente.Nome)
}

func TestCriarCliente_CPFExistente(t *testing.T) {
	repo := new(MockClienteRepository)
	service := NovoClienteService(repo)
	ctx := context.Background()

	cpfValido := "39053344705"
	repo.On("BuscarPorCPF", ctx, cpfValido).Return(&domain.Cliente{ID: "1", Nome: "Outro"}, nil)

	cliente, err := service.CriarCliente(ctx, "João", cpfValido, "joao@email.com", "11999999999")
	assert.Error(t, err)
	assert.Nil(t, cliente)
}

func TestBuscarClientePorID_Sucesso(t *testing.T) {
	repo := new(MockClienteRepository)
	service := NovoClienteService(repo)
	ctx := context.Background()

	repo.On("BuscarPorID", ctx, "1").Return(&domain.Cliente{ID: "1", Nome: "João"}, nil)

	cliente, err := service.BuscarClientePorID(ctx, "1")
	assert.NoError(t, err)
	assert.NotNil(t, cliente)
	assert.Equal(t, "João", cliente.Nome)
}

func TestAtualizarCliente_ClienteNaoEncontrado(t *testing.T) {
	repo := new(MockClienteRepository)
	service := NovoClienteService(repo)
	ctx := context.Background()

	repo.On("BuscarPorID", ctx, "1").Return((*domain.Cliente)(nil), nil)

	err := service.AtualizarCliente(ctx, &domain.Cliente{ID: "1", Nome: "Novo"})
	assert.Error(t, err)
}

func TestDeletarCliente_Sucesso(t *testing.T) {
	repo := new(MockClienteRepository)
	service := NovoClienteService(repo)
	ctx := context.Background()

	repo.On("Deletar", ctx, "1").Return(nil)

	err := service.DeletarCliente(ctx, "1")
	assert.NoError(t, err)
}
