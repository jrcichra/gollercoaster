package game

import (
	"fmt"
	"image/color"
	"math"
	"time"

	"github.com/jrcichra/gollercoaster/tile"

	"github.com/jrcichra/gollercoaster/music"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/jrcichra/gollercoaster/level"
)

//Game - Game Object
type Game struct {
	Name           string //Name of the game ("gollercoaster for now")
	windowWidth    float64
	windowHeight   float64
	window         *pixelgl.Window
	offset         pixel.Vec
	tileSize       float64
	currentLevel   *level.Level //Pointer to the currentLevel to render
	CamPos         pixel.Vec    //
	CamSpeed       float64      //Camera Speed at any given point
	CamZoom        float64      //Camera zoom level at any given point
	CamZoomSpeed   float64      //Speed of the zoom at any given point
	music          music.Music  //Music object controlling the music of the game (simple now)
	lastSelected   pixel.Vec    //Tile X,Y which was the tile that was selected on the previous frame
	batch          *pixel.Batch //Batch from which we render the screen from when we need to make a bigger draw call
	dt             float64      //Deltatime between frames
	last           time.Time    //The last time we rendered a frame
	mousePos       pixel.Vec    //Vector of the mouse position at any given time
	currentTile    *tile.Tile   //Tile we're hovering over at any given time (nil if offscreen)
	currentTilePos pixel.Vec    //X,Y of current tile
	frames         int
	second         <-chan time.Time
}

func (g *Game) playMusic() {
	g.music.LoadRandom()
	g.music.Play()
	g.music.Close()
}

func (g *Game) newWindow() {
	cfg := pixelgl.WindowConfig{
		Title:  g.Name,
		Bounds: pixel.R(0, 0, g.windowWidth, g.windowHeight),
		// VSync:  true,
	}
	var err error
	g.window, err = pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
}

func (g *Game) newLevel() {
	var l level.Level
	g.currentLevel = &l
	g.batch = l.Spawn()
}

//Logic update to handle what may change every frame in game state
func (g *Game) logicUpdate() {
	t := g.currentTile
	if t != nil {
		if g.window.JustPressed(pixelgl.MouseButton1) {
			//Click for an angle roof
			t.Clear() //Clear what's on the tile
			t.Push(g.currentLevel.SS.LeftAngleRoof)
		} else {
			//Hovering over a tile, peek to see if it's already been selected
			topSprite := t.Peek()
			if topSprite.Sprite == g.currentLevel.SS.DoubleWall.Sprite {
				//It's already been selected. Do nothing
			} else {
				//Push the DoubleWall onto the tile
				t.Push(g.currentLevel.SS.DoubleWall)
				//Pop the Doublewall from the last selected
				lt, err := g.currentLevel.GetTile(int(g.lastSelected.X), int(g.lastSelected.Y))
				if err != nil {
					fmt.Println(err)
				} else {
					lt.Pop()
				}
			}
			//Update the last selected to this one
			g.lastSelected = g.currentTilePos
		}
	}
}

//Window movements that shouldn't be costly to performance or require a change to the batch
//There might be some game state changes in here which coorilate directly with the way the window interacts with the game
func (g *Game) windowUpdate() {
	if g.window.Pressed(pixelgl.KeyLeft) || g.window.Pressed(pixelgl.KeyA) {
		g.CamPos.X -= g.CamSpeed * g.dt / g.CamZoom
	}
	if g.window.Pressed(pixelgl.KeyRight) || g.window.Pressed(pixelgl.KeyD) {
		g.CamPos.X += g.CamSpeed * g.dt / g.CamZoom
	}
	if g.window.Pressed(pixelgl.KeyDown) || g.window.Pressed(pixelgl.KeyS) {
		g.CamPos.Y -= g.CamSpeed * g.dt / g.CamZoom
	}
	if g.window.Pressed(pixelgl.KeyUp) || g.window.Pressed(pixelgl.KeyW) {
		g.CamPos.Y += g.CamSpeed * g.dt / g.CamZoom
	}
	g.CamZoom *= math.Pow(g.CamZoomSpeed, g.window.MouseScroll().Y)
	cam := pixel.IM.Scaled(g.CamPos, g.CamZoom).Moved(g.window.Bounds().Center().Sub(g.CamPos))
	g.mousePos = g.isoToCartesian(cam.Unproject(g.window.MousePosition()))

	//Set the matrix for the window for this frame
	g.window.SetMatrix(cam)

	tileX := int((g.mousePos.X + 1))
	tileY := int((g.mousePos.Y + 1))
	t, err := g.currentLevel.GetTile(tileX, tileY)
	if err != nil {
		fmt.Println(err)
	}
	g.currentTile = t
	g.currentTilePos = pixel.V(float64(tileX), float64(tileY))

}

//Updates that require a full or partial redraw of the buffer. The goal is to do as few of these updates as possible
func (g *Game) visualUpdate() {

	//Determine what changes we need to make to the batch (This will soon be optimized to what's on screen is only what's rendered)
	g.batch.Clear()
	g.renderAll()

	// isoCoords := g.cartesianToIso(g.currentTilePos)
	// mat := pixel.IM.Moved(g.offset.Add(isoCoords))
	// t.Draw(g.window, mat)

	//Redraw the batch to the screen every frame, no matter what
	g.window.Clear(color.Black)
	g.batch.Draw(g.window)
	g.window.Update()
	g.frames++

	select {
	case <-g.second:
		g.window.SetTitle(fmt.Sprintf("%s | FPS: %d", g.Name, g.frames))
		g.frames = 0
	default:
	}
}

//Run - run the game
func (g *Game) Run() {
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
	g.frames = 0
	g.second = time.Tick(time.Second)

	go g.playMusic()

	pixelgl.Run(func() {
		g.newWindow()
		g.newLevel()
		g.last = time.Now()
		g.window.Clear(color.Black)
		g.batch.Clear()
		g.renderAll()
		g.batch.Draw(g.window)
		for !g.window.Closed() {
			g.dt = time.Since(g.last).Seconds()
			g.last = time.Now()

			g.windowUpdate() //This is our input from the user
			g.logicUpdate()  //This is what we do to react based on that input / move the game forward
			g.visualUpdate() //This is what we send back to the user after we change the game state
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
	c := 0
	for x := fx; x >= xMin; x-- {
		for y := fy; y >= yMin; y-- {
			isoCoords := g.cartesianToIso(pixel.V(float64(x), float64(y)))
			mat := pixel.IM.Moved(g.offset.Add(isoCoords))
			// mat = mat.ScaledXY(g.window.Bounds().Center(), pixel.V(1, 1))
			t, err := g.currentLevel.GetTile(x, y)
			if err != nil {
				fmt.Println(err)
			} else {
				c += t.Draw(mat)
			}
		}
	}
	fmt.Println("I drew", c, "things to the batch")
}

func (g *Game) cartesianToIso(pt pixel.Vec) pixel.Vec {
	return pixel.V((pt.X-pt.Y)*(g.tileSize/2), (pt.X+pt.Y)*(g.tileSize/4))
}

func (g *Game) isoToCartesian(pt pixel.Vec) pixel.Vec {
	return pixel.V((pt.X/(g.tileSize/2)+pt.Y/(g.tileSize/4))/2, (pt.Y/(g.tileSize/4)-(pt.X/(g.tileSize/2)))/2)
}
