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
	if err = res.ParseStruct(&drawable.Header); err != nil {
		return
	}

	if err = res.Jump(drawable.Header.ModelCollection); err != nil {
		return
	}

	if err = drawable.Models.Unpack(res); err != nil {
		return
	}

	if err = res.Return(); err != nil {
		return
	}
	return
}
