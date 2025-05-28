package domain

import (
	"errors"
	"time"
	"soat-fiap/internal/core/common"
)

type Cliente struct {
	ID        string    `json:"id"`
	Nome      string    `json:"nome"`
	CPF       string    `json:"cpf"`
	Email     string    `json:"email"`
	Telefone  string    `json:"telefone"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NovoCliente(id, nome, cpf, email, telefone string) (*Cliente, error) {
	cliente := &Cliente{
		ID:        id,
		Nome:      nome,
		CPF:       cpf,
		Email:     email,
		Telefone:  telefone,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := cliente.Validar()
	if err != nil {
		return nil, err
	}

	return cliente, nil
}

func (c *Cliente) Validar() error {
	if c.Nome == "" {
		return errors.New("nome não pode ser vazio")
	}

	if !common.ValidarCPF(c.CPF) {
		return errors.New("CPF inválido")
	}

	if !common.ValidarEmail(c.Email) {
		return errors.New("email inválido")
	}

	if c.Telefone == "" {
		return errors.New("telefone não pode ser vazio")
	}

	return nil
}


