package clima

import (
	"errors"

	"github.com/ThalesSouza1998/clima-por-cep/internal/viacep"
)

var (
	ErrCEPinvalido      = errors.New("cep inválido")
	ErrCEPNaoEncontrado = errors.New("cep não encontrado")
)

// BuscaCidade resolve o município a partir do CEP.
type BuscaCidade interface {
	Cidade(cep string) (string, error)
}

// Termometro informa a temperatura atual de uma cidade, em Celsius.
type Termometro interface {
	TemperaturaAtual(cidade string) (float64, error)
}

type Servico struct {
	cidades    BuscaCidade
	termometro Termometro
}

func NovoServico(cidades BuscaCidade, termometro Termometro) *Servico {
	return &Servico{cidades: cidades, termometro: termometro}
}

func (s *Servico) PorCEP(cep string) (Clima, error) {
	if !cepValido(cep) {
		return Clima{}, ErrCEPinvalido
	}

	cidade, err := s.cidades.Cidade(cep)
	if err != nil {
		if errors.Is(err, viacep.ErrCEPDesconhecido) {
			return Clima{}, ErrCEPNaoEncontrado
		}
		return Clima{}, err
	}

	celsius, err := s.termometro.TemperaturaAtual(cidade)
	if err != nil {
		return Clima{}, err
	}

	return Clima{
		Celsius:    celsius,
		Fahrenheit: paraFahrenheit(celsius),
		Kelvin:     paraKelvin(celsius),
	}, nil
}

func cepValido(cep string) bool {
	if len(cep) != 8 {
		return false
	}
	for i := range len(cep) {
		if cep[i] < '0' || cep[i] > '9' {
			return false
		}
	}
	return true
}
