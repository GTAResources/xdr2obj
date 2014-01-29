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
var destObjFile, destMtlFile *os.File

func OpenDestFile(baseName string) error {
	var objFile, mtlFile *os.File
	var err error
	if objFile, err = os.Create(fmt.Sprintf("%v.obj", baseName)); err != nil {
		return err
	}

	if mtlFile, err = os.Create(fmt.Sprintf("%v.mtl", baseName)); err != nil {
		return err
	}

	destObjFile = objFile
	destMtlFile = mtlFile
	return nil
}

func CloseDestFile() {
	if destObjFile != nil {
		destObjFile.Close()
	}
	if destMtlFile != nil {
		destMtlFile.Close()
	}
}

func Export(drawable *drawable.Drawable) error {
	baseName := drawable.Title[0:strings.LastIndex(drawable.Title, ".")]
	objFileName := fmt.Sprintf("%v.obj", baseName)
	mtlFileName := fmt.Sprintf("%v.mtl", baseName)
	materialNames = make(map[int]string)

	var objFile, mtlFile *os.File
	var err error
	if destObjFile != nil && destMtlFile != nil {
		objFile = destObjFile
		mtlFile = destMtlFile
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

	fmt.Fprintf(objFile, "mtllib %v\n", mtlFileName)

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
	for i, geom := range model.Geometry {
		fmt.Fprintf(file, "g %v_%v\n", name, i)
		if err := exportGeometry(geom, file, name); err != nil {
			return err
		}
	}
	return nil
}

func exportGeometry(geom *drawable.Geometry, file *os.File, name string) error {
	if !geom.Vertices.Format.Supports(drawable.VertXYZ) {
		return ErrUnsupportedVertexFormat
	}

	if geom.Shader != drawable.ShaderNone {
		fmt.Fprintf(file, "usemtl %v\n", materialNames[int(geom.Shader)])
	}

	for _, vert := range geom.Vertices.Vertex {
		fmt.Fprintf(file, "v %v %v %v\n", vert.WorldCoord[0], vert.WorldCoord[1], vert.WorldCoord[2])
		if geom.Vertices.Format.Supports(drawable.VertUV0) {
			u := vert.UV0.U.Value()
			v := (-vert.UV0.V.Value()) + 1
			fmt.Fprintf(file, "vt %v %v\n", u, v)
		}
	}

	numVerts := len(geom.Vertices.Vertex)

	for _, tri := range geom.Indices.Index {
		a, b, c := -int(numVerts-int(tri.A)), -int(numVerts-int(tri.B)), -int(numVerts-int(tri.C))
		if geom.Vertices.Format.Supports(drawable.VertUV0) {
			fmt.Fprintf(file, "f %v/%v %v/%v %v/%v\n", a, a, b, b, c, c)
		} else {
			fmt.Fprintf(file, "f %v %v %v\n", a, b, c)
		}
	}
	return nil
}
