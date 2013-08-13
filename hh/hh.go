package main

import (
	// stdlib
	"fmt"
	// local
	"math"
	//"github.com/sbinet/go-hdf5/pkg/hdf5"
)

func potassiumAlphaN(voltage float64) float64 {
	if voltage != 10 {
		return 0.01 * (-voltage + 10) / (math.Exp((-voltage+10)/10) - 1)
	} else {
		return 0.1
	}

}

func potassiumBetaN(voltage float64) float64 {
	return 0.125 * math.Exp(-voltage/80)
}

func potassiumNInf(voltage float64) float64 {
	return potassiumAlphaN(voltage) / (potassiumAlphaN(voltage) + potassiumBetaN(voltage))
}

func sodiumAlphaN(voltage float64) float64 {
	if voltage != 25 {
		return 0.1 * (-voltage + 25) / (math.Exp((-voltage+25)/10) - 1)
	} else {
		return 1
	}
}

func sodiumBetaN(voltage float64) float64 {

}

func main() {

	voltage := float64(5.0)
	voltage = potassiumAlphaN(voltage)
	fmt.Println(voltage)
}
