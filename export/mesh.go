package export

import (
	"fmt"

	"github.com/Jragonmiris/mathgl"

	"github.com/tgascoigne/xdr2obj/resource/types"
)

type Mesh struct {
	Vertices []mathgl.Vec4f
	Faces    []types.Tri
}

func NewMesh() *Mesh {
	return &Mesh{
		Vertices: make([]mathgl.Vec4f, 0),
		Faces:    make([]types.Tri, 0),
	}
}

func (mesh *Mesh) AddVert(pos mathgl.Vec4f) {
	mesh.Vertices = append(mesh.Vertices, pos)
}

func (mesh *Mesh) Rel(idx int) uint16 {
	numVerts := len(mesh.Vertices)
	return uint16(numVerts + idx)
}

func (mesh *Mesh) AddFace(face types.Tri) {
	max := func(i ...uint16) int {
		max := i[0]
		for _, x := range i {
			if x > max {
				max = x
			}
		}
		return int(max)
	}
	if max(face.A, face.B, face.C) > len(mesh.Vertices) {
		panic(fmt.Sprintf("invalid vert reference: %v", face))
	}
	mesh.Faces = append(mesh.Faces, face)
}
