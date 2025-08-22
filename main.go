package main

import (
	"net/http"

	"minha-api-go/database"
	"minha-api-go/handlers"
)

func main() {
	db := database.Conectar()
	defer db.Close()

	apiInstance := &handlers.API{DB: db}

	http.HandleFunc("/contas/importar-sale", func(w http.ResponseWriter, r *http.Request) {
		println("Rota /contas/importar-sale acessada com método:", r.Method)
		if r.Method != http.MethodPost {
			http.Error(w, "Método não permitido. Use POST para importar contas do Sale", http.StatusMethodNotAllowed)
			return
		}
		apiInstance.ImportarContasSale(w, r)
	})

	porta := ":8081"
	println("Servidor iniciado em http://localhost" + porta)
	if err := http.ListenAndServe(porta, nil); err != nil {
		panic(err)
	}
}

// lsof -i :8081
// kill -9 30368
