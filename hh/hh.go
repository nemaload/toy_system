package main

import (
	"fmt"
	"math"
	//"github.com/sbinet/go-hdf5/pkg/hdf5"
)

//PURPOSE: Hold various neuronal parameters
type NeuronParameters struct {
	restVoltage                    float64 //V_rest
	lipidBilayerCapacitance        float64 //C_m
	sodiumActivationMaxConductance float64 //g-_Na
	potassiumMaxConductance        float64 //g-_K
	leakConductance                float64 //g-_l
	sodiumReversePotential         float64 //E_Na
	potassiumReversePotential      float64 //E_K
	leakReversePotential           float64 //E_l
}

//MODIFIES: np
//EFFECTS: Fills the neuronParameters struct with default values
func (np *NeuronParameters) initializeParametersWithDefaults() {
	np.restVoltage = 0
	np.lipidBilayerCapacitance = 1
	np.sodiumActivationMaxConductance = 120
	np.potassiumMaxConductance = 36
	np.leakConductance = 0.3
	np.sodiumReversePotential = 115
	np.potassiumReversePotential = -12
	np.leakReversePotential = 10.613
}

//PURPOSE: Hold various data about a single neuron and its state
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

//REQUIRES: neuron to be a properly initialized Neuron
//EFFECTS: Prints the neuron voltages of a single neuron to the screen in such a fashion that
//if stdout is redirected into a file, it is in CSV format
func (neuron *Neuron) printToCSV() {
	for i := range neuron.V_m {
		fmt.Print(neuron.V_m[i], ",")
	}
}

//MODIFIES: neuron's currentTimeStep
//EFFECTS: Instantiates a single neuron
func (neuron *Neuron) initializeNeuron(simulation *Simulation) {
	neuron.intializeVoltageArray(simulation)
	neuron.setSimulation(simulation)
	neuron.initializeDimensionlessQuantities()
	neuron.currentTimeStep = 1

}

//MODIFIES: neuron's stimulation array
//EFFECTS: Fills the stimulation array with a single 10mV stimulation period from 5ms to 30ms
func (neuron *Neuron) setSampleStimulationValues() {
	neuron.stimulation = make([]float64, len(neuron.simulation.timeArray))
	for time, currentTime := range neuron.simulation.timeArray {
		if currentTime >= 5 && currentTime <= 30 {
			neuron.stimulation[time] = float64(10)
		}
	}
}

//MODIFIES: neuron's simulation pointer
//EFFECTS: Sets the neuron.simulation pointer to the input simulation
func (neuron *Neuron) setSimulation(simulation *Simulation) {
	neuron.simulation = simulation
}

//MODIFIES: neuron's V_m array
//EFFECTS: Initializes neuron's voltage array, and sets the first element equal to the neuron's rest voltage
func (neuron *Neuron) intializeVoltageArray(simulation *Simulation) {

	neuron.V_m = make([]float64, len(simulation.timeArray))
	neuron.V_m[0] = neuron.parameters.restVoltage
}

//REQUIRES: neuron has the voltage array initialized and parameters set.
//MODIFIES: neuron's m, n, and h
//EFFECTS: Initializes the dimensionless quantities
func (neuron *Neuron) initializeDimensionlessQuantities() {
	neuron.m = neuron.sodiumMInfinity()
	neuron.n = neuron.potassiumNInfinity()
	neuron.h = neuron.sodiumHInfinity()
}

//REQUIRES: neuron is properly initialized
//MODIFIES: neuron's m, n, and h
//EFFECTS: Calculates the dimensionless quantities for the current neuron timestep
func (neuron *Neuron) calculateDimensionlessQuantities() {
	neuron.m += (sodiumAlphaM(neuron.V_m[int(neuron.currentTimeStep)-1])*(1-neuron.m) -
		sodiumBetaM(neuron.V_m[int(neuron.currentTimeStep)-1])*neuron.m) * neuron.simulation.deltaTime
	neuron.h += (sodiumAlphaH(neuron.V_m[int(neuron.currentTimeStep)-1])*(1-neuron.h) -
		sodiumBetaH(neuron.V_m[int(neuron.currentTimeStep)-1])*neuron.h) * neuron.simulation.deltaTime
	neuron.n += (potassiumAlphaN(neuron.V_m[int(neuron.currentTimeStep)-1])*(1-neuron.n) -
		potassiumBetaN(neuron.V_m[int(neuron.currentTimeStep)-1])*neuron.n) * neuron.simulation.deltaTime

}

//MODIFIES: neuron's currentTimeStep
//EFFECTS: Runs all update calculations for the neuron in its current timestep
func (neuron *Neuron) calculateSimulationStep() {
	//potassium and sodium calculation step
	neuron.calculatePotassiumAndSodiumConductance()
	neuron.calculateDimensionlessQuantities()
	neuron.calculateNewVoltage()
	neuron.currentTimeStep += 1
}

//MODIFIES: neuron's V_m array
//EFFECTS: Calculates the neuron's membrane potential at the current neuron timestep
func (n *Neuron) calculateNewVoltage() {
	n.V_m[int(n.currentTimeStep)] = n.V_m[int(n.currentTimeStep)-1]
	n.V_m[int(n.currentTimeStep)] += (n.stimulation[int(n.currentTimeStep)-1] - n.sodiumConductance*
		(n.V_m[int(n.currentTimeStep)-1]-n.parameters.sodiumReversePotential) - n.potassiumConductance*
		(n.V_m[int(n.currentTimeStep)-1]-n.parameters.potassiumReversePotential) - n.parameters.leakConductance*
		(n.V_m[int(n.currentTimeStep)-1]-n.parameters.leakReversePotential)) / n.parameters.lipidBilayerCapacitance * n.simulation.deltaTime
}

//MODIFIES: neuron's sodiumConductance and potassiumConductance
//EFFECTS: Calculates the neuron's potassium and sodium conductance and the current neuron timestep
func (neuron *Neuron) calculatePotassiumAndSodiumConductance() {
	neuron.sodiumConductance = neuron.parameters.sodiumActivationMaxConductance * neuron.h * math.Pow(neuron.m, 3.)
	neuron.potassiumConductance = neuron.parameters.potassiumMaxConductance * math.Pow(neuron.n, 4.)
}

//EFFECTS: Returns neuron's membrane potential at the current neuron timestep
func (neuron *Neuron) currentVoltage() float64 {
	return neuron.V_m[int(neuron.currentTimeStep)]
}

//EFFECTS: Calculates and returns the neuron's current first potassium rate constant value
func (neuron *Neuron) potassiumAlphaN() float64 {
	if neuron.currentVoltage() != 10 {
		return 0.01 * (-neuron.currentVoltage() + 10) / (math.Exp((-neuron.currentVoltage()+10)/10) - 1)
	}
	return 0.1
}

//EFFECTS: Calculates and returns the neuron's current second potassium rate constant value
func (neuron *Neuron) potassiumBetaN() float64 {
	return 0.125 * math.Exp(-neuron.currentVoltage()/80)
}

//EFFECTS: Calculates and returns the neuron's current steady state potassium activation value
func (neuron *Neuron) potassiumNInfinity() float64 {
	return neuron.potassiumAlphaN() /
		(neuron.potassiumAlphaN() + neuron.potassiumBetaN())
}

//EFFECTS: Calculates and returns the neuron's current first sodium rate constant value
func (neuron *Neuron) sodiumAlphaM() float64 {
	if neuron.currentVoltage() != 25 {
		return 0.1 * (-neuron.currentVoltage() + 25) / (math.Exp((-neuron.currentVoltage()+25)/10) - 1)
	}
	return 1
}

//EFFECTS: Calculates and returns the neuron's current second sodium rate constant value
func (neuron *Neuron) sodiumBetaM() float64 {
	return 4 * math.Exp(-neuron.currentVoltage()/18)
}

//EFFECTS: Calculates and returns the neuron's current sodium steady state activation value
func (neuron *Neuron) sodiumMInfinity() float64 {
	return neuron.sodiumAlphaM() /
		(neuron.sodiumAlphaM() + neuron.sodiumBetaM())
}

//EFFECTS: Calculates and returns the neuron's current first sodium inactivation rate constant value
func (neuron *Neuron) sodiumAlphaH() float64 {
	return 0.07 * math.Exp(-neuron.currentVoltage()/20)
}

//EFFECTS: Calculates and returns the neuron's current second sodium inactivation rate constant value
func (neuron *Neuron) sodiumBetaH() float64 {
	return 1 / (math.Exp((-neuron.currentVoltage()+30)/10) + 1)
}

//EFFECTS: Calculates and returns the neuron's current sodium steady state inactivation value
func (neuron *Neuron) sodiumHInfinity() float64 {
	return neuron.sodiumAlphaH() /
		(neuron.sodiumAlphaH() + neuron.sodiumBetaH())
}

//PURPOSE: Hold values and states related to the neural network simulation
type Simulation struct {
	totalSimulationTime float64
	deltaTime           float64
	timeArray           []float64
	weightMap           map[*Neuron]map[*Neuron]float64
	neuronArray         []*Neuron
}

//MODIFIES: Neurons in neuronArray
//EFFECTS: Initializes all neurons in the neuron array with params
func (simulation *Simulation) initializeNeuronArray(params NeuronParameters) {
	for _, neuron := range simulation.neuronArray {
		neuron.initializeNeuron(simulation)
		neuron.parameters = params
		neuron.setSampleStimulationValues()
	}
}

//MODIFIES: Simulation's neuronArray
//EFFECTS: Allocates new neurons, and adds them to the simulation's neuronArray
func (simulation *Simulation) addNumberofNeuronsToSimulation(neuronCount int) {
	for i := 0; i < neuronCount; i++ {
		simulation.neuronArray = append(simulation.neuronArray, new(Neuron))
	}
}

//MODIFIES: Simulation's weightMap
//EFFECTS: Sets all values in the weightMap equal to zero
//TODO: Programatically create the array using the neurons contained within the simulation
func (simulation *Simulation) initializeWeightMap() {
	for neuron1 := range simulation.weightMap {
		for neuron2 := range simulation.weightMap[neuron1] {
			simulation.weightMap[neuron1][neuron2] = 0.0
		}
	}
}

//MODIFIES: simulation's weightMap
//EFFECTS: Sets the connection weight from neuron1 to neuron2 as weight
func (simulation *Simulation) setSynapseWeightPair(neuron1, neuron2 Neuron, weight float64) {
	simulation.weightMap[&neuron1][&neuron2] = weight
}

//MODIFIES: Simulation's totalSimulationTime, deltaTime
//EFFECTS: Intializes the simulation with totalSimulationTime and deltaTime
func (simulation *Simulation) initializeSimulation(totalSimulationTime float64, deltaTime float64) {
	simulation.totalSimulationTime = totalSimulationTime
	simulation.deltaTime = deltaTime
	simulation.initializeTimeArray()

}

//MODIFIES: simulation's timeArray
//EFFECTS: Fills timeArray with values from 0 to totalSimulationTime in increments of deltaTime
func (simulation *Simulation) initializeTimeArray() {
	for timestep := float64(0); timestep < simulation.totalSimulationTime+simulation.deltaTime; timestep += simulation.deltaTime {
		simulation.timeArray = append(simulation.timeArray, timestep)
	}
}

//EFFECTS: Run the simulation
//TODO: Add support for synaptic connections
func (simulation *Simulation) runSimulation() {
	//simulation code goes here
	for timeStep := 1; timeStep < len(simulation.timeArray); timeStep++ {
		for _, neuron := range simulation.neuronArray {
			neuron.calculateSimulationStep()
		}
	}
}

func main() {

	//Simulation Parameters
	totalSimulationTime := float64(55) //Total simulation time in milliseconds
	deltaTime := float64(0.025)        //Simulation timestep in milliseconds

	var simulation Simulation
	simulation.initializeSimulation(totalSimulationTime, deltaTime)
	simulation.initializeWeightMap()
	var params NeuronParameters
	params.initializeParametersWithDefaults() //defaults are initialized
	simulation.addNumberofNeuronsToSimulation(1)
	simulation.initializeNeuronArray(params)
	simulation.runSimulation()
	simulation.neuronArray[0].printToCSV()

}
