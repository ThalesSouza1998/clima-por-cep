package clima

import (
	"errors"
	"testing"

	"github.com/ThalesSouza1998/clima-por-cep/internal/viacep"
)

type buscaFalsa struct {
	cidade string
	err    error
}

func (b buscaFalsa) Cidade(string) (string, error) {
	return b.cidade, b.err
}

type termometroFalso struct {
	celsius float64
	err     error
}

func (t termometroFalso) TemperaturaAtual(string) (float64, error) {
	return t.celsius, t.err
}

func TestServicoPorCEP(t *testing.T) {
	erroDeRede := errors.New("serviço fora do ar")

	casos := []struct {
		nome       string
		cep        string
		busca      buscaFalsa
		termometro termometroFalso
		querErro   error
		quer       Clima
	}{
		{
			nome:       "consulta normal",
			cep:        "01001000",
			busca:      buscaFalsa{cidade: "São Paulo"},
			termometro: termometroFalso{celsius: 28.5},
			quer:       Clima{Celsius: 28.5, Fahrenheit: 83.3, Kelvin: 301.65},
		},
		{
			nome:     "cep vazio",
			querErro: ErrCEPinvalido,
		},
		{
			nome:     "cep com letras",
			cep:      "0100100a",
			querErro: ErrCEPinvalido,
		},
		{
			nome:     "cep com hífen",
			cep:      "01001-000",
			querErro: ErrCEPinvalido,
		},
		{
			nome:     "cep inexistente",
			cep:      "99999999",
			busca:    buscaFalsa{err: viacep.ErrCEPDesconhecido},
			querErro: ErrCEPNaoEncontrado,
		},
		{
			nome:     "viacep fora do ar",
			cep:      "01001000",
			busca:    buscaFalsa{err: erroDeRede},
			querErro: erroDeRede,
		},
		{
			nome:       "weatherapi fora do ar",
			cep:        "01001000",
			busca:      buscaFalsa{cidade: "São Paulo"},
			termometro: termometroFalso{err: erroDeRede},
			querErro:   erroDeRede,
		},
	}

	for _, caso := range casos {
		t.Run(caso.nome, func(t *testing.T) {
			clima, err := NovoServico(caso.busca, caso.termometro).PorCEP(caso.cep)

			if !errors.Is(err, caso.querErro) {
				t.Fatalf("erro = %v; queria %v", err, caso.querErro)
			}

			if clima != caso.quer {
				t.Errorf("clima = %+v; queria %+v", clima, caso.quer)
			}
		})
	}
}
