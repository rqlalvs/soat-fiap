package presenters

import (
	"errors"
	"soat-fiap/internal/core/domain"
	"testing"
	"time"
)

func TestNewPedidoPresenter(t *testing.T) {
	presenter := NewPedidoPresenter()
	
	if presenter == nil {
		t.Error("NewPedidoPresenter() retornou nil")
	}
}

func TestPedidoPresenter_Present(t *testing.T) {
	presenter := NewPedidoPresenter()
	now := time.Now()
	clienteID := "cliente-123"
	
	pedido := &domain.Pedido{
		ID:         "pedido-123",
		ClienteID:  &clienteID,
		ValorTotal: 100.50,
		Status:     domain.StatusRecebido,
		Itens: []domain.ItemPedido{
			{
				ProdutoID:  "produto-1",
				Nome:       "Hamburguer",
				Preco:      50.00,
				Quantidade: 2,
				Observacao: "Sem cebola",
			},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
	
	result := presenter.Present(pedido)
	
	if result == nil {
		t.Fatal("Present() retornou nil")
	}
	
}

func TestPedidoPresenter_Present_NilPedido(t *testing.T) {
	presenter := NewPedidoPresenter()
	
	defer func() {
		if r := recover(); r != nil {
			t.Log("Present() com pedido nil causou panic como esperado")
		}
	}()
	
	result := presenter.Present(nil)
	
	if result != nil {
		t.Log("Present() com pedido nil retornou resultado não-nil")
	}
}

func TestPedidoPresenter_PresentList(t *testing.T) {
	presenter := NewPedidoPresenter()
	now := time.Now()
	clienteID1 := "cliente-1"
	clienteID2 := "cliente-2"
	
	pedidos := []*domain.Pedido{
		{
			ID:         "pedido-1",
			ClienteID:  &clienteID1,
			ValorTotal: 50.00,
			Status:     domain.StatusRecebido,
			Itens:      []domain.ItemPedido{},
			CreatedAt:  now,
			UpdatedAt:  now,
		},
		{
			ID:         "pedido-2",
			ClienteID:  &clienteID2,
			ValorTotal: 75.00,
			Status:     domain.StatusEmPreparacao,
			Itens:      []domain.ItemPedido{},
			CreatedAt:  now,
			UpdatedAt:  now,
		},
	}
	
	result := presenter.PresentList(pedidos)
	
	if result == nil {
		t.Fatal("PresentList() retornou nil")
	}
	
	if len(result) != len(pedidos) {
		t.Errorf("PresentList() retornou %d items, esperado %d", len(result), len(pedidos))
	}
}

func TestPedidoPresenter_PresentList_EmptyList(t *testing.T) {
	presenter := NewPedidoPresenter()
	
	pedidos := []*domain.Pedido{}
	
	result := presenter.PresentList(pedidos)
	
	if result == nil {
		t.Fatal("PresentList() retornou nil para lista vazia")
	}
	
	if len(result) != 0 {
		t.Errorf("PresentList() retornou %d items para lista vazia, esperado 0", len(result))
	}
}

func TestPedidoPresenter_PresentList_NilList(t *testing.T) {
	presenter := NewPedidoPresenter()
	
	defer func() {
		if r := recover(); r != nil {
			t.Log("PresentList() com lista nil causou panic")
		}
	}()
	
	result := presenter.PresentList(nil)
	
	if result == nil {
		t.Log("PresentList() retornou nil para lista nil")
	}
}

func TestPedidoPresenter_PresentError(t *testing.T) {
	presenter := NewPedidoPresenter()
	
	tests := []struct {
		name          string
		err           error
		expectedCode  int
		shouldContain string
	}{
		{
			name:          "Erro genérico",
			err:           errors.New("erro genérico"),
			expectedCode:  400,
			shouldContain: "erro genérico",
		},
		{
			name:          "Erro de validação",
			err:           errors.New("campo obrigatório não informado"),
			expectedCode:  400,
			shouldContain: "campo obrigatório",
		},
		{
			name:          "Erro de negócio",
			err:           errors.New("pedido não pode ser cancelado"),
			expectedCode:  400,
			shouldContain: "não pode ser cancelado",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, response := presenter.PresentError(tt.err)
			
			if code != tt.expectedCode {
				t.Errorf("PresentError() code = %d, esperado %d", code, tt.expectedCode)
			}
			
			if response == nil {
				t.Fatal("PresentError() response é nil")
			}
			
			respMap, ok := response.(map[string]string)
			if !ok {
				t.Fatal("PresentError() response não é map[string]string")
			}
			
			errorMsg, exists := respMap["error"]
			if !exists {
				t.Error("PresentError() response não contém campo 'error'")
			}
			
			if errorMsg != tt.err.Error() {
				t.Errorf("PresentError() error message = %s, esperado %s", errorMsg, tt.err.Error())
			}
		})
	}
}

func TestPedidoPresenter_PresentError_NilError(t *testing.T) {
	presenter := NewPedidoPresenter()
	
	defer func() {
		if r := recover(); r != nil {
			t.Log("PresentError() com erro nil causou panic")
		}
	}()
	
	code, response := presenter.PresentError(nil)
	
	if code != 400 {
		t.Errorf("PresentError() code = %d, esperado 400", code)
	}
	
	t.Logf("PresentError() com erro nil retornou: %v", response)
}

func TestPedidoPresenter_PresentSuccess(t *testing.T) {
	presenter := NewPedidoPresenter()
	
	tests := []struct {
		name    string
		message string
		data    interface{}
	}{
		{
			name:    "Sucesso com dados",
			message: "Operação realizada com sucesso",
			data:    map[string]string{"id": "123"},
		},
		{
			name:    "Sucesso sem dados",
			message: "Operação realizada",
			data:    nil,
		},
		{
			name:    "Sucesso com string",
			message: "Item criado",
			data:    "item-123",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, response := presenter.PresentSuccess(tt.message, tt.data)
			
			if code != 200 {
				t.Errorf("PresentSuccess() code = %d, esperado 200", code)
			}
			
			if response == nil {
				t.Fatal("PresentSuccess() response é nil")
			}
			
			respMap, ok := response.(map[string]interface{})
			if !ok {
				t.Fatal("PresentSuccess() response não é map[string]interface{}")
			}
			
			message, exists := respMap["message"]
			if !exists {
				t.Error("PresentSuccess() response não contém campo 'message'")
			}
			
			if message != tt.message {
				t.Errorf("PresentSuccess() message = %s, esperado %s", message, tt.message)
			}
			
			data, exists := respMap["data"]
			if !exists {
				t.Error("PresentSuccess() response não contém campo 'data'")
			}
			
			if tt.data == nil && data != nil {
				t.Error("PresentSuccess() data deveria ser nil")
			}
		})
	}
}

func TestPedidoPresenter_PresentCreated(t *testing.T) {
	presenter := NewPedidoPresenter()
	
	tests := []struct {
		name string
		data interface{}
	}{
		{
			name: "Criado com map",
			data: map[string]string{"id": "123", "status": "criado"},
		},
		{
			name: "Criado com string",
			data: "recurso-123",
		},
		{
			name: "Criado com nil",
			data: nil,
		},
		{
			name: "Criado com struct",
			data: struct {
				ID   string
				Nome string
			}{
				ID:   "123",
				Nome: "Teste",
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, response := presenter.PresentCreated(tt.data)
			
			if code != 201 {
				t.Errorf("PresentCreated() code = %d, esperado 201", code)
			}
			
			if response == nil {
				t.Fatal("PresentCreated() response é nil")
			}
			
			respMap, ok := response.(map[string]interface{})
			if !ok {
				t.Fatal("PresentCreated() response não é map[string]interface{}")
			}
			
			message, exists := respMap["message"]
			if !exists {
				t.Error("PresentCreated() response não contém campo 'message'")
			}
			
			expectedMessage := "Resource created successfully"
			if message != expectedMessage {
				t.Errorf("PresentCreated() message = %s, esperado %s", message, expectedMessage)
			}
			
			data, exists := respMap["data"]
			if !exists {
				t.Error("PresentCreated() response não contém campo 'data'")
			}
			
			if tt.data == nil && data != nil {
				t.Error("PresentCreated() data deveria ser nil")
			}
		})
	}
}

func TestPedidoPresenter_MultipleCalls(t *testing.T) {
	presenter := NewPedidoPresenter()
	
	err1 := errors.New("erro 1")
	err2 := errors.New("erro 2")
	
	code1, resp1 := presenter.PresentError(err1)
	code2, resp2 := presenter.PresentError(err2)
	
	if code1 != code2 {
		t.Error("Códigos diferentes para mesma operação")
	}
	
	respMap1 := resp1.(map[string]string)
	respMap2 := resp2.(map[string]string)
	
	if respMap1["error"] == respMap2["error"] {
		t.Error("Mensagens de erro não deveriam ser iguais")
	}
}

func BenchmarkPedidoPresenter_Present(b *testing.B) {
	presenter := NewPedidoPresenter()
	clienteID := "cliente-123"
	pedido := &domain.Pedido{
		ID:         "pedido-123",
		ClienteID:  &clienteID,
		ValorTotal: 100.00,
		Status:     domain.StatusRecebido,
		Itens:      []domain.ItemPedido{},
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		presenter.Present(pedido)
	}
}

func BenchmarkPedidoPresenter_PresentError(b *testing.B) {
	presenter := NewPedidoPresenter()
	err := errors.New("erro de teste")
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		presenter.PresentError(err)
	}
}

func BenchmarkPedidoPresenter_PresentSuccess(b *testing.B) {
	presenter := NewPedidoPresenter()
	data := map[string]string{"id": "123"}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		presenter.PresentSuccess("sucesso", data)
	}
}
