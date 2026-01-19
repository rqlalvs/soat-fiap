package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Cliente struct {
	ID       string `json:"id"`
	Nome     string `json:"nome"`
	CPF      string `json:"cpf"`
	Email    string `json:"email"`
	Telefone string `json:"telefone"`
}

type CriarClienteRequest struct {
	Nome     string `json:"nome"`
	CPF      string `json:"cpf"`
	Email    string `json:"email"`
	Telefone string `json:"telefone"`
}

type AtualizarClienteRequest struct {
	Nome     string `json:"nome"`
	CPF      string `json:"cpf"`
	Email    string `json:"email"`
	Telefone string `json:"telefone"`
}

type ClienteClient struct {
	baseURL    string
	httpClient *http.Client
}

func NovoClienteClient(baseURL string) *ClienteClient {
	return &ClienteClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *ClienteClient) CriarCliente(ctx context.Context, nome, cpf, email, telefone string) (*Cliente, error) {
	url := fmt.Sprintf("%s/api/v1/clientes", c.baseURL)

	reqBody := CriarClienteRequest{
		Nome:     nome,
		CPF:      cpf,
		Email:    email,
		Telefone: telefone,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("erro ao serializar requisição: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("erro ao criar cliente: status %d, body: %s", resp.StatusCode, string(body))
	}

	var cliente Cliente
	if err := json.NewDecoder(resp.Body).Decode(&cliente); err != nil {
		return nil, err
	}

	return &cliente, nil
}

func (c *ClienteClient) BuscarPorID(ctx context.Context, id string) (*Cliente, error) {
	url := fmt.Sprintf("%s/api/v1/clientes/%s", c.baseURL, id)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("erro ao buscar cliente: %s", string(body))
	}

	var cliente Cliente
	if err := json.NewDecoder(resp.Body).Decode(&cliente); err != nil {
		return nil, err
	}

	return &cliente, nil
}

func (c *ClienteClient) BuscarPorCPF(ctx context.Context, cpf string) (*Cliente, error) {
	url := fmt.Sprintf("%s/api/v1/clientes/cpf/%s", c.baseURL, cpf)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("erro ao buscar cliente por CPF: %s", string(body))
	}

	var cliente Cliente
	if err := json.NewDecoder(resp.Body).Decode(&cliente); err != nil {
		return nil, err
	}

	return &cliente, nil
}

func (c *ClienteClient) ListarClientes(ctx context.Context) ([]*Cliente, error) {
	url := fmt.Sprintf("%s/api/v1/clientes", c.baseURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("erro ao listar clientes: %s", string(body))
	}

	var clientes []*Cliente
	if err := json.NewDecoder(resp.Body).Decode(&clientes); err != nil {
		return nil, err
	}

	return clientes, nil
}

func (c *ClienteClient) AtualizarCliente(ctx context.Context, id string, nome, cpf, email, telefone string) error {
	url := fmt.Sprintf("%s/api/v1/clientes/%s", c.baseURL, id)

	reqBody := AtualizarClienteRequest{
		Nome:     nome,
		CPF:      cpf,
		Email:    email,
		Telefone: telefone,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("erro ao serializar requisição: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("erro ao atualizar cliente: status %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}

func (c *ClienteClient) DeletarCliente(ctx context.Context, id string) error {
	url := fmt.Sprintf("%s/api/v1/clientes/%s", c.baseURL, id)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("erro ao deletar cliente: status %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}

func (c *ClienteClient) ValidarCliente(ctx context.Context, clienteID string) error {
	cliente, err := c.BuscarPorID(ctx, clienteID)
	if err != nil {
		return err
	}

	if cliente == nil {
		return errors.New("cliente não encontrado")
	}

	return nil
}
