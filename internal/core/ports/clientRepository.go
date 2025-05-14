package ports

import (
	"context"

	"soat-fiap/internal/core/domain"
)

type ClienteRepository interface {
	Criar(ctx context.Context, cliente *domain.Cliente) error
	BuscarPorID(ctx context.Context, id string) (*domain.Cliente, error)
	BuscarPorCPF(ctx context.Context, cpf string) (*domain.Cliente, error)
	Listar(ctx context.Context) ([]*domain.Cliente, error)
	Atualizar(ctx context.Context, cliente *domain.Cliente) error
	Deletar(ctx context.Context, id string) error
}
