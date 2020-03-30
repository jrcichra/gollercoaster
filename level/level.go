package level

import (
	"math/rand"
	"time"

	perlin "github.com/aquilax/go-perlin"
	"github.com/jrcichra/gollercoaster/spriteset"
	"github.com/jrcichra/gollercoaster/tile"
)

//Level - all the things a level has to store
type Level struct {
	Name  string        // Name of your level
	Level [][]tile.Tile //Map of the level (2D array of tiles)
}

//Spawn - spawns a new level
func (l *Level) Spawn() {
	//Load in the sprite set
	var ss spriteset.SpriteSet
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

	lvl := make([][]tile.Tile, 0) //start with a blank level

	alpha := 2.
	beta := 2.
	n := 3
	rand.Seed(time.Now().UTC().UnixNano())
	var seed = rand.Int63n(100)
	p := perlin.NewPerlin(alpha, beta, n, seed)

	x := 20
	y := 20

	for i := 0; i < x; i++ {
		row := make([]tile.Tile, 0, x)
		for j := 0; j < y; j++ {
			var t tile.Tile
			val := p.Noise2D(float64(seed)/float64(i+1), float64(seed)/float64(j+1))
			if val < 0 {
				//We'll say it's a wall
				t.Append(ss.TallWall)
			} else {
				//It's the floor
				t.Append(ss.Floor)
			}
			row = append(row, t)
		}
		lvl = append(lvl, row)
	}
	l.Level = lvl

}
