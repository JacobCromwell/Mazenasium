package maze

import (
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
)

// Constants used by maze package
const (
	Radius = 120 // Radius of the circular maze
	TileSize = 40  // Size of each tile in the maze
)

// TileType represents different types of tiles in the maze
type TileType int

const (
	Floor TileType = iota
	Wall
	Goal
)

// Maze represents the maze grid
type Maze struct {
	grid          [][]MazeTile
	width, height int
	centerX       float64
	centerY       float64
	rotationAngle float64
	goalX, goalY  int // Track goal position
}

// MazeTile represents a single tile in the maze
type MazeTile struct {
	tileType TileType
	visited  bool
	hasItem  bool
	itemType int
}

// New creates a new maze with the specified dimensions
func New(width, height int, centerX, centerY float64) *Maze {
	goalX := width - 2  // Goal near the bottom-right corner
	goalY := height - 2

	return &Maze{
		width:         width,
		height:        height,
		centerX:       centerX,
		centerY:       centerY,
		rotationAngle: 0,
		goalX:         goalX,
		goalY:         goalY,
		grid:          createMazeGrid(width, height, goalX, goalY),
	}
}

// Create a simple maze grid with a goal
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

// Draw the maze grid
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

			// Draw different tile types
			switch m.grid[y][x].tileType {
			case Wall:
				ebitenutil.DrawRect(screen, tileX, tileY, TileSize, TileSize, color.RGBA{70, 70, 70, 255})
			case Goal:
				ebitenutil.DrawRect(screen, tileX, tileY, TileSize, TileSize, color.RGBA{200, 0, 200, 255}) // Purple goal
			default: // Floor
				ebitenutil.DrawRect(screen, tileX, tileY, TileSize, TileSize, color.RGBA{200, 200, 200, 100})
			}
		}
	}
}

// DrawCircular draws the circular maze in the corner
func (m *Maze) DrawCircular(screen *ebiten.Image, playerGridX, playerGridY int) {
	// Draw outer circle
	ebitenutil.DrawCircle(screen, m.centerX, m.centerY, Radius, color.RGBA{200, 200, 200, 100})

	// Draw a simplified representation of the maze in the circle
	cellAngle := 2 * math.Pi / float64(m.width)
	cellRadius := Radius / float64(m.height)

	for y := 0; y < m.height; y++ {
		radius := float64(y+1) * cellRadius

		for x := 0; x < m.width; x++ {
			angle := m.rotationAngle + float64(x)*cellAngle

			// Calculate position
			cellX := m.centerX + math.Cos(angle)*radius
			cellY := m.centerY + math.Sin(angle)*radius

			// Draw different cell types
			switch m.grid[y][x].tileType {
			case Wall:
				ebitenutil.DrawCircle(screen, cellX, cellY, cellRadius/2, color.RGBA{70, 70, 70, 255})
			case Goal:
				ebitenutil.DrawCircle(screen, cellX, cellY, cellRadius/2, color.RGBA{200, 0, 200, 255}) // Purple goal
			}
		}
	}

	// Draw player position in the minimap
	playerAngle := m.rotationAngle + float64(playerGridX)*cellAngle
	playerRadius := float64(playerGridY+1) * cellRadius
	playerMiniX := m.centerX + math.Cos(playerAngle)*playerRadius
	playerMiniY := m.centerY + math.Sin(playerAngle)*playerRadius
	ebitenutil.DrawCircle(screen, playerMiniX, playerMiniY, cellRadius/2, color.RGBA{0, 0, 255, 255})

	// Draw rotation controls
	ebitenutil.DebugPrintAt(screen, "Q/E: Rotate", int(m.centerX)-40, int(m.centerY)+Radius+10)
}

// RotateLeft rotates the maze left
func (m *Maze) RotateLeft() {
	m.rotationAngle -= 0.05
}

// RotateRight rotates the maze right
func (m *Maze) RotateRight() {
	m.rotationAngle += 0.05
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