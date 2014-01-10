package obj

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/tgascoigne/xdr2obj/resource/drawable"
	"github.com/tgascoigne/xdr2obj/resource/drawable/shader"
)

var (
	ErrUnsupportedVertexFormat = errors.New("unsupported vertex format")
)

const (
	BitmapFormat = "dds"
)

var materialNames map[int]string
var mergeObjFile, mergeMtlFile *os.File

func OpenMergeFile(baseName string) error {
	var objFile, mtlFile *os.File
	var err error
	if objFile, err = os.Create(fmt.Sprintf("%v.obj", baseName)); err != nil {
		return err
	}

	if mtlFile, err = os.Create(fmt.Sprintf("%v.mtl", baseName)); err != nil {
		return err
	}

	mergeObjFile = objFile
	mergeMtlFile = mtlFile
	return nil
}

func CloseMergeFile() {
	if mergeObjFile != nil {
		mergeObjFile.Close()
	}
	if mergeMtlFile != nil {
		mergeMtlFile.Close()
	}
}

func Export(drawable *drawable.Drawable) error {
	baseName := drawable.Title[0:strings.LastIndex(drawable.Title, ".")]
	objFileName := fmt.Sprintf("%v.obj", baseName)
	mtlFileName := fmt.Sprintf("%v.mtl", baseName)
	materialNames = make(map[int]string)

	var objFile, mtlFile *os.File
	var err error
	if mergeObjFile != nil && mergeMtlFile != nil {
		objFile = mergeObjFile
		mtlFile = mergeMtlFile
	} else {
		if objFile, err = os.Create(objFileName); err != nil {
			return err
		}
		defer objFile.Close()
		if mtlFile, err = os.Create(mtlFileName); err != nil {
			return err
		}
		defer mtlFile.Close()
	}

	for i, shader := range drawable.Shaders.Shaders {
		materialNames[i] = fmt.Sprintf("%v_%v", baseName, i)
		fmt.Fprintf(mtlFile, "newmtl %v\n", materialNames[i])
		if err := exportMaterial(shader, mtlFile); err != nil {
			return err
		}
	}

	for i, model := range drawable.Models.Models {
		modelName := fmt.Sprintf("%v_%v", baseName, i)
		fmt.Fprintf(objFile, "o %v\n", modelName)
		if err := exportModel(model, objFile, modelName); err != nil {
			return err
		}
	}

	return nil
}

func exportMaterial(material *shader.Shader, file *os.File) error {
	if material.DiffusePath != "" {
		fmt.Fprintf(file, "map_Kd %v.%v\n", material.DiffusePath, BitmapFormat)
	}

	if material.NormalPath != "" {
		fmt.Fprintf(file, "bump %v.%v\n", material.NormalPath, BitmapFormat)
	}

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
	if !geom.Vertices.Format.Supports(drawable.VertXYZ) {
		return ErrUnsupportedVertexFormat
	}

	if geom.Shader != drawable.ShaderNone {
		fmt.Fprintf(file, "usemtl %v\n", materialNames[int(geom.Shader)])
	}

	for _, vert := range geom.Vertices.Vertex {
		fmt.Fprintf(file, "v %v %v %v\n", vert.X, vert.Y, vert.Z)
		if geom.Vertices.Format.Supports(drawable.VertUV) {
			u := vert.U.Value()
			v := (-vert.V.Value()) + 1
			fmt.Fprintf(file, "vt %v %v\n", u, v)
		}
	}

	numVerts := len(geom.Vertices.Vertex)

	for _, tri := range geom.Indices.Index {
		a, b, c := -int(numVerts-int(tri.A)), -int(numVerts-int(tri.B)), -int(numVerts-int(tri.C))
		if geom.Vertices.Format.Supports(drawable.VertUV) {
			fmt.Fprintf(file, "f %v/%v %v/%v %v/%v\n", a, a, b, b, c, c)
		} else {
			fmt.Fprintf(file, "f %v %v %v\n", a, b, c)
		}
	}
	return nil
}
