package game

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten"
	"github.com/jrcichra/gollercoaster/music"

	"github.com/jrcichra/gollercoaster/level"
)

//Game - Game Object
type Game struct {
	Name         string //Name of the game ("gollercoaster for now")
	windowWidth  int
	windowHeight int
	tileSize     int
	currentLevel *level.Level //pointer to the currentLevel to render
	CamPosX      float64
	CamPosY      float64
	CamSpeed     float64
	CamZoom      float64
	CamZoomSpeed float64
	music        music.Music
	frames       int64
	screen       *ebiten.Image
}

func (g *Game) playMusic() {
	g.music.LoadRandom()
	g.music.Play()
	g.music.Close()
}

// update is called every frame (1/60 [s]).
func (g *Game) update(screen *ebiten.Image) error {
	g.screen = screen
	g.frames++
	// fmt.Println("Rendering frame", g.frames)
	dt := 1.0 / 60

	// Write your game's logical update.

	if ebiten.IsDrawingSkipped() {
		// When the game is running slowly, the rendering result
		// will not be adopted.
		return nil
	}

	// Write your game's rendering.

	screen.Fill(color.Black)

	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		g.CamPosX -= g.CamSpeed * dt / g.CamZoom
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		g.CamPosX += g.CamSpeed * dt / g.CamZoom
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		g.CamPosY -= g.CamSpeed * dt / g.CamZoom
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		g.CamPosY += g.CamSpeed * dt / g.CamZoom
	}
	_, sY := ebiten.Wheel()
	g.CamZoom *= math.Pow(g.CamZoomSpeed, sY)

	// cam := pixel.IM.Scaled(g.CamPos, g.CamZoom).Moved(pixel.V(float64(g.currentLevel.GetWidth()/2), float64(g.currentLevel.GetHeight()/2)).Sub(g.CamPos))
	// var identity ebiten.GeoM
	// identity.Scale(g.CamPosX,g.CamPosY).Scale(g.CamZoom)
	// if g.window.MouseScroll().Y > 0 {
	// 	g.CamPos = g.CamPos.Add(cam.Unproject(g.window.MousePosition())))

	// } else if g.window.MouseScroll().Y < 0 {
	// 	g.CamPos = g.CamPos.Sub(cam.Unproject(g.window.MousePosition())))
	// }

	// if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
	// 	mx, my := ebiten.CursorPosition()
	// 	mp := pixel.V(float64(mx), float64(my))
	// 	mousePos := g.isoToCartesian(cam.Unproject(mp))
	// 	tileX := int((mousePos.X + 1))
	// 	tileY := int((mousePos.Y + 1))
	// 	t, err := g.currentLevel.GetTile(tileX, tileY)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	} else {
	// 		t.Clear()
	// 		t.Push(g.currentLevel.SS.LeftAngleRoof)
	// 	}
	// }

	//Render the whole thing
	g.renderAll()

	return nil
}

//Run - run the game
func (g *Game) Run() {
	g.Name = "gollercoaster"
	g.windowWidth = 1280
	g.windowHeight = 720
	g.tileSize = 64
	g.CamPosX = 0
	g.CamPosY = 0
	g.CamSpeed = 500
	g.CamZoom = 1
	g.CamZoomSpeed = 1.2
	// g.music = music.Music{}
	g.frames = 0

	// go g.playMusic()

	//Create a new level
	var l level.Level
	g.currentLevel = &l
	l.Spawn()
	ebiten.SetMaxTPS(ebiten.UncappedTPS)
	if err := ebiten.Run(g.update, g.windowWidth, g.windowHeight, 1, g.Name); err != nil {
		panic(err)
	}

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

	for x := xMin; x <= fx; x++ {
		for y := yMin; y <= fy; y++ {
			xi, yi := g.cartesianToIso(x, y)
			op := &ebiten.DrawImageOptions{}
			//Translate for isometric
			op.GeoM.Translate(float64(xi), float64(yi))
			//Transfer for camera position
			op.GeoM.Translate(-g.CamPosX, g.CamPosY)
			t, err := g.currentLevel.GetTile(x, y)
			if err != nil {
				fmt.Println(err)
			} else {
				t.Draw(g.screen, op)
			}
		}
	}
}

func (g *Game) cartesianToIso(x, y int) (int, int) {
	rx := (x - y) * (g.tileSize / 2)
	ry := (x + y) * (g.tileSize / 4)
	return rx, ry
}

func (g *Game) isoToCartesian(x, y int) (int, int) {
	rx := (x/(g.tileSize/2) + y/(g.tileSize/4)) / 2
	ry := (y/(g.tileSize/4) - (x / (g.tileSize / 2))) / 2
	return rx, ry
}
