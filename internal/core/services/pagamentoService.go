package services

import (
	"context"
	"errors"
	"fmt"
	"soat-fiap/internal/core/domain"
	"soat-fiap/internal/core/ports"
	"time"

	"github.com/google/uuid"
)

type PagamentoService struct {
	pagamentoRepository ports.PagamentoRepository
	pedidoService       ports.PedidoService
}

func NovoPagamentoService(pagamentoRepository ports.PagamentoRepository, pedidoService ports.PedidoService) *PagamentoService {
	return &PagamentoService{
		pagamentoRepository: pagamentoRepository,
		pedidoService:       pedidoService,
	}
}

func (s *PagamentoService) CriarPagamento(ctx context.Context, pedidoID string, valor float64, tipo domain.TipoPagamento) (*domain.Pagamento, error) {

	pedido, err := s.pedidoService.BuscarPedidoPorID(ctx, pedidoID)
	if err != nil {
		return nil, err
	}
	if pedido == nil {
		return nil, errors.New("pedido não encontrado")
	}

	pagamentoExistente, err := s.pagamentoRepository.BuscarPorPedidoID(ctx, pedidoID)
	if err == nil && pagamentoExistente != nil {
		return nil, errors.New("já existe um pagamento para este pedido")
	}

	id := uuid.New().String()
	pagamento, err := domain.NovoPagamento(id, pedidoID, valor, tipo)
	if err != nil {
		return nil, err
	}

	if tipo == domain.TipoPagamentoPix {
		qrCode, err := s.GerarQRCodePIX(ctx, pagamento.ID)
		if err != nil {
			return nil, fmt.Errorf("erro ao gerar QR Code PIX: %w", err)
		}
		pagamento.QRCodeData = qrCode
		
		expiresAt := time.Now().Add(30 * time.Minute)
		pagamento.ExpiresAt = &expiresAt
	}

	pagamento.ExternalID = fmt.Sprintf("MP_%s", uuid.New().String()[:8])

	err = s.pagamentoRepository.Criar(ctx, pagamento)
	if err != nil {
		return nil, err
	}

	return pagamento, nil
}

func (s *PagamentoService) BuscarPagamentoPorID(ctx context.Context, id string) (*domain.Pagamento, error) {
	return s.pagamentoRepository.BuscarPorID(ctx, id)
}

func (s *PagamentoService) BuscarPagamentoPorPedidoID(ctx context.Context, pedidoID string) (*domain.Pagamento, error) {
	return s.pagamentoRepository.BuscarPorPedidoID(ctx, pedidoID)
}

func (s *PagamentoService) ProcessarWebhookPagamento(ctx context.Context, externalID string, status domain.StatusPagamento) error {

	pagamento, err := s.pagamentoRepository.BuscarPorExternalID(ctx, externalID)
	if err != nil {
		return err
	}
	if pagamento == nil {
		return errors.New("pagamento não encontrado para o ID externo: " + externalID)
	}

	switch status {
	case domain.StatusPagamentoAprovado:
		err = pagamento.Aprovar()
	case domain.StatusPagamentoRecusado:
		err = pagamento.Recusar()
	case domain.StatusPagamentoCancelado:
		err = pagamento.Cancelar()
	default:
		return errors.New("status de pagamento inválido no webhook")
	}

	if err != nil {
		return err
	}

	err = s.pagamentoRepository.Atualizar(ctx, pagamento)
	if err != nil {
		return err
	}

	if status == domain.StatusPagamentoAprovado {
		err = s.pedidoService.AtualizarStatusPedido(ctx, pagamento.PedidoID, domain.StatusEmPreparacao)
		if err != nil {
			fmt.Printf("Erro ao atualizar status do pedido %s: %v\n", pagamento.PedidoID, err)
		}
	}

	return nil
}

func (s *PagamentoService) GerarQRCodePIX(ctx context.Context, pagamentoID string) (string, error) {
	qrCodeData := fmt.Sprintf("00020126580014br.gov.bcb.pix0136%s5204000053039865802BR5925LANCHONETE SOAT FIAP6009SAO PAULO62070503***6304", pagamentoID)
	return qrCodeData, nil
}

func (s *PagamentoService) ConsultarStatusPagamento(ctx context.Context, pedidoID string) (*domain.Pagamento, error) {
	pagamento, err := s.pagamentoRepository.BuscarPorPedidoID(ctx, pedidoID)
	if err != nil {
		return nil, err
	}
	if pagamento == nil {
		return nil, errors.New("pagamento não encontrado para o pedido")
	}

	if pagamento.Tipo == domain.TipoPagamentoPix && pagamento.IsExpirado() && pagamento.Status == domain.StatusPagamentoPendente {
		err = pagamento.Cancelar()
		if err != nil {
			return nil, err
		}
		err = s.pagamentoRepository.Atualizar(ctx, pagamento)
		if err != nil {
			return nil, err
		}
	}

	return pagamento, nil
}

func (s *PagamentoService) CancelarPagamento(ctx context.Context, pagamentoID string) error {
	pagamento, err := s.pagamentoRepository.BuscarPorID(ctx, pagamentoID)
	if err != nil {
		return err
	}
	if pagamento == nil {
		return errors.New("pagamento não encontrado")
	}

	err = pagamento.Cancelar()
	if err != nil {
		return err
	}

	return s.pagamentoRepository.Atualizar(ctx, pagamento)
}

func (s *PagamentoService) ListarPagamentosExpirados(ctx context.Context) ([]*domain.Pagamento, error) {
	return s.pagamentoRepository.ListarExpirados(ctx)
}
