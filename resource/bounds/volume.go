package bounds

import (
	"fmt"
	"log"

	"github.com/tgascoigne/xdr2obj/export"
	"github.com/tgascoigne/xdr2obj/resource"
	"github.com/tgascoigne/xdr2obj/resource/types"
)

type VolumeHeader struct {
	_             uint32
	_             uint32
	_             uint32
	_             float32
	_             [4]types.Vec4 /* looks like world coords also */
	_             types.Vec4
	_             uint32
	_             uint32
	_             uint16
	_             uint16 /* some count */
	IndicesAddr   types.Ptr32
	ScaleFactor   types.Vec4
	Offset        types.Vec4
	VerticesAddr  types.Ptr32
	_             uint32
	_             uint32
	_             uint16
	VertexCount   uint16
	_             uint16
	IndexCount    uint16
	_             uint32
	_             uint32
	_             uint32
	_             types.Ptr32
	_             uint32
	_             uint32
	_             uint32
	_             uint32
	MaterialMap   types.Ptr32
	_             uint32 /* count? */
	_             uint32
	BoundBoxesPtr types.Ptr32
	_             uint32
	_             uint32
	_             uint32
}

type Volume struct {
	VolumeHeader
	VolumeInfo
	*export.Mesh
}

func (vol *Volume) Unpack(res *resource.Container) error {
	res.Parse(&vol.VolumeHeader)

	vol.Mesh = export.NewMesh()

	res.Detour(vol.VerticesAddr, func() error {
		log.Printf("verts at %x size %x\n", res.Tell(), vol.VertexCount*6)
		return vol.unpackVertices(res)
	})

	res.Detour(vol.IndicesAddr, func() error {
		log.Printf("indices at %x size %x\n", res.Tell(), vol.IndexCount*16)
		return vol.unpackFaces(res)
	})

	return nil
}

func (vol *Volume) unpackVertices(res *resource.Container) error {
	for i := 0; i < int(vol.VertexCount); i++ {
		iVec := new(types.Vec3i)
		res.Parse(iVec)
		v := types.Vec3{
			X: (float32(iVec.X) * vol.ScaleFactor.X) + vol.Offset.X,
			Y: (float32(iVec.Y) * vol.ScaleFactor.Y) + vol.Offset.Y,
			Z: (float32(iVec.Z) * vol.ScaleFactor.Z) + vol.Offset.Z,
		}
		vol.AddVert(v)
	}

	return nil
}

func (vol *Volume) unpackFaces(res *resource.Container) error {
	fixIndex := func(c uint16) uint16 {
		c &= 0x7FFF
		return c
	}

	for i := 0; i < int(vol.IndexCount); i++ {
		var polygonType, junk uint16
		res.Parse(&junk)
		res.Parse(&polygonType)
		idxValues := make([]uint16, 6)
		res.Parse(idxValues)

		debug := func(junk, polyType uint16, idxValues []uint16) {
			fmt.Printf("%.4x ", junk)
			fmt.Printf("%.4x ", polyType)
			for _, i := range idxValues {
				fmt.Printf("%.4x ", i)
			}
			fmt.Printf("\n")
		}

		polygonType &= 0xF
		if polygonType == 0x4 {
			/* what is this? */
			debug(junk, polygonType, idxValues)
		} else if polygonType == 0x3 || polygonType == 0xB {
			/* cube, 4 points specified */
			a := fixIndex(idxValues[0])
			b := fixIndex(idxValues[1])
			c := fixIndex(idxValues[2])
			d := fixIndex(idxValues[3])
			buildCube(vol.Mesh, a, b, c, d)
		} else if polygonType == 0x2 {
			/* seems to be two points */
			a := fixIndex(junk)
			b := fixIndex(idxValues[2])
			traceVerts(vol.Mesh, a, b)
		} else /* if polygonType == 0x8 */ { /* it's probably a triangle */
			a := fixIndex(idxValues[0])
			b := fixIndex(idxValues[1])
			c := fixIndex(idxValues[2])
			vol.AddFace(types.Tri{
				A: a,
				B: b,
				C: c,
			})
		}
	}
	return nil
}

func traceVerts(mesh *export.Mesh, verts ...uint16) {
	if len(verts) < 2 {
		panic("not enough verts to trace")
	}
	for i := 1; i < len(verts); i++ {
		mesh.AddVert(mesh.Vertices[verts[i]])
		mesh.AddVert(mesh.Vertices[verts[i]])
		mesh.AddVert(mesh.Vertices[verts[i-1]])
		mesh.AddFace(types.Tri{mesh.Rel(-1), mesh.Rel(-2), mesh.Rel(-3)})
	}
}

func buildCube(mesh *export.Mesh, a, b, c, d uint16) {
	A, B, C, D := mesh.Vertices[a], mesh.Vertices[b], mesh.Vertices[c], mesh.Vertices[d]

	/* create the base vertices */
	mesh.AddVert(types.Vec3{
		X: D.X,
		Y: D.Y,
		Z: B.Z,
	})
	mesh.AddVert(types.Vec3{
		X: C.X,
		Y: C.Y,
		Z: A.Z,
	})
	mesh.AddVert(types.Vec3{
		X: B.X,
		Y: B.Y,
		Z: A.Z,
	})
	mesh.AddVert(types.Vec3{
		X: A.X,
		Y: A.Y,
		Z: B.Z,
	})

	/* Create the faces */
	e, f, g, h := mesh.Rel(-1), mesh.Rel(-2), mesh.Rel(-3), mesh.Rel(-4)
	mesh.AddFace(types.Tri{
		A: a,
		B: f,
		C: e,
	})
	mesh.AddFace(types.Tri{
		A: e,
		B: f,
		C: b,
	})

	mesh.AddFace(types.Tri{
		A: f,
		B: d,
		C: b,
	})
	mesh.AddFace(types.Tri{
		A: d,
		B: h,
		C: b,
	})

	mesh.AddFace(types.Tri{
		A: d,
		B: g,
		C: h,
	})
	mesh.AddFace(types.Tri{
		A: g,
		B: c,
		C: h,
	})

	mesh.AddFace(types.Tri{
		A: g,
		B: a,
		C: c,
	})
	mesh.AddFace(types.Tri{
		A: c,
		B: a,
		C: e,
	})

}
