package dictionary

import (
	"github.com/tgascoigne/xdr2obj/resource"
	"github.com/tgascoigne/xdr2obj/resource/drawable"
)

type DrawableCollection struct {
	resource.Collection
	Drawables []*drawable.Drawable
}
