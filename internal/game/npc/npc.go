package npc

import (
	"image/color"
	"math"
	"math/rand"
)

// NPC represents a non-player character
type NPC struct {
	ID           int
	GridX, GridY int
	X, Y         float64 // Actual position for smooth movement
	DestX, DestY float64 // Destination for smooth movement
	Moving       bool
	Size         float64
	Color        color.RGBA
	HasMoved     bool    // Track if NPC has moved in current turn
}

// New creates a new NPC instance
func New(id, gridX, gridY int, size float64, color color.RGBA) *NPC {
	npc := &NPC{
		ID:       id,
		GridX:    gridX,
		GridY:    gridY,
		Size:     size,
		Color:    color,
		HasMoved: false,
	}
	
	// Set initial position
	npc.X = float64(gridX) * size
	npc.Y = float64(gridY) * size
	npc.DestX = npc.X
	npc.DestY = npc.Y
	
	return npc
}

// IsMoving checks if the NPC is currently moving
func (n *NPC) IsMoving() bool {
	return n.Moving
}

// ResetMovedStatus resets the HasMoved flag
func (n *NPC) ResetMovedStatus() {
	n.HasMoved = false
}

// UpdatePosition updates the NPC's position with smooth movement
// Returns true if the NPC has reached the destination
func (n *NPC) UpdatePosition(moveSpeed float64) bool {
	if !n.Moving {
		return false
	}
	
	dx := n.DestX - n.X
	dy := n.DestY - n.Y
	
	if math.Abs(dx) < moveSpeed && math.Abs(dy) < moveSpeed {
		// Arrived at destination
		n.X = n.DestX
		n.Y = n.DestY
		n.Moving = false
		return true
	} else {
		// Move toward destination
		if dx != 0 {
			n.X += math.Copysign(moveSpeed, dx)
		}
		if dy != 0 {
			n.Y += math.Copysign(moveSpeed, dy)
		}
		return false
	}
}

// TryMove attempts to move the NPC in a valid direction
// validMoveFn is a callback that determines if a move is valid
// Returns true if successfully moved
func (n *NPC) TryMove(validMoveFn func(x, y int) bool) bool {
	if n.Moving || n.HasMoved {
		return false // Already moving or has moved this turn
	}

	// Possible movement directions: left, right, up, down
	directions := []struct{ dx, dy int }{
		{-1, 0}, {1, 0}, {0, -1}, {0, 1},
	}

	// Shuffle directions for randomized movement
	rand.Shuffle(len(directions), func(i, j int) {
		directions[i], directions[j] = directions[j], directions[i]
	})

	// Try each direction until a valid move is found
	for _, dir := range directions {
		newGridX := n.GridX + dir.dx
		newGridY := n.GridY + dir.dy

		// Check if movement is valid using the callback
		if validMoveFn(newGridX, newGridY) {
			// Update grid position
			n.GridX = newGridX
			n.GridY = newGridY

			// Set destination for smooth movement
			n.DestX = float64(newGridX) * n.Size
			n.DestY = float64(newGridY) * n.Size
			n.Moving = true
			n.HasMoved = true
			return true
		}
	}

	// If NPC can't move in any direction, mark as moved anyway
	n.HasMoved = true
	return false
}

// Manager handles a collection of NPCs
type Manager struct {
	NPCs []*NPC
}

// NewManager creates a new NPC manager
func NewManager() *Manager {
	return &Manager{
		NPCs: make([]*NPC, 0),
	}
}

// AddNPC adds an NPC to the manager
func (m *Manager) AddNPC(npc *NPC) {
	m.NPCs = append(m.NPCs, npc)
}

// AnyMoving checks if any NPC is currently moving
func (m *Manager) AnyMoving() bool {
	for _, npc := range m.NPCs {
		if npc.Moving {
			return true
		}
	}
	return false
}

// AllMoved checks if all NPCs have moved this turn
func (m *Manager) AllMoved() bool {
	for _, npc := range m.NPCs {
		if !npc.HasMoved {
			return false
		}
	}
	return true
}

// ResetMovedStatus resets the moved status for all NPCs
func (m *Manager) ResetMovedStatus() {
	for _, npc := range m.NPCs {
		npc.ResetMovedStatus()
	}
}

// ProcessTurn processes the turn for one NPC that hasn't moved yet
// Returns true if an NPC moved
func (m *Manager) ProcessTurn(validMoveFn func(x, y int) bool) bool {
	if m.AnyMoving() {
		return false // Wait for movement to complete
	}

	// If all NPCs have moved, we're done with the turn
	if m.AllMoved() {
		return false
	}

	// Process NPCs that haven't moved yet
	for _, npc := range m.NPCs {
		if !npc.HasMoved && !npc.Moving {
			if npc.TryMove(validMoveFn) {
				return true // An NPC moved
			}
		}
	}

	return false // No NPCs could move
}

// UpdatePositions updates positions for all NPCs
// Returns a slice of NPCs that reached their destinations this frame
func (m *Manager) UpdatePositions(moveSpeed float64) []*NPC {
	arrivedNPCs := make([]*NPC, 0)
	
	for _, npc := range m.NPCs {
		if npc.UpdatePosition(moveSpeed) {
			arrivedNPCs = append(arrivedNPCs, npc)
		}
	}
	
	return arrivedNPCs
}