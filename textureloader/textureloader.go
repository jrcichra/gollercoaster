package textureloader

import (
	"image"
	_ "image/png" //support for png inputs
	"os"

	"github.com/hajimehoshi/ebiten"
)

//TextureLoader - loads sprites from a file, just keeps the file managed
type TextureLoader struct {
	file    *os.File
	dump    *ebiten.Image
	picture *ebiten.Image
}

//Open - open's a file
func (t *TextureLoader) Open(path string) error {
	var err error
	t.file, err = os.Open(path)
	if err != nil {
		return err
	}
	img, _, err := image.Decode(t.file)
	if err != nil {
		return err
	}
	t.picture, err = ebiten.NewImageFromImage(img, ebiten.FilterDefault)
	return err
}

//GetTexture - returns a sprite from the coordinates on the texture
func (t *TextureLoader) GetTexture(x, y int) *ebiten.Image {
	//I personally want 0,0 to be in the top left corner, not the bottom left
	return t.picture.SubImage(image.Rect(x*64, ((y + 1) * 64), (x+1)*64, (y * 64))).(*ebiten.Image)
}

//Close - close the file for this texture loader
func (t *TextureLoader) Close() error {
	return t.file.Close()
}
