package text

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type Text struct {
	Width int32
	Height int32
	Texture *sdl.Texture
}

// renderText renders texture from ttf font.
func New(text string, font *ttf.Font, e *sdl.Renderer,
			color sdl.Color, scale float32) (texture *sdl.Texture, err error) {

	surface, err := font.RenderUTF8Blended(text, color)
	if err != nil {
		return
	}
	defer surface.Free()

	width, height := surface.W, surface.H
	width = int32(float32(width) * scale)
	height = int32(float32(height) * scale)

	surface.W = width
	surface.H = height

	return e.CreateTextureFromSurface(surface)
}

