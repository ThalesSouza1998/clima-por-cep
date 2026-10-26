// Package viacep consulta o ViaCEP para descobrir a cidade de um CEP.
package viacep

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

var ErrCEPDesconhecido = errors.New("cep não existe na base do ViaCEP")

type Cliente struct {
	cliente      *http.Client
	enderecoBase string
}

func Novo() *Cliente {
	return &Cliente{
		cliente:      &http.Client{Timeout: 10 * time.Second},
		enderecoBase: "https://viacep.com.br",
	}
}

// Quando o CEP não existe o ViaCEP devolve {"erro": true}, sem localidade.
type endereco struct {
	Localidade string `json:"localidade"`
}

func (c *Cliente) Cidade(cep string) (string, error) {
	resp, err := c.cliente.Get(c.enderecoBase + "/ws/" + cep + "/json/")
	if err != nil {
		return "", fmt.Errorf("consultando o ViaCEP: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ViaCEP respondeu %s", resp.Status)
	}

	var dados endereco
	if err := json.NewDecoder(resp.Body).Decode(&dados); err != nil {
		return "", fmt.Errorf("lendo a resposta do ViaCEP: %w", err)
	}

	if dados.Localidade == "" {
		return "", ErrCEPDesconhecido
	}

	return dados.Localidade, nil
}
