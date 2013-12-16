package resource

import (
	"bytes"
	"encoding/binary"
	"errors"
	"log"
	"io"

	"github.com/tgascoigne/xdr2obj/util/stack"
)

const (
	rscMagic = 0x52534337
	baseSize = 0x2000
)

var ErrInvalidResource error = errors.New("invalid resource")

type Container struct {
	Magic     uint32
	sysFlags uint32
	gfxFlags uint32
	Data      []byte

	SysOffset int64
	GfxOffset int64
	jumpStack stack.Stack
	position int64
	size int64
}

func (rsc *Container) Read(p []byte) (int, error) {
	var read int64
	toRead := int64(len(p))

	if rsc.position + toRead >= rsc.size {
		return 0, io.EOF
	}

	for read < toRead && rsc.position + read < rsc.size {
		p[read] = rsc.Data[rsc.position]
		rsc.position++
		read++
	}
	return int(read), nil
}

func (rsc *Container) Seek(offset int64, whence int) (int64, error) {
	var err error
	if whence == 0 {
		partition := (offset >> 24) & 0xFF;
		part_offset := int64(offset & 0xFFFFFF);
		switch partition {
		case 0x50:
			offset = rsc.SysOffset + part_offset
			break
		case 0x60:
			offset = rsc.GfxOffset + part_offset
			break
		}
 	} else if whence == 1 {
		offset = rsc.position + offset
	} else if whence == 2 {
		offset = rsc.size - offset - 1
	}

	if offset < 0 || offset >= rsc.size {
		err = io.EOF
	}

	rsc.position = offset
	return rsc.position, err
}

func (rsc *Container) Tell() int64 {
	return rsc.position
}

func (rsc *Container) Jump(offset int64) (int64, error) {
	position := rsc.Tell()
	rsc.jumpStack.Push(&stack.Item{position})
	return rsc.Seek(offset, 0)
}

func (rsc *Container) Return() (int64, error) {
	position := rsc.jumpStack.Pop()
	return rsc.Seek(position.Value.(int64), 0)
}

func (rsc *Container) Unpack(data []byte) (err error) {
	reader := bytes.NewReader(data)

	if err = binary.Read(reader, binary.BigEndian, &rsc.Magic); err != nil {
		return
	}

	if rsc.Magic != rscMagic {
		err = ErrInvalidResource
		return
	}

	if _, err = reader.Seek(0x4, 1); err != nil {
		return
	}
	if err = binary.Read(reader, binary.BigEndian, &rsc.sysFlags); err != nil {
		return
	}
	if err = binary.Read(reader, binary.BigEndian, &rsc.gfxFlags); err != nil {
		return
	}

	rsc.SysOffset = 0x10
	rsc.GfxOffset = rsc.SysOffset + int64(getPartitionSize(rsc.sysFlags))
	rsc.jumpStack.Allocate(0xF);
	rsc.Seek(0x50000000, 0) // seek to the start of the system partition
	rsc.Data = data
	rsc.size = int64(len(data))

	log.Printf("sys: %x gfx: %x", rsc.SysOffset, rsc.GfxOffset)
	return
}

func getPartitionSize(flags uint32) uint32 {
	var base uint32 = baseSize << (flags & 0xF)
	var size uint32 = 0
	size += ((flags >> 17) & 0x7F)
	size += (((flags >> 11) & 0x3F) << 1)
	size += (((flags >> 7) & 0xF) << 2)
	size += (((flags >> 5) & 0x3) << 3)
	size += (((flags >> 4) & 0x1) << 4)
	size *= base
	for i := uint(0); i < 4; i++ {
		if (flags >> (24 + i) & 1) == 1 {
			size += (base >> (1 + i))
		}
	}
	return size
}
