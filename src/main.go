package main

import "C"

import (
	"CardGameGo/src/components/buttons/imagebutton"
	"CardGameGo/src/components/buttons/rectbutton"
	"CardGameGo/src/engine"
	"CardGameGo/src/managers/gamemanager"
	"CardGameGo/src/managers/interfaces"
	"CardGameGo/src/screens"
	"CardGameGo/src/utils"
	"errors"
	"fmt"
	"runtime"

	"github.com/veandco/go-sdl2/sdl"
)

var gameUi *gamemanager.GameUiManager
var startNewGame = true

func Draw(e *engine.Engine, screen int, args ...interface{}) error {

	switch screen {
	case screens.MainScreen:
		return drawMainScreen(e)
	case screens.GameScreen:
		return drawGameScreen(e, args)
	case screens.SettingsScreen:
		return drawSettingsScreen(e, args)
	default:
		return errors.New("draw error: unexpected error occurred")
	}
}

func drawMainScreen(e *engine.Engine) error {
	_ = e.Renderer.Clear()
	_ = e.Renderer.SetDrawColor(66, 152, 66, 1)
	_ = e.Renderer.FillRect(nil)

	//Insert Card Image
	image := e.Image.Images["cardicon"]
	w, h := e.Window.GetSize()
	err := e.Renderer.Copy(image, nil, utils.CenterTexture(image, w, h/2))
	if err != nil {
		return err
	}

	// Insert New Game Button
	color := utils.GRAY
	font, _ := e.Font.GetFont("universalfruitcake", 20)
	newGameButton := rectbutton.New("New Game", 350, 75, color, font)
	cenX, newGameButtonY := utils.GetCenterCoordinates(newGameButton.Width, newGameButton.Height, w, h)
	err = newGameButton.Draw(cenX, newGameButtonY, e.Renderer)
	if err != nil {
		return err
	}
	newGameButton.CallBack = func(...interface{}) error {
		e.CurrentScreen = screens.GameScreen
		return nil
	}
	e.Event[e.CurrentScreen].RegisterEvent(newGameButton)

	// Insert Settings Button
	settingsButton := rectbutton.New("Settings Button", 350, 75, color, font)
	err = settingsButton.Draw(cenX, newGameButtonY+100, e.Renderer)
	if err != nil {
		return err
	}
	settingsButton.CallBack = func(...interface{}) error {
		e.CurrentScreen = screens.SettingsScreen
		return nil
	}
	e.Event[e.CurrentScreen].RegisterEvent(settingsButton)
	return nil
}

func drawGameScreen(e *engine.Engine, args []interface{}) error {
	w, h := e.Window.GetSize()
	_ = e.Renderer.Clear()

	// Background
	_ = e.Renderer.SetDrawColor(168, 235, 254, 255)
	_ = e.Renderer.FillRect(nil)

	// Home Button
	image := e.Image.Images["home"]
	_, _, imageW, imageH, _ := image.Query()
	homeButton := imagebutton.New(image)
	err := homeButton.Draw(w-imageW-10, imageH, e.Renderer)
	if err != nil {
		return err
	}
	homeButton.CallBack = func(i ...interface{}) error {
		e.CurrentScreen = screens.MainScreen
		return nil
	}
	e.Event[e.CurrentScreen].RegisterEvent(homeButton)

	//Draw Card Game Rack
	hostPlayer := &interfaces.Player{Direction: utils.East}
	player1 := &interfaces.Player{Direction: utils.North}
	player2 := &interfaces.Player{Direction: utils.South}
	player3 := &interfaces.Player{Direction: utils.West}
	players := map[*interfaces.Player]bool{
		hostPlayer: true,
		player1: true,
		player2: true,
		player3: true,
	}

	dummyContext := interfaces.GameContext{
		GameId:  "hello world",
		Players: players,
		Host:    hostPlayer,
	}

	if startNewGame {
		gameUi = gamemanager.New(hostPlayer, dummyContext)
		err := gameUi.SetCurrentPlayer(hostPlayer)
		if err != nil {
			return err
		}
		gameUi.Init(e.Image, e.Event[screens.GameScreen], e.Font, e.Renderer)
		startNewGame = false
		cards := []string{"c1", "c2", "c3", "c4", "c5", "c6", "c7", "c8", "c9", "cX", "cJ", "cQ", "cK"}
		gameUi.AssignCards(cards)
	}

	err = gameUi.Draw(w, h, e.Renderer)
	if err != nil {
		return err
	}

	return nil
}

func drawSettingsScreen(e *engine.Engine, args []interface{}) error {
	w, _ := e.Window.GetSize()

	_ = e.Renderer.Clear()
	_ = e.Renderer.SetDrawColor(255, 250, 205, 255)
	_ = e.Renderer.FillRect(nil)

	image := e.Image.Images["home"]
	_, _, imageW, imageH, _ := image.Query()
	homeButton := imagebutton.New(image)
	err := homeButton.Draw(w-imageW-10, imageH, e.Renderer)
	if err != nil {
		return err
	}
	homeButton.CallBack = func(i ...interface{}) error {
		e.CurrentScreen = screens.MainScreen
		return nil
	}

	e.Event[e.CurrentScreen].RegisterEvent(homeButton)
	return nil
}

//export SDL_main
func SDL_main() {
	runtime.LockOSThread()
	//e := engine.New("Go SDL2", 480, 800)
	e := engine.New("Go SDL2", 720, 1280)

	err := e.Init()
	if err != nil {
		sdl.LogError(sdl.LOG_CATEGORY_APPLICATION, "Init: %s\n", err)
	}
	defer e.Destroy()

	e.Load()
	defer e.Unload()

	for e.Running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				e.Quit()

			case *sdl.MouseButtonEvent:
				if t.Type == sdl.MOUSEBUTTONDOWN && t.Button == sdl.BUTTON_LEFT {
					err := e.Event[e.CurrentScreen].ProcessClickEvents(t)
					if err != nil {
						fmt.Printf("ignoring event %q: %d\n", err, t.Timestamp)
					}
				}

			case *sdl.KeyboardEvent:
				if t.Keysym.Scancode == sdl.SCANCODE_ESCAPE || t.Keysym.Scancode == sdl.SCANCODE_AC_BACK {
					e.Quit()
				}
			}
		}

		err = Draw(e, e.CurrentScreen)
		if err != nil {
			fmt.Println(err)
			return
		}

		e.Renderer.Present()
		sdl.Delay(50)
	}
}

func main() {
	SDL_main()
}
