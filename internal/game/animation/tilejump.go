// internal/game/animation/tilejump.go
package animation

import (
	"math"
)

// TileAnimData holds animation data for a single tile
type TileAnimData struct {
	// Export fields so they're accessible from other packages
	StartX, StartY     float64
	EndX, EndY         float64
	CurrentX, CurrentY float64
	JumpHeight         float64
	TileSize           float64
}

// TileJumpAnimation animates tiles "jumping" during x-rotation
type TileJumpAnimation struct {
	active      bool
	duration    float64 // Total animation duration in seconds
	currentTime float64 // Current animation time
	
	// Tile positions and animation data
	tiles       []TileAnimData
	
	// Callback to apply final positions
	onComplete  func()
}

// NewTileJumpAnimation creates a new tile jump animation
func NewTileJumpAnimation(duration float64, onComplete func()) *TileJumpAnimation {
	return &TileJumpAnimation{
		active:      false,
		duration:    duration,
		currentTime: 0,
		tiles:       make([]TileAnimData, 0),
		onComplete:  onComplete,
	}
}

// SetTiles sets the tiles to animate
func (a *TileJumpAnimation) SetTiles(startPositions, endPositions [][2]int, tileSize float64) {
	a.tiles = make([]TileAnimData, len(startPositions))
	
	for i := 0; i < len(startPositions); i++ {
		a.tiles[i] = TileAnimData{
			StartX:     float64(startPositions[i][0]) * tileSize,
			StartY:     float64(startPositions[i][1]) * tileSize,
			EndX:       float64(endPositions[i][0]) * tileSize,
			EndY:       float64(endPositions[i][1]) * tileSize,
			CurrentX:   float64(startPositions[i][0]) * tileSize,
			CurrentY:   float64(startPositions[i][1]) * tileSize,
			JumpHeight: 20.0, // Jump height in pixels
			TileSize:   tileSize,
		}
	}
}

// Update updates the animation state
// Returns true when animation is complete
func (a *TileJumpAnimation) Update(deltaTime float64) bool {
	if !a.active {
		return true
	}
	
	a.currentTime += deltaTime
	progress := math.Min(a.currentTime / a.duration, 1.0)
	
	if progress >= 1.0 {
		// Animation complete
		a.active = false
		
		// Call completion callback
		if a.onComplete != nil {
			a.onComplete()
		}
		
		return true
	}
	
	// Update all tile positions
	for i := range a.tiles {
		// X movement (linear)
		a.tiles[i].CurrentX = a.tiles[i].StartX + (a.tiles[i].EndX - a.tiles[i].StartX) * progress
		
		// Y movement (parabolic jump)
		jumpProgress := math.Sin(progress * math.Pi) // 0->1->0 curve
		jumpOffset := a.tiles[i].JumpHeight * jumpProgress
		
		// Move from start to end position with jump
		a.tiles[i].CurrentY = a.tiles[i].StartY + (a.tiles[i].EndY - a.tiles[i].StartY) * progress - jumpOffset
	}
	
	return false
}

// GetTilePositions returns the current positions of all tiles
func (a *TileJumpAnimation) GetTilePositions() []TileAnimData {
	return a.tiles
}

// IsActive checks if the animation is currently running
func (a *TileJumpAnimation) IsActive() bool {
	return a.active
}

// Reset resets the animation to its initial state
func (a *TileJumpAnimation) Reset() {
	a.currentTime = 0
	a.active = true
	
	// Reset positions to starting points
	for i := range a.tiles {
		a.tiles[i].CurrentX = a.tiles[i].StartX
		a.tiles[i].CurrentY = a.tiles[i].StartY
	}
}