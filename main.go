package main

import (
	"image/color"
	"log"
	"strconv"
	"strings"

	"bismuth/internal/engine"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	GridWidth  = 40
	GridHeight = 40
	TileSize   = 20
)

type App struct {
	grid         *engine.Grid
	showHelp     bool
	selectedTool string
	currentLayer int
}

func (app *App) Update() error {
	mx, my := ebiten.CursorPosition()

	gridX := mx / TileSize
	gridY := my / TileSize
	gridZ := app.currentLayer

	if inpututil.IsKeyJustPressed(ebiten.KeyH) {
		app.showHelp = !app.showHelp
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		switch app.selectedTool {
		case "wire":
			app.grid.SetTile(gridX, gridY, gridZ, engine.Wire, engine.None)
		case "diode":
			app.grid.SetTile(gridX, gridY, gridZ, engine.Diode, engine.Up)
		case "switch":
			app.grid.SetTile(gridX, gridY, gridZ, engine.Switch, engine.None)
		case "button":
			app.grid.SetTile(gridX, gridY, gridZ, engine.Button, engine.None)
		case "light":
			app.grid.SetTile(gridX, gridY, gridZ, engine.Light, engine.None)
		case "notgate":
			app.grid.SetTile(gridX, gridY, gridZ, engine.NotGate, engine.Up)
		case "via":
			app.grid.SetTile(gridX, gridY, gridZ, engine.Via, engine.None)
		default:
			app.grid.SetTile(gridX, gridY, gridZ, engine.Empty, engine.None)
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		app.grid.ToggleTile(gridX, gridY, gridZ)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		app.grid.RotateGate(gridX, gridY, gridZ)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		app.grid.SetTile(gridX, gridY, gridZ, engine.Empty, engine.None)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		app.currentLayer++
		if app.currentLayer >= app.grid.Layers {
			app.currentLayer = 0
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		app.currentLayer--
		if app.currentLayer < 0 {
			app.currentLayer = app.grid.Layers - 1
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		app.selectedTool = "wire"
	}

	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		app.selectedTool = "diode"
	}

	if inpututil.IsKeyJustPressed(ebiten.Key3) {
		app.selectedTool = "switch"
	}

	if inpututil.IsKeyJustPressed(ebiten.Key4) {
		app.selectedTool = "button"
	}

	if inpututil.IsKeyJustPressed(ebiten.Key5) {
		app.selectedTool = "light"
	}

	if inpututil.IsKeyJustPressed(ebiten.Key6) {
		app.selectedTool = "notgate"
	}

	if inpututil.IsKeyJustPressed(ebiten.Key7) {
		app.selectedTool = "via"
	}

	app.grid.UpdatePower()

	return nil
}

func (app *App) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{20, 20, 20, 255})

	for x := 0; x < app.grid.Width; x++ {
		for y := 0; y < app.grid.Height; y++ {
			tile := app.grid.Tiles[x][y][app.currentLayer]

			if tile.Type == engine.Empty {
				continue
			}

			var c color.Color
			if tile.Type == engine.Wire {
				c = color.RGBA{100, 100, 100, 255}
				if tile.Powered {
					c = color.RGBA{255, 50, 50, 255}
				}
			} else if tile.Type == engine.Switch {
				c = color.RGBA{0, 100, 0, 255}
				if tile.Powered {
					c = color.RGBA{50, 255, 50, 255}
				}
			} else if tile.Type == engine.Light {
				c = color.RGBA{145, 145, 0, 255}
				if tile.Powered {
					c = color.RGBA{255, 255, 0, 255}
				}
			} else if tile.Type == engine.NotGate {
				c = color.RGBA{100, 0, 150, 255}
				if tile.Powered {
					c = color.RGBA{200, 50, 255, 255}
				}
			} else if tile.Type == engine.Button {
				c = color.RGBA{0, 70, 150, 255}
				if tile.Powered {
					c = color.RGBA{50, 150, 255, 255}
				}
			} else if tile.Type == engine.Diode {
				c = color.RGBA{100, 100, 50, 255}
				if tile.Powered {
					c = color.RGBA{255, 50, 0, 255}
				}
			} else if tile.Type == engine.Via {
				c = color.RGBA{50, 150, 150, 255}
				if tile.Powered {
					c = color.RGBA{100, 255, 255, 255}
				}
			}

			rectX := float32(x * TileSize)
			rectY := float32(y * TileSize)
			vector.FillRect(screen, rectX, rectY, float32(TileSize-1), float32(TileSize-1), c, true)

			if tile.Type == engine.Via {
				indSize := float32(8)
				centerX := rectX + float32(TileSize)/2.0 - indSize/2.0
				centerY := rectY + float32(TileSize)/2.0 - indSize/2.0

				holeColor := color.RGBA{0, 0, 0, 255}
				if tile.Powered {
					holeColor = color.RGBA{255, 255, 255, 255}
				}
				vector.FillRect(screen, centerX, centerY, indSize, indSize, holeColor, true)
			}

			if tile.Type == engine.NotGate || tile.Type == engine.Diode {
				indSize := float32(4)

				centerX := rectX + float32(TileSize)/2.0 - indSize/2.0
				centerY := rectY + float32(TileSize)/2.0 - indSize/2.0

				var indX, indY float32

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
					indX, indY = centerX, centerY
				}

				vector.FillRect(screen, indX, indY, indSize, indSize, color.RGBA{0, 0, 0, 255}, true)
			}
		}
	}

	if app.showHelp {
		panelWidth := float32(250)
		screenHeight := float32(GridHeight * TileSize)
		panelX := float32(GridWidth*TileSize) - panelWidth

		vector.FillRect(screen, panelX, 0, panelWidth, screenHeight, color.RGBA{30, 30, 30, 240}, true)

		legend := "=== BISMUTH CONTROLS ===\n\n" +
			"Left Click : Place Selected\n" +
			"Spacebar   : Toggle Tile\n" +
			"R          : Rotate Tile\n" +
			"Backspace  : Clear Tile\n" +
			"Up/Down    : Change Layer\n" +
			"Num 1      : Select Wire\n" +
			"Num 2      : Select Diode\n" +
			"Num 3      : Select Switch\n" +
			"Num 4      : Select Button\n" +
			"Num 5      : Select Light\n" +
			"Num 6      : Select NotGate\n" +
			"Num 7      : Select Via\n\n" +
			"Press 'H' to hide/show this menu"

		ebitenutil.DebugPrintAt(screen, legend, int(panelX)+20, 20)
	}

	bottomY := (GridHeight * TileSize) - 20
	layerInfo := " | Layer: " + strconv.Itoa(app.currentLayer)
	ebitenutil.DebugPrintAt(screen, "Tool selected > "+strings.Title(app.selectedTool)+layerInfo, 10, bottomY)
}

func (app *App) Layout(outsideWidth, outsideHeight int) (int, int) {
	return GridWidth * TileSize, GridHeight * TileSize
}

func main() {
	myApp := &App{
		grid:         engine.NewGrid(GridWidth, GridHeight, 3), // Creates 3 layers!
		selectedTool: "wire",
		currentLayer: 0,
	}

	ebiten.SetWindowSize(GridWidth*TileSize, GridHeight*TileSize)
	ebiten.SetWindowTitle("Bismuth")

	if err := ebiten.RunGame(myApp); err != nil {
		log.Fatal(err)
	}
}
