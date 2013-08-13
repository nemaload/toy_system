all: ziggurat izhekevich 
	gcc -O3 -o izhekevich izhekevich.o ziggurat.o -lhdf5 -lm 


ziggurat: ./lib/ziggurat/ziggurat.c
	gcc -O3 -c -o ziggurat.o ./lib/ziggurat/ziggurat.c 

izhekevich: izhekevich.c
	gcc -O3 -c -o izhekevich.o izhekevich.c 

clean: ziggurat izhekevich 
	rm -f ./ziggurat.o
	rm -f ./izhekevich.o
	rm -f ./izhekevich