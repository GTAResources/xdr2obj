package resource

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"

	"github.com/tgascoigne/xdr2obj/util/stack"
)

const (
	resMagic  = 0x52534337
	baseSize  = 0x2000
	stringMax = 32
)

var ErrInvalidResource error = errors.New("invalid resource")
var ErrInvalidString error = errors.New("invalid string")

type ContainerHeader struct {
	Magic    uint32
	_        uint32
	SysFlags uint32
	GfxFlags uint32
}

type Container struct {
	Header    ContainerHeader
	SysOffset int64
	GfxOffset int64

	jumpStack stack.Stack
	position  int64
	size      int64
	Data      []byte
}

func (res *Container) Unpack(data []byte) (err error) {
	reader := bytes.NewReader(data)

	header := &res.Header
	if err = binary.Read(reader, binary.BigEndian, header); err != nil {
		return
	}

	if header.Magic != resMagic {
		err = ErrInvalidResource
		return
	}

	res.SysOffset = 0x10
	res.GfxOffset = res.SysOffset + int64(getPartitionSize(header.SysFlags))
	res.jumpStack.Allocate(0xF)
	res.Seek(0x50000000, 0) // seek to the start of the system partition
	res.Data = data
	res.size = int64(len(data))
	return
}

func (res *Container) Parse(data interface{}) (err error) {
	switch data.(type) {
	case *string:
		/* find NULL */
		var i int64
		buf := make([]byte, stringMax)
		err = ErrInvalidString
		for i = 0; i < stringMax; i++ {
			if res.Data[res.position+i] == 0 {
				err = nil
				break
			}
		}
		str := data.(*string)
		copy(buf, res.Data[res.position:res.position+i])
		*str = string(buf)
		res.position += i
		return err
	default:
		return binary.Read(res, binary.BigEndian, data)
	}
	return nil
}

func (res *Container) Detour(addr Ptr32, callback func() error) (err error) {
	if err = res.Jump(addr); err != nil {
		return
	}

	defer func() {
		if err == nil {
			err = res.Return()
		} else {
			res.Return()
		}
	}()

	return callback()
}

func (res *Container) Peek(addr Ptr32, data interface{}) (err error) {
	if err = res.Jump(addr); err != nil {
		return err
	}

	if err = binary.Read(res, binary.BigEndian, data); err != nil {
		return err
	}

	if err = res.Return(); err != nil {
		return err
	}
	return
}

func (res *Container) PeekElem(addr Ptr32, element int, data interface{}) (err error) {
	return res.Peek(addr+Ptr32(element*intDataSize(data)), data)
}

/* Container Util functions */
func (res *Container) Read(p []byte) (int, error) {
	var read int64
	toRead := int64(len(p))

	if res.position+toRead >= res.size {
		return 0, io.EOF
	}

	for read < toRead && res.position+read < res.size {
		p[read] = res.Data[res.position]
		res.position++
		read++
	}
	return int(read), nil
}

func (res *Container) Skip(offset int64) (int64, error) {
	return res.Seek(offset, 1)
}

func (res *Container) Seek(offset int64, whence int) (int64, error) {
	var err error
	if whence == 0 {
		partition := (offset >> 24) & 0xFF
		part_offset := int64(offset & 0xFFFFFF)
		switch partition {
		case 0x50:
			offset = res.SysOffset + part_offset
		case 0x60:
			offset = res.GfxOffset + part_offset
		}
	} else if whence == 1 {
		offset = res.position + offset
	} else if whence == 2 {
		offset = res.size - offset - 1
	}

	if offset < 0 || offset >= res.size {
		err = io.EOF
	}

	res.position = offset
	return res.position, err
}

func (res *Container) Tell() int64 {
	return res.position
}

func (res *Container) Jump(offset Ptr32) error {
	position := res.Tell()
	res.jumpStack.Push(&stack.Item{position})
	_, err := res.Seek(int64(offset), 0)
	return err
}

func (res *Container) Return() error {
	position := res.jumpStack.Pop()
	_, err := res.Seek(position.Value.(int64), 0)
	return err
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
