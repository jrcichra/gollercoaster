package game

import (
	"fmt"
	"image/color"
	"math"
	"time"

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
	CamPos       pixel.Vec
	CamSpeed     float64
	CamZoom      float64
	CamZoomSpeed float64
}

//Run - run the game
func (g *Game) Run() {
	var err error
	g.Name = "gollercoaster"
	g.windowWidth = 1280
	g.windowHeight = 720
	g.offset = pixel.V(-400*14.5, -325*20)
	g.tileSize = 64
	g.CamPos = pixel.ZV
	g.CamSpeed = 500
	g.CamZoom = 1
	g.CamZoomSpeed = 1.2
	frames := 0
	second := time.Tick(time.Second)
	// last := time.Now()

	pixelgl.Run(func() {
		cfg := pixelgl.WindowConfig{
			Title:  g.Name,
			Bounds: pixel.R(0, 0, g.windowWidth, g.windowHeight),
			// VSync:  true,
		}
		g.window, err = pixelgl.NewWindow(cfg)
		if err != nil {
			panic(err)
		}
		//Create a new level
		var l level.Level
		g.currentLevel = &l
		batch := l.Spawn()
		batch.Clear()
		g.render()
		batch.Draw(g.window)
		// g.render()
		last := time.Now()
		for !g.window.Closed() {
			dt := time.Since(last).Seconds()
			last = time.Now()
			g.window.Clear(color.Black)
			batch.Clear()
			if g.window.Pressed(pixelgl.KeyLeft) || g.window.Pressed(pixelgl.KeyA) {
				g.CamPos.X -= g.CamSpeed * dt / g.CamZoom
			}
			if g.window.Pressed(pixelgl.KeyRight) || g.window.Pressed(pixelgl.KeyD) {
				g.CamPos.X += g.CamSpeed * dt / g.CamZoom
			}
			if g.window.Pressed(pixelgl.KeyDown) || g.window.Pressed(pixelgl.KeyS) {
				g.CamPos.Y -= g.CamSpeed * dt / g.CamZoom
			}
			if g.window.Pressed(pixelgl.KeyUp) || g.window.Pressed(pixelgl.KeyW) {
				g.CamPos.Y += g.CamSpeed * dt / g.CamZoom
			}

			//TODO: Trying to get scroll where your cursor is working. Works except for zoom compensation

			// if g.window.MouseScroll().Y > 0 {
			// 	g.CamPos = g.CamPos.Sub(g.window.Bounds().Center().Sub(g.window.MousePosition()))

			// } else if g.window.MouseScroll().Y < 0 {
			// 	g.CamPos = g.CamPos.Add(g.window.Bounds().Center().Sub(g.window.MousePosition()))

			// }

			g.CamZoom *= math.Pow(g.CamZoomSpeed, g.window.MouseScroll().Y)

			cam := pixel.IM.Scaled(g.CamPos, g.CamZoom).Moved(g.window.Bounds().Center().Sub(g.CamPos))

			g.window.SetMatrix(cam)
			g.render()
			batch.Draw(g.window)

			g.window.Update()
			frames++
			select {
			case <-second:
				g.window.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, frames))
				frames = 0
			default:
			}
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
