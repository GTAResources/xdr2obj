package shader

import (
	"log"

	"github.com/tgascoigne/xdr2obj/resource"
	"github.com/tgascoigne/xdr2obj/resource/types"
)

type GroupHeader struct {
	_                 uint32 /* vtable */
	TextureDictionary types.Ptr32
	resource.Collection
}

type Group struct {
	GroupHeader
	Shaders []*Shader
}

func (group *Group) Unpack(res *resource.Container) error {
	res.Parse(&group.GroupHeader)

	log.Printf("Reading %v shaders", group.Count)

	/* Read our shader headers */
	group.Shaders = make([]*Shader, group.Count)
	for i := range group.Shaders {
		group.Shaders[i] = new(Shader)
	}

	for i, shader := range group.Shaders {
		if err := group.JumpTo(res, i); err != nil {
			log.Printf("Error reading shader")
			return err
		}

		if err := shader.Unpack(res); err != nil {
			log.Printf("Error reading shader")
			return err
		}

		if err := res.Return(); err != nil {
			log.Printf("Error reading shader")
			return err
		}
	}

	return nil
}
