package main

import "C"

import (
	"CardGameGo/src/managers/fontmanager"
	"CardGameGo/src/managers/imgmanager"
	"CardGameGo/src/components/buttons/rectbutton"
	"CardGameGo/src/utils"
	"errors"
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const (
	winTitle  = "Go SDL2"
	winWidth  = 480
	winHeight = 800
)

//noinspection GoSnakeCaseUsage
const (
	MAIN_SCREEN = iota
	GAME_SCREEN
	SETTINGS_SCREEN
)

// Engine represents SDL engine.
type Engine struct {
	State    int
	Window   *sdl.Window
	Renderer *sdl.Renderer
	Image    *imgmanager.ImageManager
	Font     *fontmanager.FontManager
	Music    *mix.Music
	Sound    *mix.Chunk

	running   bool
}

// NewEngine returns new engine.
func NewEngine() (e *Engine) {
	e = &Engine{}
	e.running = true
	return
}

func Draw(e *Engine, screen int, args ...interface{}) error {
	switch screen {
	case MAIN_SCREEN:
		return drawMainScreen(e)
	case GAME_SCREEN:
		return drawGameSceen(e, args)
	case SETTINGS_SCREEN:
		return drawSettingsScreen(e, args)
	default:
		return errors.New("draw error: unexpected error occurred")
	}
}

func drawMainScreen(e *Engine) error {
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

	return nil
}

func drawGameSceen(e *Engine, args []interface{}) error {
	return nil
}

func drawSettingsScreen(e *Engine, args []interface{}) error {
	return nil
}

// Init initializes SDL.
func (e *Engine) Init() (err error) {
	err = sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		return
	}

	e.Font, err = fontmanager.New()
	if err != nil {
		return
	}

	err = mix.Init(mix.INIT_MP3)
	if err != nil {
		return
	}

	err = mix.OpenAudio(mix.DEFAULT_FREQUENCY, mix.DEFAULT_FORMAT, mix.DEFAULT_CHANNELS, 3072)
	if err != nil {
		return
	}

	e.Window, err = sdl.CreateWindow(winTitle, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, winWidth, winHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		return
	}

	e.Renderer, err = sdl.CreateRenderer(e.Window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		return
	}

	e.Image, err = imgmanager.New(e.Renderer)
	if err != nil {
		return
	}

	return nil
}

// Destroy destroys SDL and releases the memory.
func (e *Engine) Destroy() {
	e.Renderer.Destroy()
	e.Window.Destroy()
	mix.CloseAudio()

	img.Quit()
	mix.Quit()
	ttf.Quit()
	sdl.Quit()
}

// Running checks if loop is running.
func (e *Engine) Running() bool {
	return e.running
}

// Quit exits main loop.
func (e *Engine) Quit() {
	e.running = false
}

// Load loads resources.
func (e *Engine) Load() {
	assetDir := ""
	if runtime.GOOS != "android" {
		assetDir = filepath.Join("assets")
	}

	var err error

	err = e.Font.Load()
	if err != nil {
		sdl.LogError(sdl.LOG_CATEGORY_APPLICATION, "load font error: %s\n", err)
	}

	err = e.Image.Load()
	if err != nil {
		sdl.LogError(sdl.LOG_CATEGORY_APPLICATION, "load image error: %s\n", err)
	}

	e.Music, err = mix.LoadMUS(filepath.Join(assetDir, "music", "frantic-gameplay.mp3"))
	if err != nil {
		sdl.LogError(sdl.LOG_CATEGORY_APPLICATION, "LoadMUS: %s\n", err)
	}

	e.Sound, err = mix.LoadWAV(filepath.Join(assetDir, "sounds", "click.wav"))
	if err != nil {
		sdl.LogError(sdl.LOG_CATEGORY_APPLICATION, "LoadWAV: %s\n", err)
	}
}

// Unload unloads resources.
func (e *Engine) Unload() {

	//e.Sprite.Destroy()
	e.Font.Close()
	e.Music.Free()
	e.Sound.Free()
}

//export SDL_main
func SDL_main() {
	runtime.LockOSThread()
	e := NewEngine()

	err := e.Init()
	if err != nil {
		sdl.LogError(sdl.LOG_CATEGORY_APPLICATION, "Init: %s\n", err)
	}
	defer e.Destroy()

	e.Load()
	defer e.Unload()

	for e.Running() {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				e.Quit()

			case *sdl.MouseButtonEvent:
				if t.Type == sdl.MOUSEBUTTONDOWN && t.Button == sdl.BUTTON_LEFT {

				}

			case *sdl.KeyboardEvent:
				if t.Keysym.Scancode == sdl.SCANCODE_ESCAPE || t.Keysym.Scancode == sdl.SCANCODE_AC_BACK {
					e.Quit()
				}
			}
		}

		err = Draw(e, MAIN_SCREEN)
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
