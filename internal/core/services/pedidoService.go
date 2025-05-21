package services

import (
	"context"
	"errors"
	"soat-fiap/internal/core/domain"
	"soat-fiap/internal/core/ports"

	"github.com/google/uuid"
)

type PedidoService struct {
	pedidoRepository  ports.PedidoRepository
	produtoRepository ports.ProdutoRepository
}

func NovoPedidoService(pedidoRepository ports.PedidoRepository, produtoRepository ports.ProdutoRepository) *PedidoService {
	return &PedidoService{
		pedidoRepository:  pedidoRepository,
		produtoRepository: produtoRepository,
	}
}

func (s *PedidoService) CriarPedido(ctx context.Context, clienteID *string, itens []domain.ItemPedido) (*domain.Pedido, error) {

	for i, item := range itens {
		produto, err := s.produtoRepository.BuscarPorID(ctx, item.ProdutoID)
		if err != nil {
			return nil, err
		}
		if produto == nil {
			return nil, errors.New("produto não encontrado: " + item.ProdutoID)
		}
		if !produto.Disponivel {
			return nil, errors.New("produto não disponível: " + produto.Nome)
		}

		itens[i].Nome = produto.Nome
		itens[i].Preco = produto.Preco
	}

	id := uuid.New().String()
	pedido, err := domain.NovoPedido(id, clienteID, itens)
	if err != nil {
		return nil, err
	}

	err = s.pedidoRepository.Criar(ctx, pedido)
	if err != nil {
		return nil, err
	}

	return pedido, nil
}

func (s *PedidoService) BuscarPedidoPorID(ctx context.Context, id string) (*domain.Pedido, error) {
	return s.pedidoRepository.BuscarPorID(ctx, id)
}

func (s *PedidoService) ListarPedidos(ctx context.Context) ([]*domain.Pedido, error) {
	return s.pedidoRepository.Listar(ctx)
}

func (s *PedidoService) ListarPedidosPorStatus(ctx context.Context, status domain.StatusPedido) ([]*domain.Pedido, error) {
	if !domain.IsStatusValido(status) {
		return nil, errors.New("status inválido")
	}
	return s.pedidoRepository.ListarPorStatus(ctx, status)
}

func (s *PedidoService) ListarPedidosPorCliente(ctx context.Context, clienteID string) ([]*domain.Pedido, error) {
	return s.pedidoRepository.ListarPorCliente(ctx, clienteID)
}

func (s *PedidoService) AtualizarStatusPedido(ctx context.Context, id string, status domain.StatusPedido) error {
	if !domain.IsStatusValido(status) {
		return errors.New("status inválido")
	}

	pedido, err := s.pedidoRepository.BuscarPorID(ctx, id)
	if err != nil {
		return err
	}
	if pedido == nil {
		return errors.New("pedido não encontrado")
	}

	pedido.AtualizarStatus(status)
	return s.pedidoRepository.Atualizar(ctx, pedido)
}
