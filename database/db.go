package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	dbHost  = "localhost"
	dbPort  = 5432
	dbUser  = "api_conta"
	dbPass  = "sales123"
	dbName  = "AccountBank"
	sslMode = "disable"
)

func Conectar() *sql.DB {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPass, dbName, sslMode)

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
	CREATE TABLE IF NOT EXISTS contas (
		id SERIAL PRIMARY KEY,
		nome VARCHAR(100) NOT NULL,
		email VARCHAR(100) UNIQUE NOT NULL,
		telefone VARCHAR(20),
		cpf VARCHAR(14) UNIQUE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Erro ao criar/verificar tabela 'contas': %v", err)
	}
	fmt.Println("Tabela 'contas' verificada/criada com sucesso.")

	return db
}
