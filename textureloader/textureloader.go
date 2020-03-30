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
func (t *TextureLoader) GetTexture(minX, minY, maxX, maxY float64) *pixel.Sprite {
	return pixel.NewSprite(t.picture, pixel.R(minX, minY, maxX, maxY))
}

//Close - close the file for this texture loader
func (t *TextureLoader) Close() error {
	return t.file.Close()
}
