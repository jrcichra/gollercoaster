package game

import (
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/jrcichra/gollercoaster/level"
)

//Game - Game Object
type Game struct {
	Name         string //Name of the game ("gollercoaster for now")
	windowWidth  float64
	windowHeight float64
	window       *pixelgl.Window
	offset       pixel.Vec
	tileSize     float64
	// tileMap      map[tile.Tile]int //map of tile objects that make up the game's tiles
	currentLevel *level.Level //pointer to the currentLevel to render
}

//Run - run the game
func (g *Game) Run() {
	var err error
	g.Name = "gollercoaster"
	g.windowWidth = 1920
	g.windowHeight = 1080
	g.offset = pixel.V(-800, -825)
	g.tileSize = 64

	pixelgl.Run(func() {
		cfg := pixelgl.WindowConfig{
			Title:  g.Name,
			Bounds: pixel.R(0, 0, g.windowWidth, g.windowHeight),
			VSync:  true,
		}
		g.window, err = pixelgl.NewWindow(cfg)
		if err != nil {
			panic(err)
		}

		//Create a new level
		var l level.Level
		l.Spawn()
		g.currentLevel = &l

		// g.render()

		for !g.window.Closed() {
			// l.Spawn()
			// time.Sleep(100 * time.Millisecond)
			g.window.Clear(color.Black)
			g.render()
			g.window.Update()
		}
	})

}

// Draw level data tiles to window, from farthest to closest.
// In order to achieve the depth effect, we need to render tiles up to down, being lower
// closer to the viewer (see painter's algorithm). To do that, we need to process levelData in reverse order,
// so its first row is rendered last, as OpenGL considers its origin to be in the lower left corner of the display.
func (g *Game) render() {
	for x := len(g.currentLevel.Level) - 1; x >= 0; x-- {
		for y := len(g.currentLevel.Level[x]) - 1; y >= 0; y-- {
			isoCoords := g.cartesianToIso(pixel.V(float64(x), float64(y)))
			mat := pixel.IM.Moved(g.offset.Add(isoCoords))
			// Not really needed, just put to show bigger blocks
			mat = mat.ScaledXY(g.window.Bounds().Center(), pixel.V(.1, .1))
			g.currentLevel.Level[x][y].Draw(g.window, mat)
		}
	}
}

func (g *Game) cartesianToIso(pt pixel.Vec) pixel.Vec {
	return pixel.V((pt.X-pt.Y)*(g.tileSize/2), (pt.X+pt.Y)*(g.tileSize/4))
}
