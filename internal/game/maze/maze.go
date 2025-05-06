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
	TileSize = 40  // Size of each tile in the maze
)

// TileType represents different types of tiles in the maze
type TileType int

const (
	Floor TileType = iota
	Wall
	Goal
	// New tile types can be added here as needed
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

func New(width, height int, centerX int, centerY int) *Maze {
	goalX := width - 2  // Goal near the bottom-right corner
	goalY := height - 2

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

// Add these functions to the Maze struct
func createMazeGrid(width, height, goalX, goalY int) [][]MazeTile {
	grid := make([][]MazeTile, height)
	for y := range grid {
		grid[y] = make([]MazeTile, width)
		for x := range grid[y] {
			// Create walls around the edges and some random walls
			if x == 0 || y == 0 || x == width-1 || y == height-1 || (rand.Intn(100) < 20 && x > 1 && y > 1) {
				grid[y][x].tileType = Wall
			} else {
				grid[y][x].tileType = Floor
			}
		}
	}

	// Ensure the starting positions are not walls
	grid[1][1].tileType = Floor // Player start
	grid[3][3].tileType = Floor // NPC1 start
	grid[5][5].tileType = Floor // NPC2 start

	// Add the goal tile
	grid[goalY][goalX].tileType = Goal

	return grid
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

// Update the Draw method to show highlights
func (m *Maze) Draw(screen *ebiten.Image) {
	// Draw grid lines and tiles
	for y := 0; y < m.height; y++ {
		for x := 0; x < m.width; x++ {
			// Calculate tile position
			tileX := float64(x) * TileSize
			tileY := float64(y) * TileSize

			// Draw tile border
			borderColor := color.RGBA{100, 100, 100, 255}
			ebitenutil.DrawLine(screen, tileX, tileY, tileX+TileSize, tileY, borderColor)
			ebitenutil.DrawLine(screen, tileX, tileY, tileX, tileY+TileSize, borderColor)
			ebitenutil.DrawLine(screen, tileX+TileSize, tileY, tileX+TileSize, tileY+TileSize, borderColor)
			ebitenutil.DrawLine(screen, tileX, tileY+TileSize, tileX+TileSize, tileY+TileSize, borderColor)

			// Determine tile color based on type and highlight status
			var tileColor color.RGBA

			if m.grid[y][x].highlighted {
				// Highlighted tiles are red-tinted
				tileColor = color.RGBA{255, 100, 100, 255}
			} else {
				// Normal coloring based on tile type
				switch m.grid[y][x].tileType {
				case Wall:
					tileColor = color.RGBA{70, 70, 70, 255}
				case Goal:
					tileColor = color.RGBA{200, 0, 200, 255} // Purple goal
				default: // Floor
					tileColor = color.RGBA{200, 200, 200, 100}
				}
			}

			// Draw the tile
			ebitenutil.DrawRect(screen, tileX, tileY, TileSize, TileSize, tileColor)
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
