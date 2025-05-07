package player

import (
	"math"
)

// Constants related to player
const (
	Size = 28 // Player size (reduced from 38 to match smaller tile size)
)

// Player represents the player character
type Player struct {
	GridX, GridY int
	X, Y         float64 // Actual position for smooth movement
	DestX, DestY float64 // Destination for smooth movement
	Moving       bool
	Size         float64
}

// New creates a new player with the given initial grid position
func New(gridX, gridY int, tileSize float64) *Player {
	x := float64(gridX) * tileSize
	y := float64(gridY) * tileSize
	
	return &Player{
		GridX:  gridX,
		GridY:  gridY,
		X:      x,
		Y:      y,
		DestX:  x,
		DestY:  y,
		Moving: false,
		Size:   Size,
	}
}

// SetDestination sets a new destination for the player to move to
func (p *Player) SetDestination(gridX, gridY int, tileSize float64) {
	p.GridX = gridX
	p.GridY = gridY
	p.DestX = float64(gridX) * tileSize
	p.DestY = float64(gridY) * tileSize
	p.Moving = true
}

// Update updates the player's position with smooth movement
// Returns true if the player has arrived at the destination
func (p *Player) Update(moveSpeed float64) bool {
	if !p.Moving {
		return false
	}
	
	dx := p.DestX - p.X
	dy := p.DestY - p.Y
	
	if math.Abs(dx) < moveSpeed && math.Abs(dy) < moveSpeed {
		// Arrived at destination
		p.X = p.DestX
		p.Y = p.DestY
		p.Moving = false
		return true
	} else {
		// Move toward destination
		if dx != 0 {
			p.X += math.Copysign(moveSpeed, dx)
		}
		if dy != 0 {
			p.Y += math.Copysign(moveSpeed, dy)
		}
		return false
	}
}

// IsMoving returns whether the player is currently moving
func (p *Player) IsMoving() bool {
	return p.Moving
}

// GetGridPosition returns the current grid position of the player
func (p *Player) GetGridPosition() (int, int) {
	return p.GridX, p.GridY
}

// GetPosition returns the current pixel position of the player
func (p *Player) GetPosition() (float64, float64) {
	return p.X, p.Y
}