package bounds

import (
	"log"

	"github.com/tgascoigne/xdr2obj/resource"
	"github.com/tgascoigne/xdr2obj/resource/types"
)

type VolumeHeader struct {
	_            uint32
	_            uint32
	_            uint32
	_            float32
	_            [4]types.Vec4 /* looks like world coords also */
	_            types.Vec4
	_            uint32
	_            uint32
	_            uint16
	_            uint16 /* some count */
	IndicesAddr  types.Ptr32
	ScaleFactor  types.Vec4
	WorldCoords  types.Vec4
	VerticesAddr types.Ptr32
	_            uint32
	_            uint32
	_            uint16
	VertexCount  uint16
	_            uint16
	IndexCount   uint16
	_            uint32
	_            uint32
	_            uint32
	_            types.Ptr32
	_            uint32
	_            uint32
	_            uint32
	_            uint32
	MaterialMap  types.Ptr32
	_            uint32 /* count? */
	_            uint32
	UnkPtr       types.Ptr32
	_            uint32
	_            uint32
	_            uint32
}

type Volume struct {
	VolumeHeader
	Vertices Vertices
	Indices  Indices
}

func (vol *Volume) Unpack(res *resource.Container) error {
	res.Parse(&vol.VolumeHeader)
	vol.Vertices = make(Vertices, vol.VertexCount)
	vol.Indices = make(Indices, vol.IndexCount)

	var vertAddr types.Ptr32
	res.Peek(vol.UnkPtr, &vertAddr)

	res.Detour(vol.VerticesAddr, func() error {
		log.Printf("verts at %x size %x\n", res.Tell(), vol.VertexCount*6)
		return vol.Vertices.Unpack(res)
	})

	res.Detour(vol.IndicesAddr, func() error {
		log.Printf("indices at %x size %x\n", res.Tell(), vol.IndexCount*16)
		return vol.Indices.Unpack(res)
	})

	return nil
}
