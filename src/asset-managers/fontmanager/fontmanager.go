package fontmanager

import (
	"github.com/veandco/go-sdl2/ttf"
	"path/filepath"
	"runtime"
	"strings"
)

var LOADED_FONTS map[string]int = map[string]int{
	"universalfruitcake.ttf": 24,
}

type FontManager struct {
	Fonts map[string]*ttf.Font
}

func New() (*FontManager, error) {
	err := ttf.Init()
	if err != nil {
		return nil, err
	}
	fManager := FontManager{make(map[string]*ttf.Font)}

	return &fManager, nil
}

func (fManager *FontManager) Load() error {
	assetDir := ""
	if runtime.GOOS != "android" {
		assetDir = filepath.Join( "assets")
	}

	var err error
	for font, size := range LOADED_FONTS {
		fontName:= strings.Split(font, ".")[0]
		fManager.Fonts[fontName], err = ttf.OpenFont(filepath.Join(assetDir,"fonts", font), size)
		if err != nil {
			return err
		}
	}

	return nil
}

func (fManager *FontManager) Close() {
	for _, font := range fManager.Fonts {
		font.Close()
	}
}