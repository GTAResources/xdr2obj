package drawable

import (
	"bytes"
	"encoding/binary"
	"errors"
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
	if ofs, err := buf.Format.OffsetOf(VertXYZ); err == nil {
		reader.Seek(int64(ofs), 0)
		if err = binary.Read(reader, binary.BigEndian, &vert.WorldCoord); err != nil {
			return err
		}
	}

	if ofs, err := buf.Format.OffsetOf(VertUV); err == nil {
		reader.Seek(int64(ofs), 0)
		if err = binary.Read(reader, binary.BigEndian, &vert.UV); err != nil {
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

var (
	elementSizes = map[int]int{
		VertXYZ:   4 * 3,
		VertUnkA:  4,
		VertUnkB:  4,
		VertUnkC:  4,
		VertColor: 4,
		VertUnkD:  8,
		VertUV:    4,
	}
	ErrUnsupported error = errors.New("unsupported bit")
)

func (f VertexFormat) Supports(field int) bool {
	return (int(f) & field) != 0
}

func (f VertexFormat) OffsetOf(field int) (int, error) {
	if !f.Supports(field) {
		return 0, ErrUnsupported
	}

	size := 0
	for i := 1; i < field; i <<= 1 {
		if f.Supports(i) {
			size += elementSizes[i]
		}
	}

	return size, nil
}

func (f VertexFormat) String() string {
	return fmt.Sprintf("0x%x", int(f)) /* todo: better representation. */
}
