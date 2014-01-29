package export

import (
	"github.com/tgascoigne/xdr2obj/resource/types"
)

type Mesh struct {
	Vertices []types.Vec3
	Faces    []types.Tri
}

func NewMesh() *Mesh {
	return &Mesh{
		Vertices: make([]types.Vec3, 0),
		Faces:    make([]types.Tri, 0),
	}
}

func (mesh *Mesh) AddVert(pos types.Vec3) {
	mesh.Vertices = append(mesh.Vertices, pos)
}

func (mesh *Mesh) Rel(idx int) uint16 {
	numVerts := len(mesh.Vertices)
	return uint16(numVerts + idx)
}

func (mesh *Mesh) AddFace(face types.Tri) {
	mesh.Faces = append(mesh.Faces, face)
}
