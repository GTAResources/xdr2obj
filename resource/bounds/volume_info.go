package bounds

import (
	"github.com/tgascoigne/xdr2obj/resource"
	"github.com/tgascoigne/xdr2obj/resource/types"
)

type VolumeInfoHeader struct {
	Low  types.Vec4
	High types.Vec4
}

type VolumeInfo struct {
	VolumeInfoHeader
	Center types.Vec4
}

func (info *VolumeInfo) Unpack(res *resource.Container) error {
	res.Parse(&info.VolumeInfoHeader)

	info.Center = types.Vec4{
		X: info.Low.X + ((info.High.X - info.Low.X) / 2),
		Y: info.Low.Y + ((info.High.Y - info.Low.Y) / 2),
		Z: info.Low.Z + ((info.High.Z - info.Low.Z) / 2),
		W: info.Low.W + ((info.High.W - info.Low.W) / 2),
	}

	return nil
}
