package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/tgascoigne/xdr2obj/export/obj"
	"github.com/tgascoigne/xdr2obj/resource"
	"github.com/tgascoigne/xdr2obj/resource/dictionary"
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

	switch {
	case filepath.Ext(in_file) == ".xdr":
		exportDrawable(res)
	case filepath.Ext(in_file) == ".xdd":
		exportDrawableDictionary(res)
	}

}

func exportDrawable(res *resource.Container) {
	/* Unpack the drawable */
	drawable := new(drawable.Drawable)
	if err := drawable.Unpack(res); err != nil {
		log.Fatal(err)
	}

	/* Export it */
	if err := obj.Export(drawable); err != nil {
		log.Fatal(err)
	}
}

func exportDrawableDictionary(res *resource.Container) {
	/* Unpack the dictionary */
	dictionary := new(dictionary.Dictionary)
	if err := dictionary.Unpack(res); err != nil {
		log.Fatal(err)
	}

	/* Fix up the file names */
	for i, drawable := range dictionary.Drawables {
		drawable.Title = strings.Replace(drawable.Title, ".#dd", fmt.Sprintf("_%v.#dd", i), -1)
		log.Print(drawable.Title)
	}

	/* Export it */
	for _, drawable := range dictionary.Drawables {
		if err := obj.Export(drawable); err != nil {
			log.Printf(err.Error())
		}
	}
}
