package drawable

import (
	"log"

	"github.com/tgascoigne/xdr2obj/resource"
	"github.com/tgascoigne/xdr2obj/resource/types"
)

type ModelCollection struct {
	resource.Collection
	Models []*Model
}

type ModelHeader struct {
	_                  uint32 /* vtable */
	GeometryCollection resource.Collection
	_                  types.Ptr32 /* Ptr to vectors */
	ShaderMappings     types.Ptr32
}

type Model struct {
	Header   ModelHeader
	Geometry []*Geometry
}

func (col *ModelCollection) Unpack(res *resource.Container) error {
	res.Parse(&col.Collection)

	col.Models = make([]*Model, col.Count)
	for i := range col.Models {
		col.Models[i] = new(Model)
	}

	/* Read our model headers */
	for i, model := range col.Models {
		if err := col.Detour(res, i, func() error {
			if err := model.Unpack(res); err != nil {
				log.Printf("Error reading model")
				return err
			}
			return nil
		}); err != nil {
			return err
		}
	}

	return nil
}

func (model *Model) Unpack(res *resource.Container) error {
	res.Parse(&model.Header)

	col := &model.Header.GeometryCollection

	model.Geometry = make([]*Geometry, col.Count)
	for i := range model.Geometry {
		model.Geometry[i] = new(Geometry)
	}

	for i, geom := range model.Geometry {
		if err := col.JumpTo(res, i); err != nil {
			log.Printf("Error reading geometry")
			return err
		}

		if err := geom.Unpack(res); err != nil {
			log.Printf("Error reading geometry")
			return err
		}

		if model.Header.ShaderMappings.Valid() {
			if err := res.PeekElem(model.Header.ShaderMappings, i, &geom.Shader); err != nil {
				return err
			}
		} else {
			geom.Shader = ShaderNone
		}

		if err := res.Return(); err != nil {
			log.Printf("Error reading geometry")
			return err
		}
	}

	return nil
}
