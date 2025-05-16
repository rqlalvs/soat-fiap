package repositories

import (
	"context"
	"database/sql"
	"soat-fiap/internal/core/domain"
	"time"
)

type ClienteRepository struct {
	db *sql.DB
}

func NovoClienteRepository(db *sql.DB) *ClienteRepository {
	return &ClienteRepository{
		db: db,
	}
}

func (r *ClienteRepository) Criar(ctx context.Context, cliente *domain.Cliente) error {
	stmt, err := r.db.PrepareContext(ctx, `
		INSERT INTO clientes (id, nome, cpf, email, telefone, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		cliente.ID,
		cliente.Nome,
		cliente.CPF,
		cliente.Email,
		cliente.Telefone,
		cliente.CreatedAt.Format(time.RFC3339),
		cliente.UpdatedAt.Format(time.RFC3339),
	)

	return err
}

func (r *ClienteRepository) BuscarPorID(ctx context.Context, id string) (*domain.Cliente, error) {
	stmt, err := r.db.PrepareContext(ctx, `
		SELECT id, nome, cpf, email, telefone, created_at, updated_at
		FROM clientes
		WHERE id = ?
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var cliente domain.Cliente
	var createdAtStr, updatedAtStr string

	err = stmt.QueryRowContext(ctx, id).Scan(
		&cliente.ID,
		&cliente.Nome,
		&cliente.CPF,
		&cliente.Email,
		&cliente.Telefone,
		&createdAtStr,
		&updatedAtStr,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	cliente.CreatedAt, _ = time.Parse(time.RFC3339, createdAtStr)
	cliente.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAtStr)

	return &cliente, nil
}

func (r *ClienteRepository) BuscarPorCPF(ctx context.Context, cpf string) (*domain.Cliente, error) {
	stmt, err := r.db.PrepareContext(ctx, `
		SELECT id, nome, cpf, email, telefone, created_at, updated_at
		FROM clientes
		WHERE cpf = ?
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var cliente domain.Cliente
	var createdAtStr, updatedAtStr string

	err = stmt.QueryRowContext(ctx, cpf).Scan(
		&cliente.ID,
		&cliente.Nome,
		&cliente.CPF,
		&cliente.Email,
		&cliente.Telefone,
		&createdAtStr,
		&updatedAtStr,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	cliente.CreatedAt, _ = time.Parse(time.RFC3339, createdAtStr)
	cliente.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAtStr)

	return &cliente, nil
}

func (r *ClienteRepository) Listar(ctx context.Context) ([]*domain.Cliente, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, nome, cpf, email, telefone, created_at, updated_at
		FROM clientes
		ORDER BY nome
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clientes []*domain.Cliente

	for rows.Next() {
		var cliente domain.Cliente
		var createdAtStr, updatedAtStr string

		err := rows.Scan(
			&cliente.ID,
			&cliente.Nome,
			&cliente.CPF,
			&cliente.Email,
			&cliente.Telefone,
			&createdAtStr,
			&updatedAtStr,
		)
		if err != nil {
			return nil, err
		}

		cliente.CreatedAt, _ = time.Parse(time.RFC3339, createdAtStr)
		cliente.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAtStr)

		clientes = append(clientes, &cliente)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return clientes, nil
}

func (r *ClienteRepository) Atualizar(ctx context.Context, cliente *domain.Cliente) error {
	stmt, err := r.db.PrepareContext(ctx, `
		UPDATE clientes
		SET nome = ?, cpf = ?, email = ?, telefone = ?, updated_at = ?
		WHERE id = ?
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx,
		cliente.Nome,
		cliente.CPF,
		cliente.Email,
		cliente.Telefone,
		cliente.UpdatedAt.Format(time.RFC3339),
		cliente.ID,
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

func (r *ClienteRepository) Deletar(ctx context.Context, id string) error {
	stmt, err := r.db.PrepareContext(ctx, `
		DELETE FROM clientes
		WHERE id = ?
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, id)
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
