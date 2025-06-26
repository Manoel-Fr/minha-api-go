package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"minha-api-go/models"
)

type API struct {
	DB *sql.DB
}

func (api *API) CriarConta(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var conta models.Conta
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
	json.NewEncoder(w).Encode(models.Resposta{
		Status:   "sucesso",
		Mensagem: "Conta criada com sucesso! ID: " + string(rune(conta.ID)),
	})
}
