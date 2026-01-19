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

type Produto struct {
	ID         string  `json:"id"`
	Nome       string  `json:"nome"`
	Descricao  string  `json:"descricao"`
	Preco      float64 `json:"preco"`
	Categoria  string  `json:"categoria"`
	Disponivel bool    `json:"disponivel"`
}

type CriarProdutoRequest struct {
	Nome       string  `json:"nome"`
	Descricao  string  `json:"descricao"`
	Preco      float64 `json:"preco"`
	Categoria  string  `json:"categoria"`
	Disponivel bool    `json:"disponivel"`
}

type AtualizarProdutoRequest struct {
	Nome       string  `json:"nome"`
	Descricao  string  `json:"descricao"`
	Preco      float64 `json:"preco"`
	Categoria  string  `json:"categoria"`
	Disponivel bool    `json:"disponivel"`
}

type ProdutoClient struct {
	baseURL    string
	httpClient *http.Client
}

func NovoProdutoClient(baseURL string) *ProdutoClient {
	return &ProdutoClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *ProdutoClient) CriarProduto(ctx context.Context, nome, descricao string, preco float64, categoria string, disponivel bool) (*Produto, error) {
	url := fmt.Sprintf("%s/api/v1/produtos", c.baseURL)

	reqBody := CriarProdutoRequest{
		Nome:       nome,
		Descricao:  descricao,
		Preco:      preco,
		Categoria:  categoria,
		Disponivel: disponivel,
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
		return nil, fmt.Errorf("erro ao criar produto: status %d, body: %s", resp.StatusCode, string(body))
	}

	var produto Produto
	if err := json.NewDecoder(resp.Body).Decode(&produto); err != nil {
		return nil, err
	}

	return &produto, nil
}

func (c *ProdutoClient) BuscarPorID(ctx context.Context, id string) (*Produto, error) {
	url := fmt.Sprintf("%s/api/v1/produtos/%s", c.baseURL, id)

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
		return nil, fmt.Errorf("erro ao buscar produto: %s", string(body))
	}

	var produto Produto
	if err := json.NewDecoder(resp.Body).Decode(&produto); err != nil {
		return nil, err
	}

	return &produto, nil
}

func (c *ProdutoClient) ListarProdutos(ctx context.Context) ([]*Produto, error) {
	url := fmt.Sprintf("%s/api/v1/produtos", c.baseURL)

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
		return nil, fmt.Errorf("erro ao listar produtos: %s", string(body))
	}

	var produtos []*Produto
	if err := json.NewDecoder(resp.Body).Decode(&produtos); err != nil {
		return nil, err
	}

	return produtos, nil
}

func (c *ProdutoClient) AtualizarProduto(ctx context.Context, id string, nome, descricao string, preco float64, categoria string, disponivel bool) error {
	url := fmt.Sprintf("%s/api/v1/produtos/%s", c.baseURL, id)

	reqBody := AtualizarProdutoRequest{
		Nome:       nome,
		Descricao:  descricao,
		Preco:      preco,
		Categoria:  categoria,
		Disponivel: disponivel,
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
		return fmt.Errorf("erro ao atualizar produto: status %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}

func (c *ProdutoClient) DeletarProduto(ctx context.Context, id string) error {
	url := fmt.Sprintf("%s/api/v1/produtos/%s", c.baseURL, id)

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
		return fmt.Errorf("erro ao deletar produto: status %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}

func (c *ProdutoClient) ValidarProduto(ctx context.Context, produtoID string) (*Produto, error) {
	produto, err := c.BuscarPorID(ctx, produtoID)
	if err != nil {
		return nil, err
	}

	if produto == nil {
		return nil, errors.New("produto não encontrado")
	}

	if !produto.Disponivel {
		return nil, errors.New("produto não disponível")
	}

	return produto, nil
}
