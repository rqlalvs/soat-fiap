package ports

import (
	"context"
	"soat-fiap/internal/core/domain"
)

type PagamentoRepository interface {
	Criar(ctx context.Context, pagamento *domain.Pagamento) error
	BuscarPorID(ctx context.Context, id string) (*domain.Pagamento, error)
	BuscarPorPedidoID(ctx context.Context, pedidoID string) (*domain.Pagamento, error)
	BuscarPorExternalID(ctx context.Context, externalID string) (*domain.Pagamento, error)
	Atualizar(ctx context.Context, pagamento *domain.Pagamento) error
	ListarPorStatus(ctx context.Context, status domain.StatusPagamento) ([]*domain.Pagamento, error)
	ListarExpirados(ctx context.Context) ([]*domain.Pagamento, error)
}
