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

//MODIFIES: neuron's currentTimeStep
//EFFECTS: Instantiates a single neuron
func (neuron *Neuron) initializeNeuron(simulation *Simulation) {
	neuron.setSimulation(simulation)
	neuron.intializeVoltageArray()
	neuron.initializeStimulationArray()
	neuron.currentTimeStep = 1
	neuron.initializeDimensionlessQuantities()

}

func (neuron *Neuron) initializeStimulationArray() {
	neuron.stimulation = make([]float64, len(neuron.simulation.timeArray))
}

//MODIFIES: neuron's stimulation array
//EFFECTS: Fills the stimulation array with a single 10mV stimulation period from 5ms to 30ms
func (neuron *Neuron) setSampleStimulationValues() {
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
func (neuron *Neuron) intializeVoltageArray() {
	neuron.V_m = make([]float64, len(neuron.simulation.timeArray))
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
	neuron.m += (neuron.sodiumAlphaM()*(1-neuron.m) - //FUNCTIONS ARE CALLING CURRENT VALUES AND NOT PREVIOUS ONES
		neuron.sodiumBetaM()*neuron.m) * neuron.simulation.deltaTime
	neuron.h += (neuron.sodiumAlphaH()*(1-neuron.h) -
		neuron.sodiumBetaH()*neuron.h) * neuron.simulation.deltaTime
	neuron.n += (neuron.potassiumAlphaN()*(1-neuron.n) -
		neuron.potassiumBetaN()*neuron.n) * neuron.simulation.deltaTime

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
		(n.previousVoltage()-n.parameters.sodiumReversePotential) - n.potassiumConductance*
		(n.previousVoltage()-n.parameters.potassiumReversePotential) - n.parameters.leakConductance*
		(n.previousVoltage()-n.parameters.leakReversePotential)) / n.parameters.lipidBilayerCapacitance * n.simulation.deltaTime
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

func (neuron *Neuron) previousVoltage() float64 {
	return neuron.V_m[int(neuron.currentTimeStep)-1]
}

//EFFECTS: Calculates and returns the neuron's current first potassium rate constant value
func (neuron *Neuron) potassiumAlphaN() float64 {
	if neuron.previousVoltage() != 10 {
		return 0.01 * (-neuron.previousVoltage() + 10) / (math.Exp((-neuron.previousVoltage()+10)/10) - 1)
	}
	return 0.1
}

//EFFECTS: Calculates and returns the neuron's current second potassium rate constant value
func (neuron *Neuron) potassiumBetaN() float64 {
	return 0.125 * math.Exp(-neuron.previousVoltage()/80)
}

//EFFECTS: Calculates and returns the neuron's current steady state potassium activation value
func (neuron *Neuron) potassiumNInfinity() float64 {
	return neuron.potassiumAlphaN() /
		(neuron.potassiumAlphaN() + neuron.potassiumBetaN())
}

//EFFECTS: Calculates and returns the neuron's current first sodium rate constant value
func (neuron *Neuron) sodiumAlphaM() float64 {
	if neuron.previousVoltage() != 25 { //change to previous voltage
		return 0.1 * (-neuron.previousVoltage() + 25) / (math.Exp((-neuron.previousVoltage()+25)/10) - 1)
	}
	return 1
}

//EFFECTS: Calculates and returns the neuron's current second sodium rate constant value
func (neuron *Neuron) sodiumBetaM() float64 {
	return 4 * math.Exp(-neuron.previousVoltage()/18)
}

//EFFECTS: Calculates and returns the neuron's current sodium steady state activation value
func (neuron *Neuron) sodiumMInfinity() float64 {
	return neuron.sodiumAlphaM() /
		(neuron.sodiumAlphaM() + neuron.sodiumBetaM())
}

//EFFECTS: Calculates and returns the neuron's current first sodium inactivation rate constant value
func (neuron *Neuron) sodiumAlphaH() float64 {
	return 0.07 * math.Exp(-neuron.previousVoltage()/20)
}

//EFFECTS: Calculates and returns the neuron's current second sodium inactivation rate constant value
func (neuron *Neuron) sodiumBetaH() float64 {
	return 1 / (math.Exp((-neuron.previousVoltage()+30)/10) + 1)
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
		neuron.parameters = params
		neuron.initializeNeuron(simulation)
	}
}

//MODIFIES: Simulation's neuronArray
//EFFECTS: Allocates new neurons, and adds them to the simulation's neuronArray
func (simulation *Simulation) addNumberofNeuronsToSimulation(neuronCount int) {
	for i := 0; i < neuronCount; i++ {
		simulation.neuronArray = append(simulation.neuronArray, new(Neuron))
	}
}

func (simulation *Simulation) allocateWeightMap() {
	simulation.weightMap = make(map[*Neuron]map[*Neuron]float64, len(simulation.neuronArray))
	for _, neuron := range simulation.neuronArray {
		simulation.weightMap[neuron] = make(map[*Neuron]float64, len(simulation.neuronArray))

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
func (simulation *Simulation) setSynapseWeightPair(neuron1, neuron2 *Neuron, weight float64) {
	simulation.weightMap[neuron1][neuron2] = weight
}

//MODIFIES: simulation's neuronArray's neurons
//EFFECTS: Updates all of the stimulations based on the current voltages
func (simulation *Simulation) updateNeuronStimulationValues() {
	for _, neuron1 := range simulation.neuronArray {
		for neuron2 := range simulation.weightMap[neuron1] {
			if neuron1 == neuron2 {
				continue
			}
			neuron2.stimulation[int(neuron2.currentTimeStep)] += neuron1.V_m[int(neuron1.currentTimeStep)-1] * simulation.weightMap[neuron1][neuron2]
		}
	}
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
		simulation.updateNeuronStimulationValues()
		for _, neuron := range simulation.neuronArray {
			neuron.calculateSimulationStep()
		}
	}
}

func (simulation *Simulation) printToCSV() {
	for row, time := range simulation.timeArray {
		fmt.Print(time, ",")
		for index, neuron := range simulation.neuronArray {
			if index != len(simulation.neuronArray)-1 {
				fmt.Print(neuron.V_m[row], ",")
			} else {
				fmt.Print(neuron.V_m[row])
			}
		}
		fmt.Print("\n")
	}
}

/*type HDF5File struct {
	dims      []int
	dataspace *Dataspace
}*/

func main() {

	//Simulation Parameters
	totalSimulationTime := float64(220) //Total simulation time in milliseconds
	deltaTime := float64(0.025)         //Simulation timestep in milliseconds

	var simulation Simulation
	simulation.initializeSimulation(totalSimulationTime, deltaTime)
	var params NeuronParameters
	params.initializeParametersWithDefaults() //defaults are initialized
	simulation.addNumberofNeuronsToSimulation(3)
	simulation.initializeNeuronArray(params)
	simulation.neuronArray[0].setSampleStimulationValues()
	simulation.allocateWeightMap()
	simulation.initializeWeightMap()
	simulation.setSynapseWeightPair(simulation.neuronArray[0], simulation.neuronArray[1], 0.8)
	simulation.setSynapseWeightPair(simulation.neuronArray[1], simulation.neuronArray[2], 0.9)
	simulation.setSynapseWeightPair(simulation.neuronArray[2], simulation.neuronArray[1], 0.1)
	simulation.setSynapseWeightPair(simulation.neuronArray[2], simulation.neuronArray[0], 0.1)

	simulation.runSimulation()
	simulation.printToCSV()

}
