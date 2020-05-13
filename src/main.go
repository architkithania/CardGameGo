package main

import "C"

import (
	"CardGameGo/src/components/buttons/rectbutton"
	"CardGameGo/src/engine"
	"CardGameGo/src/screens"
	"CardGameGo/src/utils"
	"errors"
	"fmt"
	"runtime"

	"github.com/veandco/go-sdl2/sdl"
)

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

	// Insert Button
	button := rectbutton.New("New Game", 350, 75, &sdl.Color{255, 0, 0, 0})
	cenX, cenY := utils.GetCenterCoordinates(button.Width, button.Height, w, h)
	font, _ := e.Font.GetFont("universalfruitcake", 20)
	err = button.Draw(cenX, cenY, font, e.Renderer)
	if err != nil {
		return err
	}

	button.CallBack = func(...interface{}) error {
		e.CurrentScreen = screens.GameScreen
		return nil
	}

	e.Event[e.CurrentScreen].RegisterEvent(button)

	return nil
}

func drawGameScreen(e *engine.Engine, args []interface{}) error {
	_ = e.Renderer.Clear()
	_ = e.Renderer.SetDrawColor(168, 235, 254, 255)
	_ = e.Renderer.FillRect(nil)

	return nil
}

func drawSettingsScreen(e *engine.Engine, args []interface{}) error {
	return nil
}

//export SDL_main
func SDL_main() {
	runtime.LockOSThread()
	e := engine.New("Go SDL2", 480, 800)

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
