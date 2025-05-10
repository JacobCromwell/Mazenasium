// internal/game/ui/text.go
package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
)

var (
	// DefaultFont is the default font to use for text rendering
	DefaultFont font.Face = basicfont.Face7x13
	
	// FontScale is the scaling factor for text rendering (2.0 = twice as big)
	FontScale float64 = 2.0
	
	// DefaultTextColor is the default color for text (light yellow)
	DefaultTextColor = color.RGBA{255, 250, 205, 255} // Light yellow
	
	// OutlineColor is the color for text outlines
	OutlineColor = color.RGBA{0, 0, 0, 255} // Black
)

// DrawTextWithOutline draws text with a 1px outline
func DrawTextWithOutline(screen *ebiten.Image, s string, x, y int, textColor, outlineColor color.Color) {
	if s == "" {
		return // Don't try to render an empty string
	}
	
	// Calculate the text bounds
	bounds := text.BoundString(DefaultFont, s)
	w := bounds.Dx()
	h := bounds.Dy()
	
	// Add extra space for the outline
	w += 2
	h += 2
	
	// Ensure width and height are at least 1 pixel
	if w < 1 {
		w = 1
	}
	if h < 1 {
		h = 1
	}
	
	// Create a temporary image for rendering the text with outline
	textImage := ebiten.NewImage(w, h)
	
	// Calculate the base position, adjusted for the outline
	baseX := 1 // Add 1px for outline
	baseY := -bounds.Min.Y + 1 // Adjust position and add 1px for outline
	
	// Draw the outline by drawing the text multiple times at offset positions
	// Top-left
	text.Draw(textImage, s, DefaultFont, baseX-1, baseY-1, outlineColor)
	// Top
	text.Draw(textImage, s, DefaultFont, baseX, baseY-1, outlineColor)
	// Top-right
	text.Draw(textImage, s, DefaultFont, baseX+1, baseY-1, outlineColor)
	// Left
	text.Draw(textImage, s, DefaultFont, baseX-1, baseY, outlineColor)
	// Right
	text.Draw(textImage, s, DefaultFont, baseX+1, baseY, outlineColor)
	// Bottom-left
	text.Draw(textImage, s, DefaultFont, baseX-1, baseY+1, outlineColor)
	// Bottom
	text.Draw(textImage, s, DefaultFont, baseX, baseY+1, outlineColor)
	// Bottom-right
	text.Draw(textImage, s, DefaultFont, baseX+1, baseY+1, outlineColor)
	
	// Draw the main text in the center
	text.Draw(textImage, s, DefaultFont, baseX, baseY, textColor)
	
	// Draw the scaled text to the screen
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(FontScale, FontScale)
	op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(textImage, op)
}

// DrawTextColor draws text with a color but no outline
func DrawTextColor(screen *ebiten.Image, s string, x, y int, clr color.Color) {
	// Now just calls DrawTextWithOutline with outline color set to transparent
	transparent := color.RGBA{0, 0, 0, 0}
	DrawTextWithOutline(screen, s, x, y, clr, transparent)
}

// DrawText draws text with the default color and a black outline
func DrawText(screen *ebiten.Image, s string, x, y int) {
	DrawTextWithOutline(screen, s, x, y, DefaultTextColor, OutlineColor)
}

// Helper function for Go versions earlier than 1.21 which might not have max in the standard library
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}