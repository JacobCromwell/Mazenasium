package maze

import (
	"math/rand"
)

// Maze is the main interface for interacting with the maze
type Maze struct {
    State     *State
    Generator *Generator
}

// New creates a new maze with the specified dimensions
// func New(width, height int, centerX, centerY int) *Maze {
//     // Double the dimensions for a larger maze
//     width = width * 2
//     height = height * 2
    
//     // Create a generator with a random seed
//     generator := NewGenerator(rand.Int63())
    
//     // Generate the initial maze state
//     state := generator.Generate(width, height)
    
//     return &Maze{
//         State:     state,
//         Generator: generator,
//     }
// }

func New(width, height int, centerX, centerY int) *Maze {
    // Double the dimensions for a larger maze
    width = width * 2
    height = height * 2
    
    // Create a generator with a random seed
    generator := NewGenerator(rand.Int63())
    
    // Generate the initial maze state
    state := generator.Generate(width, height)
    
    return &Maze{
        State:     state,
        Generator: generator,
    }
}

// IsWall checks if the given coordinates are a wall
func (m *Maze) IsWall(x, y int) bool {
    tile := m.State.GetTile(x, y)
    return tile == nil || tile.IsWall()
}

// IsGoal checks if the given coordinates are the goal
func (m *Maze) IsGoal(x, y int) bool {
    tile := m.State.GetTile(x, y)
    return tile != nil && tile.IsGoal()
}

// IsValidMove checks if a move to the given coordinates is valid
func (m *Maze) IsValidMove(x, y int) bool {
    return m.State.IsValidMove(x, y)
}

// GetFlavorImageForTile returns the flavor image path for the tile at the given position
func (m *Maze) GetFlavorImageForTile(x, y int) string {
    tile := m.State.GetTile(x, y)
    if tile == nil {
        return ""
    }
    return tile.GetFlavorImage()
}

// HighlightXRotation highlights tiles that would be affected by X-rotation
func (m *Maze) HighlightXRotation(playerX, playerY int) {
    m.State.HighlightXRotation(playerX, playerY)
}

// ClearHighlights removes all highlighting
func (m *Maze) ClearHighlights() {
    m.State.ClearHighlights()
}

// PerformXRotate performs the rotation of tiles on the X-axis
func (m *Maze) PerformXRotate(playerX, playerY, direction int) {
    m.State.PerformXRotate(playerX, playerY, direction)
}

// GetTileSize returns the size of each tile in pixels
func (m *Maze) GetTileSize() float64 {
    return TileSize
}

// internal/game/maze/maze.go
// Add this method to the Maze struct

// CheckXRotateCollisions checks if the rotation would cause any walls to move onto entities
// entityPositions is a slice of positions where entities (player, NPCs) are located
// Returns true if there's a collision, false otherwise
func (m *Maze) CheckXRotateCollisions(playerX, playerY, direction int, entityPositions []Position) bool {
    if playerY < 0 || playerY >= m.State.Height {
        return false
    }
    
    // Create a map of entity positions for quick lookup
    entityMap := make(map[Position]bool)
    for _, pos := range entityPositions {
        entityMap[pos] = true
    }
    
    // Get a copy of the current row for rotation simulation
    row := m.State.Grid[playerY]
    
    // Simulate the rotation and check for collisions
    for x := 1; x < m.State.Width-1; x++ {
        // Skip player position
        if x == playerX {
            continue
        }
        
        // Calculate new position after rotation
        newX := x + direction
        
        // Handle wrapping within the playable area (excluding boundary walls)
        if newX <= 0 {
            newX = m.State.Width - 2
        } else if newX >= m.State.Width-1 {
            newX = 1
        }
        
        // Check if we're moving a wall onto an entity position
        if row[x].Type == Wall {
            // Check if there's an entity at the destination
            pos := Position{X: newX, Y: playerY}
            if entityMap[pos] {
                return true // Collision detected!
            }
        }
    }
    
    return false // No collisions
}