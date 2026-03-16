package engine

/*
X + something goes Right

X - something goes Left

Y + something goes DOWN (Because you are moving to the next row of pixels)

Y - something goes UP (Because you are moving back toward row 0).

Z layer depth
*/

type Grid struct {
	Width  int
	Height int
	Layers int
	Tiles  [][][]Tile // Grid > Column > Row > Layer
}

type Point struct {
	X int
	Y int
	Z int
}

// Creates a fresh grid of a specific size
func NewGrid(width, height, layers int) *Grid {
	// Make Columns(X)
	tiles := make([][][]Tile, width)

	// Rows(Y)
	for i := range tiles {
		tiles[i] = make([][]Tile, height)
		// loop to create the Z layers for every single X,Y tile
		for j := range tiles[i] {
			tiles[i][j] = make([]Tile, layers)
		}
	}

	return &Grid{
		Width:  width,
		Height: height,
		Layers: layers,
		Tiles:  tiles,
	}
}

// Replaces a component on the grid
func (g *Grid) SetTile(x, y, z int, compType ComponentType, facing Direction) {
	// Make sure its inside the grid size
	if x >= 0 && x < g.Width && y >= 0 && y < g.Height && z >= 0 && z < g.Layers {
		g.Tiles[x][y][z] = Tile{
			Type:    compType,
			Powered: false,
			Facing:  facing,
		}
	}
}

// Rotate Gate
func (g *Grid) RotateGate(x, y, z int) {
	// Make sure its inside the grid size
	if x >= 0 && x < g.Width && y >= 0 && y < g.Height && z >= 0 && z < g.Layers {
		if g.Tiles[x][y][z].Facing != None {
			switch g.Tiles[x][y][z].Facing {
			case Up:
				g.Tiles[x][y][z].Facing = Right
			case Down:
				g.Tiles[x][y][z].Facing = Left
			case Left:
				g.Tiles[x][y][z].Facing = Up
			case Right:
				g.Tiles[x][y][z].Facing = Down
			}
		}
	}
}

func (g *Grid) ToggleTile(x, y, z int) {
	if x >= 0 && x < g.Width && y >= 0 && y < g.Height && z >= 0 && z < g.Layers {
		switch g.Tiles[x][y][z].Type {
		case Switch:
			g.Tiles[x][y][z].Powered = !g.Tiles[x][y][z].Powered
		case Button:
			g.Tiles[x][y][z].Powered = true
			g.Tiles[x][y][z].Timer = 120 // sets timer to 120 frames (60 frames per second)
		}
	}
}

func (g *Grid) UpdatePower() {
	powerSource := make([]Point, 0)

	// READ PHASE
	for i := 0; i < g.Width; i++ {
		for j := 0; j < g.Height; j++ {
			for k := 0; k < g.Layers; k++ {
				if g.Tiles[i][j][k].Type == NotGate {
					g.Tiles[i][j][k].Powered = false
					switch g.Tiles[i][j][k].Facing {
					case Up:
						if j+1 < g.Height {
							if g.Tiles[i][j+1][k].Powered == false {
								g.Tiles[i][j][k].Powered = true
							}
						}
					case Down:
						if j-1 >= 0 {
							if g.Tiles[i][j-1][k].Powered == false {
								g.Tiles[i][j][k].Powered = true
							}
						}
					case Left:
						if i+1 < g.Width {
							if g.Tiles[i+1][j][k].Powered == false {
								g.Tiles[i][j][k].Powered = true
							}
						}
					case Right:
						if i-1 >= 0 {
							if g.Tiles[i-1][j][k].Powered == false {
								g.Tiles[i][j][k].Powered = true
							}
						}
					}
				}

				isSwitch := g.Tiles[i][j][k].Type == Switch
				isButton := g.Tiles[i][j][k].Type == Button
				isNotGate := g.Tiles[i][j][k].Type == NotGate

				if (isSwitch || isButton || isNotGate) && g.Tiles[i][j][k].Powered {
					powerSource = append(powerSource, Point{X: i, Y: j, Z: k})
				}
				if g.Tiles[i][j][k].Type == Button {
					if g.Tiles[i][j][k].Timer > 0 {
						g.Tiles[i][j][k].Timer--
					} else {
						g.Tiles[i][j][k].Powered = false
					}
				}
			}
		}
	}

	// WIPE PHASE
	for i := 0; i < g.Width; i++ {
		for j := 0; j < g.Height; j++ {
			for k := 0; k < g.Layers; k++ {
				switch g.Tiles[i][j][k].Type {
				case Wire, Light, Diode, Via:
					g.Tiles[i][j][k].Powered = false
				}
			}
		}
	}

	// SPREAD PHASE
	for len(powerSource) > 0 {
		currentSource := powerSource[0]
		powerSource = powerSource[1:]

		sourceTile := g.Tiles[currentSource.X][currentSource.Y][currentSource.Z]

		// If this tile is a Via, or the tile above/below is a Via, power moves vertically
		layerAbove := currentSource.Z + 1
		layerBelow := currentSource.Z - 1

		if layerAbove < g.Layers {
			targetAbove := g.Tiles[currentSource.X][currentSource.Y][layerAbove]
			if (sourceTile.Type == Via || targetAbove.Type == Via) && !targetAbove.Powered && targetAbove.Type != Empty {
				g.Tiles[currentSource.X][currentSource.Y][layerAbove].Powered = true
				powerSource = append(powerSource, Point{X: currentSource.X, Y: currentSource.Y, Z: layerAbove})
			}
		}
		if layerBelow >= 0 {
			targetBelow := g.Tiles[currentSource.X][currentSource.Y][layerBelow]
			if (sourceTile.Type == Via || targetBelow.Type == Via) && !targetBelow.Powered && targetBelow.Type != Empty {
				g.Tiles[currentSource.X][currentSource.Y][layerBelow].Powered = true
				powerSource = append(powerSource, Point{X: currentSource.X, Y: currentSource.Y, Z: layerBelow})
			}
		}

		// HORIZONTAL X/Y SPREAD
		up := Point{X: currentSource.X, Y: currentSource.Y - 1, Z: currentSource.Z}
		down := Point{X: currentSource.X, Y: currentSource.Y + 1, Z: currentSource.Z}
		left := Point{X: currentSource.X - 1, Y: currentSource.Y, Z: currentSource.Z}
		right := Point{X: currentSource.X + 1, Y: currentSource.Y, Z: currentSource.Z}

		isLaser := sourceTile.Type == NotGate || sourceTile.Type == Diode

		canGoUp := !isLaser || sourceTile.Facing == Up
		canGoDown := !isLaser || sourceTile.Facing == Down
		canGoLeft := !isLaser || sourceTile.Facing == Left
		canGoRight := !isLaser || sourceTile.Facing == Right

		if canGoUp {
			if up.X >= 0 && up.X < g.Width && up.Y >= 0 && up.Y < g.Height {
				target := &g.Tiles[up.X][up.Y][up.Z]
				if (target.Type == Wire || target.Type == Via) && !target.Powered {
					target.Powered = true
					powerSource = append(powerSource, up)
				} else if target.Type == Light && !target.Powered {
					target.Powered = true
				} else if target.Type == Diode && target.Facing == Up && !target.Powered {
					target.Powered = true
					powerSource = append(powerSource, up)
				}
			}
		}
		if canGoDown {
			if down.X >= 0 && down.X < g.Width && down.Y >= 0 && down.Y < g.Height {
				target := &g.Tiles[down.X][down.Y][down.Z]
				if (target.Type == Wire || target.Type == Via) && !target.Powered {
					target.Powered = true
					powerSource = append(powerSource, down)
				} else if target.Type == Light && !target.Powered {
					target.Powered = true
				} else if target.Type == Diode && target.Facing == Down && !target.Powered {
					target.Powered = true
					powerSource = append(powerSource, down)
				}
			}
		}
		if canGoLeft {
			if left.X >= 0 && left.X < g.Width && left.Y >= 0 && left.Y < g.Height {
				target := &g.Tiles[left.X][left.Y][left.Z]
				if (target.Type == Wire || target.Type == Via) && !target.Powered {
					target.Powered = true
					powerSource = append(powerSource, left)
				} else if target.Type == Light && !target.Powered {
					target.Powered = true
				} else if target.Type == Diode && target.Facing == Left && !target.Powered {
					target.Powered = true
					powerSource = append(powerSource, left)
				}
			}
		}
		if canGoRight {
			if right.X >= 0 && right.X < g.Width && right.Y >= 0 && right.Y < g.Height {
				target := &g.Tiles[right.X][right.Y][right.Z]
				if (target.Type == Wire || target.Type == Via) && !target.Powered {
					target.Powered = true
					powerSource = append(powerSource, right)
				} else if target.Type == Light && !target.Powered {
					target.Powered = true
				} else if target.Type == Diode && target.Facing == Right && !target.Powered {
					target.Powered = true
					powerSource = append(powerSource, right)
				}
			}
		}
	}
}
