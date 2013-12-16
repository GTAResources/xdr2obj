package resource

import (
	"encoding/binary"
)

func parseStruct(res *Container, data interface{}) (err error) {
	return binary.Read(res, binary.BigEndian, data)
}

/* Borrowed/Adapted from encoding/binary/binary.go */
func intDataSize(data interface{}) int {
	switch data := data.(type) {
	case int8, *int8, *uint8:
		return 1
	case []int8:
		return len(data)
	case []uint8:
		return len(data)
	case int16, *int16, *uint16:
		return 2
	case []int16:
		return 2 * len(data)
	case []uint16:
		return 2 * len(data)
	case int32, *int32, *uint32:
		return 4
	case []int32:
		return 4 * len(data)
	case []uint32:
		return 4 * len(data)
	case int64, *int64, *uint64:
		return 8
	case []int64:
		return 8 * len(data)
	case []uint64:
		return 8 * len(data)
	case Ptr32, *Ptr32:
		return 4
	}
	return 0
}
