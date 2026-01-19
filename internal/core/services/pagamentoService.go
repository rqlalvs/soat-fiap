package services

import (
	"context"
	"errors"
	"fmt"
	"soat-fiap/internal/core/domain"
	"soat-fiap/pkg/clients"
	"time"
)

type PagamentoService struct {
	pagamentoClient *clients.PagamentoClient
	pedidoClient    *clients.PedidoClient
}

func NovoPagamentoService(pagamentoClient *clients.PagamentoClient, pedidoClient *clients.PedidoClient) *PagamentoService {
	return &PagamentoService{
		pagamentoClient: pagamentoClient,
		pedidoClient:    pedidoClient,
	}
}

func pagamentoFromClient(p *clients.Pagamento) *domain.Pagamento {
	if p == nil {
		return nil
	}

	var expiresAt *time.Time
	if p.ExpiresAt != nil {
		t, err := time.Parse(time.RFC3339, *p.ExpiresAt)
		if err == nil {
			expiresAt = &t
		}
	}

	var processedAt *time.Time
	if p.ProcessedAt != nil {
		t, err := time.Parse(time.RFC3339, *p.ProcessedAt)
		if err == nil {
			processedAt = &t
		}
	}

	return &domain.Pagamento{
		ID:          p.ID,
		PedidoID:    p.PedidoID,
		Valor:       p.Valor,
		Status:      domain.StatusPagamento(p.Status),
		Tipo:        domain.TipoPagamento(p.Tipo),
		QRCodeData:  p.QRCodeData,
		ExternalID:  p.ExternalID,
		ExpiresAt:   expiresAt,
		ProcessedAt: processedAt,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func (s *PagamentoService) CriarPagamento(ctx context.Context, pedidoID string, valor float64, tipo domain.TipoPagamento) (*domain.Pagamento, error) {
	pedido, err := s.pedidoClient.ValidarPedido(ctx, pedidoID)
	if err != nil {
		return nil, err
	}
	if pedido == nil {
		return nil, errors.New("pedido não encontrado")
	}

	pagamentoClient, err := s.pagamentoClient.CriarPagamento(ctx, pedidoID, valor, string(tipo))
	if err != nil {
		return nil, err
	}

	return pagamentoFromClient(pagamentoClient), nil
}

func (s *PagamentoService) BuscarPagamentoPorID(ctx context.Context, id string) (*domain.Pagamento, error) {
	pagamentoClient, err := s.pagamentoClient.BuscarPorID(ctx, id)
	if err != nil {
		return nil, err
	}
	return pagamentoFromClient(pagamentoClient), nil
}

func (s *PagamentoService) BuscarPagamentoPorPedidoID(ctx context.Context, pedidoID string) (*domain.Pagamento, error) {
	pagamentoClient, err := s.pagamentoClient.BuscarPorPedidoID(ctx, pedidoID)
	if err != nil {
		return nil, err
	}
	return pagamentoFromClient(pagamentoClient), nil
}

func (s *PagamentoService) ProcessarWebhookPagamento(ctx context.Context, externalID string, status domain.StatusPagamento) error {
	// O webhook deve ser processado diretamente pelo serviço de pagamentos
	// Este método pode ser removido ou mantido como proxy se necessário
	return fmt.Errorf("webhook deve ser processado diretamente pelo serviço de pagamentos")
}

func (s *PagamentoService) ConsultarStatusPagamento(ctx context.Context, pedidoID string) (*domain.Pagamento, error) {
	pagamentoClient, err := s.pagamentoClient.ConsultarStatusPagamento(ctx, pedidoID)
	if err != nil {
		return nil, err
	}

	pagamento := pagamentoFromClient(pagamentoClient)
	if pagamento == nil {
		return nil, errors.New("pagamento não encontrado para o pedido")
	}

	return pagamento, nil
}

func (s *PagamentoService) CancelarPagamento(ctx context.Context, pagamentoID string) error {
	return s.pagamentoClient.CancelarPagamento(ctx, pagamentoID)
}

func (s *PagamentoService) ListarPagamentosExpirados(ctx context.Context) ([]*domain.Pagamento, error) {
	// Esta funcionalidade deve ser implementada no serviço de pagamentos
	// Por enquanto retornamos erro ou lista vazia
	return nil, fmt.Errorf("funcionalidade não disponível via API")
}
