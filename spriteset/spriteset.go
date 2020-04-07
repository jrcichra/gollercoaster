package spriteset

import (
	"github.com/jrcichra/gollercoaster/sprite"
	"github.com/jrcichra/gollercoaster/textureloader"
)

//SpriteSet - collection of sprites that can be stacked in a tile
type SpriteSet struct {
	TallWall      *sprite.Sprite
	Floor         *sprite.Sprite
	SmallTable    *sprite.Sprite
	LeftAngleRoof *sprite.Sprite
	Selected      *sprite.Sprite
}

//Load - loads all the sprites defined in the SpriteSet
func (s *SpriteSet) Load() {
	var t textureloader.TextureLoader
	err := t.Open("textures.png")
	if err != nil {
		panic(err)
	}
	s.TallWall = t.GetTexture(0, 0)
	s.Floor = t.GetTexture(0, 5)
	s.SmallTable = t.GetTexture(3, 5)
	s.LeftAngleRoof = t.GetTexture(1, 1)
	s.Selected = t.GetTexture(0, 7)
}
