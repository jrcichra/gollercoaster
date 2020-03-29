package sprite

import (
	"github.com/faiface/pixel"
)

//Sprite - something that comes from a file (ex. a texture) - a tile might have more than one sprite
type Sprite struct {
	Sprite *pixel.Sprite
}
