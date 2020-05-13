package imgmanager

import (
	"errors"
	"fmt"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"path/filepath"
	"runtime"
	"strings"
)

var LOADED_IMAGES = []string{
	"cardicon.png",
	"home.png",
}

type ImageManager struct{
	Images map[string]*sdl.Texture
	renderer *sdl.Renderer
}

func New(renderer *sdl.Renderer) (*ImageManager, error) {
	err := img.Init(img.INIT_PNG)
	if err != nil {
		return nil, err
	}

	imgManager := ImageManager{make(map[string]*sdl.Texture), renderer}

	return &imgManager, nil
}


func (i *ImageManager) Load() error {
	assetDir := ""
	if runtime.GOOS != "android" {
		assetDir = filepath.Join( "assets")
	}

	var err error
	for _, image := range LOADED_IMAGES {
		imageName := strings.Split(image, ".")[0]
		i.Images[imageName], err = img.LoadTexture(i.renderer, filepath.Join(assetDir, "images", image))
		if err != nil {
			return errors.New(fmt.Sprintf("image manager error: %q couldn't be loaded", image))
		}
	}

	return nil
}

func (i *ImageManager) Close() {
	for _, texture := range i.Images {
		_ = texture.Destroy()
	}
}
