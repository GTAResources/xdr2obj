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

func (col *Collection) Detour(res *Container, i int, callback func() error) error {
	addr, err := col.GetPtr(res, i)
	if err != nil {
		return err
	}

	return res.Detour(addr, callback)
}

func (col *Collection) JumpTo(res *Container, i int) error {
	addr, err := col.GetPtr(res, i)
	if err != nil {
		return err
	}

	if err = res.Jump(addr); err != nil {
		log.Printf("Error performing collection lookup")
		return err
	}

	return nil
}

func (col *Collection) GetPtr(res *Container, i int) (types.Ptr32, error) {
	var addr types.Ptr32
	if err := res.PeekElem(col.Lookup, i, &addr); err != nil {
		log.Printf("Error performing collection lookup")
		return 0, err
	}
	return addr, nil
}
