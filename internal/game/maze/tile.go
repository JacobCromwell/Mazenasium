// internal/game/maze/tile.go
package maze

// TileType represents different types of tiles in the maze
type TileType int

const (
    Floor TileType = iota
    Wall
    Goal
    SpecialTrigger // For tiles that trigger special events
    Trap           // For hazardous tiles
    // Add more types as needed
)

// Make these constants accessible to external packages
func (t TileType) String() string {
    switch t {
    case Floor:
        return "Floor"
    case Wall:
        return "Wall"
    case Goal:
        return "Goal"
    case SpecialTrigger:
        return "SpecialTrigger"
    case Trap:
        return "Trap"
    default:
        return "Unknown"
    }
}

// Tile represents a single cell in the maze
type Tile struct {
    ID          int
    Type        TileType
    FlavorImage string
    X, Y        int
    Highlighted bool
    Visited     bool // Used during maze generation
    
    // Additional properties can be added as needed
}

// NewTile creates a new tile with the specified type and position
func NewTile(id int, tileType TileType, x, y int) *Tile {
    return &Tile{
        ID:          id,
        Type:        tileType,
        X:           x,
        Y:           y,
        FlavorImage: "", // Default empty, to be set later
        Highlighted: false,
        Visited:     false,
    }
}

// IsWall checks if this tile is a wall
func (t *Tile) IsWall() bool {
    return t.Type == Wall
}

// IsGoal checks if this tile is a goal
func (t *Tile) IsGoal() bool {
    return t.Type == Goal
}

// IsFloor checks if this tile is a floor
func (t *Tile) IsFloor() bool {
    return t.Type == Floor
}

// SetFlavorImage sets the flavor image for this tile
func (t *Tile) SetFlavorImage(path string) {
    t.FlavorImage = path
}

// GetFlavorImage returns the flavor image path for this tile
func (t *Tile) GetFlavorImage() string {
    return t.FlavorImage
}