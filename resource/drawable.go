package resource

type DrawableHeader struct {
	_               uint32
	BlockMap        Ptr32
	ShaderTable     Ptr32
	SkeletonData    Ptr32
	Center          Vec4
	BoundsMin       Vec4
	BoundsMax       Vec4
	ModelCollection Ptr32
	LodCollections  [3]Ptr32
	PointMax        Vec4
}

type Drawable struct {
	Header DrawableHeader
	Models ModelCollection
}

func (drawable *Drawable) Unpack(res *Container) (err error) {
	if err = res.Parse(&drawable.Header); err != nil {
		return
	}

	err = res.Detour(drawable.Header.ModelCollection, func() error {
		return drawable.Models.Unpack(res)
	})
	if err != nil {
		return
	}

	return
}
