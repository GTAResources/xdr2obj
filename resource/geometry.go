package resource

import (
	"log"
)

type GeometryHeader struct {
	_             uint32 /* vtable */
	_             uint32
	_             uint32
	VertexBuffer  Ptr32
	_             uint32
	_             uint32
	_             uint32
	IndexBuffer   Ptr32
	_             uint32
	_             uint32
	_             uint32
	IndexCount    uint32
	FaceCount     uint32
	VertexCount   uint16
	PrimitiveType uint16
}

type Geometry struct {
	GeometryHeader
}

func (geom *Geometry) Unpack(res *Container) (err error) {
	if err = res.ParseStruct(&geom.GeometryHeader); err != nil {
		return
	}

	log.Printf("Found geometry with %v vertices %v triangles", geom.VertexCount, geom.FaceCount)
	return
}
