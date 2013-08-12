#include <stdlib.h>
#include <time.h>
#include <math.h>
#include "lib/ziggurat/ziggurat.h"

#define NUMBER_EXCITATORY_NEURONS 800
#define NUMBER_INHIBITORY_NEURONS 200
#define SIMULATION_TIME_MS 1000

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

//MODIFIES: inputArray
//EFFECTS: Fills numElements elements of inputArray with number
void fillArrayWithNumber(double *inputArray, double number, int numElements) {
	int currentArrayElementIndex;
	for (currentArrayElementIndex=0;
		currentArrayElementIndex < numElements;
		++currentArrayElementIndex) {
		inputArray[currentArrayElementIndex] = number;
	}
}

double generateNormallyDistributedRandomDouble( unsigned long int *seed,
	int ziggurat_kn[128], float ziggurat_fn[128], float ziggurat_wn[128]) {
	return (double) r4_nor(seed, ziggurat_kn, ziggurat_fn, ziggurat_wn);
}


int main(int argc, char **argv)
{
	srand(time(NULL)); //seed the random number generator
	//set up the ziggurat normally distributed pseudorandom number generator
	int ziggurat_kn[128];
	float ziggurat_fn[128];
	float ziggurat_wn[128];
	r4_nor_setup(ziggurat_kn,ziggurat_fn,ziggurat_wn);





	//declare two arrays to hold random values
	//corresponds to "Ne=800;     Ni=200;"
	double randomExcitatory[NUMBER_EXCITATORY_NEURONS];
	double randomInhibitory[NUMBER_INHIBITORY_NEURONS];
	//fill the excitatory array with random numbers between 0 and 1
	//corresponds to "re=rand(Ne,1); ri=rand(Ni,1);"
	fillArrayWithRandomNumbers(randomExcitatory, 
							   NUMBER_EXCITATORY_NEURONS);
	fillArrayWithRandomNumbers(randomInhibitory, 
							   NUMBER_INHIBITORY_NEURONS);
	//declare an array to hold the "a" in the Izhekevich model
	double a[NUMBER_EXCITATORY_NEURONS+NUMBER_INHIBITORY_NEURONS];
	
	int currentArrayElementIndex; //declare the iterator
	//corresponds to "a=[0.02*ones(Ne,1);     0.02+0.08*ri];"
	for (currentArrayElementIndex = 0;
		currentArrayElementIndex < NUMBER_EXCITATORY_NEURONS;
		++currentArrayElementIndex) {
		//fill the first NUMBER_EXCITATORY_NEURONS with 0.02
		a[currentArrayElementIndex] = 0.02;
	}
	for (currentArrayElementIndex = NUMBER_EXCITATORY_NEURONS;
		currentArrayElementIndex < 
			(NUMBER_EXCITATORY_NEURONS+ NUMBER_INHIBITORY_NEURONS);
		++currentArrayElementIndex) {
		//fill the rest of the array with the formula below
		a[currentArrayElementIndex] = 
			0.02 + 0.08 * randomInhibitory[currentArrayElementIndex];
	}

	//declare the "b" array
	double b[NUMBER_EXCITATORY_NEURONS + NUMBER_INHIBITORY_NEURONS];
	//corresponds to "b=[0.2*ones(Ne,1);      0.25-0.05*ri];"
	for (currentArrayElementIndex = 0;
		currentArrayElementIndex < NUMBER_EXCITATORY_NEURONS;
		++currentArrayElementIndex) {
		b[currentArrayElementIndex] = 0.02;
	}
	for (currentArrayElementIndex = NUMBER_EXCITATORY_NEURONS;
		currentArrayElementIndex < 
			(NUMBER_EXCITATORY_NEURONS + NUMBER_INHIBITORY_NEURONS);
		++currentArrayElementIndex) {
		b[currentArrayElementIndex] = 
			0.25 - 0.05* randomInhibitory[currentArrayElementIndex];
	}

	//declare the "c" array
	//corresponds to "c=[-65+15*re.^2;        -65*ones(Ni,1)];"
	double c[NUMBER_EXCITATORY_NEURONS + NUMBER_INHIBITORY_NEURONS];
	for (currentArrayElementIndex = 0;
		currentArrayElementIndex < NUMBER_EXCITATORY_NEURONS;
		++currentArrayElementIndex) {
		c[currentArrayElementIndex] = 
			-65 + 15 * powf(randomExcitatory[currentArrayElementIndex],2);
	}
	for (currentArrayElementIndex = 0;
		currentArrayElementIndex < 
			(NUMBER_EXCITATORY_NEURONS + NUMBER_INHIBITORY_NEURONS);
		++currentArrayElementIndex) {
		c[currentArrayElementIndex] = -65;
	}
	//corresponds to "d=[8-6*re.^2;           2*ones(Ni,1)];"
	double d[NUMBER_EXCITATORY_NEURONS + NUMBER_INHIBITORY_NEURONS];
	for (currentArrayElementIndex = 0;
		currentArrayElementIndex < NUMBER_EXCITATORY_NEURONS;
		++currentArrayElementIndex) {
		d[currentArrayElementIndex] = 
			8 - 6 * powf(randomExcitatory[currentArrayElementIndex],2);
	}
	for (currentArrayElementIndex = 0;
		currentArrayElementIndex < 
			(NUMBER_EXCITATORY_NEURONS + NUMBER_INHIBITORY_NEURONS);
		++currentArrayElementIndex) {
		c[currentArrayElementIndex] = 2;
	}

	//S is a square 2D array
	//corresponds to S=[0.5*rand(Ne+Ni,Ne),  -rand(Ne+Ni,Ni)];
	double S[NUMBER_INHIBITORY_NEURONS + NUMBER_EXCITATORY_NEURONS]
	[NUMBER_EXCITATORY_NEURONS + NUMBER_INHIBITORY_NEURONS];

	int currentRowIndex;
	int currentColumnIndex;
	//fill the first NUMBER_EXCITATORY_NEURONS rows with 0.5*random numbers
	for (currentRowIndex = 0;
		currentRowIndex < NUMBER_EXCITATORY_NEURONS;
		++currentRowIndex) {
		for (currentColumnIndex = 0;
			currentColumnIndex < 
				(NUMBER_EXCITATORY_NEURONS + NUMBER_INHIBITORY_NEURONS);
			++currentColumnIndex) {
			S[currentRowIndex][currentColumnIndex] = 
				0.5 * generateRandomDouble();
		}
	}

	//fill the next NUMBER_INHIBITORY_NEURONS rows with -random numbers
	for (currentRowIndex = NUMBER_EXCITATORY_NEURONS;
		currentRowIndex < 
			(NUMBER_EXCITATORY_NEURONS + NUMBER_INHIBITORY_NEURONS);
		++currentRowIndex) {
		for (currentColumnIndex = 0;
			currentColumnIndex < 
				(NUMBER_EXCITATORY_NEURONS + NUMBER_INHIBITORY_NEURONS);
			++currentColumnIndex) {
			S[currentRowIndex][currentColumnIndex] = 
				-1.0 * generateRandomDouble();
		}
	}

	//generate initial values of v
	//corresponds to "v=-65*ones(Ne+Ni,1);    % Initial values of v"
	double v[NUMBER_EXCITATORY_NEURONS + NUMBER_INHIBITORY_NEURONS];
	for (currentArrayElementIndex = 0;
		currentArrayElementIndex < 
			(NUMBER_EXCITATORY_NEURONS + NUMBER_INHIBITORY_NEURONS);
		++currentArrayElementIndex) {
		v[currentArrayElementIndex] = -65;
	}
	//fill u
	//corresponds to "u=b.*v;    % Initial values of u"
	double u[NUMBER_EXCITATORY_NEURONS + NUMBER_INHIBITORY_NEURONS];
	for (currentArrayElementIndex = 0;
		currentArrayElementIndex <
			(NUMBER_EXCITATORY_NEURONS + NUMBER_INHIBITORY_NEURONS);
		++currentArrayElementIndex) {
		u[currentArrayElementIndex] = 
			b[currentArrayElementIndex] * v[currentArrayElementIndex];
	}

	//double firings[]; 
	int currentTime;

	for (currentTime = 0; currentTime < SIMULATION_TIME_MS; ++currentTime) {
		

	}









	//use powf(a,b) to do float exponentiation

	//construct "a" matrix

	//take out fired
	//v and u values are relevant

}

