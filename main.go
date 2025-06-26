package main

import (
	"log"
	"net/http"

	"minha-api-go/database"
	"minha-api-go/handlers"
)

func main() {
	db := database.Conectar()
	defer db.Close()

	apiInstance := &handlers.API{DB: db}

	http.HandleFunc("/conta", apiInstance.CriarConta)

	log.Println("API rodando em http://localhost:8888")
	log.Fatal(http.ListenAndServe(":8888", nil))
}
