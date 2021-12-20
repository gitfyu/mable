package chat

import "errors"

var errInvalidColor = errors.New("invalid color")

// Color represents a color code used in a Msg
type Color uint8

const (
	// ColorReset resets the color to the default. Which value that is depends on the context of the Msg.
	ColorReset Color = iota
	ColorBlack
	ColorDarkBlue
	ColorDarkGreen
	ColorDarkCyan
	ColorDarkRed
	ColorPurple
	ColorGold
	ColorGray
	ColorDarkGray
	ColorBlue
	ColorGreen
	ColorCyan
	ColorRed
	ColorPink
	ColorYellow
	ColorWhite
)

var colorNames = [][]byte{
	[]byte("reset"),
	[]byte("black"),
	[]byte("dark_blue"),
	[]byte("dark_green"),
	[]byte("dark_aqua"),
	[]byte("dark_red"),
	[]byte("dark_purple"),
	[]byte("gold"),
	[]byte("gray"),
	[]byte("dark_gray"),
	[]byte("blue"),
	[]byte("green"),
	[]byte("aqua"),
	[]byte("red"),
	[]byte("light_purple"),
	[]byte("yellow"),
	[]byte("white"),
}

func (c Color) MarshalText() ([]byte, error) {
	if c < 0 || int(c) > len(colorNames) {
		return nil, errInvalidColor
	}

	return colorNames[c], nil
}
