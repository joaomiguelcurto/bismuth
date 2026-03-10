package main

import (
	"image/color"
	"log"

	"bismuth/internal/engine"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	GridWidth  = 40
	GridHeight = 40
	TileSize   = 20 // each tile is x pixels on the screen
)

type App struct {
	grid *engine.Grid
}

func (app *App) Update() error {
	// get mouse pos
	mx, my := ebiten.CursorPosition()

	// convert pixel to grid cords
	gridX := mx / TileSize
	gridY := my / TileSize

	// LCLICK place wire
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		app.grid.SetTile(gridX, gridY, engine.Wire)
	}

	// RCLICK place switch
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		app.grid.SetTile(gridX, gridY, engine.Switch)
	}

	// SPACE toggle switch on/off
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		app.grid.ToggleSwitch(gridX, gridY)
	}

	app.grid.UpdatePower()

	return nil
}

// Draw runs right after update
func (app *App) Draw(screen *ebiten.Image) {
	// paint the whole background dark gray
	screen.Fill(color.RGBA{20, 20, 20, 255})

	// loop through the 2d grid
	for x := 0; x < app.grid.Width; x++ {
		for y := 0; y < app.grid.Height; y++ {
			tile := app.grid.Tiles[x][y]

			if tile.Type == engine.Empty {
				continue // skip empty tiles
			}

			// pick color depending on what the tile is
			var c color.Color
			if tile.Type == engine.Wire {
				c = color.RGBA{100, 100, 100, 255} // Unpowered wire (Gray)
				if tile.Powered {
					c = color.RGBA{255, 50, 50, 255} // Powered wire (Red)
				}
			} else if tile.Type == engine.Switch {
				c = color.RGBA{0, 100, 0, 255} // Unpowered switch (Dark Green)
				if tile.Powered {
					c = color.RGBA{50, 255, 50, 255} // Powered switch (Bright Green)
				}
			}

			// draw the tile
			rectX := float32(x * TileSize)
			rectY := float32(y * TileSize)
			// the -1 exists so it exists a 1 pixel gap between tiles
			vector.FillRect(screen, rectX, rectY, float32(TileSize-1), float32(TileSize-1), c, true)
		}
	}
}

// tells Ebitengine how big the window should be
func (app *App) Layout(outsideWidth, outsideHeight int) (int, int) {
	return GridWidth * TileSize, GridHeight * TileSize
}

func main() {
	// create app and initialize grid
	myApp := &App{
		grid: engine.NewGrid(GridWidth, GridHeight),
	}

	// window config
	ebiten.SetWindowSize(GridWidth*TileSize, GridHeight*TileSize)
	ebiten.SetWindowTitle("Bismuth")

	// starts game loop
	if err := ebiten.RunGame(myApp); err != nil {
		log.Fatal(err)
	}
}
