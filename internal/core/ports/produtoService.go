package ports

import (
	"context"
	"soat-fiap/internal/core/domain"
)

type ProdutoService interface {
	CriarProduto(ctx context.Context, nome, descricao string, preco float64, categoria domain.Categoria) (*domain.Produto, error)
	BuscarProdutoPorID(ctx context.Context, id string) (*domain.Produto, error)
	ListarProdutos(ctx context.Context) ([]*domain.Produto, error)
	ListarProdutosPorCategoria(ctx context.Context, categoria domain.Categoria) ([]*domain.Produto, error)
	AtualizarProduto(ctx context.Context, produto *domain.Produto) error
	DeletarProduto(ctx context.Context, id string) error
}
