package domain

import (
	"errors"
	"time"
)

type StatusPedido string

const (
	StatusRecebido     StatusPedido = "RECEBIDO"
	StatusEmPreparacao StatusPedido = "EM_PREPARACAO"
	StatusPronto       StatusPedido = "PRONTO"
	StatusFinalizado   StatusPedido = "FINALIZADO"
)

type ItemPedido struct {
	ProdutoID  string  `json:"produto_id"`
	Nome       string  `json:"nome"`
	Preco      float64 `json:"preco"`
	Quantidade int     `json:"quantidade"`
	Observacao string  `json:"observacao,omitempty"`
}

type Pedido struct {
	ID         string       `json:"id"`
	ClienteID  *string      `json:"cliente_id,omitempty"`
	Itens      []ItemPedido `json:"itens"`
	ValorTotal float64      `json:"valor_total"`
	Status     StatusPedido `json:"status"`
	CreatedAt  time.Time    `json:"created_at"`
	UpdatedAt  time.Time    `json:"updated_at"`
}

func NovoPedido(id string, clienteID *string, itens []ItemPedido) (*Pedido, error) {
	pedido := &Pedido{
		ID:        id,
		ClienteID: clienteID,
		Itens:     itens,
		Status:    StatusRecebido,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := pedido.Validar(); err != nil {
		return nil, err
	}

	pedido.CalcularValorTotal()

	return pedido, nil
}

func (p *Pedido) Validar() error {
	if len(p.Itens) == 0 {
		return errors.New("pedido deve ter pelo menos um item")
	}

	for _, item := range p.Itens {
		if item.ProdutoID == "" {
			return errors.New("produto ID não pode ser vazio")
		}
		if item.Quantidade <= 0 {
			return errors.New("quantidade deve ser maior que zero")
		}
		if item.Preco <= 0 {
			return errors.New("preço do item deve ser maior que zero")
		}
	}

	return nil
}

func (p *Pedido) CalcularValorTotal() {
	total := 0.0
	for _, item := range p.Itens {
		total += item.Preco * float64(item.Quantidade)
	}
	p.ValorTotal = total
}

func (p *Pedido) AtualizarStatus(novoStatus StatusPedido) {
	p.Status = novoStatus
	p.UpdatedAt = time.Now()
}

func IsStatusValido(status StatusPedido) bool {
	switch status {
	case StatusRecebido, StatusEmPreparacao, StatusPronto, StatusFinalizado:
		return true
	default:
		return false
	}
}
