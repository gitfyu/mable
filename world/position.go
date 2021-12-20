package world

type Pos struct {
	X     float64
	Y     float64
	Z     float64
	Yaw   float32
	Pitch float32
}

func NewPos(x, y, z float64, yaw, pitch float32) Pos {
	return Pos{
		X:     x,
		Y:     y,
		Z:     z,
		Yaw:   yaw,
		Pitch: pitch,
	}
}
