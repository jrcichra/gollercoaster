package sprite

import "github.com/hajimehoshi/ebiten"

//Sprite - something that comes from a file (ex. a texture) - a tile might have more than one sprite
type Sprite struct {
	Sprite *ebiten.Image
}
