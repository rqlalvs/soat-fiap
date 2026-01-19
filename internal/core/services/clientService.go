package services

import (
	"context"
	"soat-fiap/internal/core/domain"
	"soat-fiap/pkg/clients"
	"time"
)

type ClienteService struct {
	clienteClient *clients.ClienteClient
}

func NovoClienteService(clienteClient *clients.ClienteClient) *ClienteService {
	return &ClienteService{
		clienteClient: clienteClient,
	}
}

func clienteFromClient(c *clients.Cliente) *domain.Cliente {
	if c == nil {
		return nil
	}
	return &domain.Cliente{
		ID:        c.ID,
		Nome:      c.Nome,
		CPF:       c.CPF,
		Email:     c.Email,
		Telefone:  c.Telefone,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func clienteToClient(c *domain.Cliente) *clients.Cliente {
	if c == nil {
		return nil
	}
	return &clients.Cliente{
		ID:       c.ID,
		Nome:     c.Nome,
		CPF:      c.CPF,
		Email:    c.Email,
		Telefone: c.Telefone,
	}
}

func (s *ClienteService) CriarCliente(ctx context.Context, nome, cpf, email, telefone string) (*domain.Cliente, error) {
	clienteClient, err := s.clienteClient.CriarCliente(ctx, nome, cpf, email, telefone)
	if err != nil {
		return nil, err
	}

	return clienteFromClient(clienteClient), nil
}

func (s *ClienteService) BuscarClientePorID(ctx context.Context, id string) (*domain.Cliente, error) {
	clienteClient, err := s.clienteClient.BuscarPorID(ctx, id)
	if err != nil {
		return nil, err
	}
	return clienteFromClient(clienteClient), nil
}

func (s *ClienteService) BuscarClientePorCPF(ctx context.Context, cpf string) (*domain.Cliente, error) {
	clienteClient, err := s.clienteClient.BuscarPorCPF(ctx, cpf)
	if err != nil {
		return nil, err
	}
	return clienteFromClient(clienteClient), nil
}

func (s *ClienteService) ListarClientes(ctx context.Context) ([]*domain.Cliente, error) {
	clientesClient, err := s.clienteClient.ListarClientes(ctx)
	if err != nil {
		return nil, err
	}

	clientes := make([]*domain.Cliente, 0, len(clientesClient))
	for _, c := range clientesClient {
		clientes = append(clientes, clienteFromClient(c))
	}

	return clientes, nil
}

func (s *ClienteService) AtualizarCliente(ctx context.Context, cliente *domain.Cliente) error {
	err := cliente.Validar()
	if err != nil {
		return err
	}

	return s.clienteClient.AtualizarCliente(ctx, cliente.ID, cliente.Nome, cliente.CPF, cliente.Email, cliente.Telefone)
}

func (s *ClienteService) DeletarCliente(ctx context.Context, id string) error {
	return s.clienteClient.DeletarCliente(ctx, id)
}
