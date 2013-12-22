package obj

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/tgascoigne/xdr2obj/resource/drawable"
)

var (
	ErrUnsupportedVertexFormat = errors.New("unsupported vertex format")
)

func Export(drawable *drawable.Drawable) error {
	baseName := drawable.Title[0 : strings.LastIndex(drawable.Title, ".")-1]
	objFileName := fmt.Sprintf("%v.obj", baseName)

	var objFile *os.File
	var err error
	if objFile, err = os.Create(objFileName); err != nil {
		return err
	}
	defer objFile.Close()

	for i, model := range drawable.Models.Models {
		modelName := fmt.Sprintf("%v_%v", baseName, i)
		fmt.Fprintf(objFile, "o %v\n", modelName)
		if err := exportModel(model, objFile, modelName); err != nil {
			return err
		}
	}

	log.Printf("Wrote %v", objFileName)
	return nil
}

func exportModel(model *drawable.Model, file *os.File, name string) error {
	indexBase := uint16(1)

	for i, geom := range model.Geometry {
		fmt.Fprintf(file, "g %v_%v\n", name, i)
		if err := exportGeometry(geom, file, name, indexBase); err != nil {
			return err
		}

		indexBase += geom.VertexCount
	}
	return nil
}

// idxBase is added to each index specified. Used for grouping geometry properly
func exportGeometry(geom *drawable.Geometry, file *os.File, name string, idxBase uint16) error {
	if !geom.Vertices.Format.HasXYZ() {
		return ErrUnsupportedVertexFormat
	}

	for _, vert := range geom.Vertices.Vertex {
		/* todo: Do we need the W here? Is it even a W coord? */
		fmt.Fprintf(file, "v %v %v %v\n", vert.X, vert.Y, vert.Z)
		if geom.Vertices.Format.HasUV0() {
			u := vert.U.Value()
			v := (-vert.V.Value()) + 1
			fmt.Fprintf(file, "vt %v %v\n", u, v)
		}
	}

	for _, tri := range geom.Indices.Index {
		a, b, c := idxBase+tri.A, idxBase+tri.B, idxBase+tri.C
		if geom.Vertices.Format.HasUV0() {
			fmt.Fprintf(file, "f %v/%v %v/%v %v/%v\n", a, a, b, b, c, c)
		} else {
			fmt.Fprintf(file, "f %v %v %v\n", a, b, c)
		}
	}
	return nil
}
