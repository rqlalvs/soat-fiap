package ports

import (
	"context"

	"soat-fiap/internal/core/domain"
)

type ClienteService interface {
	CriarCliente(ctx context.Context, nome, cpf, email, telefone string) (*domain.Cliente, error)
	BuscarClientePorID(ctx context.Context, id string) (*domain.Cliente, error)
	BuscarClientePorCPF(ctx context.Context, cpf string) (*domain.Cliente, error)
	ListarClientes(ctx context.Context) ([]*domain.Cliente, error)
	AtualizarCliente(ctx context.Context, cliente *domain.Cliente) error
	DeletarCliente(ctx context.Context, id string) error
}
