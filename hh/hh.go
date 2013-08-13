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
	}
	return 0.1
}

//calculate the other potassium rate constant
func potassiumBetaN(voltage float64) float64 {
	return 0.125 * math.Exp(-voltage/80)
}

//calculate potassium steady state activation value
func potassiumNInfinity(voltage float64) float64 {
	return potassiumAlphaN(voltage) /
		(potassiumAlphaN(voltage) + potassiumBetaN(voltage))
}

//calculate first sodium rate constant
func sodiumAlphaM(voltage float64) float64 {
	if voltage != 25 {
		return 0.1 * (-voltage + 25) / (math.Exp((-voltage+25)/10) - 1)
	}
	return 1
}

//calculate second sodium rate constant
func sodiumBetaM(voltage float64) float64 {
	return 4 * math.Exp(-voltage/18)
}

//calculate sodium steady state activation value
func sodiumMInfinity(voltage float64) float64 {
	return sodiumAlphaM(voltage) /
		(sodiumAlphaM(voltage) + sodiumBetaM(voltage))
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
	return sodiumAlphaH(voltage) /
		(sodiumAlphaH(voltage) + sodiumBetaH(voltage))
}

func main() {

	//Simulation Parameters
	totalSimulationTime := float64(55) //Total simulation time in milliseconds
	deltaTime := float64(0.025)        //Simulation timestep in milliseconds

	var timeArray []float64 //the array to hold timesteps
	timeArray = make([]float64, int(totalSimulationTime/deltaTime),
		int(totalSimulationTime/deltaTime))
	for timestep := float64(0); timestep <
		(totalSimulationTime+deltaTime)/deltaTime; timestep++ {
		//fill the array with a number of timesteps
		timeArray = append(timeArray, timestep*deltaTime)
	}

	//Hodgkin Huxley model parameters
	restVoltage := float64(0)                      //V_rest
	lipidBilayerCapacitance := float64(1)          //C_m
	sodiumActivationMaxConductance := float64(120) //g-_Na
	potassiumMaxConductance := float64(36)         //g-_K
	leakConductance := float64(0.3)                //g-_l
	//These values aren't used until the simulation
	var sodiumConductance, potassiumConductance float64 //g_Na and g_K
	sodiumReversePotential := float64(115)              //E_Na
	potassiumReversePotential := float64(-12)           //E_K
	leakReversePotential := float64(10.613)             //E_l
	//declare and calculate activation and inactivation dimensionless quantities
	m := sodiumMInfinity(restVoltage)
	n := potassiumNInfinity(restVoltage)
	h := sodiumHInfinity(restVoltage)
	//make an array of the calculated voltages
	V_m := make([]float64,
		len(timeArray), len(timeArray))
	//set the first timestep equal to the rest voltage
	V_m[0] = restVoltage
	//make an an array of stimuli over time 
	//replace this later with other neurons in the network
	stimulusValues := make([]float64, len(timeArray), len(timeArray))
	for time, currentTime := range timeArray {
		//arbitrary stimulation times
		if currentTime >= 10 && currentTime <= 20 {
			stimulusValues[time] = float64(20) //some arbitrary stimulus value
		}
	}
	//the main simulation loop
	for timeStep := range timeArray {
		if timeStep == 0 {
			//Skip the very first step of simulation
			continue
		}

		sodiumConductance = sodiumActivationMaxConductance * h * math.Pow(m, 3)
		potassiumConductance = potassiumMaxConductance * math.Pow(n, 4)

		//Update the activation/inactivation dimensionless quantities
		m += (sodiumAlphaM(V_m[timeStep-1])*(1-m) -
			sodiumBetaM(V_m[timeStep-1])*m) * deltaTime
		h += (sodiumAlphaH(V_m[timeStep-1])*(1-h) -
			sodiumBetaH(V_m[timeStep-1])*h) * deltaTime
		n += (potassiumAlphaN(V_m[timeStep-1])*(1-n) -
			potassiumBetaN(V_m[timeStep-1])*n) * deltaTime

		//Calculate the new membrane potential
		//first, set the voltage to the old voltage
		V_m[timeStep] = V_m[timeStep-1]
		//then, update it with the model equation
		V_m[timeStep] += (stimulusValues[timeStep-1] - sodiumConductance*
			(V_m[timeStep-1]-sodiumReversePotential) - potassiumConductance*
			(V_m[timeStep-1]-potassiumReversePotential) - leakConductance*
			V_m[timeStep-1] - leakReversePotential) / lipidBilayerCapacitance
		V_m[timeStep] *= deltaTime
		fmt.Print(V_m[timeStep], ",")

	}

}
