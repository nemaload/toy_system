package main

import (
	"fmt"
	"math"
	//"github.com/sbinet/go-hdf5/pkg/hdf5"
)

type NeuronParameters struct {
	//these declare constant values
	restVoltage                    float64 //V_rest
	lipidBilayerCapacitance        float64 //C_m
	sodiumActivationMaxConductance float64 //g-_Na
	potassiumMaxConductance        float64 //g-_K
	leakConductance                float64 //g-_l
	sodiumReversePotential         float64 //E_Na
	potassiumReversePotential      float64 //E_K
	leakReversePotential           float64 //E_l
}

func (np NeuronParameters) initializeParametersWithDefaults() {
	np.restVoltage = 0
	np.lipidBilayerCapacitance = 1
	np.sodiumActivationMaxConductance = 120
	np.potassiumMaxConductance = 36
	np.leakConductance = 0.3
	//sodium conductance and potassium conductance are set during the simulation
	np.sodiumReversePotential = 115
	np.potassiumReversePotential = -12
	np.leakReversePotential = 10.613
}

type Neuron struct {
	parameters           NeuronParameters
	sodiumConductance    float64
	potassiumConductance float64
	//activation and inactivation dimensionless quantities
	m               float64
	n               float64
	h               float64
	stimulation     []float64
	V_m             []float64
	currentTimeStep float64
	simulation      *Simulation
}

func (neuron Neuron) initializeNeuron(simulation Simulation) {
	neuron.intializeVoltageArray(simulation)
	neuron.setSimulation(simulation)
	neuron.initializeDimensionlessQuantities()

}

func (neuron Neuron) setSampleStimulationValues() {
	neuron.stimulation = make([]float64, len(neuron.simulation.timeArray))
	for time, currentTime := range neuron.simulation.timeArray {
		if currentTime >= 5 && currentTime <= 30 {
			neuron.stimulation[time] = float64(10)
		}
	}
}

func (neuron Neuron) setSimulation(simulation Simulation) {
	neuron.simulation = &simulation
}

func (neuron Neuron) intializeVoltageArray(simulation Simulation) {
	neuron.V_m = make([]float64, len(simulation.timeArray))
	neuron.V_m[0] = neuron.parameters.restVoltage
}

func (neuron Neuron) initializeDimensionlessQuantities() {
	neuron.m = neuron.sodiumMInfinity()
	neuron.n = neuron.potassiumNInfinity()
	neuron.h = neuron.sodiumHInfinity()
}
func (neuron Neuron) calculateDimensionlessQuantities() {
	neuron.m += (sodiumAlphaM(neuron.V_m[int(neuron.currentTimeStep)-1])*(1-neuron.m) -
		sodiumBetaM(neuron.V_m[int(neuron.currentTimeStep)-1])*neuron.m) * neuron.simulation.deltaTime
	neuron.h += (sodiumAlphaH(neuron.V_m[timeStep-1])*(1-neuron.h) -
		sodiumBetaH(neuron.V_m[int(neuron.currentTimeStep)-1])*neuron.h) * neuron.simulation.deltaTime
	neuron.n += (potassiumAlphaN(neuron.V_m[timeStep-1])*(1-neuron.n) -
		potassiumBetaN(neuron.V_m[int(neuron.currentTimeStep)-1])*neuron.n) * neuron.simulation.deltaTime

}
func (neuron Neuron) calculateSimulationStep() {
	//potassium and sodium calculation step
	neuron.calculatePotassiumAndSodiumConductance()
	neuron.calculateDimensionlessQuantities()
	neuron.calculateNewVoltage()
}

func (n Neuron) calculateNewVoltage() {
	n.V_m[n.currentTimeStep] = n.V_m[n.currentTimeStep-1]
	n.V_m[n.currentTimeStep] += (n.stimulation[n.currentTimeStep-1] - n.sodiumConductance*
		(n.V_m[n.currentTimeStep-1]-n.parameters.sodiumReversePotential) - n.potassiumConductance*
		(n.V_m[n.currentTimeStep-1]-n.parameters.potassiumReversePotential) - n.parameters.leakConductance*
		(n.V_m[n.currentTimeStep-1]-n.parameters.leakReversePotential)) / n.parameters.lipidBilayerCapacitance * neuron.simulation.deltaTime
}

func (neuron Neuron) calculatePotassiumAndSodiumConductance() {
	neuron.sodiumConductance = neuron.parameters.sodiumActivationMaxConductance * neuron.h * math.Pow(neuron.m, 3.)
	neuron.potassiumConductance = neuron.parameters.potassiumMaxConductance * math.Pow(neuron.n, 4.)
}

func (neuron Neuron) currentVoltage() {
	return neuron.V_m[neuron.currentTimeStep]
}

func (neuron Neuron) potassiumAlphaN() {
	if neuron.currentVoltage() != 10 {
		return 0.01 * (-neuron.currentVoltage() + 10) / (math.Exp((-neuron.currentVoltage()+10)/10) - 1)
	}
	return 0.1
}

func (neuron Neuron) potassiumBetaN() {
	return 0.125 * math.Exp(-neuron.currentVoltage()/80)
}

func (neuron Neuron) potassiumNInfinity() {
	return neuron.potassiumAlphaN() /
		(neuron.potassiumAlphaN() + neuron.potassiumBetaN())
}

func (neuron Neuron) sodiumAlphaM() {
	if neuron.currentVoltage() != 25 {
		return 0.1 * (-neuron.currentVoltage() + 25) / (math.Exp((-neuron.currentVoltage()+25)/10) - 1)
	}
	return 1
}

func (neuron Neuron) sodiumBetaM() {
	return 4 * math.Exp(-neuron.currentVoltage()/18)
}

func (neuron Neuron) sodiumMInfinity() {
	return neuron.sodiumAlphaM() /
		(neuron.sodiumAlphaM() + neuron.sodiumBetaM())
}

func (neuron Neuron) sodiumAlphaH() {
	return 0.07 * math.Exp(-neuron.currentVoltage()/20)
}

func (neuron Neuron) sodiumBetaH() {
	return 1 / (math.Exp((-neuron.currentVoltage()+30)/10) + 1)
}

func (neuron Neuron) sodiumHInfinity() {
	return neuron.sodiumAlphaH() /
		(neuron.sodiumAlphaH() + neuron.sodiumBetaH())
}

type Simulation struct {
	totalSimulationTime float64
	deltaTime           float64
	timeArray           []float64
	weightMap           map[*Neuron]map[*Neuron]float64
	neuronArray         []*Neuron
}

func (simulation Simulation) addNeuronToSimulation(neuron Neuron) {
	simulation.neuronArray = append(simulation.neuronArray, &neuron)
}

func (simulation Simulation) initializeWeightMap() {
	for neuron1 := range simulation.weightMap {
		for neuron2 := range simulation.weightMap[neuron1] {
			simulation.weightMap[neuron1][neuron2] = 0.0
		}
	}
}

func (simulation Simulation) setSynapseWeightPair(neuron1, neuron2 Neuron, weight float64) {
	simulation.weightMap[neuron1][neuron2] = weight
	simulation.weightMap[neuron2][neuron1] = weight
}

func (simulation Simulation) initializeSimulation(totalSimulationTime float64, deltaTime float64) {
	simulation.totalSimulationTime = totalSimulationTime
	simulation.deltaTime = deltaTime
	simulation.initializeTimeArray()

}

func (simulation Simulation) initializeTimeArray() {
	for timestep := float64(0); timestep < simulation.totalSimulationTime+simulation.deltaTime; timestep += simulation.deltaTime {
		simulation.timeArray = append(simulation.timeArray, timestep)
	}
}

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
	for timestep := float64(0); timestep <
		totalSimulationTime+deltaTime; timestep += deltaTime {
		timeArray = append(timeArray, timestep)
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
	V_m := make([]float64, len(timeArray))
	//set the first timestep equal to the rest voltage
	V_m[0] = restVoltage
	//make an an array of stimuli over time 
	//replace this later with other neurons in the network
	stimulusValues := make([]float64, len(timeArray))
	for time, currentTime := range timeArray {
		//arbitrary stimulation times
		if currentTime >= 5 && currentTime <= 30 {
			stimulusValues[time] = float64(10) //some arbitrary stimulus value
		}
	}
	for timeStep := 1; timeStep < len(timeArray); timeStep++ {

		sodiumConductance = sodiumActivationMaxConductance * h * math.Pow(m, 3.)
		potassiumConductance = potassiumMaxConductance * math.Pow(n, 4.)

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
		//update it with the model equation
		V_m[timeStep] += (stimulusValues[timeStep-1] - sodiumConductance*
			(V_m[timeStep-1]-sodiumReversePotential) - potassiumConductance*
			(V_m[timeStep-1]-potassiumReversePotential) - leakConductance*
			(V_m[timeStep-1]-leakReversePotential)) / lipidBilayerCapacitance * deltaTime
		fmt.Print(V_m[timeStep], ",")

	}

}
