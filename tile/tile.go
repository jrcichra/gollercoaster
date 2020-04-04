package tile

import (
	"github.com/hajimehoshi/ebiten"
)

//Tile - basic tile object (made up of many sprites)
type Tile struct {
	sprites    []*ebiten.Image
	temppushes int
}

//TempPush - push a new sprite on but pop after a draw call
func (t *Tile) TempPush(s *ebiten.Image) {
	t.Push(s)
	t.temppushes++
}

//Push - push a new sprite on this tile
func (t *Tile) Push(s *ebiten.Image) {
	t.sprites = append(t.sprites, s)
}

//Clear - clear a sprite's elements (make a new array and let go garbage collect the old one)
func (t *Tile) Clear() {
	t.sprites = make([]*ebiten.Image, 0)
}

//Pop - Pops the last sprite off the tile - unexpected drawing behavior when you pop off the last
func (t *Tile) Pop() *ebiten.Image {
	if len(t.sprites) == 1 {
		return nil
	}
	s := t.sprites[len(t.sprites)-1]
	t.sprites = t.sprites[:len(t.sprites)-1]
	return s
}

//Draw - draws the tile
func (t *Tile) Draw(win *ebiten.Image, options *ebiten.DrawImageOptions) {
	//To draw a tile, you need to render the sprites in order from furthest to closest
	//Here I'm assuming 0 is furthest and N is closest
	for _, s := range t.sprites {
		win.DrawImage(s, options)
	}

	//Pop off for every temp push
	for i := 0; i < t.temppushes; i++ {
		t.Pop()
		t.temppushes--
	}

}
