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

func TestNovoPagamentoRepository(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Erro ao criar mock: %v", err)
	}
	defer db.Close()

	repo := NovoPagamentoRepository(db)

	if repo == nil {
		t.Fatal("NovoPagamentoRepository() retornou nil")
	}

	if repo.db == nil {
		t.Error("db não foi inicializado")
	}
}

func TestPagamentoRepository_Criar(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Erro ao criar mock: %v", err)
	}
	defer db.Close()

	repo := NovoPagamentoRepository(db)
	now := time.Now()
	expiresAt := now.Add(30 * time.Minute)

	tests := []struct {
		name      string
		pagamento *domain.Pagamento
		mockSetup func()
		wantErr   bool
	}{
		{
			name: "Sucesso ao criar pagamento",
			pagamento: &domain.Pagamento{
				ID:         "pag-123",
				PedidoID:   "pedido-123",
				Valor:      100.50,
				Status:     domain.StatusPagamentoPendente,
				Tipo:       domain.TipoPagamentoPix,
				QRCodeData: "qr-code",
				ExternalID: "ext-123",
				PaymentURL: "https://payment.url",
				ExpiresAt:  &expiresAt,
				CreatedAt:  now,
				UpdatedAt:  now,
			},
			mockSetup: func() {
				mock.ExpectPrepare("INSERT INTO pagamentos").
					ExpectExec().
					WithArgs("pag-123", "pedido-123", 100.50, domain.StatusPagamentoPendente,
						domain.TipoPagamentoPix, "qr-code", "ext-123", "https://payment.url",
						sqlmock.AnyArg(), nil, sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "Erro ao preparar statement",
			pagamento: &domain.Pagamento{
				ID:        "pag-123",
				CreatedAt: now,
				UpdatedAt: now,
			},
			mockSetup: func() {
				mock.ExpectPrepare("INSERT INTO pagamentos").
					WillReturnError(errors.New("erro ao preparar"))
			},
			wantErr: true,
		},
		{
			name: "Erro ao executar insert",
			pagamento: &domain.Pagamento{
				ID:        "pag-123",
				CreatedAt: now,
				UpdatedAt: now,
			},
			mockSetup: func() {
				mock.ExpectPrepare("INSERT INTO pagamentos").
					ExpectExec().
					WillReturnError(errors.New("erro ao executar"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			err := repo.Criar(context.Background(), tt.pagamento)

			if (err != nil) != tt.wantErr {
				t.Errorf("Criar() erro = %v, wantErr %v", err, tt.wantErr)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Expectativas não atendidas: %v", err)
			}
		})
	}
}

func TestPagamentoRepository_BuscarPorID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Erro ao criar mock: %v", err)
	}
	defer db.Close()

	repo := NovoPagamentoRepository(db)

	tests := []struct {
		name      string
		id        string
		mockSetup func()
		wantNil   bool
		wantErr   bool
	}{
		{
			name: "Sucesso ao buscar pagamento",
			id:   "pag-123",
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"id", "pedido_id", "valor", "status", "tipo", "qr_code_data", "external_id", "payment_url", "expires_at", "processed_at", "created_at", "updated_at"}).
					AddRow("pag-123", "pedido-123", 100.50, "PENDENTE", "PIX", "qr-code", "ext-123", "https://url", nil, nil,
						"2024-01-01 10:00:00", "2024-01-01 10:00:00")

				mock.ExpectPrepare("SELECT (.+) FROM pagamentos WHERE id").
					ExpectQuery().
					WithArgs("pag-123").
					WillReturnRows(rows)
			},
			wantNil: false,
			wantErr: false,
		},
		{
			name: "Pagamento não encontrado",
			id:   "pag-999",
			mockSetup: func() {
				mock.ExpectPrepare("SELECT (.+) FROM pagamentos WHERE id").
					ExpectQuery().
					WithArgs("pag-999").
					WillReturnError(sql.ErrNoRows)
			},
			wantNil: true,
			wantErr: false,
		},
		{
			name: "Erro ao buscar pagamento",
			id:   "pag-123",
			mockSetup: func() {
				mock.ExpectPrepare("SELECT (.+) FROM pagamentos WHERE id").
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

			pagamento, err := repo.BuscarPorID(context.Background(), tt.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("BuscarPorID() erro = %v, wantErr %v", err, tt.wantErr)
			}

			if (pagamento == nil) != tt.wantNil {
				t.Errorf("BuscarPorID() pagamento = %v, wantNil %v", pagamento, tt.wantNil)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Expectativas não atendidas: %v", err)
			}
		})
	}
}

func TestPagamentoRepository_BuscarPorPedidoID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Erro ao criar mock: %v", err)
	}
	defer db.Close()

	repo := NovoPagamentoRepository(db)

	tests := []struct {
		name      string
		pedidoID  string
		mockSetup func()
		wantNil   bool
		wantErr   bool
	}{
		{
			name:     "Sucesso ao buscar por pedido",
			pedidoID: "pedido-123",
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"id", "pedido_id", "valor", "status", "tipo", "qr_code_data", "external_id", "payment_url", "expires_at", "processed_at", "created_at", "updated_at"}).
					AddRow("pag-123", "pedido-123", 100.50, "PENDENTE", "PIX", "qr-code", "ext-123", "https://url", nil, nil,
						"2024-01-01 10:00:00", "2024-01-01 10:00:00")

				mock.ExpectPrepare("SELECT (.+) FROM pagamentos WHERE pedido_id").
					ExpectQuery().
					WithArgs("pedido-123").
					WillReturnRows(rows)
			},
			wantNil: false,
			wantErr: false,
		},
		{
			name:     "Pagamento não encontrado",
			pedidoID: "pedido-999",
			mockSetup: func() {
				mock.ExpectPrepare("SELECT (.+) FROM pagamentos WHERE pedido_id").
					ExpectQuery().
					WithArgs("pedido-999").
					WillReturnError(sql.ErrNoRows)
			},
			wantNil: true,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			pagamento, err := repo.BuscarPorPedidoID(context.Background(), tt.pedidoID)

			if (err != nil) != tt.wantErr {
				t.Errorf("BuscarPorPedidoID() erro = %v, wantErr %v", err, tt.wantErr)
			}

			if (pagamento == nil) != tt.wantNil {
				t.Errorf("BuscarPorPedidoID() pagamento = %v, wantNil %v", pagamento, tt.wantNil)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Expectativas não atendidas: %v", err)
			}
		})
	}
}

func TestPagamentoRepository_BuscarPorExternalID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Erro ao criar mock: %v", err)
	}
	defer db.Close()

	repo := NovoPagamentoRepository(db)

	tests := []struct {
		name       string
		externalID string
		mockSetup  func()
		wantNil    bool
		wantErr    bool
	}{
		{
			name:       "Sucesso ao buscar por external ID",
			externalID: "ext-123",
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"id", "pedido_id", "valor", "status", "tipo", "qr_code_data", "external_id", "payment_url", "expires_at", "processed_at", "created_at", "updated_at"}).
					AddRow("pag-123", "pedido-123", 100.50, "PENDENTE", "PIX", "qr-code", "ext-123", "https://url", nil, nil,
						"2024-01-01 10:00:00", "2024-01-01 10:00:00")

				mock.ExpectPrepare("SELECT (.+) FROM pagamentos WHERE external_id").
					ExpectQuery().
					WithArgs("ext-123").
					WillReturnRows(rows)
			},
			wantNil: false,
			wantErr: false,
		},
		{
			name:       "External ID não encontrado",
			externalID: "ext-999",
			mockSetup: func() {
				mock.ExpectPrepare("SELECT (.+) FROM pagamentos WHERE external_id").
					ExpectQuery().
					WithArgs("ext-999").
					WillReturnError(sql.ErrNoRows)
			},
			wantNil: true,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			pagamento, err := repo.BuscarPorExternalID(context.Background(), tt.externalID)

			if (err != nil) != tt.wantErr {
				t.Errorf("BuscarPorExternalID() erro = %v, wantErr %v", err, tt.wantErr)
			}

			if (pagamento == nil) != tt.wantNil {
				t.Errorf("BuscarPorExternalID() pagamento = %v, wantNil %v", pagamento, tt.wantNil)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Expectativas não atendidas: %v", err)
			}
		})
	}
}

func TestPagamentoRepository_Atualizar(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Erro ao criar mock: %v", err)
	}
	defer db.Close()

	repo := NovoPagamentoRepository(db)
	now := time.Now()

	tests := []struct {
		name      string
		pagamento *domain.Pagamento
		mockSetup func()
		wantErr   bool
	}{
		{
			name: "Sucesso ao atualizar pagamento",
			pagamento: &domain.Pagamento{
				ID:        "pag-123",
				PedidoID:  "pedido-123",
				Valor:     100.50,
				Status:    domain.StatusPagamentoAprovado,
				UpdatedAt: now,
			},
			mockSetup: func() {
				mock.ExpectPrepare("UPDATE pagamentos").
					ExpectExec().
					WithArgs(100.50, domain.StatusPagamentoAprovado, sqlmock.AnyArg(), sqlmock.AnyArg(),
						sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
						sqlmock.AnyArg(), "pag-123").
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name: "Erro ao atualizar",
			pagamento: &domain.Pagamento{
				ID:        "pag-123",
				UpdatedAt: now,
			},
			mockSetup: func() {
				mock.ExpectPrepare("UPDATE pagamentos").
					ExpectExec().
					WillReturnError(errors.New("erro no banco"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			err := repo.Atualizar(context.Background(), tt.pagamento)

			if (err != nil) != tt.wantErr {
				t.Errorf("Atualizar() erro = %v, wantErr %v", err, tt.wantErr)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Expectativas não atendidas: %v", err)
			}
		})
	}
}

func TestPagamentoRepository_ListarPorStatus(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Erro ao criar mock: %v", err)
	}
	defer db.Close()

	repo := NovoPagamentoRepository(db)

	tests := []struct {
		name      string
		status    domain.StatusPagamento
		mockSetup func()
		wantLen   int
		wantErr   bool
	}{
		{
			name:   "Sucesso ao listar por status",
			status: domain.StatusPagamentoPendente,
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"id", "pedido_id", "valor", "status", "tipo", "qr_code_data", "external_id", "payment_url", "expires_at", "processed_at", "created_at", "updated_at"}).
					AddRow("pag-1", "pedido-1", 100.0, "PENDENTE", "PIX", "qr", "ext-1", "url", nil, nil,
						"2024-01-01 10:00:00", "2024-01-01 10:00:00").
					AddRow("pag-2", "pedido-2", 200.0, "PENDENTE", "PIX", "qr", "ext-2", "url", nil, nil,
						"2024-01-01 10:00:00", "2024-01-01 10:00:00")

				mock.ExpectPrepare("SELECT (.+) FROM pagamentos WHERE status").
					ExpectQuery().
					WithArgs(domain.StatusPagamentoPendente).
					WillReturnRows(rows)
			},
			wantLen: 2,
			wantErr: false,
		},
		{
			name:   "Lista vazia",
			status: domain.StatusPagamentoAprovado,
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"id", "pedido_id", "valor", "status", "tipo", "qr_code_data", "external_id", "payment_url", "expires_at", "processed_at", "created_at", "updated_at"})

				mock.ExpectPrepare("SELECT (.+) FROM pagamentos WHERE status").
					ExpectQuery().
					WithArgs(domain.StatusPagamentoAprovado).
					WillReturnRows(rows)
			},
			wantLen: 0,
			wantErr: false,
		},
		{
			name:   "Erro ao listar",
			status: domain.StatusPagamentoPendente,
			mockSetup: func() {
				mock.ExpectPrepare("SELECT (.+) FROM pagamentos WHERE status").
					ExpectQuery().
					WillReturnError(errors.New("erro no banco"))
			},
			wantLen: 0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			pagamentos, err := repo.ListarPorStatus(context.Background(), tt.status)

			if (err != nil) != tt.wantErr {
				t.Errorf("ListarPorStatus() erro = %v, wantErr %v", err, tt.wantErr)
			}

			if len(pagamentos) != tt.wantLen {
				t.Errorf("ListarPorStatus() len = %d, wantLen %d", len(pagamentos), tt.wantLen)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Expectativas não atendidas: %v", err)
			}
		})
	}
}

func TestPagamentoRepository_ListarExpirados(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Erro ao criar mock: %v", err)
	}
	defer db.Close()

	repo := NovoPagamentoRepository(db)

	tests := []struct {
		name      string
		mockSetup func()
		wantLen   int
		wantErr   bool
	}{
		{
			name: "Sucesso ao listar expirados",
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"id", "pedido_id", "valor", "status", "tipo", "qr_code_data", "external_id", "payment_url", "expires_at", "processed_at", "created_at", "updated_at"}).
					AddRow("pag-1", "pedido-1", 100.0, "PENDENTE", "PIX", "qr", "ext-1", "url", "2024-01-01 09:00:00", nil,
						"2024-01-01 10:00:00", "2024-01-01 10:00:00")

				mock.ExpectPrepare("SELECT (.+) FROM pagamentos (.+) WHERE expires_at IS NOT NULL").
					ExpectQuery().
					WithArgs(domain.StatusPagamentoPendente).
					WillReturnRows(rows)
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "Lista vazia",
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"id", "pedido_id", "valor", "status", "tipo", "qr_code_data", "external_id", "payment_url", "expires_at", "processed_at", "created_at", "updated_at"})

				mock.ExpectPrepare("SELECT (.+) FROM pagamentos (.+) WHERE expires_at IS NOT NULL").
					ExpectQuery().
					WithArgs(domain.StatusPagamentoPendente).
					WillReturnRows(rows)
			},
			wantLen: 0,
			wantErr: false,
		},
		{
			name: "Erro ao listar",
			mockSetup: func() {
				mock.ExpectPrepare("SELECT (.+) FROM pagamentos (.+) WHERE expires_at IS NOT NULL").
					ExpectQuery().
					WillReturnError(errors.New("erro no banco"))
			},
			wantLen: 0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			pagamentos, err := repo.ListarExpirados(context.Background())

			if (err != nil) != tt.wantErr {
				t.Errorf("ListarExpirados() erro = %v, wantErr %v", err, tt.wantErr)
			}

			if len(pagamentos) != tt.wantLen {
				t.Errorf("ListarExpirados() len = %d, wantLen %d", len(pagamentos), tt.wantLen)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Expectativas não atendidas: %v", err)
			}
		})
	}
}
