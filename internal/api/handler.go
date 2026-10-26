// Package api expõe o serviço de clima via HTTP.
package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/ThalesSouza1998/clima-por-cep/internal/clima"
)

// Previsor é tudo o que o handler precisa do serviço de clima.
type Previsor interface {
	PorCEP(cep string) (clima.Clima, error)
}

type Handler struct {
	previsor Previsor
}

func Novo(previsor Previsor) *Handler {
	return &Handler{previsor: previsor}
}

func (h *Handler) Clima(w http.ResponseWriter, r *http.Request) {
	previsao, err := h.previsor.PorCEP(r.URL.Query().Get("cep"))
	if err != nil {
		codigo, mensagem := respostaDeErro(err)
		escrever(w, codigo, map[string]string{"error": mensagem})
		return
	}

	escrever(w, http.StatusOK, previsao)
}

// As mensagens ficam em inglês porque o contrato do desafio pede assim.
func respostaDeErro(err error) (int, string) {
	switch {
	case errors.Is(err, clima.ErrCEPinvalido):
		return http.StatusUnprocessableEntity, "invalid zipcode"
	case errors.Is(err, clima.ErrCEPNaoEncontrado):
		return http.StatusNotFound, "can not find zipcode"
	default:
		log.Printf("erro inesperado ao consultar o clima: %v", err)
		return http.StatusInternalServerError, "internal error"
	}
}

func escrever(w http.ResponseWriter, codigo int, corpo any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(codigo)

	if err := json.NewEncoder(w).Encode(corpo); err != nil {
		log.Printf("erro ao escrever a resposta: %v", err)
	}
}
