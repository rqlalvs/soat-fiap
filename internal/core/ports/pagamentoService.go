package ports

import (
	"context"
	"soat-fiap/internal/core/domain"
)

type PagamentoService interface {
	CriarPagamento(ctx context.Context, pedidoID string, valor float64, tipo domain.TipoPagamento) (*domain.Pagamento, error)
	BuscarPagamentoPorID(ctx context.Context, id string) (*domain.Pagamento, error)
	BuscarPagamentoPorPedidoID(ctx context.Context, pedidoID string) (*domain.Pagamento, error)
	ProcessarWebhookPagamento(ctx context.Context, externalID string, status domain.StatusPagamento) error
	GerarQRCodePIX(ctx context.Context, pagamentoID string) (string, error)
	ConsultarStatusPagamento(ctx context.Context, pedidoID string) (*domain.Pagamento, error)
	CancelarPagamento(ctx context.Context, pagamentoID string) error
	ListarPagamentosExpirados(ctx context.Context) ([]*domain.Pagamento, error)
}
