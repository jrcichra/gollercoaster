package tile

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/jrcichra/gollercoaster/sprite"
)

//Tile - basic tile object (made up of many sprites)
type Tile struct {
	sprites []*sprite.Sprite
}

//Append - put a new sprite on this tile (on top)
func (t *Tile) Append(s *sprite.Sprite) {
	t.sprites = append(t.sprites, s)
}

//Draw - draws the tile
func (t *Tile) Draw(win *pixelgl.Window, mat pixel.Matrix) {
	//To draw a tile, you need to render the sprites in order from furthest to closest
	//Here I'm assuming 0 is furthest and N is closest
	for _, s := range t.sprites {
		s.Sprite.Draw(win, mat)
	}
}
