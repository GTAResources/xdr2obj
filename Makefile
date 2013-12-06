all: xdr2obj
	
xdr2obj: main.c
	gcc -o $@ $< -std=gnu99