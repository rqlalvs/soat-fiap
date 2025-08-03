package repositories

import (
	"context"
	"database/sql"
	"soat-fiap/internal/core/domain"
	"time"
)

type PagamentoRepository struct {
	db *sql.DB
}

func NovoPagamentoRepository(db *sql.DB) *PagamentoRepository {
	return &PagamentoRepository{
		db: db,
	}
}

func (r *PagamentoRepository) Criar(ctx context.Context, pagamento *domain.Pagamento) error {
	stmt, err := r.db.PrepareContext(ctx, `
		INSERT INTO pagamentos (id, pedido_id, valor, status, tipo, qr_code_data, external_id, payment_url, expires_at, processed_at, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var expiresAt, processedAt *string
	if pagamento.ExpiresAt != nil {
		expiresAtStr := pagamento.ExpiresAt.Format("2006-01-02 15:04:05")
		expiresAt = &expiresAtStr
	}
	if pagamento.ProcessedAt != nil {
		processedAtStr := pagamento.ProcessedAt.Format("2006-01-02 15:04:05")
		processedAt = &processedAtStr
	}

	_, err = stmt.ExecContext(ctx,
		pagamento.ID,
		pagamento.PedidoID,
		pagamento.Valor,
		pagamento.Status,
		pagamento.Tipo,
		pagamento.QRCodeData,
		pagamento.ExternalID,
		pagamento.PaymentURL,
		expiresAt,
		processedAt,
		pagamento.CreatedAt.Format("2006-01-02 15:04:05"),
		pagamento.UpdatedAt.Format("2006-01-02 15:04:05"),
	)

	return err
}

func (r *PagamentoRepository) BuscarPorID(ctx context.Context, id string) (*domain.Pagamento, error) {
	stmt, err := r.db.PrepareContext(ctx, `
		SELECT id, pedido_id, valor, status, tipo, qr_code_data, external_id, payment_url, expires_at, processed_at, created_at, updated_at
		FROM pagamentos WHERE id = ?
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, id)
	return r.scanPagamento(row)
}

func (r *PagamentoRepository) BuscarPorPedidoID(ctx context.Context, pedidoID string) (*domain.Pagamento, error) {
	stmt, err := r.db.PrepareContext(ctx, `
		SELECT id, pedido_id, valor, status, tipo, qr_code_data, external_id, payment_url, expires_at, processed_at, created_at, updated_at
		FROM pagamentos WHERE pedido_id = ?
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, pedidoID)
	return r.scanPagamento(row)
}

func (r *PagamentoRepository) BuscarPorExternalID(ctx context.Context, externalID string) (*domain.Pagamento, error) {
	stmt, err := r.db.PrepareContext(ctx, `
		SELECT id, pedido_id, valor, status, tipo, qr_code_data, external_id, payment_url, expires_at, processed_at, created_at, updated_at
		FROM pagamentos WHERE external_id = ?
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, externalID)
	return r.scanPagamento(row)
}

func (r *PagamentoRepository) Atualizar(ctx context.Context, pagamento *domain.Pagamento) error {
	stmt, err := r.db.PrepareContext(ctx, `
		UPDATE pagamentos 
		SET valor = ?, status = ?, tipo = ?, qr_code_data = ?, external_id = ?, payment_url = ?, expires_at = ?, processed_at = ?, updated_at = ?
		WHERE id = ?
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var expiresAt, processedAt *string
	if pagamento.ExpiresAt != nil {
		expiresAtStr := pagamento.ExpiresAt.Format("2006-01-02 15:04:05")
		expiresAt = &expiresAtStr
	}
	if pagamento.ProcessedAt != nil {
		processedAtStr := pagamento.ProcessedAt.Format("2006-01-02 15:04:05")
		processedAt = &processedAtStr
	}

	_, err = stmt.ExecContext(ctx,
		pagamento.Valor,
		pagamento.Status,
		pagamento.Tipo,
		pagamento.QRCodeData,
		pagamento.ExternalID,
		pagamento.PaymentURL,
		expiresAt,
		processedAt,
		pagamento.UpdatedAt.Format("2006-01-02 15:04:05"),
		pagamento.ID,
	)

	return err
}

func (r *PagamentoRepository) ListarPorStatus(ctx context.Context, status domain.StatusPagamento) ([]*domain.Pagamento, error) {
	stmt, err := r.db.PrepareContext(ctx, `
		SELECT id, pedido_id, valor, status, tipo, qr_code_data, external_id, payment_url, expires_at, processed_at, created_at, updated_at
		FROM pagamentos WHERE status = ?
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanPagamentos(rows)
}

func (r *PagamentoRepository) ListarExpirados(ctx context.Context) ([]*domain.Pagamento, error) {
	stmt, err := r.db.PrepareContext(ctx, `
		SELECT id, pedido_id, valor, status, tipo, qr_code_data, external_id, payment_url, expires_at, processed_at, created_at, updated_at
		FROM pagamentos 
		WHERE expires_at IS NOT NULL AND expires_at < NOW() AND status = ?
		ORDER BY expires_at ASC
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, domain.StatusPagamentoPendente)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanPagamentos(rows)
}

func (r *PagamentoRepository) scanPagamento(row *sql.Row) (*domain.Pagamento, error) {
	var pagamento domain.Pagamento
	var expiresAt, processedAt sql.NullString
	var createdAtStr, updatedAtStr string

	err := row.Scan(
		&pagamento.ID,
		&pagamento.PedidoID,
		&pagamento.Valor,
		&pagamento.Status,
		&pagamento.Tipo,
		&pagamento.QRCodeData,
		&pagamento.ExternalID,
		&pagamento.PaymentURL,
		&expiresAt,
		&processedAt,
		&createdAtStr,
		&updatedAtStr,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// Parse das datas
	pagamento.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAtStr)
	if err != nil {
		return nil, err
	}

	pagamento.UpdatedAt, err = time.Parse("2006-01-02 15:04:05", updatedAtStr)
	if err != nil {
		return nil, err
	}

	if expiresAt.Valid {
		expiresAtTime, err := time.Parse("2006-01-02 15:04:05", expiresAt.String)
		if err != nil {
			return nil, err
		}
		pagamento.ExpiresAt = &expiresAtTime
	}

	if processedAt.Valid {
		processedAtTime, err := time.Parse("2006-01-02 15:04:05", processedAt.String)
		if err != nil {
			return nil, err
		}
		pagamento.ProcessedAt = &processedAtTime
	}

	return &pagamento, nil
}

func (r *PagamentoRepository) scanPagamentos(rows *sql.Rows) ([]*domain.Pagamento, error) {
	var pagamentos []*domain.Pagamento

	for rows.Next() {
		var pagamento domain.Pagamento
		var expiresAt, processedAt sql.NullString
		var createdAtStr, updatedAtStr string

		err := rows.Scan(
			&pagamento.ID,
			&pagamento.PedidoID,
			&pagamento.Valor,
			&pagamento.Status,
			&pagamento.Tipo,
			&pagamento.QRCodeData,
			&pagamento.ExternalID,
			&pagamento.PaymentURL,
			&expiresAt,
			&processedAt,
			&createdAtStr,
			&updatedAtStr,
		)

		if err != nil {
			return nil, err
		}

		// Parse das datas
		pagamento.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAtStr)
		if err != nil {
			return nil, err
		}

		pagamento.UpdatedAt, err = time.Parse("2006-01-02 15:04:05", updatedAtStr)
		if err != nil {
			return nil, err
		}

		if expiresAt.Valid {
			expiresAtTime, err := time.Parse("2006-01-02 15:04:05", expiresAt.String)
			if err != nil {
				return nil, err
			}
			pagamento.ExpiresAt = &expiresAtTime
		}

		if processedAt.Valid {
			processedAtTime, err := time.Parse("2006-01-02 15:04:05", processedAt.String)
			if err != nil {
				return nil, err
			}
			pagamento.ProcessedAt = &processedAtTime
		}

		pagamentos = append(pagamentos, &pagamento)
	}

	return pagamentos, rows.Err()
}
