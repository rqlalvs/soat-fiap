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

type ItemPedido struct {
	ProdutoID string  `json:"produto_id"`
	Nome      string  `json:"nome"`
	Preco     float64 `json:"preco"`
	Quantidade int    `json:"quantidade"`
}

type Pedido struct {
	ID         string       `json:"id"`
	ClienteID  *string      `json:"cliente_id,omitempty"`
	Itens      []ItemPedido `json:"itens"`
	ValorTotal float64     `json:"valor_total"`
	Status     string       `json:"status"`
	CreatedAt  string       `json:"created_at,omitempty"`
	UpdatedAt  string       `json:"updated_at,omitempty"`
}

type CriarPedidoRequest struct {
	ClienteID *string      `json:"cliente_id,omitempty"`
	Itens     []ItemPedido `json:"itens"`
}

type AtualizarStatusRequest struct {
	Status string `json:"status"`
}

type PedidoClient struct {
	baseURL    string
	httpClient *http.Client
}

func NovoPedidoClient(baseURL string) *PedidoClient {
	return &PedidoClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *PedidoClient) CriarPedido(ctx context.Context, clienteID *string, itens []ItemPedido) (*Pedido, error) {
	url := fmt.Sprintf("%s/api/v1/pedidos", c.baseURL)

	reqBody := CriarPedidoRequest{
		ClienteID: clienteID,
		Itens:     itens,
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
		return nil, fmt.Errorf("erro ao criar pedido: status %d, body: %s", resp.StatusCode, string(body))
	}

	var pedido Pedido
	if err := json.NewDecoder(resp.Body).Decode(&pedido); err != nil {
		return nil, err
	}

	return &pedido, nil
}

func (c *PedidoClient) BuscarPorID(ctx context.Context, id string) (*Pedido, error) {
	url := fmt.Sprintf("%s/api/v1/pedidos/%s", c.baseURL, id)

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
		return nil, fmt.Errorf("erro ao buscar pedido: %s", string(body))
	}

	var pedido Pedido
	if err := json.NewDecoder(resp.Body).Decode(&pedido); err != nil {
		return nil, err
	}

	return &pedido, nil
}

func (c *PedidoClient) ListarPedidos(ctx context.Context) ([]*Pedido, error) {
	url := fmt.Sprintf("%s/api/v1/pedidos", c.baseURL)

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
		return nil, fmt.Errorf("erro ao listar pedidos: %s", string(body))
	}

	var pedidos []*Pedido
	if err := json.NewDecoder(resp.Body).Decode(&pedidos); err != nil {
		return nil, err
	}

	return pedidos, nil
}

func (c *PedidoClient) ListarPedidosPorStatus(ctx context.Context, status string) ([]*Pedido, error) {
	url := fmt.Sprintf("%s/api/v1/pedidos?status=%s", c.baseURL, status)

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
		return nil, fmt.Errorf("erro ao listar pedidos por status: %s", string(body))
	}

	var pedidos []*Pedido
	if err := json.NewDecoder(resp.Body).Decode(&pedidos); err != nil {
		return nil, err
	}

	return pedidos, nil
}

func (c *PedidoClient) ListarPedidosPorCliente(ctx context.Context, clienteID string) ([]*Pedido, error) {
	url := fmt.Sprintf("%s/api/v1/pedidos?cliente_id=%s", c.baseURL, clienteID)

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
		return nil, fmt.Errorf("erro ao listar pedidos por cliente: %s", string(body))
	}

	var pedidos []*Pedido
	if err := json.NewDecoder(resp.Body).Decode(&pedidos); err != nil {
		return nil, err
	}

	return pedidos, nil
}

func (c *PedidoClient) AtualizarStatusPedido(ctx context.Context, id string, status string) error {
	url := fmt.Sprintf("%s/api/v1/pedidos/%s/status", c.baseURL, id)

	reqBody := AtualizarStatusRequest{
		Status: status,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("erro ao serializar requisição: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, url, bytes.NewBuffer(jsonData))
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
		return fmt.Errorf("erro ao atualizar status do pedido: status %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}

func (c *PedidoClient) ValidarPedido(ctx context.Context, pedidoID string) (*Pedido, error) {
	pedido, err := c.BuscarPorID(ctx, pedidoID)
	if err != nil {
		return nil, err
	}

	if pedido == nil {
		return nil, errors.New("pedido não encontrado")
	}

	return pedido, nil
}
