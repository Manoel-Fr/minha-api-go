package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func Conectar() (*sql.DB, error) {
	connStr := os.Getenv("DATABASE_URL")

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

	return db, nil
}
