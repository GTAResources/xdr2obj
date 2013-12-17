package obj

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/tgascoigne/xdr2obj/resource"
)

func Export(drawable *resource.Drawable) (err error) {
	baseName := drawable.Title[0 : strings.LastIndex(drawable.Title, ".")-1]
	objFileName := fmt.Sprintf("%v.obj", baseName)
	objFile, err := os.Create(objFileName)
	if err != nil {
		return err
	}
	defer func() {
		objFile.Close()
	}()

	for i, model := range drawable.Models.Models {
		modelName := fmt.Sprintf("%v_%v", baseName, i)
		fmt.Fprintf(objFile, "o %v\n", modelName)
		if err = exportModel(&model, objFile, modelName); err != nil {
			return err
		}
	}

	log.Printf("Wrote %v", objFileName)
	return nil
}

func exportModel(model *resource.Model, file *os.File, name string) (err error) {
	indexBase := uint16(1)
	for i := range model.Geometry {
		geom := &model.Geometry[i]
		fmt.Fprintf(file, "g %v_%v\n", name, i)
		if err = exportGeometry(geom, file, name, indexBase); err != nil {
			return err
		}

		indexBase += geom.VertexCount
	}
	return nil
}

// idxBase is added to each index specified. Used for grouping geometry properly
func exportGeometry(geom *resource.Geometry, file *os.File, name string, idxBase uint16) (err error) {
	for _, vert := range geom.Vertices.Vertex {
		/* todo: Do we need the W here? Is it even a W coord? */
		fmt.Fprintf(file, "v %v %v %v\n", vert.X, vert.Y, vert.Z)
	}

	for _, tri := range geom.Indices.Index {
		fmt.Fprintf(file, "f %v %v %v\n", idxBase+tri.A, idxBase+tri.B, idxBase+tri.C)
	}
	return nil
}
