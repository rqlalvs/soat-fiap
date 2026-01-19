package services

import (
	"context"
	"soat-fiap/internal/core/domain"
	"soat-fiap/pkg/clients"
	"time"
)

type ProdutoService struct {
	produtoClient *clients.ProdutoClient
}

func NovoProdutoService(produtoClient *clients.ProdutoClient) *ProdutoService {
	return &ProdutoService{
		produtoClient: produtoClient,
	}
}

func produtoFromClient(p *clients.Produto) *domain.Produto {
	if p == nil {
		return nil
	}
	return &domain.Produto{
		ID:         p.ID,
		Nome:       p.Nome,
		Descricao:  p.Descricao,
		Preco:      p.Preco,
		Categoria:   domain.Categoria(p.Categoria),
		Disponivel: p.Disponivel,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

func (s *ProdutoService) CriarProduto(ctx context.Context, nome, descricao string, preco float64, categoria domain.Categoria) (*domain.Produto, error) {
	produtoClient, err := s.produtoClient.CriarProduto(ctx, nome, descricao, preco, string(categoria), true)
	if err != nil {
		return nil, err
	}

	return produtoFromClient(produtoClient), nil
}

func (s *ProdutoService) BuscarProdutoPorID(ctx context.Context, id string) (*domain.Produto, error) {
	produtoClient, err := s.produtoClient.BuscarPorID(ctx, id)
	if err != nil {
		return nil, err
	}
	return produtoFromClient(produtoClient), nil
}

func (s *ProdutoService) ListarProdutos(ctx context.Context) ([]*domain.Produto, error) {
	produtosClient, err := s.produtoClient.ListarProdutos(ctx)
	if err != nil {
		return nil, err
	}

	produtos := make([]*domain.Produto, 0, len(produtosClient))
	for _, p := range produtosClient {
		produtos = append(produtos, produtoFromClient(p))
	}

	return produtos, nil
}

func (s *ProdutoService) ListarProdutosPorCategoria(ctx context.Context, categoria domain.Categoria) ([]*domain.Produto, error) {
	produtos, err := s.ListarProdutos(ctx)
	if err != nil {
		return nil, err
	}

	var produtosFiltrados []*domain.Produto
	for _, p := range produtos {
		if p.Categoria == categoria {
			produtosFiltrados = append(produtosFiltrados, p)
		}
	}

	return produtosFiltrados, nil
}

func (s *ProdutoService) AtualizarProduto(ctx context.Context, produto *domain.Produto) error {
	err := produto.Validar()
	if err != nil {
		return err
	}

	return s.produtoClient.AtualizarProduto(ctx, produto.ID, produto.Nome, produto.Descricao, produto.Preco, string(produto.Categoria), produto.Disponivel)
}

func (s *ProdutoService) DeletarProduto(ctx context.Context, id string) error {
	return s.produtoClient.DeletarProduto(ctx, id)
}
