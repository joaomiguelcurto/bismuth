package engine

type ComponentType int

// iota auto assigns numbers -> empty = 0; wire = 1; switch = 0
const (
	Empty ComponentType = iota
	Wire
	Switch
	Light
)

// 1x1 tile
type Tile struct {
	Type    ComponentType
	Powered bool
}
