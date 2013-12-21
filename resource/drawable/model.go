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
	MaterialLookup     types.Ptr32
}

type Model struct {
	Header   ModelHeader
	Geometry []*Geometry
}

func (col *ModelCollection) Unpack(res *resource.Container) (err error) {
	if err = res.Parse(&col.Collection); err != nil {
		return err
	}

	log.Printf("Reading %v models", col.Count)

	/* Read our model headers */
	col.Models = make([]*Model, col.Count)
	for i := range col.Models {
		col.Models[i] = new(Model)
	}

	for i, model := range col.Models {
		if err = col.JumpTo(res, i); err != nil {
			log.Printf("Error reading model")
			continue
		}

		if err = model.Unpack(res); err != nil {
			log.Printf("Error reading model")
			continue
		}

		if err = res.Return(); err != nil {
			log.Printf("Error reading model")
			continue
		}
	}

	return
}

func (model *Model) Unpack(res *resource.Container) (err error) {
	if err = res.Parse(&model.Header); err != nil {
		return err
	}

	col := &model.Header.GeometryCollection

	log.Printf("Reading %v geometries", col.Count)

	model.Geometry = make([]*Geometry, col.Count)
	for i := range model.Geometry {
		model.Geometry[i] = new(Geometry)
	}

	for i, geom := range model.Geometry {
		if err = col.JumpTo(res, i); err != nil {
			log.Printf("Error reading geometry")
			continue
		}

		if err = geom.Unpack(res); err != nil {
			log.Printf("Error reading geometry")
			continue
		}

		if err = res.Return(); err != nil {
			log.Printf("Error reading geometry")
			continue
		}
	}

	return
}
