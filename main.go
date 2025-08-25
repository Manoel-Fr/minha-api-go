package main

import (
	"log"
	"net/http"
	"os"

	"minha-api-go/database"
	"minha-api-go/handlers"
)

func main() {
	db, err := database.Conectar()
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}
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

	// rota de saúde (útil para Render)
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	// pega porta do ambiente do Render
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081" // fallback local
	}

	log.Println("Servidor iniciado na porta " + port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}

// lsof -i :8081
// kill -9 87388
