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
			color sdl.Color, size int) (*sdl.Texture, error) {

	surface, err := font.RenderUTF8Blended(text, color)
	if err != nil {
		return nil, err
	}
	defer surface.Free()

	return e.CreateTextureFromSurface(surface)
}

