package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

const (
	dbHost = "localhost"
	dbPort = 3306 // Porta padrão do MySQL
	dbUser = "root"
	dbPass = "2202data"
	dbName = "SalesforceTest"
)

func Conectar() *sql.DB {
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		dbUser, dbPass, dbHost, dbPort, dbName)

	db, err := sql.Open("mysql", connStr)
	if err != nil {
		log.Fatalf("Erro ao configurar a conexão com o banco de dados: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("ERRO: Não foi possível conectar ao banco de dados. Verifique MySQL (rodando?), Firewall, Host/Porta, Usuário/Senha. Detalhes: %v", err)
	}
	fmt.Println("Conexão com o banco de dados estabelecida com sucesso!")

	createTableSQL := `
		CREATE TABLE IF NOT EXISTS Accounts (
			 Id INT PRIMARY KEY AUTO_INCREMENT,
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
