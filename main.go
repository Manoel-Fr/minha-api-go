package main

import (
	"database/sql" // Adicionado para operações com banco de dados
	"encoding/json"
	"fmt" // Adicionado para formatação da string de conexão
	"log"
	"net/http"

	_ "github.com/lib/pq" // Driver PostgreSQL (o underscore é importante!)
)

// Conta representa a estrutura dos dados de uma conta para o JSON e para o banco de dados
type Conta struct {
	ID       int    `json:"id"` // Adicionado ID para representar a chave primária do banco
	Nome     string `json:"nome"`
	Email    string `json:"email"`
	Telefone string `json:"telefone"`
	Cpf      string `json:"cpf"`
}

// Resposta é a estrutura padrão para as respostas da API
type Resposta struct {
	Status   string `json:"status"`
	Mensagem string `json:"mensagem"`
}

// Definimos uma struct que conterá a conexão com o banco de dados.
// Isso nos permite passar a conexão para as funções de handler de forma mais limpa.
type API struct {
	DB *sql.DB
}

// criarConta agora é um método de API, recebendo acesso ao DB
func (api *API) criarConta(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var conta Conta
	err := json.NewDecoder(r.Body).Decode(&conta)
	if err != nil {
		log.Printf("Erro ao decodificar JSON: %v", err)
		http.Error(w, "Erro ao processar JSON", http.StatusBadRequest)
		return
	}

	if conta.Nome == "" || conta.Email == "" {
		http.Error(w, "Nome e Email são obrigatórios", http.StatusBadRequest)
		return
	}

	insertSQL := `INSERT INTO contas (nome, email, telefone, cpf) VALUES ($1, $2, $3, $4) RETURNING id;`

	err = api.DB.QueryRow(insertSQL, conta.Nome, conta.Email, conta.Telefone, conta.Cpf).Scan(&conta.ID)
	if err != nil {
		log.Printf("Erro ao inserir conta no banco de dados: %v", err)
		http.Error(w, "Erro interno ao criar conta", http.StatusInternalServerError)
		return
	}

	log.Printf("Conta criada e salva no DB: ID=%d, Nome=%s, Email=%s", conta.ID, conta.Nome, conta.Email)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(Resposta{
		Status:   "sucesso",
		Mensagem: fmt.Sprintf("Conta criada com sucesso! ID: %d", conta.ID),
	})
}

func main() {

	const (
		dbHost  = "localhost"
		dbPort  = 5432
		dbUser  = "api_conta"
		dbPass  = "sales123"
		dbName  = "AccountBank"
		sslMode = "disable"
	)

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPass, dbName, sslMode)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Erro ao configurar a conexão com o banco de dados: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("ERRO: Não foi possível conectar ao banco de dados. Verifique PostgreSQL (rodando?), Firewall, Host/Porta, Usuário/Senha. Detalhes: %v", err)
	}
	fmt.Println("Conexão com o banco de dados estabelecida com sucesso!")

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS contas (
		id SERIAL PRIMARY KEY,
		nome VARCHAR(100) NOT NULL,
		email VARCHAR(100) UNIQUE NOT NULL, -- Email como UNIQUE para evitar duplicatas
		telefone VARCHAR(20),
		cpf VARCHAR(14) UNIQUE, -- CPF também como UNIQUE
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Erro ao criar/verificar tabela 'contas': %v", err)
	}
	fmt.Println("Tabela 'contas' verificada/criada com sucesso.")

	apiInstance := &API{DB: db}

	http.HandleFunc("/conta", apiInstance.criarConta)

	log.Println("API rodando em http://localhost:8888")
	log.Fatal(http.ListenAndServe(":8888", nil))
}
