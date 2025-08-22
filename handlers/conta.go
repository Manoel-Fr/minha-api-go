package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"minha-api-go/models"
)

type API struct {
	DB *sql.DB
}

// Função utilitária para responder erros de forma padronizada
func respondError(w http.ResponseWriter, status int, mensagem string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(models.Resposta{
		Status:   "erro",
		Mensagem: mensagem,
	})
}

// Função utilitária para responder sucesso
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func (api *API) ListarContas(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Método não permitido")
		return
	}

	rows, err := api.DB.Query("SELECT Id, Name, Phone, Website, Industry FROM Accounts")
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Erro ao buscar contas")
		return
	}
	defer rows.Close()

	var contas []models.Conta
	for rows.Next() {
		var conta models.Conta
		err := rows.Scan(&conta.ID, &conta.Name, &conta.Phone, &conta.Website, &conta.Industry)
		if err != nil {
			fmt.Printf("Erro ao ler conta do banco: %v\n", err)
			continue
		}
		contas = append(contas, conta)
	}

	// Verifica se houve algum erro durante a iteração
	if err = rows.Err(); err != nil {
		fmt.Printf("Erro ao iterar sobre as contas: %v\n", err)
		respondError(w, http.StatusInternalServerError, "Erro ao listar contas")
		return
	}

	fmt.Printf("Total de contas encontradas: %d\n", len(contas))
	respondJSON(w, http.StatusOK, contas)
}

func (api *API) CriarConta(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Método não permitido. Use POST para criar conta")
		return
	}

	var conta models.Conta
	if err := json.NewDecoder(r.Body).Decode(&conta); err != nil {
		respondError(w, http.StatusBadRequest, "Erro ao processar JSON")
		return
	}

	if conta.Name == "" {
		respondError(w, http.StatusBadRequest, "Name é obrigatório")
		return
	}

	insertSQL := `INSERT INTO Accounts (Name, Phone, Website, Industry) VALUES (?, ?, ?, ?)`
	result, err := api.DB.Exec(insertSQL, conta.Name, conta.Phone, conta.Website, conta.Industry)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Erro interno ao criar conta")
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Erro interno ao obter ID da conta criada")
		return
	}
	conta.ID = int(id)

	respondJSON(w, http.StatusCreated, models.Resposta{
		Status:   "sucesso",
		Mensagem: fmt.Sprintf("Conta criada com sucesso! ID: %d", conta.ID),
	})
}
