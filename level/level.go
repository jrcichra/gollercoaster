package level

import (
	"errors"
	"math/rand"
	"strconv"
	"time"

	"github.com/jrcichra/gollercoaster/spriteset"
	"github.com/jrcichra/gollercoaster/tile"
)

//Level - all the things a level has to store
type Level struct {
	Name  string         // Name of your level
	level [][]*tile.Tile //Map of the level (2D array of tiles)
	SS    *spriteset.SpriteSet
}

//GetTile - get the attributes of this tile
func (l *Level) GetTile(x, y int) (*tile.Tile, error) {
	if x >= 0 && y >= 0 && x < l.GetWidth() && y < l.GetHeight() {
		return l.level[x][y], nil
	}
	return nil, errors.New("Tile index out of range: x=" + strconv.Itoa(x) + " y=" + strconv.Itoa(y))
}

//GetWidth - returns the width of the level
func (l *Level) GetWidth() int {
	return len(l.level)
}

//GetHeight - returns the height of the level
func (l *Level) GetHeight() int {
	return len(l.level[0])
}

//Spawn - spawns a new level
func (l *Level) Spawn() {
	//Load in the sprite set
	var ss spriteset.SpriteSet
	l.SS = &ss
	ss.Load()

	// //Floor Tile
	// var f tile.Tile
	// f.Append(ss.Floor)

	// //Table with a floor under it
	// var t tile.Tile
	// t.Append(ss.Floor)
	// t.Append(ss.BigTable)

	// //Wall
	// var w tile.Tile
	// w.Append(ss.TallWall)

	//Hardcoded level for now
	// l.Level = [][]tile.Tile{
	// 	{f, f, f, f, f, f, t}, // This row will be rendered in the lower left part of the screen (closer to the viewer)
	// 	{w, f, f, f, f, w, w},
	// 	{w, f, f, f, f, w, w},
	// 	{w, f, f, f, f, w, w},
	// 	{w, f, f, f, f, w, w},
	// 	{w, w, w, w, w, w, w}, // And this in the upper right
	// }

	lvl := make([][]*tile.Tile, 0) //start with a blank level

	rand.Seed(time.Now().UTC().UnixNano())

	x := 200
	y := 200

	for i := 0; i < x; i++ {
		row := make([]*tile.Tile, 0, x)
		for j := 0; j < y; j++ {
			t := &tile.Tile{}
			val := rand.Float64()
			if val < .3 {
				t.Push(ss.TallWall)
			} else if val < .6 {
				t.Push(ss.Floor)
				t.Push(ss.SmallTable)
			} else {
				//It's the floor
				t.Push(ss.Floor)
			}
			row = append(row, t)
		}
		lvl = append(lvl, row)
	}
	l.level = lvl
}
