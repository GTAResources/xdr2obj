package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/tgascoigne/xdr2obj/export/obj"
	"github.com/tgascoigne/xdr2obj/resource"
	"github.com/tgascoigne/xdr2obj/resource/dictionary"
	"github.com/tgascoigne/xdr2obj/resource/drawable"
	"github.com/tgascoigne/xdr2obj/resource/frag"
)

var (
	SupportedExtensions = []string{".xdr", ".xdd", ".xft"}
)

func main() {
	var mergeFile = flag.String("merge", "", "The basename of a file to merge all output to")
	flag.Parse()

	log.SetFlags(0)

	if *mergeFile != "" {
		obj.OpenMergeFile(*mergeFile)
		defer obj.CloseMergeFile()
	}

	converted := 0
	defer func() {
		log.Printf("Converted %v models\n", converted)
	}()

	if flag.NArg() > 0 {
		for _, s := range flag.Args() {
			func() {
				defer func() {
					if err := recover(); err != nil {
						log.Printf("Unable to convert %v\n", s)
						panic(err)
					}
				}()
				log.Printf("Converting %v..\n", s)
				processModel(s)
				log.Printf("done\n")
				converted++
			}()
		}
		return
	}

	/* Convert all supported models in directory */
	files, err := ioutil.ReadDir("./")
	if err != nil {
		panic(err)
	}

	for _, f := range files {
		for _, ext := range SupportedExtensions {
			if strings.Contains(f.Name(), ext) {
				func() {
					defer func() {
						if err := recover(); err != nil {
							log.Printf("Unable to convert %v\n", f.Name())
							panic(err)
						}
					}()
					filePath := fmt.Sprintf("./%v", f.Name())
					log.Printf("Converting %v..\n", filePath)
					processModel(filePath)
					log.Printf("done\n")
					converted++
				}()
			}
		}
	}

}

func processModel(inFile string) {
	var data []byte
	var err error

	/* Read the file */
	if data, err = ioutil.ReadFile(inFile); err != nil {
		panic(err)
	}

	/* Unpack the container */
	res := new(resource.Container)
	if err = res.Unpack(data); err != nil {
		panic(err)
	}

	switch {
	case filepath.Ext(inFile) == ".xdr":
		exportDrawable(res)
	case filepath.Ext(inFile) == ".xdd":
		exportDrawableDictionary(res)
	case filepath.Ext(inFile) == ".xft":
		exportFragType(res, filepath.Base(inFile))
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
	}

	/* Export it */
	for _, drawable := range dictionary.Drawables {
		if err := obj.Export(drawable); err != nil {
			log.Printf(err.Error())
		}
	}
}

func exportFragType(res *resource.Container, title string) {
	/* Unpack the frag type */
	frag := new(frag.FragType)
	if err := frag.Unpack(res); err != nil {
		log.Fatal(err)
	}

	/* Drawables inside frag files dont seem to be named properly. */
	frag.Drawable.Title = title

	/* Export it */
	if err := obj.Export(&frag.Drawable); err != nil {
		log.Fatal(err)
	}
}
