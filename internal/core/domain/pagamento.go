package domain

import (
	"errors"
	"time"
)

type StatusPagamento string

const (
	StatusPagamentoPendente  StatusPagamento = "PENDENTE"
	StatusPagamentoAprovado  StatusPagamento = "APROVADO"
	StatusPagamentoRecusado  StatusPagamento = "RECUSADO"
	StatusPagamentoCancelado StatusPagamento = "CANCELADO"
)

type TipoPagamento string

const (
	TipoPagamentoPix        TipoPagamento = "PIX"
	TipoPagamentoCartao     TipoPagamento = "CARTAO"
	TipoPagamentoDinheiro   TipoPagamento = "DINHEIRO"
)

type Pagamento struct {
	ID              string          `json:"id"`
	PedidoID        string          `json:"pedido_id"`
	Valor           float64         `json:"valor"`
	Status          StatusPagamento `json:"status"`
	Tipo            TipoPagamento   `json:"tipo"`
	QRCodeData      string          `json:"qr_code_data,omitempty"`
	ExternalID      string          `json:"external_id,omitempty"`
	PaymentURL      string          `json:"payment_url,omitempty"`
	ExpiresAt       *time.Time      `json:"expires_at,omitempty"`
	ProcessedAt     *time.Time      `json:"processed_at,omitempty"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

func NovoPagamento(id, pedidoID string, valor float64, tipo TipoPagamento) (*Pagamento, error) {
	pagamento := &Pagamento{
		ID:        id,
		PedidoID:  pedidoID,
		Valor:     valor,
		Status:    StatusPagamentoPendente,
		Tipo:      tipo,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := pagamento.Validar(); err != nil {
		return nil, err
	}

	return pagamento, nil
}

func (p *Pagamento) Validar() error {
	if p.PedidoID == "" {
		return errors.New("pedido_id não pode ser vazio")
	}

	if p.Valor <= 0 {
		return errors.New("valor deve ser maior que zero")
	}

	if !IsStatusPagamentoValido(p.Status) {
		return errors.New("status de pagamento inválido")
	}

	if !IsTipoPagamentoValido(p.Tipo) {
		return errors.New("tipo de pagamento inválido")
	}

	return nil
}

func (p *Pagamento) Aprovar() error {
	if p.Status != StatusPagamentoPendente {
		return errors.New("apenas pagamentos pendentes podem ser aprovados")
	}

	p.Status = StatusPagamentoAprovado
	now := time.Now()
	p.ProcessedAt = &now
	p.UpdatedAt = now

	return nil
}

func (p *Pagamento) Recusar() error {
	if p.Status != StatusPagamentoPendente {
		return errors.New("apenas pagamentos pendentes podem ser recusados")
	}

	p.Status = StatusPagamentoRecusado
	now := time.Now()
	p.ProcessedAt = &now
	p.UpdatedAt = now

	return nil
}

func (p *Pagamento) Cancelar() error {
	if p.Status == StatusPagamentoAprovado {
		return errors.New("pagamentos aprovados não podem ser cancelados")
	}

	p.Status = StatusPagamentoCancelado
	now := time.Now()
	p.ProcessedAt = &now
	p.UpdatedAt = now

	return nil
}

func (p *Pagamento) IsExpirado() bool {
	if p.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*p.ExpiresAt)
}

func IsStatusPagamentoValido(status StatusPagamento) bool {
	switch status {
	case StatusPagamentoPendente, StatusPagamentoAprovado, StatusPagamentoRecusado, StatusPagamentoCancelado:
		return true
	default:
		return false
	}
}

func IsTipoPagamentoValido(tipo TipoPagamento) bool {
	switch tipo {
	case TipoPagamentoPix, TipoPagamentoCartao, TipoPagamentoDinheiro:
		return true
	default:
		return false
	}
}
