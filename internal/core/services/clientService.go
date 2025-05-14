package services

import (
	"context"
	"errors"
	"soat-fiap/internal/core/domain"
	"soat-fiap/internal/core/ports"
	"time"

	"github.com/google/uuid"
)

type ClienteService struct {
	repository ports.ClienteRepository
}

func NovoClienteService(repository ports.ClienteRepository) *ClienteService {
	return &ClienteService{
		repository: repository,
	}
}

func (s *ClienteService) CriarCliente(ctx context.Context, nome, cpf, email, telefone string) (*domain.Cliente, error) {

	clienteExistente, err := s.repository.BuscarPorCPF(ctx, cpf)
	if err == nil && clienteExistente != nil {
		return nil, errors.New("já existe um cliente com este CPF")
	}

	id := uuid.New().String()

	cliente, err := domain.NovoCliente(id, nome, cpf, email, telefone)
	if err != nil {
		return nil, err
	}

	err = s.repository.Criar(ctx, cliente)
	if err != nil {
		return nil, err
	}

	return cliente, nil
}

func (s *ClienteService) BuscarClientePorID(ctx context.Context, id string) (*domain.Cliente, error) {
	return s.repository.BuscarPorID(ctx, id)
}

func (s *ClienteService) BuscarClientePorCPF(ctx context.Context, cpf string) (*domain.Cliente, error) {
	if !domain.ValidarCPF(cpf) {
		return nil, errors.New("CPF inválido")
	}
	return s.repository.BuscarPorCPF(ctx, cpf)
}

func (s *ClienteService) ListarClientes(ctx context.Context) ([]*domain.Cliente, error) {
	return s.repository.Listar(ctx)
}

func (s *ClienteService) AtualizarCliente(ctx context.Context, cliente *domain.Cliente) error {
	clienteExistente, err := s.repository.BuscarPorID(ctx, cliente.ID)
	if err != nil {
		return err
	}
	if clienteExistente == nil {
		return errors.New("cliente não encontrado")
	}

	if clienteExistente.CPF != cliente.CPF {
		clienteComMesmoCPF, err := s.repository.BuscarPorCPF(ctx, cliente.CPF)
		if err == nil && clienteComMesmoCPF != nil && clienteComMesmoCPF.ID != cliente.ID {
			return errors.New("já existe outro cliente com este CPF")
		}
	}

	err = cliente.Validar()
	if err != nil {
		return err
	}

	cliente.UpdatedAt = time.Now()

	return s.repository.Atualizar(ctx, cliente)
}

func (s *ClienteService) DeletarCliente(ctx context.Context, id string) error {
	return s.repository.Deletar(ctx, id)
}
