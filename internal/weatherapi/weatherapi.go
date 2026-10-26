// Package weatherapi consulta a temperatura atual na WeatherAPI.
package weatherapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type Cliente struct {
	cliente      *http.Client
	chave        string
	enderecoBase string
}

func Novo(chave string) *Cliente {
	return &Cliente{
		cliente:      &http.Client{Timeout: 10 * time.Second},
		chave:        chave,
		enderecoBase: "https://api.weatherapi.com",
	}
}

type resposta struct {
	Current struct {
		Celsius float64 `json:"temp_c"`
	} `json:"current"`
}

func (c *Cliente) TemperaturaAtual(cidade string) (float64, error) {
	parametros := url.Values{}
	parametros.Set("key", c.chave)
	parametros.Set("q", cidade)
	parametros.Set("aqi", "no")

	resp, err := c.cliente.Get(c.enderecoBase + "/v1/current.json?" + parametros.Encode())
	if err != nil {
		return 0, fmt.Errorf("consultando a WeatherAPI: %w", semURL(err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("WeatherAPI respondeu %s", resp.Status)
	}

	var dados resposta
	if err := json.NewDecoder(resp.Body).Decode(&dados); err != nil {
		return 0, fmt.Errorf("lendo a resposta da WeatherAPI: %w", err)
	}

	return dados.Current.Celsius, nil
}

// semURL remove a URL do erro para não levar a chave da API para logs e traces.
func semURL(err error) error {
	for {
		var erroDeURL *url.Error
		if !errors.As(err, &erroDeURL) {
			return err
		}
		if erroDeURL.Err == nil {
			return errors.New("falha de rede")
		}
		err = erroDeURL.Err
	}
}
