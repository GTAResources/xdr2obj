package bounds

import (
	_ "log"

	"github.com/tgascoigne/xdr2obj/resource"
	"github.com/tgascoigne/xdr2obj/resource/types"
)

type Vertices []*Vertex

type Vertex struct {
	X, Y, Z int16
}

func (v Vertices) Unpack(res *resource.Container) error {
	for i := 0; i < cap(v); i++ {
		v[i] = new(Vertex)
		res.Parse(v[i])
		v[i].Y, v[i].Z = v[i].Z, v[i].Y
	}

	return nil
}

type Index struct {
	_ float32
	types.Tri
	_ types.Tri /* Face relations. Adjacency based collision? */
}

type Indices []*Index

func (idx Indices) Unpack(res *resource.Container) error {
	fixIndex := func(c uint16) uint16 {
		c &= 0xFFF
		return c
	}

	for i := 0; i < cap(idx); i++ {
		idx[i] = new(Index)
		res.Parse(idx[i])
		idx[i].A = fixIndex(idx[i].A)
		idx[i].B = fixIndex(idx[i].B)
		idx[i].C = fixIndex(idx[i].C)
	}
	return nil
}
