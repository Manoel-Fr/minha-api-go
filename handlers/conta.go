package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"minha-api-go/models"
)

// API estrutura que contém a conexão com o banco de dados
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

// ImportarContasSale importa contas do Sale para o banco de dados
func (api *API) ImportarContasSale(w http.ResponseWriter, r *http.Request) {
	// Configurar headers CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Content-Type", "application/json")

	// Log da requisição
	fmt.Printf("Método da requisição: %s\n", r.Method)
	fmt.Printf("Headers da requisição: %v\n", r.Header)

	// Tratar requisição OPTIONS (preflight)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Verificar método
	if r.Method != "POST" {
		respondError(w, http.StatusMethodNotAllowed, "Método não permitido. Use POST")
		return
	}

	// Verificar Content-Type
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		respondError(w, http.StatusBadRequest, "Content-Type deve ser application/json")
		return
	}

	// Estrutura para receber múltiplas contas
	var contasSale []models.Conta
	if err := json.NewDecoder(r.Body).Decode(&contasSale); err != nil {
		respondError(w, http.StatusBadRequest, "Erro ao processar JSON das contas do Sale")
		return
	}

	if len(contasSale) == 0 {
		respondError(w, http.StatusBadRequest, "Nenhuma conta para importar")
		return
	}

	// Iniciar uma transação para garantir a consistência dos dados
	tx, err := api.DB.Begin()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Erro ao iniciar transação")
		return
	}

	// Preparar a query de inserção
	insertSQL := `INSERT INTO Accounts (Name, Phone, Website, Industry) VALUES ($1, $2, $3, $4)`
	stmt, err := tx.Prepare(insertSQL)
	if err != nil {
		tx.Rollback()
		respondError(w, http.StatusInternalServerError, "Erro ao preparar inserção")
		return
	}
	defer stmt.Close()

	contasImportadas := 0
	var erros []string

	// Processar cada conta
	for i, conta := range contasSale {
		if conta.Name == "" {
			erros = append(erros, fmt.Sprintf("Conta %d: Nome é obrigatório", i+1))
			continue
		}

		_, err := stmt.Exec(conta.Name, conta.Phone, conta.Website, conta.Industry)
		if err != nil {
			erros = append(erros, fmt.Sprintf("Erro ao importar conta %d: %v", i+1, err))
			continue
		}
		contasImportadas++
	}

	// Se houver erros, fazer rollback e retornar os erros
	if len(erros) > 0 {
		tx.Rollback()
		respondJSON(w, http.StatusBadRequest, models.Resposta{
			Status:   "erro",
			Mensagem: fmt.Sprintf("Erros na importação: %v", erros),
		})
		return
	}

	// Commit da transação
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		respondError(w, http.StatusInternalServerError, "Erro ao finalizar importação")
		return
	}

	// Resposta de sucesso
	respondJSON(w, http.StatusCreated, models.Resposta{
		Status:   "sucesso",
		Mensagem: fmt.Sprintf("Importação concluída com sucesso! %d contas importadas", contasImportadas),
	})
}
