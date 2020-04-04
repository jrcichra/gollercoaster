package game

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
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
	op           *ebiten.DrawImageOptions
}

func (g *Game) playMusic() {
	g.music.LoadRandom()
	g.music.Play()
	g.music.Close()
}

// update is called every frame (1/60 [s]).
func (g *Game) update(screen *ebiten.Image) error {
	// fmt.Println("Rendering frame", g.frames)
	dt := 1.0 / 60

	// Write your game's logical update.

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

	//Get the cursor position
	mx, my := ebiten.CursorPosition()
	//Offset for center
	fmx := float64(mx) - float64(g.windowWidth)/2.0
	fmy := float64(my) - float64(g.windowHeight)/2.0
	// x, y := float64(mx)+float64(g.windowWidth/2.0), float64(my)+float64(g.windowHeight/2.0)
	//Translate it to game coordinates
	x, y := (float64(fmx/g.CamZoom) + g.CamPosX), float64(fmy/g.CamZoom)-g.CamPosY

	//Do a half tile mouse shift because of our perspective
	x -= .5 * float64(g.tileSize)
	y -= .5 * float64(g.tileSize)
	//Convert isometric
	imx, imy := g.isoToCartesian(x, y)

	tileX := int(imx)
	tileY := int(imy)
	t, err := g.currentLevel.GetTile(tileX, tileY)
	if err != nil {
		fmt.Println(err)
	} else {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			t.Clear()
			t.Push(g.currentLevel.SS.LeftAngleRoof)
		} else {
			t.TempPush(g.currentLevel.SS.Selected)
		}
	}

	if ebiten.IsDrawingSkipped() {
		// When the game is running slowly, the rendering result
		// will not be adopted.
		fmt.Println("WARNING: We skipped a frame")
		return nil
	}

	// Write your game's rendering.
	g.render(screen)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS %f, FPS %f", ebiten.CurrentTPS(), ebiten.CurrentFPS()))

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
	// go g.playMusic()
	g.op = &ebiten.DrawImageOptions{}
	//Create a new level
	l := &level.Level{}
	g.currentLevel = l
	l.Spawn()
	// ebiten.SetMaxTPS(ebiten.UncappedTPS)
	// ebiten.SetVsyncEnabled(false)
	if err := ebiten.Run(g.update, g.windowWidth, g.windowHeight, 1, g.Name); err != nil {
		panic(err)
	}

}

func (g *Game) render(screen *ebiten.Image) {

	for x := 0; x <= g.currentLevel.GetWidth()-1; x++ {
		for y := 0; y <= g.currentLevel.GetHeight()-1; y++ {
			xi, yi := g.cartesianToIso(float64(x), float64(y))
			g.op.GeoM.Reset()
			//Translate for isometric
			g.op.GeoM.Translate(float64(xi), float64(yi))
			//Translate for camera position
			g.op.GeoM.Translate(-g.CamPosX, g.CamPosY)
			//Scale for camera zoom
			g.op.GeoM.Scale(g.CamZoom, g.CamZoom)
			//Translate for center of screen offset
			g.op.GeoM.Translate(float64(g.windowWidth/2.0), float64(g.windowHeight/2.0))
			t, err := g.currentLevel.GetTile(x, y)
			if err != nil {
				fmt.Println(err)
			} else {
				t.Draw(screen, g.op)
			}
		}
	}
}

func (g *Game) cartesianToIso(x, y float64) (float64, float64) {
	rx := (x - y) * float64(g.tileSize/2)
	ry := (x + y) * float64(g.tileSize/4)
	return rx, ry
}

func (g *Game) isoToCartesian(x, y float64) (float64, float64) {
	rx := (x/float64(g.tileSize/2) + y/float64(g.tileSize/4)) / 2
	ry := (y/float64(g.tileSize/4) - (x / float64(g.tileSize/2))) / 2
	return rx, ry
}
