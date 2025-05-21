package repositories

import (
	"context"
	"database/sql"
	"soat-fiap/internal/core/domain"
	"time"
)

type PedidoRepository struct {
	db *sql.DB
}

func NovoPedidoRepository(db *sql.DB) *PedidoRepository {
	return &PedidoRepository{
		db: db,
	}
}

func (r *PedidoRepository) Criar(ctx context.Context, pedido *domain.Pedido) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO pedidos (id, cliente_id, valor_total, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		pedido.ID,
		pedido.ClienteID,
		pedido.ValorTotal,
		pedido.Status,
		pedido.CreatedAt.Format(time.RFC3339),
		pedido.UpdatedAt.Format(time.RFC3339),
	)
	if err != nil {
		return err
	}

	stmtItem, err := tx.PrepareContext(ctx, `
		INSERT INTO pedido_itens (pedido_id, produto_id, nome, preco, quantidade, observacao)
		VALUES (?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmtItem.Close()

	for _, item := range pedido.Itens {
		_, err = stmtItem.ExecContext(ctx,
			pedido.ID,
			item.ProdutoID,
			item.Nome,
			item.Preco,
			item.Quantidade,
			item.Observacao,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *PedidoRepository) BuscarPorID(ctx context.Context, id string) (*domain.Pedido, error) {
	stmt, err := r.db.PrepareContext(ctx, `
		SELECT id, cliente_id, valor_total, status, created_at, updated_at
		FROM pedidos
		WHERE id = ?
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var pedido domain.Pedido
	var createdAtStr, updatedAtStr string
	var clienteID sql.NullString

	err = stmt.QueryRowContext(ctx, id).Scan(
		&pedido.ID,
		&clienteID,
		&pedido.ValorTotal,
		&pedido.Status,
		&createdAtStr,
		&updatedAtStr,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if clienteID.Valid {
		pedido.ClienteID = &clienteID.String
	}

	pedido.CreatedAt, _ = time.Parse(time.RFC3339, createdAtStr)
	pedido.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAtStr)

	itens, err := r.buscarItensPorPedidoID(ctx, id)
	if err != nil {
		return nil, err
	}

	pedido.Itens = itens

	return &pedido, nil
}

func (r *PedidoRepository) Listar(ctx context.Context) ([]*domain.Pedido, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, cliente_id, valor_total, status, created_at, updated_at
		FROM pedidos
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.processarResultados(ctx, rows)
}

func (r *PedidoRepository) ListarPorStatus(ctx context.Context, status domain.StatusPedido) ([]*domain.Pedido, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, cliente_id, valor_total, status, created_at, updated_at
		FROM pedidos
		WHERE status = ?
		ORDER BY created_at ASC
	`, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.processarResultados(ctx, rows)
}

func (r *PedidoRepository) ListarPorCliente(ctx context.Context, clienteID string) ([]*domain.Pedido, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, cliente_id, valor_total, status, created_at, updated_at
		FROM pedidos
		WHERE cliente_id = ?
		ORDER BY created_at DESC
	`, clienteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.processarResultados(ctx, rows)
}

func (r *PedidoRepository) Atualizar(ctx context.Context, pedido *domain.Pedido) error {
	stmt, err := r.db.PrepareContext(ctx, `
		UPDATE pedidos
		SET status = ?, updated_at = ?
		WHERE id = ?
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx,
		pedido.Status,
		pedido.UpdatedAt.Format(time.RFC3339),
		pedido.ID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *PedidoRepository) processarResultados(ctx context.Context, rows *sql.Rows) ([]*domain.Pedido, error) {
	var pedidos []*domain.Pedido

	for rows.Next() {
		var pedido domain.Pedido
		var createdAtStr, updatedAtStr string
		var clienteID sql.NullString

		err := rows.Scan(
			&pedido.ID,
			&clienteID,
			&pedido.ValorTotal,
			&pedido.Status,
			&createdAtStr,
			&updatedAtStr,
		)
		if err != nil {
			return nil, err
		}

		if clienteID.Valid {
			pedido.ClienteID = &clienteID.String
		}

		pedido.CreatedAt, _ = time.Parse(time.RFC3339, createdAtStr)
		pedido.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAtStr)

		itens, err := r.buscarItensPorPedidoID(ctx, pedido.ID)
		if err != nil {
			return nil, err
		}

		pedido.Itens = itens
		pedidos = append(pedidos, &pedido)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return pedidos, nil
}

func (r *PedidoRepository) buscarItensPorPedidoID(ctx context.Context, pedidoID string) ([]domain.ItemPedido, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT produto_id, nome, preco, quantidade, observacao
		FROM pedido_itens
		WHERE pedido_id = ?
		ORDER BY produto_id
	`, pedidoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var itens []domain.ItemPedido

	for rows.Next() {
		var item domain.ItemPedido
		var observacao sql.NullString

		err := rows.Scan(
			&item.ProdutoID,
			&item.Nome,
			&item.Preco,
			&item.Quantidade,
			&observacao,
		)
		if err != nil {
			return nil, err
		}

		if observacao.Valid {
			item.Observacao = observacao.String
		}

		itens = append(itens, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return itens, nil
}
