package sqlite

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func ConectarDB(arquivo string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", arquivo)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	if err = iniciarTabelas(db); err != nil {
		return nil, err
	}

	return db, nil
}

func iniciarTabelas(db *sql.DB) error {

	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS clientes (
			id TEXT PRIMARY KEY,
			nome TEXT NOT NULL,
			cpf TEXT NOT NULL UNIQUE,
			email TEXT NOT NULL,
			telefone TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_clientes_cpf ON clientes(cpf);
	`)

	if err != nil {
		log.Printf("Erro ao criar tabelas: %v", err)
		return err
	}

	return nil
}
