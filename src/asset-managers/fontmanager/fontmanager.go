package fontmanager

import (
	"fmt"
	"github.com/veandco/go-sdl2/ttf"
	"path/filepath"
	"runtime"
	"strings"
)

var PRE_LOADED_FONTS map[string]int = map[string]int{
	"universalfruitcake.ttf": 24,
}

type FontManager struct {
	fonts map[struct{string;int}]*ttf.Font
}

func New() (*FontManager, error) {
	err := ttf.Init()
	if err != nil {
		return nil, err
	}
	fManager := FontManager{make(map[struct{string;int}]*ttf.Font)}

	return &fManager, nil
}

func (fManager *FontManager) GetFont(font string, size int) (*ttf.Font, bool) {
	assetDir := ""
	if runtime.GOOS != "android" {
		assetDir = filepath.Join( "assets")
	}

	if val, ok := fManager.fonts[struct{string;int}{font,size}]; ok {
		return val, true
	}

	var err error
	fontPack, err := ttf.OpenFont(filepath.Join(assetDir, "fonts", font + ".ttf"), size)

	if err != nil {
		fmt.Println(err)
		return fManager.fonts[struct{string;int}{"universalfruitcake", 24}], false
	}

	fManager.fonts[struct{string;int}{font,size}] = fontPack
	return fontPack, true
}

func (fManager *FontManager) Load() error {
	assetDir := ""
	if runtime.GOOS != "android" {
		assetDir = filepath.Join( "assets")
	}

	var err error
	var key struct{string;int}
	for font, size := range PRE_LOADED_FONTS {
		fontName:= strings.Split(font, ".")[0]
		key = struct{string;int}{fontName,size}
		fManager.fonts[key], err = ttf.OpenFont(filepath.Join(assetDir,"fonts", font), size)
		if err != nil {
			return err
		}
	}

	return nil
}

func (fManager *FontManager) Close() {
	for _, font := range fManager.fonts {
		font.Close()
	}
}