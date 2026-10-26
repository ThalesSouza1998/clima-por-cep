// Package clima junta as consultas de CEP e de temperatura e devolve o
// resultado já convertido para as três escalas usadas na resposta.
package clima

// Clima é o corpo devolvido pela API.
type Clima struct {
	Celsius    float64 `json:"temp_C"`
	Fahrenheit float64 `json:"temp_F"`
	Kelvin     float64 `json:"temp_K"`
}
