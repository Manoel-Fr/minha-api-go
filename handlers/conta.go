package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"minha-api-go/models"
)

type API struct {
	DB *sql.DB
}

func respondError(w http.ResponseWriter, status int, mensagem string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(models.Resposta{
		Status:   "erro",
		Mensagem: mensagem,
	})
}

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func (api *API) ImportarContasSale(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method != "POST" {
		respondError(w, http.StatusMethodNotAllowed, "Método não permitido. Use POST")
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
		respondError(w, http.StatusBadRequest, "Content-Type deve ser application/json")
		return
	}

	var contasSale []models.Conta
	if err := json.NewDecoder(r.Body).Decode(&contasSale); err != nil {
		respondError(w, http.StatusBadRequest, "Erro ao processar JSON das contas do Sale")
		return
	}
	if len(contasSale) == 0 {
		respondError(w, http.StatusBadRequest, "Nenhuma conta para importar")
		return
	}

	tx, err := api.DB.Begin()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Erro ao iniciar transação")
		return
	}

	stmt, err := tx.Prepare(`INSERT INTO Accounts (Name, Phone, Website, Industry) VALUES ($1, $2, $3, $4)`)
	if err != nil {
		tx.Rollback()
		respondError(w, http.StatusInternalServerError, "Erro ao preparar inserção")
		return
	}
	defer stmt.Close()

	contasImportadas := 0
	var erros []string

	for _, conta := range contasSale {
		if conta.Name == "" {
			erros = append(erros, "Nome é obrigatório")
			continue
		}
		if _, err := stmt.Exec(conta.Name, conta.Phone, conta.Website, conta.Industry); err != nil {
			erros = append(erros, err.Error())
			continue
		}
		contasImportadas++
	}

	if len(erros) > 0 {
		tx.Rollback()
		respondJSON(w, http.StatusBadRequest, models.Resposta{
			Status:   "erro",
			Mensagem: "Erros na importação",
		})
		return
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		respondError(w, http.StatusInternalServerError, "Erro ao finalizar importação")
		return
	}

	respondJSON(w, http.StatusCreated, models.Resposta{
		Status:   "sucesso",
		Mensagem: "Importação concluída com sucesso!",
	})
}
