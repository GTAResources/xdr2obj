package drawable

import (
	"github.com/tgascoigne/xdr2obj/resource"
	"github.com/tgascoigne/xdr2obj/resource/types"
)

type DrawableHeader struct {
	_               uint32
	BlockMap        types.Ptr32
	ShaderTable     types.Ptr32
	SkeletonData    types.Ptr32
	Center          types.Vec4
	BoundsMin       types.Vec4
	BoundsMax       types.Vec4
	ModelCollection types.Ptr32
	LodCollections  [3]types.Ptr32
	PointMax        types.Vec4
	_               [6]uint32
	_               types.Ptr32
	Title           types.Ptr32
}

type Drawable struct {
	Header DrawableHeader
	Models ModelCollection
	Title  string
}

func (drawable *Drawable) Unpack(res *resource.Container) error {
	if err := res.Parse(&drawable.Header); err != nil {
		return err
	}

	if err := res.Detour(drawable.Header.ModelCollection, func() error {
		return drawable.Models.Unpack(res)
	}); err != nil {
		return err
	}

	if err := res.Detour(drawable.Header.Title, func() error {
		return res.Parse(&drawable.Title)
	}); err != nil {
		return err
	}

	return nil
}
