package gamemanager

import (
	"CardGameGo/src/components/buttons/imagebutton"
	"CardGameGo/src/components/buttons/rectbutton"
	"CardGameGo/src/managers/eventmanager"
	"CardGameGo/src/managers/fontmanager"
	"CardGameGo/src/managers/imgmanager"
	"CardGameGo/src/managers/interfaces"
	"CardGameGo/src/utils"
	"errors"
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"strconv"
)

var allCards = make(map[string]*imagebutton.ImageButton)
var playButton *rectbutton.RectangularButton = nil
var claimButton *rectbutton.RectangularButton = nil
var playerIcon *rectbutton.RectangularButton = nil
var claimedHandsText *rectbutton.RectangularButton = nil

var cardYPosition int32

type GameUiManager struct {
	GameId string

	Players       map[*interfaces.Player]bool
	Host          *interfaces.Player
	DevicePlayer  *interfaces.Player
	CurrentPlayer *interfaces.Player

	Cards       []string
	PlayedCards []string
	DeviceTurn  bool
	GameStarted bool

	selectedCard string
	claimedHands int
}

func callBackGenerator(ui *GameUiManager, cardName string) func(...interface{}) error {
	return func(...interface{}) error {
		if ui.selectedCard == cardName {
			ui.selectedCard = ""
		} else {
			ui.selectedCard = cardName
		}
		return nil
	}
}

func (ui *GameUiManager) Init(manager *imgmanager.ImageManager,
	eventManager *eventmanager.EventManager, fontManager *fontmanager.FontManager, renderer *sdl.Renderer) {

	// Init card image buttons
	cardNames := []string{"c1", "c2", "c3", "c4", "c5", "c6", "c7", "c8", "c9", "cX", "cJ", "cQ", "cK",
		"s1", "s2", "s3", "s4", "s5", "s6", "s7", "s8", "s9", "sX", "sJ", "sQ", "sK",
		"h1", "h2", "h3", "h4", "h5", "h6", "h7", "h8", "h9", "hX", "hJ", "hQ", "hK",
		"d1", "d2", "d3", "d4", "d5", "d6", "d7", "d8", "d9", "dX", "dJ", "dQ", "dK"}

	for _, card := range cardNames {
		allCards[card] = imagebutton.New(GetCard(card, manager))
	}

	for key, card := range allCards {
		card.CallBack = callBackGenerator(ui, key)
	}

	//for i := len(cardNames) - 1; i >= 0; i-- {
	//	eventManager.RegisterEvent(allCards[cardNames[i]])
	//}
	for _, card := range cardNames {
		eventManager.RegisterEvent(allCards[card])
	}

	// Init play card button
	font, _ := fontManager.GetFont("universalfruitcake", 20)
	playButton = rectbutton.New("Play", 200, 100, utils.GREEN, font)
	playButton.CallBack = func(inter ...interface{}) error {
		if ui.selectedCard == "" {
			return nil
		}
		ui.PlayedCards[ui.DevicePlayer.Direction] = ui.selectedCard
		ui.removeSelectedCard()
		return nil
	}
	eventManager.RegisterEvent(playButton)

	// Init claim button
	claimButton = rectbutton.New("Claim", 200, 100, utils.GREEN, font)
	claimButton.CallBack = func(i ...interface{}) error {
		ui.claimedHands++
		return nil
	}
	eventManager.RegisterEvent(claimButton)

	// Init Player Icons
	playerIcon = rectbutton.New("", 150, 150, utils.SILVER, font)
	playerIcon.CallBack = func(i ...interface{}) error {
		return nil
	}

	// Init claimed hands text
	claimedHandsText = rectbutton.New("Claimed: " + strconv.Itoa(ui.claimedHands), 150, 50, &sdl.Color{168, 235, 254, 255}, font)
}

func New(devicePlayer *interfaces.Player, context interfaces.GameContext) *GameUiManager {
	ui := GameUiManager{
		GameId:        context.GameId,
		Players:       context.Players,
		Host:          context.Host,
		DevicePlayer:  devicePlayer,
		CurrentPlayer: nil,
		Cards:         nil,
		PlayedCards:   []string{"", "", "", ""},
		DeviceTurn:    false,
		GameStarted:   false,
		selectedCard:  "",
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

func (ui *GameUiManager) AssignCards(cards []string) {
	ui.Cards = cards
}

func (ui *GameUiManager) Draw(
	winWidth, winHeight int32,
	renderer *sdl.Renderer,
) error {

	_, firstCardY, err := ui.drawCardRack(winWidth, winHeight, renderer)
	if err != nil {
		return err
	}

	if ui.CurrentPlayer == ui.DevicePlayer {
		err = ui.drawPlayButton(firstCardY-100, renderer)
		if err != nil {
			return err
		}
	}

	err = ui.drawClaimButton(winWidth, firstCardY-100, renderer)
	if err != nil {
		return err
	}

	err = ui.drawOpponentsAndPlayedCards(winWidth, winHeight, renderer)
	if err != nil {
		return err
	}

	err = ui.drawClaimedHands(renderer)

	return nil
}

func (ui *GameUiManager) drawCardRack(w, h int32, renderer *sdl.Renderer) (int32, int32, error) {

	if len(ui.Cards) == 0 {
		return 0, cardYPosition, nil
	}

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
		return 0, 0, err
	}

	imageW := allCards[ui.Cards[0]].Width

	intervals := generateCenteredIntervals(w, imageW, len(ui.Cards), 45)

	for i, e := range intervals {
		if ui.Cards[i] == ui.selectedCard {
			err = allCards[ui.Cards[i]].Draw(e, h-rectHeight-100, renderer)
		} else {
			err = allCards[ui.Cards[i]].Draw(e, h-rectHeight, renderer)
		}

		if err != nil {
			return 0, 0, err
		}
	}

	// Explicitly draw the cards that are not used off the screen as they were invisibly taking on
	// default values at the top edge
	shownCards := make(map[string]bool)
	for _, card := range ui.Cards {
		shownCards[card] = true
	}
	for key, value := range allCards {
		if !shownCards[key] {
			err := value.Draw(w, h, renderer)
			if err != nil {
				return 0, 0, nil
			}
		}
	}

	cardYPosition = h - rectHeight

	return intervals[0], cardYPosition, nil
}

func (ui *GameUiManager) drawPlayButton(firstCardY int32, renderer *sdl.Renderer) error {

	if ui.selectedCard == "" {
		playButton.Color = utils.SILVER
	} else {
		playButton.Color = utils.GREEN
	}
	err := playButton.Draw(20, firstCardY-125, renderer)
	if err != nil {
		return err
	}

	return nil
}

func (ui *GameUiManager) drawOpponentsAndPlayedCards(w, h int32, renderer *sdl.Renderer) error {
	numDirections := len(utils.DirectionOrder)

	var startIndex int
	for startIndex = 0; startIndex < numDirections; startIndex++ {
		if ui.DevicePlayer.Direction == utils.DirectionOrder[startIndex] {
			break
		}
	}

	localOrder := []int{utils.DirectionOrder[(startIndex+1)%numDirections],
		utils.DirectionOrder[(startIndex+2)%numDirections],
		utils.DirectionOrder[(startIndex+3)%numDirections]}

	// Draw Left Player
	for count, orderElem := range localOrder {
		for player, _ := range ui.Players {
			if player.Direction == orderElem {
				switch count {
				case 0:
					// Draw Left
					x, y := int32(15), h/2-50
					imageX, imageY := x+playerIcon.Width+15, y-50
					err := ui.drawPlayerIcon(player, x, y, renderer)
					if err != nil {
						return err
					}
					err = ui.drawPlayedCard(player, imageX, imageY, renderer)
					if err != nil {
						return err
					}
				case 1:
					// Draw Top
					x, y := w/2-playerIcon.Width/2, int32(150)
					leftX, leftY := int32(15), h/2-50
					imageX, imageY := leftX+playerIcon.Width+110, leftY-200
					err := ui.drawPlayerIcon(player, x, y, renderer)
					if err != nil {
						return err
					}
					err = ui.drawPlayedCard(player, imageX, imageY, renderer)
					if err != nil {
						return err
					}
				case 2:
					// Draw Right
					x, y := w-playerIcon.Width-15, h/2-50
					imageX, imageY := x-allCards["c1"].Width-15, y-50
					err := ui.drawPlayerIcon(player, x, y, renderer)
					if err != nil {
						return err
					}
					err = ui.drawPlayedCard(player, imageX, imageY, renderer)
					if err != nil {
						return err
					}
				}
			}
		}
	}

	// Draw player played card if any
	if playedCard := ui.PlayedCards[ui.DevicePlayer.Direction]; playedCard != "" {
		leftX, leftY := int32(15), h/2-50
		imageX, imageY := leftX+playerIcon.Width+110, leftY+playerIcon.Height/2
		err := ui.drawPlayedCard(ui.DevicePlayer, imageX, imageY, renderer)
		if err != nil {
			return err
		}
	}

	return nil
}

func (ui *GameUiManager) drawPlayerIcon(player *interfaces.Player, x, y int32, renderer *sdl.Renderer) error {
	playerIcon.BtnText = utils.DirectionToString(player.Direction)
	if player == ui.CurrentPlayer {
		playerIcon.Color = utils.BRIGHT_GREEN
	} else {
		playerIcon.Color = utils.SILVER
	}
	err := playerIcon.Draw(x, y, renderer)
	if err != nil {
		return err
	}

	return nil
}

func (ui *GameUiManager) drawPlayedCard(player *interfaces.Player, imageX int32, imageY int32, renderer *sdl.Renderer) error {
	if playedCard := ui.PlayedCards[player.Direction]; playedCard != "" {
		err := allCards[playedCard].Draw(imageX, imageY, renderer)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ui *GameUiManager) drawClaimButton(winWidth, cardY int32, renderer *sdl.Renderer) error {
	err := claimButton.Draw(winWidth-claimButton.Width-20, cardY-125, renderer)
	if err != nil {
		return err
	}

	return nil
}

func (ui *GameUiManager) drawClaimedHands(renderer *sdl.Renderer) error {
	claimedHandsText.BtnText = "Claimed: " + strconv.Itoa(ui.claimedHands)
	err := claimedHandsText.Draw(50, 50, renderer)
	if err != nil {
		return err
	}

	return nil
}

func (ui *GameUiManager) removeSelectedCard() {
	newCards := make([]string, len(ui.Cards) - 1)
	index := 0
	for _, e := range ui.Cards {
		if e != ui.selectedCard {
			newCards[index] = e
			index++
		}
	}
	ui.Cards = newCards
}

func GetCard(card string, manager *imgmanager.ImageManager) *sdl.Texture {
	return manager.Images["cards/fronts/"+card]
}

func generateCenteredIntervals(width, cardWidth int32, count int, delta int32) []int32 {
	cardSpace := (int32(count)-1)*delta + cardWidth
	start := (width - cardSpace) / 2
	return generateIntervals(start, count, delta)
}

func generateIntervals(start int32, count int, delta int32) []int32 {
	interval := make([]int32, count, count)
	interval[0] = start
	for i := 1; i < count; i++ {
		interval[i] = interval[i-1] + delta
	}
	return interval
}
