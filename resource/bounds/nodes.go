package bounds

import (
	"fmt"
	"log"
	"os"

	"github.com/tgascoigne/xdr2obj/resource"
	"github.com/tgascoigne/xdr2obj/resource/types"
)

type NodesHeader struct {
	_           uint32
	_           types.Ptr32
	_           uint32
	_           uint32
	WorldCoords [4]types.Vec4
	BoundingBox types.Vec4
	BoundsTable types.Ptr32
	Mat4        types.Ptr32 /* transformation matrix? */
	_           types.Ptr32
	_           types.Ptr32
	_           types.Ptr32
	_           types.Ptr32
	Count       uint16
	Size        uint16
	_           types.Ptr32
}

type Nodes struct {
	NodesHeader
	Volumes []*Volume
}

func (nodes *Nodes) Unpack(res *resource.Container) error {
	res.Parse(&nodes.NodesHeader)

	nodes.Volumes = make([]*Volume, nodes.Size)
	volCollection := resource.Collection{
		Lookup: nodes.BoundsTable,
		Count:  nodes.Count,
		Size:   nodes.Size,
	}

	for i := 0; i < int(nodes.Count); i++ {
		volCollection.Detour(res, i, func() error {
			log.Printf("volume %v at %v\n", i, res.Tell())
			nodes.Volumes[i] = new(Volume)
			return nodes.Volumes[i].Unpack(res)
		})
	}

	/* test export - simple obj */
	for i, vol := range nodes.Volumes {
		file, err := os.Create(fmt.Sprintf("volume_%v.obj", i))
		if err != nil {
			panic(err)
		}
		defer file.Close()

		numVerts := int(vol.VertexCount)

		log.Printf("%v verts\n", numVerts)
		log.Printf("%v faces\n", len(vol.Indices))

		scale := vol.ScaleFactor

		fmt.Fprintf(file, "o volume_%v\n", i)
		for _, v := range vol.Vertices {
			x, y, z := float32(v.X)*scale.Y, float32(v.Y)*scale.Y, float32(v.Z)*scale.Y
			fmt.Fprintf(file, "v %v %v %v\n", x, y, z)
		}

		for _, i := range vol.Indices {
			a, b, c := i.A+1, i.B+1, i.C+1
			fmt.Fprintf(file, "f %v %v %v\n", a, b, c)
		}
	}

	return nil
}