package weatherapi

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestTemperaturaAtual(t *testing.T) {
	servidor := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/current.json" {
			t.Errorf("caminho = %q; queria /v1/current.json", r.URL.Path)
		}

		parametros := r.URL.Query()
		if parametros.Get("key") != "chave-de-teste" {
			t.Errorf("key = %q", parametros.Get("key"))
		}
		if parametros.Get("q") != "São Paulo" {
			t.Errorf("q = %q; queria São Paulo", parametros.Get("q"))
		}
		if parametros.Get("aqi") != "no" {
			t.Errorf("aqi = %q; queria no", parametros.Get("aqi"))
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"current":{"temp_c":28.5}}`))
	}))
	defer servidor.Close()

	cliente := Novo("chave-de-teste")
	cliente.cliente = servidor.Client()
	cliente.enderecoBase = servidor.URL

	temperatura, err := cliente.TemperaturaAtual("São Paulo")
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if temperatura != 28.5 {
		t.Errorf("temperatura = %v; queria 28.5", temperatura)
	}
}

func TestTemperaturaAtualRespostaNaoOK(t *testing.T) {
	servidor := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "indisponível", http.StatusInternalServerError)
	}))
	defer servidor.Close()

	cliente := Novo("chave-de-teste")
	cliente.cliente = servidor.Client()
	cliente.enderecoBase = servidor.URL

	if _, err := cliente.TemperaturaAtual("São Paulo"); err == nil {
		t.Fatal("esperava um erro para resposta 500")
	}
}

type transporteComErro struct{}

func (transporteComErro) RoundTrip(*http.Request) (*http.Response, error) {
	return nil, &url.Error{
		Op:  "Get",
		URL: "https://api.weatherapi.com/v1/current.json?key=segredo",
		Err: errors.New("falha de rede"),
	}
}

func TestErroDeRedeNaoVazaAChave(t *testing.T) {
	cliente := Novo("segredo")
	cliente.cliente = &http.Client{Transport: transporteComErro{}}

	_, err := cliente.TemperaturaAtual("São Paulo")
	if err == nil {
		t.Fatal("esperava um erro")
	}
	if strings.Contains(err.Error(), "segredo") {
		t.Errorf("o erro vazou a chave da API: %v", err)
	}
}
