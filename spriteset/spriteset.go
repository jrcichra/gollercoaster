package spriteset

import (
	"github.com/faiface/pixel"
	"github.com/jrcichra/gollercoaster/sprite"
	"github.com/jrcichra/gollercoaster/textureloader"
)

const (
	TallWall      = iota
	Floor         = iota
	SmallTable    = iota
	LeftAngleRoof = iota
)

//SpriteSet - collection of sprites that can be stacked in a tile
type SpriteSet struct {
	TallWall      *sprite.Sprite
	Floor         *sprite.Sprite
	SmallTable    *sprite.Sprite
	LeftAngleRoof *sprite.Sprite
	DoubleWall    *sprite.Sprite
}

//Load - loads all the sprites defined in the SpriteSet
func (s *SpriteSet) Load() *pixel.Batch {
	var t textureloader.TextureLoader
	batch, err := t.Open("textures.png")
	if err != nil {
		panic(err)
	}
	var tallWall sprite.Sprite
	tallWall.Sprite = t.GetTexture(0, 0)
	s.TallWall = &tallWall
	var floor sprite.Sprite
	floor.Sprite = t.GetTexture(0, 5)
	s.Floor = &floor
	var smallTable sprite.Sprite
	smallTable.Sprite = t.GetTexture(3, 5)
	s.SmallTable = &smallTable
	var leftAngleRoof sprite.Sprite
	leftAngleRoof.Sprite = t.GetTexture(1, 1)
	s.LeftAngleRoof = &leftAngleRoof
	var doubleWall sprite.Sprite
	doubleWall.Sprite = t.GetTexture(0, 7)
	s.DoubleWall = &doubleWall
	//Return the batch we should write to
	return batch
}
