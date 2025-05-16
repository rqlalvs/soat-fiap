package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func ConectarMySQL(host, port, user, password, dbname string) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, password, host, port, dbname)
	db, err := sql.Open("mysql", dsn)
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
	queries := []string{
		`CREATE TABLE IF NOT EXISTS clientes (
			id VARCHAR(36) PRIMARY KEY,
			nome VARCHAR(100) NOT NULL,
			cpf VARCHAR(11) NOT NULL UNIQUE,
			email VARCHAR(100) NOT NULL,
			telefone VARCHAR(20) NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			INDEX idx_clientes_cpf (cpf)
		)`,
		`CREATE TABLE IF NOT EXISTS produtos (
			id VARCHAR(36) PRIMARY KEY,
			nome VARCHAR(100) NOT NULL,
			descricao TEXT NOT NULL,
			preco DECIMAL(10,2) NOT NULL,
			categoria VARCHAR(20) NOT NULL,
			disponivel BOOLEAN NOT NULL DEFAULT TRUE,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			INDEX idx_produtos_categoria (categoria)
		)`,
	}

	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			log.Printf("Erro ao criar tabela: %v", err)
			return err
		}
	}

	return nil
}
