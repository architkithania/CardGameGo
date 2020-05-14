package gamemanager

import (
	"CardGameGo/src/components/buttons/imagebutton"
	"CardGameGo/src/managers/eventmanager"
	"CardGameGo/src/managers/fontmanager"
	"CardGameGo/src/managers/imgmanager"
	"CardGameGo/src/managers/interfaces"
	"CardGameGo/src/utils"
	"errors"
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
)

type GameUiManager struct {
	GameId string

	Players       map[*interfaces.Player]bool
	Host          *interfaces.Player
	DevicePlayer  *interfaces.Player
	CurrentPlayer *interfaces.Player

	Cards       []rune
	PlayedCard  rune
	DeviceTurn  bool
	GameStarted bool
}

func New(devicePlayer *interfaces.Player, context interfaces.GameContext) *GameUiManager {
	ui := GameUiManager{
		GameId:        context.GameId,
		Players:       context.Players,
		Host:          context.Host,
		DevicePlayer:  devicePlayer,
		CurrentPlayer: nil,
		Cards:         nil,
		PlayedCard:    0,
		DeviceTurn:    false,
		GameStarted:   false,
	}

	return &ui
}

func (ui *GameUiManager) AddNewPlayer(player *interfaces.Player) {
	ui.Players[player] = true
}

func (ui *GameUiManager) RemovePlayer(player *interfaces.Player) {
	delete(ui.Players, player)
}

func (ui *GameUiManager) SetCurrentPlayer(player *interfaces.Player) error {
	if _, ok := ui.Players[player]; !ok {
		return errors.New(fmt.Sprintf("change player error: %q not found in game", player.Name))
	}
	ui.CurrentPlayer = player
	return nil
}

func (ui *GameUiManager) StartGame() {
	ui.GameStarted = true
}

func (ui *GameUiManager) Draw(
	winWidth, winHeight int32,
	fontManager *fontmanager.FontManager,
	eventManager *eventmanager.EventManager,
	imageManager *imgmanager.ImageManager,
	renderer *sdl.Renderer,
	) error {

	return ui.drawCardRack(winWidth, winHeight, imageManager, renderer)
}

func (ui *GameUiManager) drawCardRack(w, h int32, imageManager *imgmanager.ImageManager, renderer *sdl.Renderer) error {
	rectHeight := int32(utils.Percent(h, 20))

	rect := sdl.Rect{
		X: 0,
		Y: h - rectHeight,
		W: w,
		H: rectHeight,
	}

	_ = renderer.SetDrawColor(255, 0, 0, 255)
	err := renderer.FillRect(&rect)
	if err != nil {
		return err
	}

	image := imageManager.Images["cards/c01"]
	_, _, imageW, _, _ := image.Query()

	c01 := imagebutton.New(image)

	intervals := generateCenteredIntervals(w, imageW, 13, 45)
	for _, e := range intervals {
		err = c01.Draw(e, h - rectHeight, renderer)
		if err != nil {
			return err
		}
	}

	return nil
}

func generateCenteredIntervals(width, cardWidth int32, count int, delta int32) []int32 {
	cardSpace := (int32(count) - 1)*delta + cardWidth
	start := (width - cardSpace)/2
	return generateIntervals(start, count, delta)
}

func generateIntervals(start int32, count int, delta int32) []int32 {
	interval := make([]int32, count, count)
	interval[0] = start
	for i := 1; i < count; i++ {
		interval[i] = interval[i - 1] + delta
	}
	return interval
}
