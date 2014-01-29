package bounds

import (
	"fmt"
	"log"
	"os"

	"github.com/tgascoigne/xdr2obj/resource"
	"github.com/tgascoigne/xdr2obj/resource/types"
)

type NodesHeader struct {
	_            uint32
	_            types.Ptr32
	_            uint32
	_            uint32
	WorldCoords  [4]types.Vec4
	BoundingBox  types.Vec4
	BoundsTable  types.Ptr32
	VolumeMatrix types.Ptr32
	_            types.Ptr32
	VolumeInfo   types.Ptr32
	_            types.Ptr32
	_            types.Ptr32
	Count        uint16
	Capacity     uint16
	_            types.Ptr32
}

type Nodes struct {
	NodesHeader
	Volumes []*Volume
}

func (nodes *Nodes) Unpack(res *resource.Container) error {
	res.Parse(&nodes.NodesHeader)

	var err error

	nodes.Volumes = make([]*Volume, nodes.Capacity)
	volCollection := resource.PointerCollection{
		Addr:     nodes.BoundsTable,
		Count:    nodes.Count,
		Capacity: nodes.Capacity,
	}

	volInfoCollection := resource.Collection{
		Addr:     nodes.VolumeInfo,
		Count:    nodes.Count,
		Capacity: nodes.Capacity,
	}

	err = volInfoCollection.For(res, func(i int) error {
		nodes.Volumes[i] = new(Volume)
		if err := nodes.Volumes[i].VolumeInfo.Unpack(res); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	err = volCollection.For(res, func(i int) error {
		log.Printf("Volume %v at %v\n", i, res.Tell())
		return nodes.Volumes[i].Unpack(res)
	})
	if err != nil {
		return err
	}

	/* test export - simple obj */
	for i, vol := range nodes.Volumes {
		file, err := os.Create(fmt.Sprintf("volume_%v.obj", i))
		if err != nil {
			panic(err)
		}
		defer file.Close()

		fmt.Fprintf(file, "o volume_%v\n", i)

		for _, v := range vol.Vertices {
			x, y, z := float32(v.X), float32(v.Y), float32(v.Z)
			fmt.Fprintf(file, "v %v %v %v\n", x, z, y)
		}

		for _, i := range vol.Faces {
			a, b, c := i.A+1, i.B+1, i.C+1
			fmt.Fprintf(file, "f %v %v %v\n", a, b, c)
		}
	}

	return nil
}
