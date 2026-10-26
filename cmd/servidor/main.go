package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ThalesSouza1998/clima-por-cep/internal/api"
	"github.com/ThalesSouza1998/clima-por-cep/internal/clima"
	"github.com/ThalesSouza1998/clima-por-cep/internal/viacep"
	"github.com/ThalesSouza1998/clima-por-cep/internal/weatherapi"
	"github.com/joho/godotenv"
)

func main() {
	// O .env é só um atalho para rodar localmente; no Cloud Run as variáveis
	// chegam pelo ambiente e a ausência do arquivo não é problema.
	_ = godotenv.Load()

	chave := os.Getenv("WEATHER_API_KEY")
	if chave == "" {
		log.Fatal("WEATHER_API_KEY não definida")
	}

	servico := clima.NovoServico(viacep.Novo(), weatherapi.Novo(chave))
	handler := api.Novo(servico)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", handler.Clima)

	porta := os.Getenv("PORT")
	if porta == "" {
		porta = "8080"
	}

	log.Printf("escutando na porta %s", porta)
	if err := http.ListenAndServe(":"+porta, mux); err != nil {
		log.Fatal(err)
	}
}
