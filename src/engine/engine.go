package engine

import (
	"CardGameGo/src/managers/eventmanager"
	"CardGameGo/src/managers/fontmanager"
	"CardGameGo/src/managers/imgmanager"
	"CardGameGo/src/screens"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
	"path/filepath"
	"runtime"
)

var (
	winTitle string
	winWidth int32
	winHeight int32
)

type Engine struct {
	State    int
	Window   *sdl.Window
	Renderer *sdl.Renderer
	Image    *imgmanager.ImageManager
	Font     *fontmanager.FontManager
	Event    map[int]*eventmanager.EventManager
	Music    *mix.Music
	Sound    *mix.Chunk

	CurrentScreen int
	Running       bool
}

// NewEngine returns new engine.
func New(title string, width, height int32) (e *Engine) {
	winTitle = title
	winWidth = width
	winHeight = height

	e = &Engine{}
	e.Running = true
	return
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

	e.Event = make(map[int]*eventmanager.EventManager)
	for _, screen := range screens.Screens {
		e.Event[screen] = eventmanager.New(screen)
	}

	e.Image, err = imgmanager.New(e.Renderer)
	if err != nil {
		return
	}

	e.CurrentScreen = screens.MainScreen

	return nil
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

// Unload unloads resources.
func (e *Engine) Unload() {

	//e.Sprite.Destroy()
	e.Font.Close()
	e.Music.Free()
	e.Sound.Free()
}

// Quit exits main loop.
func (e *Engine) Quit() {
	e.Running = false
}

