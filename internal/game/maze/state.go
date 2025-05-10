package maze

// State manages the current state of the maze
type State struct {
    Grid      [][]*Tile
    Width     int
    Height    int
    GoalX     int
    GoalY     int
}

// NewState creates a new maze state with the given dimensions
func NewState(width, height int) *State {
    grid := make([][]*Tile, height)
    for y := range grid {
        grid[y] = make([]*Tile, width)
        for x := range grid[y] {
            // Initialize with wall tiles by default
            tileID := y*width + x
            grid[y][x] = NewTile(tileID, Wall, x, y)
        }
    }
    
    return &State{
        Grid:   grid,
        Width:  width,
        Height: height,
        GoalX:  -1, // To be set later
        GoalY:  -1, // To be set later
    }
}

// GetTile returns the tile at the specified position
func (s *State) GetTile(x, y int) *Tile {
    if x < 0 || x >= s.Width || y < 0 || y >= s.Height {
        return nil
    }
    return s.Grid[y][x]
}

// SetTileType sets the type of the tile at the specified position
func (s *State) SetTileType(x, y int, tileType TileType) {
    if x < 0 || x >= s.Width || y < 0 || y >= s.Height {
        return
    }
    s.Grid[y][x].Type = tileType
    
    // If setting a goal tile, update the goal position
    if tileType == Goal {
        s.GoalX = x
        s.GoalY = y
    }
}

// IsValidMove checks if a move to the given coordinates is valid
func (s *State) IsValidMove(x, y int) bool {
    tile := s.GetTile(x, y)
    return tile != nil && !tile.IsWall()
}

// HighlightXRotation highlights tiles that would be affected by X-rotation
func (s *State) HighlightXRotation(playerX, playerY int) {
    // Clear any existing highlights first
    s.ClearHighlights()

    // Highlight all tiles in the same row as the player (excluding the player's position)
    if playerY < 0 || playerY >= s.Height {
        return
    }
    
    for x := 0; x < s.Width; x++ {
        // Skip player position
        if x == playerX {
            continue
        }

        // Only highlight tiles that can be rotated (not walls at the edge)
        if x > 0 && x < s.Width-1 {
            s.Grid[playerY][x].Highlighted = true
        }
    }
}

// ClearHighlights removes all highlighting
func (s *State) ClearHighlights() {
    for y := 0; y < s.Height; y++ {
        for x := 0; x < s.Width; x++ {
            if s.Grid[y][x] != nil {
                s.Grid[y][x].Highlighted = false
            }
        }
    }
}

// PerformXRotate performs the rotation of tiles on the X-axis
func (s *State) PerformXRotate(playerX, playerY, direction int) {
    if playerY < 0 || playerY >= s.Height {
        return
    }
    
    // Create a copy of the current row for rotation
    tempRow := make([]*Tile, s.Width)
    for x := 0; x < s.Width; x++ {
        tempRow[x] = s.Grid[playerY][x]
    }

    // Perform the rotation for each tile (except boundary walls)
    for x := 1; x < s.Width-1; x++ {
        // Skip the player's position
        if x == playerX {
            continue
        }

        // Calculate new position
        newX := x + direction

        // Handle wrapping within the playable area (excluding boundary walls)
        if newX <= 0 {
            newX = s.Width - 2
        } else if newX >= s.Width-1 {
            newX = 1
        }

        // Move the tile
        s.Grid[playerY][newX] = tempRow[x]
        // Update the tile's position
        s.Grid[playerY][newX].X = newX
    }

    // Clear highlights after rotation
    s.ClearHighlights()
}