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

	/* Read our shader headers */
	group.Shaders = make([]*Shader, group.Count)
	for i := range group.Shaders {
		group.Shaders[i] = new(Shader)
	}

	for i, shader := range group.Shaders {
		if err := group.Detour(res, i, func() error {
			if err := shader.Unpack(res); err != nil {
				return err
			}
			return nil
		}); err != nil {
			log.Printf("Error reading shader %v\n", i)
			return err
		}
	}

	return nil
}
