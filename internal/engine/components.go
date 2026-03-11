package engine

type ComponentType int
type Direction int

// iota auto assigns numbers -> empty = 0; wire = 1; switch = 0
const (
	Empty ComponentType = iota
	Wire
	Switch
	Light
	NotGate
)

const (
	None Direction = iota
	Up
	Down
	Left
	Right
)

// 1x1 tile
type Tile struct {
	Type    ComponentType
	Powered bool
	Facing  Direction
}
