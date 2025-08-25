package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/lib/pq"
)

var (
	dbHost = getEnvOrDefault("DB_HOST", "localhost")
	dbPort = getEnvIntOrDefault("DB_PORT", 5432)
	dbUser = getEnvOrDefault("DB_USER", "salesuser")
	dbPass = getEnvOrDefault("DB_PASS", "2202data")
	dbName = getEnvOrDefault("DB_NAME", "salesdb")
)

func Conectar() *sql.DB {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Erro ao configurar a conexão com o banco de dados: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("ERRO: Não foi possível conectar ao banco de dados. Verifique PostgreSQL (rodando?), Firewall, Host/Porta, Usuário/Senha. Detalhes: %v", err)
	}
	fmt.Println("Conexão com o banco de dados estabelecida com sucesso!")

	createTableSQL := `
		CREATE TABLE IF NOT EXISTS Accounts (
			Id SERIAL PRIMARY KEY,
			Name VARCHAR(255) NOT NULL,
			Phone VARCHAR(20),
			Website VARCHAR(255),
			Industry VARCHAR(100)
		);
	`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Erro ao executar consulta: %v", err)
	}
	fmt.Println("Consulta executada com sucesso.")

	return db
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
