package spriteset

import (
	"github.com/jrcichra/gollercoaster/sprite"
	"github.com/jrcichra/gollercoaster/textureloader"
)

//SpriteSet - collection of sprites that can be stacked in a tile
type SpriteSet struct {
	TallWall *sprite.Sprite
	Floor    *sprite.Sprite
	BigTable *sprite.Sprite
	// SmallTable *pixel.Sprite
}

//Load - loads all the sprites defined in the SpriteSet
func (s *SpriteSet) Load() {
	var t textureloader.TextureLoader
	err := t.Open("castle.png")
	if err != nil {
		panic(err)
	}
	var tallWall sprite.Sprite
	tallWall.Sprite = t.GetTexture(0, 448, 64, 512)
	s.TallWall = &tallWall
	var floor sprite.Sprite
	floor.Sprite = t.GetTexture(0, 128, 64, 192)
	s.Floor = &floor
	var bigTable sprite.Sprite
	bigTable.Sprite = t.GetTexture(64*3, 128, 64*2, 192)
	s.BigTable = &bigTable
}
