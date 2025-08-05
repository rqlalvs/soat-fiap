package services

import (
	"context"
	"errors"
	"fmt"
	"soat-fiap/internal/core/domain"
	"soat-fiap/internal/core/ports"

	"github.com/google/uuid"
)

type PedidoService struct {
	pedidoRepository    ports.PedidoRepository
	produtoRepository   ports.ProdutoRepository
	pagamentoRepository ports.PagamentoRepository
}

func NovoPedidoService(pedidoRepository ports.PedidoRepository, produtoRepository ports.ProdutoRepository, pagamentoRepository ports.PagamentoRepository) *PedidoService {
	return &PedidoService{
		pedidoRepository:    pedidoRepository,
		produtoRepository:   produtoRepository,
		pagamentoRepository: pagamentoRepository,
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

	err = s.validarTransicaoStatus(ctx, pedido, status)
	if err != nil {
		return err
	}

	pedido.AtualizarStatus(status)
	return s.pedidoRepository.Atualizar(ctx, pedido)
}

func (s *PedidoService) validarTransicaoStatus(ctx context.Context, pedido *domain.Pedido, novoStatus domain.StatusPedido) error {

	if novoStatus == domain.StatusEmPreparacao && pedido.Status == domain.StatusRecebido {
		pagamento, err := s.pagamentoRepository.BuscarPorPedidoID(ctx, pedido.ID)
		if err != nil {
			return fmt.Errorf("erro ao verificar pagamento: %w", err)
		}
		
		if pagamento == nil {
			return errors.New("não é possível iniciar preparação: pagamento não encontrado")
		}
		
		if pagamento.Status != domain.StatusPagamentoAprovado {
			return fmt.Errorf("não é possível iniciar preparação: pagamento não aprovado (status: %s)", pagamento.Status)
		}
	}
	
	if novoStatus != domain.StatusCancelado {
		if !s.isTransicaoValida(pedido.Status, novoStatus) {
			return fmt.Errorf("transição de status inválida: de %s para %s", pedido.Status, novoStatus)
		}
	}
	
	return nil
}

func (s *PedidoService) isTransicaoValida(statusAtual, novoStatus domain.StatusPedido) bool {
	transicoesValidas := map[domain.StatusPedido][]domain.StatusPedido{
		domain.StatusRecebido: {
			domain.StatusEmPreparacao,
			domain.StatusCancelado,
		},
		domain.StatusEmPreparacao: {
			domain.StatusPronto,
			domain.StatusCancelado,
		},
		domain.StatusPronto: {
			domain.StatusFinalizado,
		},
		domain.StatusFinalizado: {},
		domain.StatusCancelado:  {},
	}
	
	statusPermitidos, existe := transicoesValidas[statusAtual]
	if !existe {
		return false
	}
	
	for _, statusPermitido := range statusPermitidos {
		if statusPermitido == novoStatus {
			return true
		}
	}
	
	return false
}
