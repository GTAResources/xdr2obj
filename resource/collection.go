package resource

import (
	"log"

	"github.com/tgascoigne/xdr2obj/resource/types"
)

type Collection struct {
	Lookup types.Ptr32
	Count  uint16
	Size   uint16
}

func (col *Collection) JumpTo(res *Container, i int) (err error) {
	var start types.Ptr32
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
