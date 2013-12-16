package drawable

import (
	"encoding/binary"
	"log"

	"github.com/tgascoigne/xdr2obj/resource"
	"github.com/tgascoigne/xdr2obj/util"
)

type Drawable struct {
	_               uint32
	blockMap        util.Ptr32
	shaderTable     util.Ptr32
	skeletonData    util.Ptr32
	center          util.Vec4
	boundsMin       util.Vec4
	boundsMax       util.Vec4
	modelCollection util.Ptr32
	lodCollections  [3]util.Ptr32
	pointMax        util.Vec4
}

func (drawable *Drawable) Unpack(res *resource.Container) (err error) {
	if err = binary.Read(res, binary.BigEndian, drawable); err != nil {
		return
	}
	log.Print(drawable)
	/*	if _, err = res.Skip(0x4); err != nil { /* Skip vtable * /
			return
		}
		if err = binary.Read(res, binary.BigEndian, &drawable.blockMap); err != nil {
			return
		}*/

	return
}
