#include <stdlib.h>
#include <time.h>

#define NUMBER_EXCITATORY_NEURONS 800
#define NUMBER_INHIBITORY_NEURONS 200

//test comment

srand(time(NULL)); //seed the random number generator

double randomExcitatoryNeurons[NUMBER_EXCITATORY_NEURONS];
double (*randomExcitatoryNeurons)[NUMBER_EXCITATORY_NEURONS];
double (*randomInhibitoryNeurons)[NUMBER_INHIBITORY_NEURONS];

double (*amatrix)[NUMBER_EXCITATORY_NEURONS + NUMBER_INHIBITORY_NEURONS];
double (*bmatrix)[]

void fillArrayWithRandomNumbers(double inputArray, int arrayLength) {
	for (int i=0; i < arrayLength; ++i) inputArray[i] = (double)rand() / (double)RAND_MAX;
}

//fill the excitatory neurons with random numbers between 0 and 1
fillArrayWithRandomNumbers(&randomExcitatoryNeurons, NUMBER_EXCITATORY_NEURONS);
fillArrayWithRandomNumbers(&randomInhibitoryNeurons, NUMBER_INHIBITORY_NEURONS);

//construct "a" matrix

//take out fired
//v and u values are relevant


int r = rand
