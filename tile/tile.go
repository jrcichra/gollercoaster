package tile

import (
	"github.com/faiface/pixel"
	"github.com/jrcichra/gollercoaster/sprite"
)

//Tile - basic tile object (made up of many sprites)
type Tile struct {
	sprites []*sprite.Sprite
	Batch   *pixel.Batch //Batch this tile writes to
}

//Peek - look at the top sprite
func (t *Tile) Peek() *sprite.Sprite {
	return t.sprites[len(t.sprites)-1]
}

//Push - push a new sprite on this tile
func (t *Tile) Push(s *sprite.Sprite) {
	t.sprites = append(t.sprites, s)
}

//Clear - clear a sprite's elements (make a new array and let go garbage collect the old one)
func (t *Tile) Clear() {
	t.sprites = make([]*sprite.Sprite, 0)
}

//Pop - Pops the last sprite off the tile - unexpected drawing behavior when you pop off the last
func (t *Tile) Pop() *sprite.Sprite {
	if len(t.sprites) == 1 {
		return nil
	}
	s := t.sprites[len(t.sprites)-1]
	t.sprites = t.sprites[:len(t.sprites)-1]
	return s
}

//Draw - draws the tile
func (t *Tile) Draw(mat pixel.Matrix) int {
	//To draw a tile, you need to render the sprites in order from furthest to closest
	//Here I'm assuming 0 is furthest and N is closest
	c := 0
	for _, s := range t.sprites {
		s.Sprite.Draw(t.Batch, mat)
		c++
	}
	return c
}
