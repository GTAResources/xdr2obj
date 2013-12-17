package resource

import (
	"log"
)

type ModelCollection struct {
	Collection
	Models []Model
}

type ModelHeader struct {
	_                  uint32 /* vtable */
	GeometryCollection Collection
	_                  Ptr32 /* Ptr to vectors */
	MaterialLookup     Ptr32
}

type Model struct {
	Header   ModelHeader
	Geometry []Geometry
}

func (col *ModelCollection) Unpack(res *Container) (err error) {
	if err = res.Parse(&col.Collection); err != nil {
		return err
	}

	log.Printf("Reading %v models", col.Count)

	/* Read our model headers */
	col.Models = make([]Model, col.Count)
	for i := range col.Models {
		if err = col.JumpTo(res, i); err != nil {
			log.Printf("Error reading model")
			continue
		}

		if err = col.Models[i].Unpack(res); err != nil {
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

func (model *Model) Unpack(res *Container) (err error) {
	if err = res.Parse(&model.Header); err != nil {
		return err
	}

	col := &model.Header.GeometryCollection

	log.Printf("Reading %v geometries", col.Count)

	model.Geometry = make([]Geometry, col.Count)
	for i := range model.Geometry {
		if err = col.JumpTo(res, i); err != nil {
			log.Printf("Error reading geometry")
			continue
		}

		if err = model.Geometry[i].Unpack(res); err != nil {
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
