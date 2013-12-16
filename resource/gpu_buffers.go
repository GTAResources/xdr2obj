package resource

import (
	"bytes"
	"encoding/binary"
)

type VertexHeader struct {
	_      uint32 /* vtable */
	Stride uint16
	_      uint16
	Buffer Ptr32
	Count  uint32
	_      Ptr32 /* Buffer */
	_      uint32
	Info   Ptr32
	_      Ptr32
}

type VertexInfo struct {
	Format uint32
	_      uint16 /* Stride */
	_      uint16 /* correlated with Stride */
	_      uint32 /* usually 0xAA111111 */
	_      uint32 /* usually 0x1199a996 */
}

type Vertex struct {
	/* This is the only information we can make use of for now */
	Position  Vec4
	TexCoords Vec2
}

type VertexBuffer struct {
	VertexHeader
	VertexInfo
	Vertex []Vertex
}

func (buf *VertexBuffer) Unpack(res *Container) (err error) {
	if err = res.Parse(&buf.VertexHeader); err != nil {
		return
	}

	err = res.Detour(buf.Info, func() error {
		return res.Parse(&buf.VertexInfo)
	})
	if err != nil {
		return
	}

	buf.Vertex = make([]Vertex, buf.Count)

	err = res.Detour(buf.Buffer, func() error {
		buffer := make([]byte, buf.Stride)
		reader := bytes.NewReader(buffer)

		for _, v := range buf.Vertex {
			/* Read the vertex into our local buffer */
			if size, err := res.Read(buffer); uint16(size) != buf.Stride || err != nil {
				return err
			}

			/* Parse out the info we can */
			if err = binary.Read(reader, binary.BigEndian, &v.Position); err != nil {
				return err
			}
			reader.Seek(0, 0)
		}
		return nil
	})
	return
}

type IndexHeader struct {
	_      uint32 /* vtable */
	Count  uint32
	Buffer Ptr32
	Info   Ptr32
}

type IndexBuffer struct {
	IndexHeader
	Index  []Tri
	Stride int /* todo: is this referenced in the geom? */
}

func (buf *IndexBuffer) Unpack(res *Container) (err error) {
	buf.Stride = 3 * 2 // 3*uint16
	if err = res.Parse(&buf.IndexHeader); err != nil {
		return
	}

	buf.Index = make([]Tri, buf.Count/3)

	err = res.Detour(buf.Buffer, func() error {
		buffer := make([]byte, buf.Stride)
		reader := bytes.NewReader(buffer)

		for _, i := range buf.Index {
			/* Read the index into our local buffer */
			if size, err := res.Read(buffer); size != buf.Stride || err != nil {
				return err
			}

			/* Parse out the info we can */
			if err = binary.Read(reader, binary.BigEndian, &i); err != nil {
				return err
			}
			reader.Seek(0, 0)
		}
		return nil
	})
	return
}