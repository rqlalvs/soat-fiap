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

func TestNovoProdutoRepository(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Erro ao criar mock: %v", err)
	}
	defer db.Close()

	repo := NovoProdutoRepository(db)

	if repo == nil {
		t.Fatal("NovoProdutoRepository() retornou nil")
	}

	if repo.db == nil {
		t.Error("db não foi inicializado")
	}
}

func TestProdutoRepository_Criar(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Erro ao criar mock: %v", err)
	}
	defer db.Close()

	repo := NovoProdutoRepository(db)
	now := time.Now()

	tests := []struct {
		name      string
		produto   *domain.Produto
		mockSetup func()
		wantErr   bool
	}{
		{
			name: "Sucesso ao criar produto",
			produto: &domain.Produto{
				ID:         "123",
				Nome:       "Hamburguer",
				Descricao:  "Hamburguer artesanal",
				Preco:      25.90,
				Categoria:  domain.CategoriaLanche,
				Disponivel: true,
				CreatedAt:  now,
				UpdatedAt:  now,
			},
			mockSetup: func() {
				mock.ExpectPrepare("INSERT INTO produtos").
					ExpectExec().
					WithArgs("123", "Hamburguer", "Hamburguer artesanal", 25.90, domain.CategoriaLanche, true,
						sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "Erro ao preparar statement",
			produto: &domain.Produto{
				ID:        "123",
				Nome:      "Hamburguer",
				CreatedAt: now,
				UpdatedAt: now,
			},
			mockSetup: func() {
				mock.ExpectPrepare("INSERT INTO produtos").
					WillReturnError(errors.New("erro ao preparar"))
			},
			wantErr: true,
		},
		{
			name: "Erro ao executar insert",
			produto: &domain.Produto{
				ID:        "123",
				Nome:      "Hamburguer",
				CreatedAt: now,
				UpdatedAt: now,
			},
			mockSetup: func() {
				mock.ExpectPrepare("INSERT INTO produtos").
					ExpectExec().
					WillReturnError(errors.New("erro ao executar"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			err := repo.Criar(context.Background(), tt.produto)

			if (err != nil) != tt.wantErr {
				t.Errorf("Criar() erro = %v, wantErr %v", err, tt.wantErr)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Expectativas não atendidas: %v", err)
			}
		})
	}
}

func TestProdutoRepository_BuscarPorID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Erro ao criar mock: %v", err)
	}
	defer db.Close()

	repo := NovoProdutoRepository(db)

	tests := []struct {
		name      string
		id        string
		mockSetup func()
		wantNil   bool
		wantErr   bool
	}{
		{
			name: "Sucesso ao buscar produto",
			id:   "123",
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"id", "nome", "descricao", "preco", "categoria", "disponivel", "created_at", "updated_at"}).
					AddRow("123", "Hamburguer", "Hamburguer artesanal", 25.90, "LANCHE", true,
						"2024-01-01 10:00:00", "2024-01-01 10:00:00")

				mock.ExpectPrepare("SELECT (.+) FROM produtos WHERE id").
					ExpectQuery().
					WithArgs("123").
					WillReturnRows(rows)
			},
			wantNil: false,
			wantErr: false,
		},
		{
			name: "Produto não encontrado",
			id:   "999",
			mockSetup: func() {
				mock.ExpectPrepare("SELECT (.+) FROM produtos WHERE id").
					ExpectQuery().
					WithArgs("999").
					WillReturnError(sql.ErrNoRows)
			},
			wantNil: true,
			wantErr: false,
		},
		{
			name: "Erro ao buscar produto",
			id:   "123",
			mockSetup: func() {
				mock.ExpectPrepare("SELECT (.+) FROM produtos WHERE id").
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

			produto, err := repo.BuscarPorID(context.Background(), tt.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("BuscarPorID() erro = %v, wantErr %v", err, tt.wantErr)
			}

			if (produto == nil) != tt.wantNil {
				t.Errorf("BuscarPorID() produto = %v, wantNil %v", produto, tt.wantNil)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Expectativas não atendidas: %v", err)
			}
		})
	}
}

func TestProdutoRepository_Listar(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Erro ao criar mock: %v", err)
	}
	defer db.Close()

	repo := NovoProdutoRepository(db)

	tests := []struct {
		name      string
		mockSetup func()
		wantLen   int
		wantErr   bool
	}{
		{
			name: "Sucesso ao listar produtos",
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"id", "nome", "descricao", "preco", "categoria", "disponivel", "created_at", "updated_at"}).
					AddRow("1", "Produto 1", "Desc 1", 10.0, "LANCHE", true,
						"2024-01-01 10:00:00", "2024-01-01 10:00:00").
					AddRow("2", "Produto 2", "Desc 2", 20.0, "BEBIDA", true,
						"2024-01-01 10:00:00", "2024-01-01 10:00:00")

				mock.ExpectQuery("SELECT (.+) FROM produtos ORDER BY nome").
					WillReturnRows(rows)
			},
			wantLen: 2,
			wantErr: false,
		},
		{
			name: "Lista vazia",
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"id", "nome", "descricao", "preco", "categoria", "disponivel", "created_at", "updated_at"})

				mock.ExpectQuery("SELECT (.+) FROM produtos ORDER BY nome").
					WillReturnRows(rows)
			},
			wantLen: 0,
			wantErr: false,
		},
		{
			name: "Erro ao listar",
			mockSetup: func() {
				mock.ExpectQuery("SELECT (.+) FROM produtos ORDER BY nome").
					WillReturnError(errors.New("erro no banco"))
			},
			wantLen: 0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			produtos, err := repo.Listar(context.Background())

			if (err != nil) != tt.wantErr {
				t.Errorf("Listar() erro = %v, wantErr %v", err, tt.wantErr)
			}

			if len(produtos) != tt.wantLen {
				t.Errorf("Listar() len = %d, wantLen %d", len(produtos), tt.wantLen)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Expectativas não atendidas: %v", err)
			}
		})
	}
}

func TestProdutoRepository_ListarPorCategoria(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Erro ao criar mock: %v", err)
	}
	defer db.Close()

	repo := NovoProdutoRepository(db)

	tests := []struct {
		name      string
		categoria domain.Categoria
		mockSetup func()
		wantLen   int
		wantErr   bool
	}{
		{
			name:      "Sucesso ao listar por categoria",
			categoria: domain.CategoriaLanche,
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"id", "nome", "descricao", "preco", "categoria", "disponivel", "created_at", "updated_at"}).
					AddRow("1", "Hamburguer", "Desc", 25.90, "LANCHE", true,
						"2024-01-01 10:00:00", "2024-01-01 10:00:00")

				mock.ExpectQuery("SELECT (.+) FROM produtos WHERE categoria").
					WithArgs(domain.CategoriaLanche).
					WillReturnRows(rows)
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name:      "Categoria sem produtos",
			categoria: domain.CategoriaSobremesa,
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"id", "nome", "descricao", "preco", "categoria", "disponivel", "created_at", "updated_at"})

				mock.ExpectQuery("SELECT (.+) FROM produtos WHERE categoria").
					WithArgs(domain.CategoriaSobremesa).
					WillReturnRows(rows)
			},
			wantLen: 0,
			wantErr: false,
		},
		{
			name:      "Erro ao listar por categoria",
			categoria: domain.CategoriaLanche,
			mockSetup: func() {
				mock.ExpectQuery("SELECT (.+) FROM produtos WHERE categoria").
					WillReturnError(errors.New("erro no banco"))
			},
			wantLen: 0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			produtos, err := repo.ListarPorCategoria(context.Background(), tt.categoria)

			if (err != nil) != tt.wantErr {
				t.Errorf("ListarPorCategoria() erro = %v, wantErr %v", err, tt.wantErr)
			}

			if len(produtos) != tt.wantLen {
				t.Errorf("ListarPorCategoria() len = %d, wantLen %d", len(produtos), tt.wantLen)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Expectativas não atendidas: %v", err)
			}
		})
	}
}

func TestProdutoRepository_Atualizar(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Erro ao criar mock: %v", err)
	}
	defer db.Close()

	repo := NovoProdutoRepository(db)
	now := time.Now()

	tests := []struct {
		name      string
		produto   *domain.Produto
		mockSetup func()
		wantErr   bool
	}{
		{
			name: "Sucesso ao atualizar produto",
			produto: &domain.Produto{
				ID:         "123",
				Nome:       "Hamburguer Atualizado",
				Descricao:  "Nova descrição",
				Preco:      30.00,
				Categoria:  domain.CategoriaLanche,
				Disponivel: true,
				UpdatedAt:  now,
			},
			mockSetup: func() {
				mock.ExpectPrepare("UPDATE produtos").
					ExpectExec().
					WithArgs("Hamburguer Atualizado", "Nova descrição", 30.00, domain.CategoriaLanche, true,
						sqlmock.AnyArg(), "123").
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name: "Produto não encontrado",
			produto: &domain.Produto{
				ID:        "999",
				Nome:      "Produto",
				UpdatedAt: now,
			},
			mockSetup: func() {
				mock.ExpectPrepare("UPDATE produtos").
					ExpectExec().
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			wantErr: true,
		},
		{
			name: "Erro ao atualizar",
			produto: &domain.Produto{
				ID:        "123",
				Nome:      "Produto",
				UpdatedAt: now,
			},
			mockSetup: func() {
				mock.ExpectPrepare("UPDATE produtos").
					ExpectExec().
					WillReturnError(errors.New("erro no banco"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			err := repo.Atualizar(context.Background(), tt.produto)

			if (err != nil) != tt.wantErr {
				t.Errorf("Atualizar() erro = %v, wantErr %v", err, tt.wantErr)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Expectativas não atendidas: %v", err)
			}
		})
	}
}

func TestProdutoRepository_Deletar(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Erro ao criar mock: %v", err)
	}
	defer db.Close()

	repo := NovoProdutoRepository(db)

	tests := []struct {
		name      string
		id        string
		mockSetup func()
		wantErr   bool
	}{
		{
			name: "Sucesso ao deletar produto",
			id:   "123",
			mockSetup: func() {
				mock.ExpectPrepare("DELETE FROM produtos").
					ExpectExec().
					WithArgs("123").
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name: "Produto não encontrado",
			id:   "999",
			mockSetup: func() {
				mock.ExpectPrepare("DELETE FROM produtos").
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
				mock.ExpectPrepare("DELETE FROM produtos").
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
