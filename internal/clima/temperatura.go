package clima

import "math"

// O contrato do desafio pede o Fahrenheit com uma casa decimal e o Kelvin
// com duas, então arredondamos logo na conversão.
func paraFahrenheit(celsius float64) float64 {
	return math.Round((celsius*1.8+32)*10) / 10
}

func paraKelvin(celsius float64) float64 {
	return math.Round((celsius+273.15)*100) / 100
}
