// internal/game/maze/position.go
package maze

// Position represents a position in the grid
type Position struct {
    X, Y int
}

// NewPosition creates a new position
func NewPosition(x, y int) Position {
    return Position{X: x, Y: y}
}