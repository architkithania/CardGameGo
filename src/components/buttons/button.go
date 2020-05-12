package buttons

import "github.com/veandco/go-sdl2/sdl"

type Button interface {
	click() error
	getRectClickedState() *sdl.Rect
	getRectNormalState() *sdl.Rect
}
