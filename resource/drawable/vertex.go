package drawable

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/tgascoigne/xdr2obj/resource"
	"github.com/tgascoigne/xdr2obj/resource/types"
)

type Vertex struct {
	/* This is the only information we can make use of for now */
	types.WorldCoord /* Position: Vertex.[X,Y,Z,W] */
	types.UV         /* UV: Vertex.[U,V] */
}

func (vert *Vertex) Unpack(res *resource.Container, buf *VertexBuffer) error {
	buffer := make([]byte, buf.Stride)
	reader := bytes.NewReader(buffer)

	/* Read the vertex into our local buffer */
	if size, err := res.Read(buffer); uint16(size) != buf.Stride || err != nil {
		return err
	}

	/* Parse out the info we can */
	if buf.Format.HasXYZ() {
		reader.Seek(0, 0)
		if err := binary.Read(reader, binary.BigEndian, &vert.WorldCoord); err != nil {
			return err
		}
	}

	if buf.Format.HasUV0() {
		reader.Seek(0x14, 0)
		if err := binary.Read(reader, binary.BigEndian, &vert.UV); err != nil {
			return err
		}
	}

	return nil
}

type VertexFormat uint32

const (
	FmtHasXYZ = (1 << 0)
	FmtHasUV0 = (1 << 4) /* I suspect there's probably multitexturing support */
)

func (f VertexFormat) HasXYZ() bool {
	return (uint32(f) & FmtHasXYZ) != 0
}

func (f VertexFormat) HasUV0() bool {
	return (uint32(f) & FmtHasXYZ) != 0
}

func (f VertexFormat) String() string {
	return fmt.Sprintf("0x%x XYZ: %v UV0: %v", uint32(f), f.HasXYZ(), f.HasUV0())
}