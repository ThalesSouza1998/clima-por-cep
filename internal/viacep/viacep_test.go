package viacep

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCidade(t *testing.T) {
	servidor := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/ws/01001000/json/" {
			t.Errorf("caminho = %q; queria /ws/01001000/json/", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"localidade":"São Paulo"}`))
	}))
	defer servidor.Close()

	cliente := &Cliente{cliente: servidor.Client(), enderecoBase: servidor.URL}

	cidade, err := cliente.Cidade("01001000")
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if cidade != "São Paulo" {
		t.Errorf("cidade = %q; queria São Paulo", cidade)
	}
}

func TestCidadeInexistente(t *testing.T) {
	servidor := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"erro":true}`))
	}))
	defer servidor.Close()

	cliente := &Cliente{cliente: servidor.Client(), enderecoBase: servidor.URL}

	_, err := cliente.Cidade("99999999")
	if !errors.Is(err, ErrCEPDesconhecido) {
		t.Fatalf("erro = %v; queria %v", err, ErrCEPDesconhecido)
	}
}

func TestCidadeRespostaNaoOK(t *testing.T) {
	servidor := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "indisponível", http.StatusInternalServerError)
	}))
	defer servidor.Close()

	cliente := &Cliente{cliente: servidor.Client(), enderecoBase: servidor.URL}

	if _, err := cliente.Cidade("01001000"); err == nil {
		t.Fatal("esperava um erro para resposta 500")
	}
}
