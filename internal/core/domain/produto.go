package domain

import (
	"errors"
	"time"
)

type Categoria string

const (
	CategoriaLanche         Categoria = "LANCHE"
	CategoriaAcompanhamento Categoria = "ACOMPANHAMENTO"
	CategoriaBebida         Categoria = "BEBIDA"
	CategoriaSobremesa      Categoria = "SOBREMESA"
)

type Produto struct {
	ID         string    `json:"id"`
	Nome       string    `json:"nome"`
	Descricao  string    `json:"descricao"`
	Preco      float64   `json:"preco"`
	Categoria  Categoria `json:"categoria"`
	Disponivel bool      `json:"disponivel"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func NovoProduto(id, nome, descricao string, preco float64, categoria Categoria) (*Produto, error) {
	produto := &Produto{
		ID:         id,
		Nome:       nome,
		Descricao:  descricao,
		Preco:      preco,
		Categoria:  categoria,
		Disponivel: true,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := produto.Validar(); err != nil {
		return nil, err
	}

	return produto, nil
}

func (p *Produto) Validar() error {
	if p.Nome == "" {
		return errors.New("nome não pode ser vazio")
	}

	if p.Descricao == "" {
		return errors.New("descrição não pode ser vazia")
	}

	if p.Preco <= 0 {
		return errors.New("preço deve ser maior que zero")
	}

	if !IsCategoriaValida(p.Categoria) {
		return errors.New("categoria inválida")
	}

	return nil
}

func IsCategoriaValida(categoria Categoria) bool {
	switch categoria {
	case CategoriaLanche, CategoriaAcompanhamento, CategoriaBebida, CategoriaSobremesa:
		return true
	default:
		return false
	}
}
