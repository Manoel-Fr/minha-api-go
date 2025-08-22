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

	http.HandleFunc("/contas", func(w http.ResponseWriter, r *http.Request) {
		println("Rota /contas acessada com método:", r.Method)
		switch r.Method {
		case http.MethodGet:
			apiInstance.ListarContas(w, r)
		case http.MethodPost:
			apiInstance.CriarConta(w, r)
		default:
			http.Error(w, "Método não permitido. Use GET para listar ou POST para criar", http.StatusMethodNotAllowed)
			return
		}
	})

	porta := ":8081"
	println("Servidor iniciado em http://localhost" + porta)
	if err := http.ListenAndServe(porta, nil); err != nil {
		panic(err)
	}
}

// lsof -i :8081
// kill -9 <18344>
