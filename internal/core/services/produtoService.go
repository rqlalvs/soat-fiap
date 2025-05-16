package services

import (
	"context"
	"errors"
	"soat-fiap/internal/core/domain"
	"soat-fiap/internal/core/ports"
	"time"

	"github.com/google/uuid"
)

type ProdutoService struct {
	repository ports.ProdutoRepository
}

func NovoProdutoService(repository ports.ProdutoRepository) *ProdutoService {
	return &ProdutoService{
		repository: repository,
	}
}

func (s *ProdutoService) CriarProduto(ctx context.Context, nome, descricao string, preco float64, categoria domain.Categoria) (*domain.Produto, error) {
	id := uuid.New().String()

	produto, err := domain.NovoProduto(id, nome, descricao, preco, categoria)
	if err != nil {
		return nil, err
	}

	err = s.repository.Criar(ctx, produto)
	if err != nil {
		return nil, err
	}

	return produto, nil
}

func (s *ProdutoService) BuscarProdutoPorID(ctx context.Context, id string) (*domain.Produto, error) {
	return s.repository.BuscarPorID(ctx, id)
}

func (s *ProdutoService) ListarProdutos(ctx context.Context) ([]*domain.Produto, error) {
	return s.repository.Listar(ctx)
}

func (s *ProdutoService) ListarProdutosPorCategoria(ctx context.Context, categoria domain.Categoria) ([]*domain.Produto, error) {
	if !domain.IsCategoriaValida(categoria) {
		return nil, errors.New("categoria inválida")
	}
	return s.repository.ListarPorCategoria(ctx, categoria)
}

func (s *ProdutoService) AtualizarProduto(ctx context.Context, produto *domain.Produto) error {
	produtoExistente, err := s.repository.BuscarPorID(ctx, produto.ID)
	if err != nil {
		return err
	}
	if produtoExistente == nil {
		return errors.New("produto não encontrado")
	}

	err = produto.Validar()
	if err != nil {
		return err
	}

	produto.UpdatedAt = time.Now()

	return s.repository.Atualizar(ctx, produto)
}

func (s *ProdutoService) DeletarProduto(ctx context.Context, id string) error {
	return s.repository.Deletar(ctx, id)
}
