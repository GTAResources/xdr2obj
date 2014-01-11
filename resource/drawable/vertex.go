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
	if buf.Format.Supports(VertXYZ) {
		if err := binary.Read(reader, binary.BigEndian, &vert.WorldCoord); err != nil {
			return err
		}
	}

	if buf.Format.Supports(VertUnkA) {
		junk := make([]byte, 4)
		reader.Read(junk)
	}

	if buf.Format.Supports(VertUnkB) {
		junk := make([]byte, 4)
		reader.Read(junk)
	}

	if buf.Format.Supports(VertUnkC) {
		junk := make([]byte, 4)
		reader.Read(junk)
	}

	if buf.Format.Supports(VertColor) {
		junk := make([]byte, 4)
		reader.Read(junk)
	}

	if buf.Format.Supports(VertUnkD) {
		junk := make([]byte, 8)
		reader.Read(junk)
	}

	if buf.Format.Supports(VertUV) {
		if err := binary.Read(reader, binary.BigEndian, &vert.UV); err != nil {
			return err
		}
	}

	return nil
}

type VertexFormat uint32

const (
	VertXYZ   = (1 << 0) /* + 4*3 */
	VertUnkA  = (1 << 1) /* + 4 */
	VertUnkB  = (1 << 2) /* + 4 */
	VertUnkC  = (1 << 3) /* + 4 */
	VertColor = (1 << 4) /* + 4 */
	VertUnkD  = (1 << 5) /* + 8 */
	VertUV    = (1 << 6) /* + 4 */
)

func (f VertexFormat) Supports(field int) bool {
	return (int(f) & field) != 0
}

func (f VertexFormat) String() string {
	return fmt.Sprintf("0x%x", int(f)) /* todo: better representation. */
}
