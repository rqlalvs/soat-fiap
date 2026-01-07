package repositories

import (
	"context"
	"database/sql"
	"errors"
	"soat-fiap/internal/core/domain"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestNovoPedidoRepository(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Erro ao criar mock: %v", err)
	}
	defer db.Close()

	repo := NovoPedidoRepository(db)

	if repo == nil {
		t.Fatal("NovoPedidoRepository() retornou nil")
	}

	if repo.db == nil {
		t.Error("db não foi inicializado")
	}
}

func TestPedidoRepository_Criar(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Erro ao criar mock: %v", err)
	}
	defer db.Close()

	repo := NovoPedidoRepository(db)
	now := time.Now()
	clienteID := "cliente-123"

	tests := []struct {
		name      string
		pedido    *domain.Pedido
		mockSetup func()
		wantErr   bool
	}{
		{
			name: "Sucesso ao criar pedido com itens",
			pedido: &domain.Pedido{
				ID:         "pedido-123",
				ClienteID:  &clienteID,
				ValorTotal: 50.00,
				Status:     domain.StatusRecebido,
				Itens: []domain.ItemPedido{
					{
						ProdutoID:  "prod-1",
						Nome:       "Hamburguer",
						Preco:      25.00,
						Quantidade: 2,
						Observacao: "Sem cebola",
					},
				},
				CreatedAt: now,
				UpdatedAt: now,
			},
			mockSetup: func() {
				mock.ExpectBegin()
				mock.ExpectPrepare("INSERT INTO pedidos").
					ExpectExec().
					WithArgs("pedido-123", "cliente-123", 50.00, domain.StatusRecebido,
						sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectPrepare("INSERT INTO pedido_itens").
					ExpectExec().
					WithArgs("pedido-123", "prod-1", "Hamburguer", 25.00, 2, "Sem cebola").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Erro ao iniciar transação",
			pedido: &domain.Pedido{
				ID:        "pedido-123",
				CreatedAt: now,
				UpdatedAt: now,
			},
			mockSetup: func() {
				mock.ExpectBegin().WillReturnError(errors.New("erro ao iniciar transação"))
			},
			wantErr: true,
		},
		{
			name: "Erro ao inserir pedido",
			pedido: &domain.Pedido{
				ID:        "pedido-123",
				CreatedAt: now,
				UpdatedAt: now,
			},
			mockSetup: func() {
				mock.ExpectBegin()
				mock.ExpectPrepare("INSERT INTO pedidos").
					ExpectExec().
					WillReturnError(errors.New("erro ao inserir"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "Erro ao inserir item",
			pedido: &domain.Pedido{
				ID:     "pedido-123",
				Status: domain.StatusRecebido,
				Itens: []domain.ItemPedido{
					{ProdutoID: "prod-1", Nome: "Item", Preco: 10.0, Quantidade: 1},
				},
				CreatedAt: now,
				UpdatedAt: now,
			},
			mockSetup: func() {
				mock.ExpectBegin()
				mock.ExpectPrepare("INSERT INTO pedidos").
					ExpectExec().
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectPrepare("INSERT INTO pedido_itens").
					ExpectExec().
					WillReturnError(errors.New("erro ao inserir item"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			err := repo.Criar(context.Background(), tt.pedido)

			if (err != nil) != tt.wantErr {
				t.Errorf("Criar() erro = %v, wantErr %v", err, tt.wantErr)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Expectativas não atendidas: %v", err)
			}
		})
	}
}

func TestPedidoRepository_BuscarPorID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Erro ao criar mock: %v", err)
	}
	defer db.Close()

	repo := NovoPedidoRepository(db)

	tests := []struct {
		name      string
		id        string
		mockSetup func()
		wantNil   bool
		wantErr   bool
	}{
		{
			name: "Sucesso ao buscar pedido com itens",
			id:   "pedido-123",
			mockSetup: func() {
				pedidoRows := sqlmock.NewRows([]string{"id", "cliente_id", "valor_total", "status", "created_at", "updated_at"}).
					AddRow("pedido-123", "cliente-123", 50.00, "RECEBIDO",
						"2024-01-01 10:00:00", "2024-01-01 10:00:00")

				mock.ExpectPrepare("SELECT (.+) FROM pedidos WHERE id").
					ExpectQuery().
					WithArgs("pedido-123").
					WillReturnRows(pedidoRows)

				itensRows := sqlmock.NewRows([]string{"produto_id", "nome", "preco", "quantidade", "observacao"}).
					AddRow("prod-1", "Hamburguer", 25.00, 2, "Sem cebola")

				mock.ExpectQuery("SELECT (.+) FROM pedido_itens WHERE pedido_id").
					WithArgs("pedido-123").
					WillReturnRows(itensRows)
			},
			wantNil: false,
			wantErr: false,
		},
		{
			name: "Pedido não encontrado",
			id:   "pedido-999",
			mockSetup: func() {
				mock.ExpectPrepare("SELECT (.+) FROM pedidos WHERE id").
					ExpectQuery().
					WithArgs("pedido-999").
					WillReturnError(sql.ErrNoRows)
			},
			wantNil: true,
			wantErr: false,
		},
		{
			name: "Erro ao buscar pedido",
			id:   "pedido-123",
			mockSetup: func() {
				mock.ExpectPrepare("SELECT (.+) FROM pedidos WHERE id").
					ExpectQuery().
					WillReturnError(errors.New("erro no banco"))
			},
			wantNil: false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			pedido, err := repo.BuscarPorID(context.Background(), tt.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("BuscarPorID() erro = %v, wantErr %v", err, tt.wantErr)
			}

			if (pedido == nil) != tt.wantNil {
				t.Errorf("BuscarPorID() pedido = %v, wantNil %v", pedido, tt.wantNil)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Expectativas não atendidas: %v", err)
			}
		})
	}
}

func TestPedidoRepository_Listar(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Erro ao criar mock: %v", err)
	}
	defer db.Close()

	repo := NovoPedidoRepository(db)

	tests := []struct {
		name      string
		mockSetup func()
		wantLen   int
		wantErr   bool
	}{
		{
			name: "Sucesso ao listar pedidos",
			mockSetup: func() {
				pedidosRows := sqlmock.NewRows([]string{"id", "cliente_id", "valor_total", "status", "created_at", "updated_at"}).
					AddRow("pedido-1", "cliente-1", 50.00, "RECEBIDO",
						"2024-01-01 10:00:00", "2024-01-01 10:00:00").
					AddRow("pedido-2", nil, 30.00, "PRONTO",
						"2024-01-01 10:00:00", "2024-01-01 10:00:00")

				mock.ExpectQuery("SELECT (.+) FROM pedidos ORDER BY created_at DESC").
					WillReturnRows(pedidosRows)

				itensRows := sqlmock.NewRows([]string{"produto_id", "nome", "preco", "quantidade", "observacao"}).
					AddRow("prod-1", "Item 1", 25.00, 2, nil)

				mock.ExpectQuery("SELECT (.+) FROM pedido_itens WHERE pedido_id").
					WithArgs("pedido-1").
					WillReturnRows(itensRows)

				itensRows2 := sqlmock.NewRows([]string{"produto_id", "nome", "preco", "quantidade", "observacao"}).
					AddRow("prod-2", "Item 2", 30.00, 1, nil)

				mock.ExpectQuery("SELECT (.+) FROM pedido_itens WHERE pedido_id").
					WithArgs("pedido-2").
					WillReturnRows(itensRows2)
			},
			wantLen: 2,
			wantErr: false,
		},
		{
			name: "Lista vazia",
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"id", "cliente_id", "valor_total", "status", "created_at", "updated_at"})

				mock.ExpectQuery("SELECT (.+) FROM pedidos ORDER BY created_at DESC").
					WillReturnRows(rows)
			},
			wantLen: 0,
			wantErr: false,
		},
		{
			name: "Erro ao listar",
			mockSetup: func() {
				mock.ExpectQuery("SELECT (.+) FROM pedidos ORDER BY created_at DESC").
					WillReturnError(errors.New("erro no banco"))
			},
			wantLen: 0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			pedidos, err := repo.Listar(context.Background())

			if (err != nil) != tt.wantErr {
				t.Errorf("Listar() erro = %v, wantErr %v", err, tt.wantErr)
			}

			if len(pedidos) != tt.wantLen {
				t.Errorf("Listar() len = %d, wantLen %d", len(pedidos), tt.wantLen)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Expectativas não atendidas: %v", err)
			}
		})
	}
}

func TestPedidoRepository_ListarPorStatus(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Erro ao criar mock: %v", err)
	}
	defer db.Close()

	repo := NovoPedidoRepository(db)

	tests := []struct {
		name      string
		status    domain.StatusPedido
		mockSetup func()
		wantLen   int
		wantErr   bool
	}{
		{
			name:   "Sucesso ao listar por status",
			status: domain.StatusRecebido,
			mockSetup: func() {
				pedidosRows := sqlmock.NewRows([]string{"id", "cliente_id", "valor_total", "status", "created_at", "updated_at"}).
					AddRow("pedido-1", "cliente-1", 50.00, "RECEBIDO",
						"2024-01-01 10:00:00", "2024-01-01 10:00:00")

				mock.ExpectQuery("SELECT (.+) FROM pedidos WHERE status").
					WithArgs(domain.StatusRecebido).
					WillReturnRows(pedidosRows)

				itensRows := sqlmock.NewRows([]string{"produto_id", "nome", "preco", "quantidade", "observacao"}).
					AddRow("prod-1", "Item 1", 25.00, 2, nil)

				mock.ExpectQuery("SELECT (.+) FROM pedido_itens WHERE pedido_id").
					WithArgs("pedido-1").
					WillReturnRows(itensRows)
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name:   "Nenhum pedido com status",
			status: domain.StatusFinalizado,
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"id", "cliente_id", "valor_total", "status", "created_at", "updated_at"})

				mock.ExpectQuery("SELECT (.+) FROM pedidos WHERE status").
					WithArgs(domain.StatusFinalizado).
					WillReturnRows(rows)
			},
			wantLen: 0,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			pedidos, err := repo.ListarPorStatus(context.Background(), tt.status)

			if (err != nil) != tt.wantErr {
				t.Errorf("ListarPorStatus() erro = %v, wantErr %v", err, tt.wantErr)
			}

			if len(pedidos) != tt.wantLen {
				t.Errorf("ListarPorStatus() len = %d, wantLen %d", len(pedidos), tt.wantLen)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Expectativas não atendidas: %v", err)
			}
		})
	}
}

func TestPedidoRepository_ListarPorCliente(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Erro ao criar mock: %v", err)
	}
	defer db.Close()

	repo := NovoPedidoRepository(db)

	tests := []struct {
		name      string
		clienteID string
		mockSetup func()
		wantLen   int
		wantErr   bool
	}{
		{
			name:      "Sucesso ao listar por cliente",
			clienteID: "cliente-123",
			mockSetup: func() {
				pedidosRows := sqlmock.NewRows([]string{"id", "cliente_id", "valor_total", "status", "created_at", "updated_at"}).
					AddRow("pedido-1", "cliente-123", 50.00, "RECEBIDO",
						"2024-01-01 10:00:00", "2024-01-01 10:00:00")

				mock.ExpectQuery("SELECT (.+) FROM pedidos WHERE cliente_id").
					WithArgs("cliente-123").
					WillReturnRows(pedidosRows)

				itensRows := sqlmock.NewRows([]string{"produto_id", "nome", "preco", "quantidade", "observacao"}).
					AddRow("prod-1", "Item 1", 25.00, 2, nil)

				mock.ExpectQuery("SELECT (.+) FROM pedido_itens WHERE pedido_id").
					WithArgs("pedido-1").
					WillReturnRows(itensRows)
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name:      "Cliente sem pedidos",
			clienteID: "cliente-999",
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"id", "cliente_id", "valor_total", "status", "created_at", "updated_at"})

				mock.ExpectQuery("SELECT (.+) FROM pedidos WHERE cliente_id").
					WithArgs("cliente-999").
					WillReturnRows(rows)
			},
			wantLen: 0,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			pedidos, err := repo.ListarPorCliente(context.Background(), tt.clienteID)

			if (err != nil) != tt.wantErr {
				t.Errorf("ListarPorCliente() erro = %v, wantErr %v", err, tt.wantErr)
			}

			if len(pedidos) != tt.wantLen {
				t.Errorf("ListarPorCliente() len = %d, wantLen %d", len(pedidos), tt.wantLen)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Expectativas não atendidas: %v", err)
			}
		})
	}
}

func TestPedidoRepository_Atualizar(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Erro ao criar mock: %v", err)
	}
	defer db.Close()

	repo := NovoPedidoRepository(db)
	now := time.Now()

	tests := []struct {
		name      string
		pedido    *domain.Pedido
		mockSetup func()
		wantErr   bool
	}{
		{
			name: "Sucesso ao atualizar pedido",
			pedido: &domain.Pedido{
				ID:        "pedido-123",
				Status:    domain.StatusEmPreparacao,
				UpdatedAt: now,
			},
			mockSetup: func() {
				mock.ExpectPrepare("UPDATE pedidos").
					ExpectExec().
					WithArgs(domain.StatusEmPreparacao, sqlmock.AnyArg(), "pedido-123").
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name: "Pedido não encontrado",
			pedido: &domain.Pedido{
				ID:        "pedido-999",
				Status:    domain.StatusEmPreparacao,
				UpdatedAt: now,
			},
			mockSetup: func() {
				mock.ExpectPrepare("UPDATE pedidos").
					ExpectExec().
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			wantErr: true,
		},
		{
			name: "Erro ao atualizar",
			pedido: &domain.Pedido{
				ID:        "pedido-123",
				Status:    domain.StatusEmPreparacao,
				UpdatedAt: now,
			},
			mockSetup: func() {
				mock.ExpectPrepare("UPDATE pedidos").
					ExpectExec().
					WillReturnError(errors.New("erro no banco"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			err := repo.Atualizar(context.Background(), tt.pedido)

			if (err != nil) != tt.wantErr {
				t.Errorf("Atualizar() erro = %v, wantErr %v", err, tt.wantErr)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Expectativas não atendidas: %v", err)
			}
		})
	}
}
