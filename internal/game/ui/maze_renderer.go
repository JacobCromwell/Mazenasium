// internal/game/ui/maze_renderer.go
package ui

import (
    "github.com/hajimehoshi/ebiten/v2"
    "github.com/hajimehoshi/ebiten/v2/ebitenutil"
    "image/color"
    
    "github.com/JacobCromwell/Mazenasium/internal/game/maze"
)

// DrawMaze renders the maze grid on the screen
func DrawMaze(screen *ebiten.Image, mazeObj *maze.Maze, offsetX, offsetY float64) {
    // For each tile in the maze state
    for y := 0; y < mazeObj.State.Height; y++ {
        for x := 0; x < mazeObj.State.Width; x++ {
            tile := mazeObj.State.GetTile(x, y)
            if tile == nil {
                continue
            }
            
            // Calculate tile position
            tileX := float64(x) * maze.TileSize + offsetX
            tileY := float64(y) * maze.TileSize + offsetY
            
            // Determine tile color based on type
            var tileColor color.RGBA
            switch tile.Type {
            case maze.Wall:
                tileColor = color.RGBA{70, 70, 70, 255}
            case maze.Goal:
                tileColor = color.RGBA{200, 0, 200, 255} // Purple goal
            default: // Floor
                tileColor = color.RGBA{200, 200, 200, 100}
            }
            
            // Draw the tile
            ebitenutil.DrawRect(screen, tileX, tileY, maze.TileSize, maze.TileSize, tileColor)
            
            // Draw highlighted tile with a 2px red outline instead of filling
            if tile.Highlighted {
                // Draw outline around the highlighted tile
                highlightColor := color.RGBA{255, 0, 0, 255} // Red outline
                
                // Draw 2px outlines
                ebitenutil.DrawRect(screen, tileX, tileY, maze.TileSize, 2, highlightColor) // Top
                ebitenutil.DrawRect(screen, tileX, tileY, 2, maze.TileSize, highlightColor) // Left
                ebitenutil.DrawRect(screen, tileX+maze.TileSize-2, tileY, 2, maze.TileSize, highlightColor) // Right
                ebitenutil.DrawRect(screen, tileX, tileY+maze.TileSize-2, maze.TileSize, 2, highlightColor) // Bottom
            }
            
            // Draw tile border
            borderColor := color.RGBA{100, 100, 100, 255}
            ebitenutil.DrawLine(screen, tileX, tileY, tileX+maze.TileSize, tileY, borderColor)
            ebitenutil.DrawLine(screen, tileX, tileY, tileX, tileY+maze.TileSize, borderColor)
            ebitenutil.DrawLine(screen, tileX+maze.TileSize, tileY, tileX+maze.TileSize, tileY+maze.TileSize, borderColor)
            ebitenutil.DrawLine(screen, tileX, tileY+maze.TileSize, tileX+maze.TileSize, tileY+maze.TileSize, borderColor)
        }
    }
}