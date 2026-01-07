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

func TestNovoClienteRepository(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Erro ao criar mock: %v", err)
	}
	defer db.Close()

	repo := NovoClienteRepository(db)

	if repo == nil {
		t.Fatal("NovoClienteRepository() retornou nil")
	}

	if repo.db == nil {
		t.Error("db não foi inicializado")
	}
}

func TestClienteRepository_Criar(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Erro ao criar mock: %v", err)
	}
	defer db.Close()

	repo := NovoClienteRepository(db)
	now := time.Now()

	tests := []struct {
		name      string
		cliente   *domain.Cliente
		mockSetup func()
		wantErr   bool
	}{
		{
			name: "Sucesso ao criar cliente",
			cliente: &domain.Cliente{
				ID:        "123",
				Nome:      "João Silva",
				CPF:       "12345678909",
				Email:     "joao@example.com",
				Telefone:  "11999999999",
				CreatedAt: now,
				UpdatedAt: now,
			},
			mockSetup: func() {
				mock.ExpectPrepare("INSERT INTO clientes").
					ExpectExec().
					WithArgs("123", "João Silva", "12345678909", "joao@example.com", "11999999999",
						sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "Erro ao preparar statement",
			cliente: &domain.Cliente{
				ID:        "123",
				Nome:      "João Silva",
				CPF:       "12345678909",
				Email:     "joao@example.com",
				CreatedAt: now,
				UpdatedAt: now,
			},
			mockSetup: func() {
				mock.ExpectPrepare("INSERT INTO clientes").
					WillReturnError(errors.New("erro ao preparar"))
			},
			wantErr: true,
		},
		{
			name: "Erro ao executar insert",
			cliente: &domain.Cliente{
				ID:        "123",
				Nome:      "João Silva",
				CPF:       "12345678909",
				Email:     "joao@example.com",
				CreatedAt: now,
				UpdatedAt: now,
			},
			mockSetup: func() {
				mock.ExpectPrepare("INSERT INTO clientes").
					ExpectExec().
					WillReturnError(errors.New("erro ao executar"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			err := repo.Criar(context.Background(), tt.cliente)

			if (err != nil) != tt.wantErr {
				t.Errorf("Criar() erro = %v, wantErr %v", err, tt.wantErr)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Expectativas não atendidas: %v", err)
			}
		})
	}
}

func TestClienteRepository_BuscarPorID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Erro ao criar mock: %v", err)
	}
	defer db.Close()

	repo := NovoClienteRepository(db)

	tests := []struct {
		name      string
		id        string
		mockSetup func()
		wantNil   bool
		wantErr   bool
	}{
		{
			name: "Sucesso ao buscar cliente",
			id:   "123",
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"id", "nome", "cpf", "email", "telefone", "created_at", "updated_at"}).
					AddRow("123", "João Silva", "12345678909", "joao@example.com", "11999999999",
						"2024-01-01 10:00:00", "2024-01-01 10:00:00")

				mock.ExpectPrepare("SELECT (.+) FROM clientes WHERE id").
					ExpectQuery().
					WithArgs("123").
					WillReturnRows(rows)
			},
			wantNil: false,
			wantErr: false,
		},
		{
			name: "Cliente não encontrado",
			id:   "999",
			mockSetup: func() {
				mock.ExpectPrepare("SELECT (.+) FROM clientes WHERE id").
					ExpectQuery().
					WithArgs("999").
					WillReturnError(sql.ErrNoRows)
			},
			wantNil: true,
			wantErr: false,
		},
		{
			name: "Erro ao buscar cliente",
			id:   "123",
			mockSetup: func() {
				mock.ExpectPrepare("SELECT (.+) FROM clientes WHERE id").
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

			cliente, err := repo.BuscarPorID(context.Background(), tt.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("BuscarPorID() erro = %v, wantErr %v", err, tt.wantErr)
			}

			if (cliente == nil) != tt.wantNil {
				t.Errorf("BuscarPorID() cliente = %v, wantNil %v", cliente, tt.wantNil)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Expectativas não atendidas: %v", err)
			}
		})
	}
}

func TestClienteRepository_BuscarPorCPF(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Erro ao criar mock: %v", err)
	}
	defer db.Close()

	repo := NovoClienteRepository(db)

	tests := []struct {
		name      string
		cpf       string
		mockSetup func()
		wantNil   bool
		wantErr   bool
	}{
		{
			name: "Sucesso ao buscar por CPF",
			cpf:  "12345678909",
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"id", "nome", "cpf", "email", "telefone", "created_at", "updated_at"}).
					AddRow("123", "João Silva", "12345678909", "joao@example.com", "11999999999",
						"2024-01-01 10:00:00", "2024-01-01 10:00:00")

				mock.ExpectPrepare("SELECT (.+) FROM clientes WHERE cpf").
					ExpectQuery().
					WithArgs("12345678909").
					WillReturnRows(rows)
			},
			wantNil: false,
			wantErr: false,
		},
		{
			name: "CPF não encontrado",
			cpf:  "99999999999",
			mockSetup: func() {
				mock.ExpectPrepare("SELECT (.+) FROM clientes WHERE cpf").
					ExpectQuery().
					WithArgs("99999999999").
					WillReturnError(sql.ErrNoRows)
			},
			wantNil: true,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			cliente, err := repo.BuscarPorCPF(context.Background(), tt.cpf)

			if (err != nil) != tt.wantErr {
				t.Errorf("BuscarPorCPF() erro = %v, wantErr %v", err, tt.wantErr)
			}

			if (cliente == nil) != tt.wantNil {
				t.Errorf("BuscarPorCPF() cliente = %v, wantNil %v", cliente, tt.wantNil)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Expectativas não atendidas: %v", err)
			}
		})
	}
}

func TestClienteRepository_Listar(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Erro ao criar mock: %v", err)
	}
	defer db.Close()

	repo := NovoClienteRepository(db)

	tests := []struct {
		name      string
		mockSetup func()
		wantLen   int
		wantErr   bool
	}{
		{
			name: "Sucesso ao listar clientes",
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"id", "nome", "cpf", "email", "telefone", "created_at", "updated_at"}).
					AddRow("1", "Cliente 1", "11111111111", "c1@example.com", "11111111111",
						"2024-01-01 10:00:00", "2024-01-01 10:00:00").
					AddRow("2", "Cliente 2", "22222222222", "c2@example.com", "22222222222",
						"2024-01-01 10:00:00", "2024-01-01 10:00:00")

				mock.ExpectQuery("SELECT (.+) FROM clientes ORDER BY nome").
					WillReturnRows(rows)
			},
			wantLen: 2,
			wantErr: false,
		},
		{
			name: "Lista vazia",
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"id", "nome", "cpf", "email", "telefone", "created_at", "updated_at"})

				mock.ExpectQuery("SELECT (.+) FROM clientes ORDER BY nome").
					WillReturnRows(rows)
			},
			wantLen: 0,
			wantErr: false,
		},
		{
			name: "Erro ao listar",
			mockSetup: func() {
				mock.ExpectQuery("SELECT (.+) FROM clientes ORDER BY nome").
					WillReturnError(errors.New("erro no banco"))
			},
			wantLen: 0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			clientes, err := repo.Listar(context.Background())

			if (err != nil) != tt.wantErr {
				t.Errorf("Listar() erro = %v, wantErr %v", err, tt.wantErr)
			}

			if len(clientes) != tt.wantLen {
				t.Errorf("Listar() len = %d, wantLen %d", len(clientes), tt.wantLen)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Expectativas não atendidas: %v", err)
			}
		})
	}
}

func TestClienteRepository_Atualizar(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Erro ao criar mock: %v", err)
	}
	defer db.Close()

	repo := NovoClienteRepository(db)
	now := time.Now()

	tests := []struct {
		name      string
		cliente   *domain.Cliente
		mockSetup func()
		wantErr   bool
	}{
		{
			name: "Sucesso ao atualizar cliente",
			cliente: &domain.Cliente{
				ID:        "123",
				Nome:      "João Silva Atualizado",
				CPF:       "12345678909",
				Email:     "joao.novo@example.com",
				Telefone:  "11988888888",
				UpdatedAt: now,
			},
			mockSetup: func() {
				mock.ExpectPrepare("UPDATE clientes").
					ExpectExec().
					WithArgs("João Silva Atualizado", "12345678909", "joao.novo@example.com", "11988888888",
						sqlmock.AnyArg(), "123").
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name: "Cliente não encontrado",
			cliente: &domain.Cliente{
				ID:        "999",
				Nome:      "Cliente",
				UpdatedAt: now,
			},
			mockSetup: func() {
				mock.ExpectPrepare("UPDATE clientes").
					ExpectExec().
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			wantErr: true,
		},
		{
			name: "Erro ao atualizar",
			cliente: &domain.Cliente{
				ID:        "123",
				Nome:      "Cliente",
				UpdatedAt: now,
			},
			mockSetup: func() {
				mock.ExpectPrepare("UPDATE clientes").
					ExpectExec().
					WillReturnError(errors.New("erro no banco"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			err := repo.Atualizar(context.Background(), tt.cliente)

			if (err != nil) != tt.wantErr {
				t.Errorf("Atualizar() erro = %v, wantErr %v", err, tt.wantErr)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Expectativas não atendidas: %v", err)
			}
		})
	}
}

func TestClienteRepository_Deletar(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Erro ao criar mock: %v", err)
	}
	defer db.Close()

	repo := NovoClienteRepository(db)

	tests := []struct {
		name      string
		id        string
		mockSetup func()
		wantErr   bool
	}{
		{
			name: "Sucesso ao deletar cliente",
			id:   "123",
			mockSetup: func() {
				mock.ExpectPrepare("DELETE FROM clientes").
					ExpectExec().
					WithArgs("123").
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name: "Cliente não encontrado",
			id:   "999",
			mockSetup: func() {
				mock.ExpectPrepare("DELETE FROM clientes").
					ExpectExec().
					WithArgs("999").
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			wantErr: true,
		},
		{
			name: "Erro ao deletar",
			id:   "123",
			mockSetup: func() {
				mock.ExpectPrepare("DELETE FROM clientes").
					ExpectExec().
					WillReturnError(errors.New("erro no banco"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			err := repo.Deletar(context.Background(), tt.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("Deletar() erro = %v, wantErr %v", err, tt.wantErr)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Expectativas não atendidas: %v", err)
			}
		})
	}
}
