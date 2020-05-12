package rectbutton

import "github.com/veandco/go-sdl2/sdl"

type RectangularButton struct {
	text string
	buttonColor *sdl.Color
	textColor *sdl.Color
	width int32
	height int32
}
