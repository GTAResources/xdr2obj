package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/tgascoigne/xdr2obj/export/obj"
	"github.com/tgascoigne/xdr2obj/resource"
	"github.com/tgascoigne/xdr2obj/resource/drawable"
)

func main() {
	var data []byte
	var err error

	/* Read the file */
	in_file := os.Args[1]
	if data, err = ioutil.ReadFile(in_file); err != nil {
		log.Fatal(err)
	}

	/* Unpack the container */
	res := new(resource.Container)
	if err = res.Unpack(data); err != nil {
		log.Fatal(err)
	}

	/* Unpack the drawable at 0x10 */
	drawable := new(drawable.Drawable)
	if err = drawable.Unpack(res); err != nil {
		log.Fatal(err)
	}

	/* Export it */
	if err = obj.Export(drawable); err != nil {
		log.Fatal(err)
	}
}
