package resource

import (
	"log"
)

type Collection struct {
	Lookup Ptr32
	Count  uint16
	Size   uint16
}

func (col *Collection) JumpTo(res *Container, i int) (err error) {
	var start Ptr32
	if err = res.PeekElem(col.Lookup, i, &start); err != nil {
		log.Printf("Error performing collection lookup")
		return
	}

	if err = res.Jump(start); err != nil {
		log.Printf("Error performing collection lookup")
		return
	}
	return
}
