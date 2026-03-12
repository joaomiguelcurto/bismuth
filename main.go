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
		app.grid.SetTile(gridX, gridY, engine.Wire, engine.None)
	}

	// RCLICK place switch
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		app.grid.SetTile(gridX, gridY, engine.Switch, engine.None)
	}

	// SPACE toggle switch/button on/off
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		app.grid.ToggleSwitch(gridX, gridY)
		app.grid.PressButton(gridX, gridY)
	}

	// R rotate gates (left/right/up/down)
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		app.grid.RotateGate(gridX, gridY)
	}

	// NUM1 place light
	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		app.grid.SetTile(gridX, gridY, engine.Light, engine.None)
	}

	// NUM2 place NotGate
	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		app.grid.SetTile(gridX, gridY, engine.NotGate, engine.Up)
	}

	// NUM3 place button
	if inpututil.IsKeyJustPressed(ebiten.Key3) {
		app.grid.SetTile(gridX, gridY, engine.Button, engine.None)
	}

	// NUM4 place button
	if inpututil.IsKeyJustPressed(ebiten.Key4) {
		app.grid.SetTile(gridX, gridY, engine.Diode, engine.Up)
	}

	// BACKSPACE sets tile to empty
	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		app.grid.SetTile(gridX, gridY, engine.Empty, engine.None)
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
			} else if tile.Type == engine.Light {
				c = color.RGBA{145, 145, 0, 255} // Unpowered light (Dark Yellow)
				if tile.Powered {
					c = color.RGBA{255, 255, 0, 255} // Powered light (Yellow)
				}
			} else if tile.Type == engine.NotGate {
				c = color.RGBA{100, 0, 150, 255} // Unpowered gate (Dark Purple)
				if tile.Powered {
					c = color.RGBA{200, 50, 255, 255} // Powered gate (Bright Neon Purple)
				}
			} else if tile.Type == engine.Button {
				c = color.RGBA{0, 70, 150, 255} // Unpowered button (Dark Blue)
				if tile.Powered {
					c = color.RGBA{50, 150, 255, 255} // Powered button (Bright Blue)
				}
			} else if tile.Type == engine.Diode {
				c = color.RGBA{100, 100, 50, 255} // Unpowered diode (Gray)
				if tile.Powered {
					c = color.RGBA{255, 50, 0, 255} // Powered diode (Red)
				}
			}

			// draw the tile
			rectX := float32(x * TileSize)
			rectY := float32(y * TileSize)
			// the -1 exists so it exists a 1 pixel gap between tiles
			vector.FillRect(screen, rectX, rectY, float32(TileSize-1), float32(TileSize-1), c, true)

			if tile.Type == engine.NotGate || tile.Type == engine.Diode {
				indSize := float32(4) // A tiny 4x4 pixel dot

				// Find the mathematical center of this specific tile
				centerX := rectX + float32(TileSize)/2.0 - indSize/2.0
				centerY := rectY + float32(TileSize)/2.0 - indSize/2.0

				var indX, indY float32

				// Move the dot to the correct edge based on the Facing direction
				switch tile.Facing {
				case engine.Up:
					indX, indY = centerX, rectY
				case engine.Down:
					indX, indY = centerX, rectY+float32(TileSize)-indSize-1
				case engine.Left:
					indX, indY = rectX, centerY
				case engine.Right:
					indX, indY = rectX+float32(TileSize)-indSize-1, centerY
				default:
					// If it's engine.None, we just draw it perfectly in the middle
					indX, indY = centerX, centerY
				}

				// Draw the tiny black dot to show where the power comes OUT
				vector.FillRect(screen, indX, indY, indSize, indSize, color.RGBA{0, 0, 0, 255}, true)
			}
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
