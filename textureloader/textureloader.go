package textureloader

import (
	"image"
	_ "image/png" //support for png inputs
	"os"

	"github.com/faiface/pixel"
)

//TextureLoader - loads sprites from a file, just keeps the file managed
type TextureLoader struct {
	file    *os.File
	dump    pixel.Picture
	picture pixel.Picture
}

//Open - open's a file
func (t *TextureLoader) Open(path string) (*pixel.Batch, error) {
	var err error
	t.file, err = os.Open(path)
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(t.file)
	if err != nil {
		return nil, err
	}
	t.picture = pixel.PictureDataFromImage(img)
	return pixel.NewBatch(&pixel.TrianglesData{}, t.picture), nil
}

//GetTexture - returns a sprite from the coordinates on the texture
func (t *TextureLoader) GetTexture(x, y int) *pixel.Sprite {
	//I personally want 0,0 to be in the top left corner, not the bottom left
	return pixel.NewSprite(t.picture, pixel.R(float64(x*64), float64(512-((y+1)*64)), float64((x+1)*64), float64(512-(y*64))))
}

//Close - close the file for this texture loader
func (t *TextureLoader) Close() error {
	return t.file.Close()
}
