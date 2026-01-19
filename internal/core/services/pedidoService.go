package services

import (
	"context"
	"errors"
	"fmt"
	"soat-fiap/internal/core/domain"
	"soat-fiap/pkg/clients"
	"time"
)

type PedidoService struct {
	pedidoClient    *clients.PedidoClient
	produtoClient   *clients.ProdutoClient
	clienteClient   *clients.ClienteClient
	pagamentoClient *clients.PagamentoClient
}

func NovoPedidoService(
	pedidoClient *clients.PedidoClient,
	produtoClient *clients.ProdutoClient,
	clienteClient *clients.ClienteClient,
	pagamentoClient *clients.PagamentoClient,
) *PedidoService {
	return &PedidoService{
		pedidoClient:    pedidoClient,
		produtoClient:   produtoClient,
		clienteClient:   clienteClient,
		pagamentoClient: pagamentoClient,
	}
}

func itemPedidoFromClient(i clients.ItemPedido) domain.ItemPedido {
	return domain.ItemPedido{
		ProdutoID:  i.ProdutoID,
		Nome:       i.Nome,
		Preco:      i.Preco,
		Quantidade: i.Quantidade,
	}
}

func itemPedidoToClient(i domain.ItemPedido) clients.ItemPedido {
	return clients.ItemPedido{
		ProdutoID:  i.ProdutoID,
		Nome:       i.Nome,
		Preco:      i.Preco,
		Quantidade: i.Quantidade,
	}
}

func pedidoFromClient(p *clients.Pedido) *domain.Pedido {
	if p == nil {
		return nil
	}

	itens := make([]domain.ItemPedido, 0, len(p.Itens))
	for _, item := range p.Itens {
		itens = append(itens, itemPedidoFromClient(item))
	}

	return &domain.Pedido{
		ID:         p.ID,
		ClienteID:  p.ClienteID,
		Itens:      itens,
		ValorTotal: p.ValorTotal,
		Status:     domain.StatusPedido(p.Status),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

func (s *PedidoService) CriarPedido(ctx context.Context, clienteID *string, itens []domain.ItemPedido) (*domain.Pedido, error) {
	if clienteID != nil && *clienteID != "" {
		if err := s.clienteClient.ValidarCliente(ctx, *clienteID); err != nil {
			return nil, fmt.Errorf("erro ao validar cliente: %w", err)
		}
	}

	itensClient := make([]clients.ItemPedido, 0, len(itens))
	for _, item := range itens {
		produto, err := s.produtoClient.ValidarProduto(ctx, item.ProdutoID)
		if err != nil {
			return nil, err
		}

		item.Nome = produto.Nome
		item.Preco = produto.Preco
		itensClient = append(itensClient, itemPedidoToClient(item))
	}

	pedidoClient, err := s.pedidoClient.CriarPedido(ctx, clienteID, itensClient)
	if err != nil {
		return nil, err
	}

	return pedidoFromClient(pedidoClient), nil
}

func (s *PedidoService) BuscarPedidoPorID(ctx context.Context, id string) (*domain.Pedido, error) {
	pedidoClient, err := s.pedidoClient.BuscarPorID(ctx, id)
	if err != nil {
		return nil, err
	}
	return pedidoFromClient(pedidoClient), nil
}

func (s *PedidoService) ListarPedidos(ctx context.Context) ([]*domain.Pedido, error) {
	pedidosClient, err := s.pedidoClient.ListarPedidos(ctx)
	if err != nil {
		return nil, err
	}

	pedidos := make([]*domain.Pedido, 0, len(pedidosClient))
	for _, p := range pedidosClient {
		pedidos = append(pedidos, pedidoFromClient(p))
	}

	return pedidos, nil
}

func (s *PedidoService) ListarPedidosPorStatus(ctx context.Context, status domain.StatusPedido) ([]*domain.Pedido, error) {
	if !domain.IsStatusValido(status) {
		return nil, errors.New("status inválido")
	}

	pedidosClient, err := s.pedidoClient.ListarPedidosPorStatus(ctx, string(status))
	if err != nil {
		return nil, err
	}

	pedidos := make([]*domain.Pedido, 0, len(pedidosClient))
	for _, p := range pedidosClient {
		pedidos = append(pedidos, pedidoFromClient(p))
	}

	return pedidos, nil
}

func (s *PedidoService) ListarPedidosPorCliente(ctx context.Context, clienteID string) ([]*domain.Pedido, error) {
	pedidosClient, err := s.pedidoClient.ListarPedidosPorCliente(ctx, clienteID)
	if err != nil {
		return nil, err
	}

	pedidos := make([]*domain.Pedido, 0, len(pedidosClient))
	for _, p := range pedidosClient {
		pedidos = append(pedidos, pedidoFromClient(p))
	}

	return pedidos, nil
}

func (s *PedidoService) AtualizarStatusPedido(ctx context.Context, id string, status domain.StatusPedido) error {
	if !domain.IsStatusValido(status) {
		return errors.New("status inválido")
	}

	pedido, err := s.pedidoClient.BuscarPorID(ctx, id)
	if err != nil {
		return err
	}
	if pedido == nil {
		return errors.New("pedido não encontrado")
	}

	pedidoDomain := pedidoFromClient(pedido)
	err = s.validarTransicaoStatus(ctx, pedidoDomain, status)
	if err != nil {
		return err
	}

	return s.pedidoClient.AtualizarStatusPedido(ctx, id, string(status))
}

func (s *PedidoService) validarTransicaoStatus(ctx context.Context, pedido *domain.Pedido, novoStatus domain.StatusPedido) error {

	if novoStatus == domain.StatusEmPreparacao && pedido.Status == domain.StatusRecebido {
		pagamentoClient, err := s.pagamentoClient.BuscarPorPedidoID(ctx, pedido.ID)
		if err != nil {
			return fmt.Errorf("erro ao verificar pagamento: %w", err)
		}
		
		if pagamentoClient == nil {
			return errors.New("não é possível iniciar preparação: pagamento não encontrado")
		}
		
		if pagamentoClient.Status != string(domain.StatusPagamentoAprovado) {
			return fmt.Errorf("não é possível iniciar preparação: pagamento não aprovado (status: %s)", pagamentoClient.Status)
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
