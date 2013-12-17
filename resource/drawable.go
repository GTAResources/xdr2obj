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
	_               [6]uint32
	_               Ptr32
	Title           Ptr32
}

type Drawable struct {
	Header DrawableHeader
	Models ModelCollection
	Title  string
}

func (drawable *Drawable) Unpack(res *Container) (err error) {
	if err = res.Parse(&drawable.Header); err != nil {
		return err
	}

	err = res.Detour(drawable.Header.ModelCollection, func() error {
		return drawable.Models.Unpack(res)
	})
	if err != nil {
		return err
	}

	err = res.Detour(drawable.Header.Title, func() error {
		return res.Parse(&drawable.Title)
	})
	if err != nil {
		return err
	}

	return
}
