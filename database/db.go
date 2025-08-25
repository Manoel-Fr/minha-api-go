package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func Conectar() (*sql.DB, error) {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		return nil, fmt.Errorf("DATABASE_URL não está configurada")
	}

	// Garante que SSL seja usado
	if connStr[len(connStr)-1] != '?' {
		connStr += "?sslmode=require"
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("erro ao testar conexão com o banco: %v", err)
	}

	// Cria tabela se não existir
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS Accounts (
			Id SERIAL PRIMARY KEY,
			Name VARCHAR(255) NOT NULL,
			Phone VARCHAR(20),
			Website VARCHAR(255),
			Industry VARCHAR(100)
		);
	`)
	if err != nil {
		return nil, err
	}

	return db, nil
}
