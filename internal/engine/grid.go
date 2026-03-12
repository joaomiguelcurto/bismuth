package engine

/*
X + something goes Right

X - something goes Left

Y + something goes DOWN (Because you are moving to the next row of pixels)

Y - something goes UP (Because you are moving back toward row 0).
*/

type Grid struct {
	Width  int
	Height int
	Tiles  [][]Tile // Grid > Column > Row
}

type Point struct {
	X int
	Y int
}

// Creates a fresh grid of a specific size
func NewGrid(width, height int) *Grid {
	// Make Columns(X)
	tiles := make([][]Tile, width)

	// Rows(Y)
	for i := range tiles {
		tiles[i] = make([]Tile, height)
	}

	return &Grid{
		Width:  width,
		Height: height,
		Tiles:  tiles,
	}
}

// Replaces a component on the grid
func (g *Grid) SetTile(x, y int, compType ComponentType, facing Direction) {
	// Make sure its inside the grid size
	if x >= 0 && x < g.Width && y >= 0 && y < g.Height {
		g.Tiles[x][y] = Tile{
			Type:    compType,
			Powered: false,
			Facing:  facing,
		}
	}
}

// Rotate Gate
func (g *Grid) RotateGate(x, y int) {
	// Make sure its inside the grid size
	if x >= 0 && x < g.Width && y >= 0 && y < g.Height {
		if g.Tiles[x][y].Facing != None {
			switch g.Tiles[x][y].Facing {
			case Up:
				g.Tiles[x][y].Facing = Right
			case Down:
				g.Tiles[x][y].Facing = Left
			case Left:
				g.Tiles[x][y].Facing = Up
			case Right:
				g.Tiles[x][y].Facing = Down
			}
		}
	}
}

func (g *Grid) ToggleTile(x, y int) {
	if x >= 0 && x < g.Width && y >= 0 && y < g.Height {
		switch g.Tiles[x][y].Type {
		case Switch:
			g.Tiles[x][y].Powered = !g.Tiles[x][y].Powered
		case Button:
			g.Tiles[x][y].Powered = !g.Tiles[x][y].Powered
			g.Tiles[x][y].Timer = 120
		}
	}
}

func (g *Grid) UpdatePower() {
	powerSource := make([]Point, 0)
	for i := 0; i < g.Width; i++ {
		for j := 0; j < g.Height; j++ {
			if g.Tiles[i][j].Type == NotGate {
				g.Tiles[i][j].Powered = false
				switch g.Tiles[i][j].Facing {
				case Up:
					if j+1 < g.Height {
						if g.Tiles[i][j+1].Powered == false {
							g.Tiles[i][j].Powered = true
						}
					}
				case Down:
					if j-1 >= 0 {
						if g.Tiles[i][j-1].Powered == false {
							g.Tiles[i][j].Powered = true
						}
					}
				case Left:
					if i+1 < g.Width {
						if g.Tiles[i+1][j].Powered == false {
							g.Tiles[i][j].Powered = true
						}
					}
				case Right:
					if i-1 >= 0 {
						if g.Tiles[i-1][j].Powered == false {
							g.Tiles[i][j].Powered = true
						}
					}
				}
			}

			isSwitch := g.Tiles[i][j].Type == Switch
			isButton := g.Tiles[i][j].Type == Button
			isNotGate := g.Tiles[i][j].Type == NotGate

			if (isSwitch || isButton || isNotGate) && g.Tiles[i][j].Powered {
				powerSource = append(powerSource, Point{X: i, Y: j})
			}
			if g.Tiles[i][j].Type == Button {
				if g.Tiles[i][j].Timer > 0 {
					g.Tiles[i][j].Timer--
				} else {
					g.Tiles[i][j].Powered = false
				}
			}
		}
	}

	for i := 0; i < g.Width; i++ {
		for j := 0; j < g.Height; j++ {
			switch g.Tiles[i][j].Type {
			case Wire:
				g.Tiles[i][j].Powered = false
			case Light:
				g.Tiles[i][j].Powered = false
			case Diode:
				g.Tiles[i][j].Powered = false
			}
		}
	}

	for len(powerSource) > 0 {
		currentSource := powerSource[0]
		powerSource = powerSource[1:]
		up := Point{X: currentSource.X, Y: currentSource.Y - 1}
		down := Point{X: currentSource.X, Y: currentSource.Y + 1}
		left := Point{X: currentSource.X - 1, Y: currentSource.Y}
		right := Point{X: currentSource.X + 1, Y: currentSource.Y}

		sourceTile := g.Tiles[currentSource.X][currentSource.Y]

		isLaser := sourceTile.Type == NotGate || sourceTile.Type == Diode

		canGoUp := !isLaser || sourceTile.Facing == Up
		canGoDown := !isLaser || sourceTile.Facing == Down
		canGoLeft := !isLaser || sourceTile.Facing == Left
		canGoRight := !isLaser || sourceTile.Facing == Right

		if canGoUp {
			if up.X >= 0 && up.X < g.Width && up.Y >= 0 && up.Y < g.Height {
				if g.Tiles[up.X][up.Y].Type == Wire && !g.Tiles[up.X][up.Y].Powered {
					g.Tiles[up.X][up.Y].Powered = true
					powerSource = append(powerSource, up)
				} else if g.Tiles[up.X][up.Y].Type == Light && !g.Tiles[up.X][up.Y].Powered {
					g.Tiles[up.X][up.Y].Powered = true
				} else if g.Tiles[up.X][up.Y].Type == Diode && g.Tiles[up.X][up.Y].Facing == Up && !g.Tiles[up.X][up.Y].Powered {
					g.Tiles[up.X][up.Y].Powered = true
					powerSource = append(powerSource, up)
				}
			}
		}
		if canGoDown {
			if down.X >= 0 && down.X < g.Width && down.Y >= 0 && down.Y < g.Height {
				if g.Tiles[down.X][down.Y].Type == Wire && !g.Tiles[down.X][down.Y].Powered {
					g.Tiles[down.X][down.Y].Powered = true
					powerSource = append(powerSource, down)
				} else if g.Tiles[down.X][down.Y].Type == Light && !g.Tiles[down.X][down.Y].Powered {
					g.Tiles[down.X][down.Y].Powered = true
				} else if g.Tiles[down.X][down.Y].Type == Diode && g.Tiles[down.X][down.Y].Facing == Down && !g.Tiles[down.X][down.Y].Powered {
					g.Tiles[down.X][down.Y].Powered = true
					powerSource = append(powerSource, down)
				}
			}
		}
		if canGoLeft {
			if left.X >= 0 && left.X < g.Width && left.Y >= 0 && left.Y < g.Height {
				if g.Tiles[left.X][left.Y].Type == Wire && !g.Tiles[left.X][left.Y].Powered {
					g.Tiles[left.X][left.Y].Powered = true
					powerSource = append(powerSource, left)
				} else if g.Tiles[left.X][left.Y].Type == Light && !g.Tiles[left.X][left.Y].Powered {
					g.Tiles[left.X][left.Y].Powered = true
				} else if g.Tiles[left.X][left.Y].Type == Diode && g.Tiles[left.X][left.Y].Facing == Left && !g.Tiles[left.X][left.Y].Powered {
					g.Tiles[left.X][left.Y].Powered = true
					powerSource = append(powerSource, left)
				}
			}
		}
		if canGoRight {
			if right.X >= 0 && right.X < g.Width && right.Y >= 0 && right.Y < g.Height {
				if g.Tiles[right.X][right.Y].Type == Wire && !g.Tiles[right.X][right.Y].Powered {
					g.Tiles[right.X][right.Y].Powered = true
					powerSource = append(powerSource, right)
				} else if g.Tiles[right.X][right.Y].Type == Light && !g.Tiles[right.X][right.Y].Powered {
					g.Tiles[right.X][right.Y].Powered = true
				} else if g.Tiles[right.X][right.Y].Type == Diode && g.Tiles[right.X][right.Y].Facing == Right && !g.Tiles[right.X][right.Y].Powered {
					g.Tiles[right.X][right.Y].Powered = true
					powerSource = append(powerSource, right)
				}
			}
		}
	}
}
