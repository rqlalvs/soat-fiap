package repositories

import (
	"context"
	"database/sql"
	"soat-fiap/internal/core/domain"
	"time"
)

type ProdutoRepository struct {
	db *sql.DB
}

func NovoProdutoRepository(db *sql.DB) *ProdutoRepository {
	return &ProdutoRepository{
		db: db,
	}
}

func (r *ProdutoRepository) Criar(ctx context.Context, produto *domain.Produto) error {
	stmt, err := r.db.PrepareContext(ctx, `
		INSERT INTO produtos (id, nome, descricao, preco, categoria, disponivel, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		produto.ID,
		produto.Nome,
		produto.Descricao,
		produto.Preco,
		produto.Categoria,
		produto.Disponivel,
		produto.CreatedAt.Format(time.RFC3339),
		produto.UpdatedAt.Format(time.RFC3339),
	)

	return err
}

func (r *ProdutoRepository) BuscarPorID(ctx context.Context, id string) (*domain.Produto, error) {
	stmt, err := r.db.PrepareContext(ctx, `
		SELECT id, nome, descricao, preco, categoria, disponivel, created_at, updated_at
		FROM produtos
		WHERE id = ?
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var produto domain.Produto
	var createdAtStr, updatedAtStr string

	err = stmt.QueryRowContext(ctx, id).Scan(
		&produto.ID,
		&produto.Nome,
		&produto.Descricao,
		&produto.Preco,
		&produto.Categoria,
		&produto.Disponivel,
		&createdAtStr,
		&updatedAtStr,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	produto.CreatedAt, _ = time.Parse(time.RFC3339, createdAtStr)
	produto.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAtStr)

	return &produto, nil
}

func (r *ProdutoRepository) Listar(ctx context.Context) ([]*domain.Produto, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, nome, descricao, preco, categoria, disponivel, created_at, updated_at
		FROM produtos
		ORDER BY nome
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var produtos []*domain.Produto

	for rows.Next() {
		var produto domain.Produto
		var createdAtStr, updatedAtStr string

		err := rows.Scan(
			&produto.ID,
			&produto.Nome,
			&produto.Descricao,
			&produto.Preco,
			&produto.Categoria,
			&produto.Disponivel,
			&createdAtStr,
			&updatedAtStr,
		)
		if err != nil {
			return nil, err
		}

		produto.CreatedAt, _ = time.Parse(time.RFC3339, createdAtStr)
		produto.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAtStr)

		produtos = append(produtos, &produto)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return produtos, nil
}

func (r *ProdutoRepository) ListarPorCategoria(ctx context.Context, categoria domain.Categoria) ([]*domain.Produto, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, nome, descricao, preco, categoria, disponivel, created_at, updated_at
		FROM produtos
		WHERE categoria = ?
		ORDER BY nome
	`, categoria)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var produtos []*domain.Produto

	for rows.Next() {
		var produto domain.Produto
		var createdAtStr, updatedAtStr string

		err := rows.Scan(
			&produto.ID,
			&produto.Nome,
			&produto.Descricao,
			&produto.Preco,
			&produto.Categoria,
			&produto.Disponivel,
			&createdAtStr,
			&updatedAtStr,
		)
		if err != nil {
			return nil, err
		}

		produto.CreatedAt, _ = time.Parse(time.RFC3339, createdAtStr)
		produto.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAtStr)

		produtos = append(produtos, &produto)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return produtos, nil
}

func (r *ProdutoRepository) Atualizar(ctx context.Context, produto *domain.Produto) error {
	stmt, err := r.db.PrepareContext(ctx, `
		UPDATE produtos
		SET nome = ?, descricao = ?, preco = ?, categoria = ?, disponivel = ?, updated_at = ?
		WHERE id = ?
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx,
		produto.Nome,
		produto.Descricao,
		produto.Preco,
		produto.Categoria,
		produto.Disponivel,
		produto.UpdatedAt.Format(time.RFC3339),
		produto.ID,
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

func (r *ProdutoRepository) Deletar(ctx context.Context, id string) error {
	stmt, err := r.db.PrepareContext(ctx, `
		DELETE FROM produtos
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
