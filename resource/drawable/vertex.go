package drawable

import (
	"bytes"
	"encoding/binary"
	"fmt"
	_ "log"

	"github.com/tgascoigne/xdr2obj/resource"
	"github.com/tgascoigne/xdr2obj/resource/types"
)

type Vertex struct {
	/* This is the only information we can make use of for now */
	types.WorldCoord          /* Position: Vertex.[X,Y,Z,W] */
	UV0, UV1         types.UV /* UV: Vertex.[U,V] */
}

func (vert *Vertex) Unpack(res *resource.Container, buf *VertexBuffer) error {
	buffer := make([]byte, buf.Stride)
	reader := bytes.NewReader(buffer)

	/* Read the vertex into our local buffer */
	if size, err := res.Read(buffer); uint16(size) != buf.Stride || err != nil {
		return err
	}

	offset := 0

	/* Parse out the info we can */
	if buf.Format.Supports(VertXYZ) {
		if err := binary.Read(reader, binary.BigEndian, &vert.WorldCoord); err != nil {
			return err
		}
		offset += (4 * 3)
	}

	if buf.Format.Supports(VertUnkA) {
		junk := make([]byte, 4)
		reader.Read(junk)
		offset += 4
	}

	if buf.Format.Supports(VertUnkB) {
		junk := make([]byte, 4)
		reader.Read(junk)
		offset += 4
	}

	if buf.Format.Supports(VertUnkC) {
		junk := make([]byte, 4)
		reader.Read(junk)
		offset += 4
	}

	if buf.Format.Supports(VertColor) {
		junk := make([]byte, 4)
		reader.Read(junk)
		offset += 4
	}

	if buf.Format.Supports(VertUnkD) {
		junk := make([]byte, 4)
		reader.Read(junk)
		offset += 4
	}

	if buf.Format.Supports(VertUV0) {
		//		log.Printf("uv0 at %v", offset)
		if err := binary.Read(reader, binary.BigEndian, &vert.UV0); err != nil {
			return err
		}
		offset += 4
	}

	if buf.Format.Supports(VertUV1) {
		//		log.Printf("uv1 at %v", offset)
		if err := binary.Read(reader, binary.BigEndian, &vert.UV1); err != nil {
			return err
		}
		offset += 4
	}

	if buf.Format.Supports(VertUnkX) {
		junk := make([]byte, 4)
		reader.Read(junk)
		offset += 4
	}
	return nil
}

type VertexFormat uint32

const (
	VertXYZ   = (1 << 0)  /* + 4*3 */
	VertUnkA  = (1 << 1)  /* + 4 */
	VertUnkB  = (1 << 2)  /* + 4 */
	VertUnkC  = (1 << 3)  /* + 4 */
	VertColor = (1 << 4)  /* + 4 */
	VertUnkD  = (1 << 5)  /* + 4 */
	VertUV0   = (1 << 6)  /* + 4 */
	VertUV1   = (1 << 7)  /* + 4 */
	VertUnkX  = (1 << 15) /* + 4 */
)

func (f VertexFormat) Supports(field int) bool {
	return (int(f) & field) != 0
}

func (f VertexFormat) String() string {
	return fmt.Sprintf("0x%x", int(f)) /* todo: better representation. */
}
