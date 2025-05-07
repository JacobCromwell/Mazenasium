// internal/game/maze/maze.go - updated maze generation

package maze

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
)

// Constants used by maze package
const (
	Radius   = 120 // Radius of the circular maze
	TileSize = 30  // Size of each tile in the maze - reduced to fit more tiles
)

// TileType represents different types of tiles in the maze
type TileType int

const (
	Floor TileType = iota
	Wall
	Goal
	// New tile types can be added here as needed
)

// Direction represents movement directions for maze generation
type Direction int

const (
	North Direction = iota
	East
	South
	West
)

// Maze represents the maze grid
type Maze struct {
	grid          [][]MazeTile
	width, height int
	centerX       int
	centerY       int
	goalX, goalY  int // Track goal position
}

// Enhanced MazeTile struct with rotation properties
type MazeTile struct {
	tileType    TileType
	visited     bool
	hasItem     bool
	itemType    int
	highlighted bool // For showing which tiles will be affected by rotation
}

// Position represents a position in the grid
type Position struct {
	X, Y int
}

func New(width, height int, centerX int, centerY int) *Maze {
	// Increase maze size - double the number of tiles in both dimensions
	width = width * 2
	height = height * 2
	
	// Choose a random location for the goal that's not too close to the start
	var goalX, goalY int
	for {
		goalX = width - 2 - rand.Intn(width/3)  // Goal in the right third of the maze
		goalY = height - 2 - rand.Intn(height/3)  // Goal in the bottom third of the maze
		
		// Ensure the goal isn't too close to the start (Manhattan distance)
		if abs(goalX-1) + abs(goalY-1) >= (width + height)/2 {
			break
		}
	}

	return &Maze{
		width:         width,
		height:        height,
		centerX:       centerX,
		centerY:       centerY,
		goalX:         goalX,
		goalY:         goalY,
		grid:          createMazeGrid(width, height, goalX, goalY),
	}
}

// Helper function for absolute value
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// Add these functions to the Maze struct
func createMazeGrid(width, height, goalX, goalY int) [][]MazeTile {
	// Initialize grid with all walls
	grid := make([][]MazeTile, height)
	for y := range grid {
		grid[y] = make([]MazeTile, width)
		for x := range grid[y] {
			grid[y][x].tileType = Wall
		}
	}

	// Use randomized depth-first search algorithm to create pathways
	// This creates a perfect maze (exactly one path between any two points)
	generatePathways(grid, 1, 1, width, height)
	
	// Add some random additional paths to make it more interesting
	// This breaks the "perfect maze" property and creates multiple possible paths
	addRandomPaths(grid, width, height)

	// Ensure the starting positions are clear
	grid[1][1].tileType = Floor // Player start
	grid[3][3].tileType = Floor // NPC1 start
	grid[5][5].tileType = Floor // NPC2 start

	// Make sure there's a path to the goal
	ensurePathToGoal(grid, 1, 1, goalX, goalY, width, height)
	
	// Add the goal tile
	grid[goalY][goalX].tileType = Goal

	return grid
}

// Generate the initial maze using a randomized depth-first search
func generatePathways(grid [][]MazeTile, startX, startY, width, height int) {
	// Initialize visited grid
	visited := make([][]bool, height)
	for y := range visited {
		visited[y] = make([]bool, width)
	}

	// Stack for backtracking
	stack := []Position{{startX, startY}}
	
	// Mark the starting position as a floor and visited
	grid[startY][startX].tileType = Floor
	visited[startY][startX] = true

	// Directions: North, East, South, West
	dx := []int{0, 1, 0, -1}
	dy := []int{-1, 0, 1, 0}

	// Keep going until the stack is empty
	for len(stack) > 0 {
		// Get the current position
		current := stack[len(stack)-1]
		
		// Find unvisited neighbors
		neighbors := []int{}
		for d := 0; d < 4; d++ {
			nx, ny := current.X + dx[d]*2, current.Y + dy[d]*2
			
			// Check if the neighbor is valid and unvisited
			if nx >= 1 && nx < width-1 && ny >= 1 && ny < height-1 && !visited[ny][nx] {
				neighbors = append(neighbors, d)
			}
		}
		
		if len(neighbors) > 0 {
			// Randomly choose a neighbor
			d := neighbors[rand.Intn(len(neighbors))]
			nx, ny := current.X + dx[d]*2, current.Y + dy[d]*2
			
			// Carve a passage by setting both the neighbor and the wall in between to floor
			grid[current.Y + dy[d]][current.X + dx[d]].tileType = Floor
			grid[ny][nx].tileType = Floor
			
			// Mark as visited and push to stack
			visited[ny][nx] = true
			stack = append(stack, Position{nx, ny})
		} else {
			// No unvisited neighbors, backtrack
			stack = stack[:len(stack)-1]
		}
	}
}

// Add some random additional paths to make the maze more interesting
func addRandomPaths(grid [][]MazeTile, width, height int) {
	// Number of random paths to add (adjustable)
	extraPaths := (width + height) / 3
	
	for i := 0; i < extraPaths; i++ {
		// Pick a random wall that's not on the border
		x, y := 0, 0
		for {
			x = rand.Intn(width-2) + 1
			y = rand.Intn(height-2) + 1
			
			if grid[y][x].tileType == Wall && 
			   x > 0 && x < width-1 && y > 0 && y < height-1 {
				break
			}
		}
		
		// Count adjacent floor tiles
		floorCount := 0
		if grid[y-1][x].tileType == Floor { floorCount++ }
		if grid[y+1][x].tileType == Floor { floorCount++ }
		if grid[y][x-1].tileType == Floor { floorCount++ }
		if grid[y][x+1].tileType == Floor { floorCount++ }
		
		// Only remove walls that connect two different passages
		// This creates loops in the maze
		if floorCount >= 2 {
			grid[y][x].tileType = Floor
		}
	}
}

// Ensure there's a path from start to goal
func ensurePathToGoal(grid [][]MazeTile, startX, startY, goalX, goalY, width, height int) {
	// Use breadth-first search to check if there's a path
	if hasPath(grid, startX, startY, goalX, goalY, width, height) {
		return
	}
	
	// If no path exists, create one
	currentX, currentY := startX, startY
	
	// Move toward the goal with a slight randomness
	for currentX != goalX || currentY != goalY {
		// Decide whether to move in X or Y direction
		moveX := rand.Intn(2) == 0
		
		if moveX && currentX != goalX {
			// Move in X direction
			dx := 1
			if currentX > goalX {
				dx = -1
			}
			grid[currentY][currentX + dx].tileType = Floor
			currentX += dx
		} else if currentY != goalY {
			// Move in Y direction
			dy := 1
			if currentY > goalY {
				dy = -1
			}
			grid[currentY + dy][currentX].tileType = Floor
			currentY += dy
		}
		
		// Set current position to floor
		grid[currentY][currentX].tileType = Floor
	}
}

// Check if there's a path from start to goal
func hasPath(grid [][]MazeTile, startX, startY, goalX, goalY, width, height int) bool {
	// Initialize visited grid
	visited := make([][]bool, height)
	for y := range visited {
		visited[y] = make([]bool, width)
	}
	
	// Queue for BFS
	queue := []Position{{startX, startY}}
	visited[startY][startX] = true
	
	// Directions: North, East, South, West
	dx := []int{0, 1, 0, -1}
	dy := []int{-1, 0, 1, 0}
	
	// BFS
	for len(queue) > 0 {
		// Get the current position
		current := queue[0]
		queue = queue[1:]
		
		// Check if we reached the goal
		if current.X == goalX && current.Y == goalY {
			return true
		}
		
		// Check all four directions
		for d := 0; d < 4; d++ {
			nx, ny := current.X + dx[d], current.Y + dy[d]
			
			// Check if the neighbor is valid, a floor, and unvisited
			if nx >= 0 && nx < width && ny >= 0 && ny < height && 
			   grid[ny][nx].tileType != Wall && !visited[ny][nx] {
				visited[ny][nx] = true
				queue = append(queue, Position{nx, ny})
			}
		}
	}
	
	return false
}

// HighlightXRotation highlights tiles that would be affected by X-rotation
func (m *Maze) HighlightXRotation(playerX, playerY int) {
	// Clear any existing highlights first
	m.ClearHighlights()

	// Highlight all tiles in the same row as the player (excluding the player's position)
	for x := 0; x < m.width; x++ {
		// Skip player position
		if x == playerX {
			continue
		}

		// Only highlight tiles that can be rotated (not walls at the edge)
		if x > 0 && x < m.width-1 {
			m.grid[playerY][x].highlighted = true
		}
	}
}

// ClearHighlights removes all highlighting
func (m *Maze) ClearHighlights() {
	for y := 0; y < m.height; y++ {
		for x := 0; x < m.width; x++ {
			m.grid[y][x].highlighted = false
		}
	}
}

// CheckXRotateCollisions checks if the rotation would cause any walls to move onto entities
// entityPositions is a slice of positions where entities (player, NPCs) are located
// Returns true if there's a collision, false otherwise
func (m *Maze) CheckXRotateCollisions(playerX, playerY, direction int, entityPositions []Position) bool {
	// Make a copy of the grid for simulation
	tempGrid := make([][]MazeTile, m.height)
	for y := range m.grid {
		tempGrid[y] = make([]MazeTile, m.width)
		copy(tempGrid[y], m.grid[y])
	}
	
	// Create a map of entity positions for quick lookup
	entityMap := make(map[Position]bool)
	for _, pos := range entityPositions {
		entityMap[pos] = true
	}
	
	// Get a copy of the current row for rotation simulation
	tempRow := make([]MazeTile, m.width)
	for x := 0; x < m.width; x++ {
		tempRow[x] = m.grid[playerY][x]
	}
	
	// Simulate the rotation and check for collisions
	for x := 1; x < m.width-1; x++ {
		// Skip player position
		if x == playerX {
			continue
		}
		
		// Calculate new position after rotation
		newX := x + direction
		
		// Handle wrapping within the playable area (excluding boundary walls)
		if newX <= 0 {
			newX = m.width - 2
		} else if newX >= m.width-1 {
			newX = 1
		}
		
		// Check if we're moving a wall onto an entity position
		if tempRow[x].tileType == Wall {
			// Check if there's an entity at the destination
			pos := Position{newX, playerY}
			if entityMap[pos] {
				return true // Collision detected!
			}
		}
	}
	
	return false // No collisions
}

// PerformXRotate performs the actual rotation of tiles on the X-axis
func (m *Maze) PerformXRotate(playerX, playerY, direction int) {
	// Create a copy of the current row for rotation
	tempRow := make([]MazeTile, m.width)
	for x := 0; x < m.width; x++ {
		tempRow[x] = m.grid[playerY][x]
	}

	// Perform the rotation for each tile (except boundary walls)
	for x := 1; x < m.width-1; x++ {
		// Skip the player's position
		if x == playerX {
			continue
		}

		// Calculate new position
		newX := x + direction

		// Handle wrapping within the playable area (excluding boundary walls)
		if newX <= 0 {
			newX = m.width - 2
		} else if newX >= m.width-1 {
			newX = 1
		}

		// Move the tile
		m.grid[playerY][newX] = tempRow[x]
	}

	// Clear highlights after rotation
	m.ClearHighlights()
}

// Update the Draw method to show highlights as outlines
func (m *Maze) Draw(screen *ebiten.Image) {
	// Draw grid lines and tiles
	for y := 0; y < m.height; y++ {
		for x := 0; x < m.width; x++ {
			// Calculate tile position
			tileX := float64(x) * TileSize
			tileY := float64(y) * TileSize

			// Determine tile color based on type
			var tileColor color.RGBA
			switch m.grid[y][x].tileType {
			case Wall:
				tileColor = color.RGBA{70, 70, 70, 255}
			case Goal:
				tileColor = color.RGBA{200, 0, 200, 255} // Purple goal
			default: // Floor
				tileColor = color.RGBA{200, 200, 200, 100}
			}

			// Draw the tile
			ebitenutil.DrawRect(screen, tileX, tileY, TileSize, TileSize, tileColor)
			
			// Draw highlighted tile with a 2px red outline instead of filling
			if m.grid[y][x].highlighted {
				// Draw 2px red outline around the highlighted tile
				highlightColor := color.RGBA{255, 0, 0, 255} // Red outline
				
				// Top outline
				ebitenutil.DrawRect(screen, tileX, tileY, TileSize, 2, highlightColor)
				// Left outline
				ebitenutil.DrawRect(screen, tileX, tileY, 2, TileSize, highlightColor)
				// Right outline
				ebitenutil.DrawRect(screen, tileX+TileSize-2, tileY, 2, TileSize, highlightColor)
				// Bottom outline
				ebitenutil.DrawRect(screen, tileX, tileY+TileSize-2, TileSize, 2, highlightColor)
			}
			
			// Draw tile border
			borderColor := color.RGBA{100, 100, 100, 255}
			ebitenutil.DrawLine(screen, tileX, tileY, tileX+TileSize, tileY, borderColor)
			ebitenutil.DrawLine(screen, tileX, tileY, tileX, tileY+TileSize, borderColor)
			ebitenutil.DrawLine(screen, tileX+TileSize, tileY, tileX+TileSize, tileY+TileSize, borderColor)
			ebitenutil.DrawLine(screen, tileX, tileY+TileSize, tileX+TileSize, tileY+TileSize, borderColor)
		}
	}
}

// IsWall checks if the given coordinates are a wall
func (m *Maze) IsWall(x, y int) bool {
	if x < 0 || x >= m.width || y < 0 || y >= m.height {
		return true
	}
	return m.grid[y][x].tileType == Wall
}

// IsGoal checks if the given coordinates are the goal
func (m *Maze) IsGoal(x, y int) bool {
	if x < 0 || x >= m.width || y < 0 || y >= m.height {
		return false
	}
	return m.grid[y][x].tileType == Goal
}

// IsValidMove checks if a move to the given coordinates is valid
func (m *Maze) IsValidMove(x, y int) bool {
	return x >= 0 && x < m.width &&
		y >= 0 && y < m.height &&
		m.grid[y][x].tileType != Wall
}