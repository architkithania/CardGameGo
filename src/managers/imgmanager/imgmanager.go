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

	"cards/fronts/c1.png",
	"cards/fronts/c1.png",
	"cards/fronts/c2.png",
	"cards/fronts/c3.png",
	"cards/fronts/c4.png",
	"cards/fronts/c5.png",
	"cards/fronts/c6.png",
	"cards/fronts/c7.png",
	"cards/fronts/c8.png",
	"cards/fronts/c9.png",
	"cards/fronts/cX.png",
	"cards/fronts/cJ.png",
	"cards/fronts/cQ.png",
	"cards/fronts/cK.png",
	"cards/fronts/h1.png",
	"cards/fronts/h2.png",
	"cards/fronts/h3.png",
	"cards/fronts/h4.png",
	"cards/fronts/h5.png",
	"cards/fronts/h6.png",
	"cards/fronts/h7.png",
	"cards/fronts/h8.png",
	"cards/fronts/h9.png",
	"cards/fronts/hX.png",
	"cards/fronts/hJ.png",
	"cards/fronts/hQ.png",
	"cards/fronts/hK.png",
	"cards/fronts/d1.png",
	"cards/fronts/d2.png",
	"cards/fronts/d3.png",
	"cards/fronts/d4.png",
	"cards/fronts/d5.png",
	"cards/fronts/d6.png",
	"cards/fronts/d7.png",
	"cards/fronts/d8.png",
	"cards/fronts/d9.png",
	"cards/fronts/dX.png",
	"cards/fronts/dJ.png",
	"cards/fronts/dQ.png",
	"cards/fronts/dK.png",
	"cards/fronts/s1.png",
	"cards/fronts/s2.png",
	"cards/fronts/s3.png",
	"cards/fronts/s4.png",
	"cards/fronts/s5.png",
	"cards/fronts/s6.png",
	"cards/fronts/s7.png",
	"cards/fronts/s8.png",
	"cards/fronts/s9.png",
	"cards/fronts/sX.png",
	"cards/fronts/sJ.png",
	"cards/fronts/sQ.png",
	"cards/fronts/sK.png",
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
