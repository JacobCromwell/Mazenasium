// internal/game/ui/layout.go
package ui

type SectionType int

const (
    MazeSection SectionType = iota
    FlavorSection
)

type Section struct {
    Type   SectionType
    Rect   Rect
    Border bool
    Title  string
}

type Rect struct {
    X, Y, Width, Height int
}

type LayoutManager struct {
    Sections map[SectionType]Section
    ScreenWidth, ScreenHeight int
}

func NewLayoutManager(screenWidth, screenHeight int) *LayoutManager {
    // Updated to put maze on the left and flavor image on the right
    mazeSection := Section{
        Type: MazeSection,
        Rect: Rect{
            X: screenWidth / 2, // Left side of the screen
            Y: 0,
            Width: screenWidth / 2,
            Height: screenHeight,
        },
        Border: true,
        Title: "Maze",
    }
    
    flavorSection := Section{
        Type: FlavorSection,
        Rect: Rect{
            X: 0, // Right side of screen
            Y: 0,
            Width: screenWidth - (screenWidth / 2),
            Height: screenHeight,
        },
        Border: true,
        Title: "View",
    }
    
    sections := make(map[SectionType]Section)
    sections[MazeSection] = mazeSection
    sections[FlavorSection] = flavorSection
    
    return &LayoutManager{
        Sections: sections,
        ScreenWidth: screenWidth,
        ScreenHeight: screenHeight,
    }
}

func (l *LayoutManager) GetSection(sectionType SectionType) Section {
    return l.Sections[sectionType]
}

// Additional methods to adjust layout if needed