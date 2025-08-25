package database

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
)

func Conectar() (*sql.DB, error) {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		return nil, sql.ErrConnDone // se não tiver variável, retorna erro simples
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
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
