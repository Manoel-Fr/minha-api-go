package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	dbHost = "localhost"
	dbPort = 5432 // Porta padrão do postgres
	dbUser = "salesuser"
	dbPass = "2202data"
	dbName = "salesdb"
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
