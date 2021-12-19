package world

type Position struct {
	X     float64
	Y     float64
	Z     float64
	Yaw   float32
	Pitch float32
}

func NewPos(x, y, z float64, yaw, pitch float32) Position {
	return Position{
		X:     x,
		Y:     y,
		Z:     z,
		Yaw:   yaw,
		Pitch: pitch,
	}
}
