package resource

import (
	"log"
)

type ModelCollectionHeader struct {
	Lookup Ptr32
	Count  uint16
	Size   uint16
}

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
	if err = res.ParseStruct(&col.Collection); err != nil {
		return err
	}

	log.Printf("Reading %v models", col.Count)

	/* Read our model headers */
	col.Models = make([]Model, col.Count)
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

func (model *Model) Unpack(res *Container) (err error) {
	if err = res.ParseStruct(&model.Header); err != nil {
		return err
	}

	col := &model.Header.GeometryCollection

	log.Printf("Reading %v geometries", col.Count)

	model.Geometry = make([]Geometry, col.Count)
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
