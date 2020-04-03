package game

import (
	"fmt"
	"image/color"
	"math"
	"time"

	"github.com/jrcichra/gollercoaster/music"

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
	currentLevel *level.Level //pointer to the currentLevel to render
	CamPos       pixel.Vec
	CamSpeed     float64
	CamZoom      float64
	CamZoomSpeed float64
	music        music.Music
}

func (g *Game) playMusic() {
	g.music.LoadRandom()
	g.music.Play()
	g.music.Close()
}

//Run - run the game
func (g *Game) Run() {
	var err error
	g.Name = "gollercoaster"
	g.windowWidth = 1280
	g.windowHeight = 720
	g.offset = pixel.ZV
	g.tileSize = 64
	g.CamPos = pixel.ZV
	g.CamSpeed = 500
	g.CamZoom = 1
	g.CamZoomSpeed = 1.2
	g.music = music.Music{}
	frames := 0
	second := time.Tick(time.Second)

	go g.playMusic()

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
		// g.render()
		last := time.Now()
		redraw := true
		g.window.Clear(color.Black)
		batch.Clear()
		g.renderAll()
		batch.Draw(g.window)
		for !g.window.Closed() {
			dt := time.Since(last).Seconds()
			last = time.Now()

			if g.window.Pressed(pixelgl.KeyLeft) || g.window.Pressed(pixelgl.KeyA) {
				g.CamPos.X -= g.CamSpeed * dt / g.CamZoom
				redraw = true
			}
			if g.window.Pressed(pixelgl.KeyRight) || g.window.Pressed(pixelgl.KeyD) {
				g.CamPos.X += g.CamSpeed * dt / g.CamZoom
				redraw = true
			}
			if g.window.Pressed(pixelgl.KeyDown) || g.window.Pressed(pixelgl.KeyS) {
				g.CamPos.Y -= g.CamSpeed * dt / g.CamZoom
				redraw = true
			}
			if g.window.Pressed(pixelgl.KeyUp) || g.window.Pressed(pixelgl.KeyW) {
				g.CamPos.Y += g.CamSpeed * dt / g.CamZoom
				redraw = true
			}
			g.CamZoom *= math.Pow(g.CamZoomSpeed, g.window.MouseScroll().Y)
			cam := pixel.IM.Scaled(g.CamPos, g.CamZoom).Moved(g.window.Bounds().Center().Sub(g.CamPos))

			// if g.window.MouseScroll().Y > 0 {
			// 	g.CamPos = g.CamPos.Add(cam.Unproject(g.window.MousePosition())))

			// } else if g.window.MouseScroll().Y < 0 {
			// 	g.CamPos = g.CamPos.Sub(cam.Unproject(g.window.MousePosition())))
			// }

			if g.window.MouseScroll().Y != 0 {
				redraw = true
			}

			if g.window.JustPressed(pixelgl.MouseButton1) {
				mousePos := g.isoToCartesian(cam.Unproject(g.window.MousePosition()))
				tileX := int((mousePos.X + 1))
				tileY := int((mousePos.Y + 1))
				t, err := g.currentLevel.GetTile(tileX, tileY)
				if err != nil {
					fmt.Println(err)
				} else {
					t.Clear()
					t.Push(g.currentLevel.SS.LeftAngleRoof)
					isoCoords := g.cartesianToIso(pixel.V(float64(tileX), float64(tileY)))
					mat := pixel.IM.Moved(g.offset.Add(isoCoords))
					t.Draw(g.window, mat)
					g.renderUpdate(tileX, tileY)
					redraw = true
				}
			}

			g.window.SetMatrix(cam)
			if redraw {
				redraw = false
				g.window.Clear(color.Black)
				batch.Draw(g.window)
			}

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
func (g *Game) renderAll() {
	g.render(g.currentLevel.GetWidth()-1, g.currentLevel.GetHeight()-1, 0, 0)
}

//TODO needs better logic to stop framedips far from the front
func (g *Game) renderUpdate(fx, fy int) {
	g.render(fx+1, fy+1, 0, 0)
}

func (g *Game) render(fx, fy, xMin, yMin int) {
	for x := fx; x >= xMin; x-- {
		for y := fy; y >= yMin; y-- {
			isoCoords := g.cartesianToIso(pixel.V(float64(x), float64(y)))
			mat := pixel.IM.Moved(g.offset.Add(isoCoords))
			// mat = mat.ScaledXY(g.window.Bounds().Center(), pixel.V(1, 1))
			t, err := g.currentLevel.GetTile(x, y)
			if err != nil {
				fmt.Println(err)
			} else {
				t.Draw(g.window, mat)
			}
		}
	}
}

func (g *Game) cartesianToIso(pt pixel.Vec) pixel.Vec {
	return pixel.V((pt.X-pt.Y)*(g.tileSize/2), (pt.X+pt.Y)*(g.tileSize/4))
}

func (g *Game) isoToCartesian(pt pixel.Vec) pixel.Vec {
	// x := (pt.X/(g.tileSize/2) + pt.Y/(g.tileSize/4)) / 2
	// y := (pt.Y/(g.tileSize/4) - (pt.X / (g.tileSize / 2))) / 2
	return pixel.V((pt.X/(g.tileSize/2)+pt.Y/(g.tileSize/4))/2, (pt.Y/(g.tileSize/4)-(pt.X/(g.tileSize/2)))/2)
}
