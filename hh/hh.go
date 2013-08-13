package main

import (
	"fmt"
	"math"
	//"github.com/sbinet/go-hdf5/pkg/hdf5"
)

//calculate the potassium rate constant
func potassiumAlphaN(voltage float64) float64 {
	if voltage != 10 {
		return 0.01 * (-voltage + 10) / (math.Exp((-voltage+10)/10) - 1)
	} else {
		return 0.1
	}
}

//calculate the other potassium rate constant
func potassiumBetaN(voltage float64) float64 {
	return 0.125 * math.Exp(-voltage/80)
}

//calculate potassium steady state activation value
func potassiumNInfinity(voltage float64) float64 {
	return potassiumAlphaN(voltage) / (potassiumAlphaN(voltage) + potassiumBetaN(voltage))
}

//calculate first sodium rate constant
func sodiumAlphaM(voltage float64) float64 {
	if voltage != 25 {
		return 0.1 * (-voltage + 25) / (math.Exp((-voltage+25)/10) - 1)
	} else {
		return 1
	}
}

//calculate second sodium rate constant
func sodiumBetaM(voltage float64) float64 {
	return 4 * math.Exp(-voltage/18)
}

//calculate sodium steady state activation value
func sodiumMInfinity(voltage float64) float64 {
	return sodiumAlphaM(voltage) / (sodiumAlphaM(voltage) + sodiumBetaM(voltage))
}

//calculate first sodium inactivation rate constant
func sodiumAlphaH(voltage float64) float64 {
	return 0.07 * math.Exp(-voltage/20)
}

//calculate second sodium inactivation rate constant
func sodiumBetaH(voltage float64) float64 {
	return 1 / (math.Exp((-voltage+30)/10) + 1)
}

//calculate sodium steady state inactivation value
func sodiumHInfinity(voltage float64) float64 {
	return sodiumAlphaH(voltage) / (sodiumAlphaH(voltage) + sodiumBetaH(voltage))
}

func main() {

	//Simulation Parameters
	totalSimulationTime := float64(55) //Total simulation time in milliseconds
	deltaTime := float64(0.025)        //Simulation timestep in milliseconds

	var timeArray []float64 //the array to hold timesteps
	timeArray = make([]float64, int(totalSimulationTime/deltaTime), int(totalSimulationTime/deltaTime))
	for timestep := float64(0); timestep < (totalSimulationTime+deltaTime)/deltaTime; timestep++ {
		timeArray = append(timeArray, timestep*deltaTime) //fill the array with a number of timesteps
	}

	//Hodgkin Huxley model parameters
	restVoltage := float64(0)                      //V_rest
	lipidBilayerCapacitance := float64(1)          //C_m
	sodiumActivationMaxConductance := float64(120) //g-_Na
	potassiumMaxConductance := float64(36)         //g-_K
	leakConductance := float64(0.3)                //g-_l
	sodiumVoltagesSource := float64(115)           //E_Na
	potassiumVoltageSource := float64(-12)         //E_K
	leakReversePotential := float64(10.613)        //E_l
	//declare and calculate activation and inactivation dimensionless quantities
	m := sodiumMInfinity(restVoltage)
	n := potassiumNInfinity(restVoltage)
	h := sodiumHInfinity(restVoltage)
	//make an array of the calculated voltages
	membraneRestingPotentialDifferentials := make([]float64, len(timeArray), len(timeArray))
	//set the first timestep equal to the rest voltage
	membraneRestingPotentialDifferentials[0] = restVoltage
	//make an an array of stimuli over time, replace this later with other neurons in the network
	stimulusValues := make([]float64, len(timeArray), len(timeArray))
	fmt.Println("Array length is", len(timeArray))
	for time, currentTime := range timeArray {
		if currentTime >= 10 && currentTime <= 20 {
			stimulusValues[time] = float64(20)
		}
	}

}
