package database

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

func Conectar() (*sql.DB, error) {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		return nil, fmt.Errorf("variável de ambiente DATABASE_URL não está definida")
	}

	// Garante que SSL seja usado
	if !strings.HasSuffix(connStr, "?") {
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
