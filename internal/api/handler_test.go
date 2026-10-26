package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/ThalesSouza1998/clima-por-cep/internal/clima"
)

type previsorFalso struct {
	previsao clima.Clima
	err      error
}

func (p previsorFalso) PorCEP(string) (clima.Clima, error) {
	return p.previsao, p.err
}

func TestHandlerClima(t *testing.T) {
	casos := []struct {
		nome       string
		url        string
		previsor   previsorFalso
		querStatus int
		querCorpo  map[string]any
	}{
		{
			nome:       "sucesso",
			url:        "/?cep=01001000",
			previsor:   previsorFalso{previsao: clima.Clima{Celsius: 28.5, Fahrenheit: 83.3, Kelvin: 301.65}},
			querStatus: http.StatusOK,
			querCorpo:  map[string]any{"temp_C": 28.5, "temp_F": 83.3, "temp_K": 301.65},
		},
		{
			nome:       "cep inválido",
			url:        "/?cep=123",
			previsor:   previsorFalso{err: clima.ErrCEPinvalido},
			querStatus: http.StatusUnprocessableEntity,
			querCorpo:  map[string]any{"error": "invalid zipcode"},
		},
		{
			nome:       "cep não encontrado",
			url:        "/?cep=99999999",
			previsor:   previsorFalso{err: clima.ErrCEPNaoEncontrado},
			querStatus: http.StatusNotFound,
			querCorpo:  map[string]any{"error": "can not find zipcode"},
		},
		{
			nome:       "erro inesperado",
			url:        "/?cep=01001000",
			previsor:   previsorFalso{err: errors.New("falha inesperada")},
			querStatus: http.StatusInternalServerError,
			querCorpo:  map[string]any{"error": "internal error"},
		},
	}

	for _, caso := range casos {
		t.Run(caso.nome, func(t *testing.T) {
			requisicao := httptest.NewRequest(http.MethodGet, caso.url, nil)
			gravador := httptest.NewRecorder()

			Novo(caso.previsor).Clima(gravador, requisicao)

			if gravador.Code != caso.querStatus {
				t.Fatalf("status = %d; queria %d", gravador.Code, caso.querStatus)
			}

			if tipo := gravador.Header().Get("Content-Type"); tipo != "application/json" {
				t.Errorf("Content-Type = %q; queria application/json", tipo)
			}

			var corpo map[string]any
			if err := json.Unmarshal(gravador.Body.Bytes(), &corpo); err != nil {
				t.Fatalf("resposta não é um JSON válido: %v", err)
			}

			if !reflect.DeepEqual(corpo, caso.querCorpo) {
				t.Errorf("corpo = %v; queria %v", corpo, caso.querCorpo)
			}
		})
	}
}
