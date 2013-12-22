package types

import (
	"math"
)

type Ptr32 uint32

type float16 uint16

func (i float16) Value() float32 {
	/* Lovingly adapted from http://stackoverflow.com/a/15118210 */
	t1 := uint32(i & 0x7fff)
	t2 := uint32(i & 0x8000)
	t3 := uint32(i & 0x7c00)
	t1 <<= 13
	t2 <<= 16
	t1 += 0x38000000
	if t3 == 0 {
		t1 = 0
	}
	t1 |= t2
	return math.Float32frombits(t1)
}

type UV struct {
	U float16
	V float16
}

type Vec2 struct {
	X float32
	Y float32
}

type Vec3 struct {
	X float32
	Y float32
	Z float32
}

type Vec4 struct {
	X float32
	Y float32
	Z float32
	W float32
}

type WorldCoord Vec4

type Tri struct {
	A uint16
	B uint16
	C uint16
}
