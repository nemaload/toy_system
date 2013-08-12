#include <stdlib.h>
#include <time.h>
#include <math.h>

#define NUMBER_EXCITATORY_NEURONS 800
#define NUMBER_INHIBITORY_NEURONS 200

//convenience functions

//EFFECTS: Returns a random number between 0 and 1 
double generateRandomDouble() {
	return (double)rand() / (double)RAND_MAX;
}

//MODIFIES: inputArray
//EFFECTS: Fills inputArray completely with random numbers between 0 and 1.
void fillArrayWithRandomNumbers(double *inputArray, int arrayLength) {
	int currentArrayElementIndex;
	for (currentArrayElementIndex=0; 
		currentArrayElementIndex < arrayLength; 
		++currentArrayElementIndex) {
		inputArray[currentArrayElementIndex] = generateRandomDouble();		
	} 
}

int main(int argc, char **argv)
{
	srand(time(NULL)); //seed the random number generator
	//declare two arrays to hold the neurons
	double randomExcitatoryNeurons[NUMBER_EXCITATORY_NEURONS];
	double randomInhibitoryNeurons[NUMBER_INHIBITORY_NEURONS];
	//fill the excitatory neurons with random numbers between 0 and 1
	fillArrayWithRandomNumbers(randomExcitatoryNeurons, NUMBER_EXCITATORY_NEURONS);
	fillArrayWithRandomNumbers(randomInhibitoryNeurons, NUMBER_INHIBITORY_NEURONS);

	
	//use powf(a,b) to do float exponentiation

	//construct "a" matrix

	//take out fired
	//v and u values are relevant

}

