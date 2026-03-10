package engine

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
func (g *Grid) SetTile(x, y int, compType ComponentType) {
	// Make sure its inside the grid size
	if x >= 0 && y < g.Width && y >= 0 && y < g.Height {
		g.Tiles[x][y] = Tile{
			Type:    compType,
			Powered: false,
		}
	}
}

func (g *Grid) ToggleSwitch(x, y int) {
	if x >= 0 && x < g.Width && y >= 0 && y < g.Height {
		// Check if its an actual switch
		if g.Tiles[x][y].Type == Switch {
			g.Tiles[x][y].Powered = !g.Tiles[x][y].Powered
		}
	}
}

func (g *Grid) UpdatePower() {
	powerSources := make([]Point, 0)
	for i := 0; i < g.Width; i++ {
		for j := 0; j < g.Height; j++ {
			switch g.Tiles[i][j].Type {
			case Wire:
				g.Tiles[i][j].Powered = false
			case Light:
				g.Tiles[i][j].Powered = false
			}
			if g.Tiles[i][j].Type == Switch && g.Tiles[i][j].Powered {
				powerSources = append(powerSources, Point{X: i, Y: j})
			}
		}
	}
	for len(powerSources) > 0 {
		currentSource := powerSources[0]
		powerSources = powerSources[1:]
		up := Point{X: currentSource.X, Y: currentSource.Y + 1}
		down := Point{X: currentSource.X, Y: currentSource.Y - 1}
		left := Point{X: currentSource.X - 1, Y: currentSource.Y}
		right := Point{X: currentSource.X + 1, Y: currentSource.Y}

		if up.X >= 0 && up.X < g.Width && up.Y >= 0 && up.Y < g.Height {
			if g.Tiles[up.X][up.Y].Type == Wire && !g.Tiles[up.X][up.Y].Powered {
				g.Tiles[up.X][up.Y].Powered = true
				powerSources = append(powerSources, up)
			} else if g.Tiles[up.X][up.Y].Type == Light && !g.Tiles[up.X][up.Y].Powered {
				g.Tiles[up.X][up.Y].Powered = true
			}
		}
		if down.X >= 0 && down.X < g.Width && down.Y >= 0 && down.Y < g.Height {
			if g.Tiles[down.X][down.Y].Type == Wire && !g.Tiles[down.X][down.Y].Powered {
				g.Tiles[down.X][down.Y].Powered = true
				powerSources = append(powerSources, down)
			} else if g.Tiles[down.X][down.Y].Type == Light && !g.Tiles[down.X][down.Y].Powered {
				g.Tiles[down.X][down.Y].Powered = true
			}
		}
		if left.X >= 0 && left.X < g.Width && left.Y >= 0 && left.Y < g.Height {
			if g.Tiles[left.X][left.Y].Type == Wire && !g.Tiles[left.X][left.Y].Powered {
				g.Tiles[left.X][left.Y].Powered = true
				powerSources = append(powerSources, left)
			} else if g.Tiles[left.X][left.Y].Type == Light && !g.Tiles[left.X][left.Y].Powered {
				g.Tiles[left.X][left.Y].Powered = true
			}
		}
		if right.X >= 0 && right.X < g.Width && right.Y >= 0 && right.Y < g.Height {
			if g.Tiles[right.X][right.Y].Type == Wire && !g.Tiles[right.X][right.Y].Powered {
				g.Tiles[right.X][right.Y].Powered = true
				powerSources = append(powerSources, right)
			} else if g.Tiles[right.X][right.Y].Type == Light && !g.Tiles[right.X][right.Y].Powered {
				g.Tiles[right.X][right.Y].Powered = true
			}
		}
	}
}
