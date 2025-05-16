package ports

import (
	"context"
	"soat-fiap/internal/core/domain"
)

type ProdutoRepository interface {
	Criar(ctx context.Context, produto *domain.Produto) error
	BuscarPorID(ctx context.Context, id string) (*domain.Produto, error)
	Listar(ctx context.Context) ([]*domain.Produto, error)
	ListarPorCategoria(ctx context.Context, categoria domain.Categoria) ([]*domain.Produto, error)
	Atualizar(ctx context.Context, produto *domain.Produto) error
	Deletar(ctx context.Context, id string) error
}
