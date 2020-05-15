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
	X       int32
	Y       int32
	Color   *sdl.Color
	Font    *ttf.Font

	CallBack func(...interface{}) error
}

func New(text string, width, height int32, color *sdl.Color, font *ttf.Font) *RectangularButton {
	button := &RectangularButton{
		BtnText: text,
		Width:   width,
		Height:  height,
		Color:   color,
		Font:    font,
	}

	return button
}

func (btn *RectangularButton) Draw(x, y int32, renderer *sdl.Renderer) error {
	rect := sdl.Rect{
		X: x,
		Y: y,
		W: btn.Width,
		H: btn.Height,
	}

	_ = renderer.SetDrawColor(btn.Color.R, btn.Color.G, btn.Color.B, btn.Color.A)
	_ = renderer.FillRect(&rect)
	textTexture, _ := text.New(btn.BtnText, btn.Font, renderer, sdl.Color{})
	defer textTexture.Destroy()

	_, _, tW, tH, _ := textTexture.Query()
	cenX, cenY := utils.GetCenterCoordinates(tW, tH, btn.Width, btn.Height)

	textRect := &sdl.Rect{
		X: x + cenX,
		Y: y + cenY,
		W: tW,
		H: tH,
	}

	btn.X = textRect.X
	btn.Y = textRect.Y

	return renderer.Copy(textTexture, nil, textRect)
}

func (btn *RectangularButton) GetX() int32 {
	return btn.X
}

func (btn *RectangularButton) GetY() int32 {
	return btn.Y
}

func (btn *RectangularButton) GetWidth() int32 {
	return btn.Width
}

func (btn *RectangularButton) GetHeight() int32 {
	return btn.Height
}

func (btn *RectangularButton) RunCallback(i ...interface{}) error {
	return btn.CallBack(i)
}
