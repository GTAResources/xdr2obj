package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/tgascoigne/xdr2obj/resource"
)

func main() {
	var data []byte
	var err error

	defer func() {
		if err := recover(); err != nil {
			log.Fatalf("Unrecoverable error: %s", err)
		}
	}()

	in_file := os.Args[1]
	if data, err = ioutil.ReadFile(in_file); err != nil {
		log.Fatal(err)
	}

	var res resource.Container
	if err = res.Unpack(data); err != nil {
		log.Fatal(err)
	}
}
