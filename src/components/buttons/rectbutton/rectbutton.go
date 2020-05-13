package rectbutton

import (
	"CardGameGo/src/components/text"
	"CardGameGo/src/utils"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type RectangularButton struct {
	BtnText string
	Width   int32
	Height  int32
	Color   *sdl.Color

	CallBack func(...interface{}) error
}

func New(text string, width, height int32, color *sdl.Color) *RectangularButton {
	button := &RectangularButton{
		BtnText: text,
		Width:   width,
		Height:  height,
		Color:   color,
	}

	return button
}

func (btn *RectangularButton) Draw(x, y int32, font *ttf.Font, renderer *sdl.Renderer) error {
	rect := sdl.Rect{
		X: x,
		Y: y,
		W: btn.Width,
		H: btn.Height,
	}

	_ = renderer.SetDrawColor(243, 241, 239, 1)
	_ = renderer.FillRect(&rect)
	textTexture, _ := text.New(btn.BtnText, font, renderer, sdl.Color{})
	_, _, tW, tH, _ := textTexture.Query()
	cenX, cenY := utils.GetCenterCoordinates(tW, tH, btn.Width, btn.Height)

	rect1 := &sdl.Rect{
		X: x + cenX,
		Y: y + cenY,
		W: tW,
		H: tH,
	}
	rect.Y = btn.Height/2
	return renderer.Copy(textTexture, nil, rect1)
}

func (btn *RectangularButton) Click() error {
	return btn.CallBack()
}
