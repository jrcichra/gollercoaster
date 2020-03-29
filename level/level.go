package level

import (
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

	//Floor Tile
	var f tile.Tile
	f.Append(ss.Floor)

	//Table with a floor under it
	var t tile.Tile
	t.Append(ss.Floor)
	t.Append(ss.BigTable)

	//Wall
	var w tile.Tile
	w.Append(ss.TallWall)

	//Hardcoded level for now
	l.Level = [][]tile.Tile{
		{f, f, f, f, f, f, t}, // This row will be rendered in the lower left part of the screen (closer to the viewer)
		{w, f, f, f, f, w, w},
		{w, f, f, f, f, w, w},
		{w, f, f, f, f, w, w},
		{w, f, f, f, f, w, w},
		{w, w, w, w, w, w, w}, // And this in the upper right
	}
}
