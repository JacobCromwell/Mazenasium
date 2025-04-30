// pkg/trivia/renderer.go
package trivia

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Renderer handles drawing trivia questions to the screen
type Renderer struct {
	backgroundColor color.RGBA
	textColor       color.RGBA
	correctColor    color.RGBA
	incorrectColor  color.RGBA
	padding         int
	width           int
	height          int
}

// NewRenderer creates a new trivia UI renderer
func NewRenderer(width, height int) *Renderer {
	return &Renderer{
		backgroundColor: color.RGBA{50, 50, 80, 240},
		textColor:       color.RGBA{255, 255, 255, 255},
		correctColor:    color.RGBA{0, 255, 0, 255},
		incorrectColor:  color.RGBA{255, 0, 0, 255},
		padding:         20,
		width:           width,
		height:          height,
	}
}

// Render draws the trivia question UI to the screen
func (r *Renderer) Render(screen *ebiten.Image, manager *Manager, x, y int) {
	question := manager.GetCurrentQuestion()
	if question == nil {
		return
	}

	// Calculate dimensions
	boxWidth := r.width - (r.padding * 2)
	boxHeight := r.height - (r.padding * 2)

	// Draw the background
	ebitenutil.DrawRect(screen, float64(x), float64(y), float64(r.width), float64(r.height), r.backgroundColor)

	// Draw the question
	ebitenutil.DebugPrintAt(screen, question.Text, x+r.padding, y+r.padding)

	// Draw the options
	optionY := y + r.padding + 50
	for i, option := range question.Options {
		optionText := fmt.Sprintf("%d: %s", i+1, option)
		ebitenutil.DebugPrintAt(screen, optionText, x+r.padding, optionY)
		optionY += 30
	}

	// Draw instructions
	instructY := y + r.height - r.padding - 30
	ebitenutil.DebugPrintAt(screen, "Press 1-4 to answer", x+r.padding, instructY)

	// If the question has been answered, show the result
	if manager.IsAnswered() && manager.GetLastResult() != nil {
		result := manager.GetLastResult()
		resultText := "Incorrect!"
		resultColor := r.incorrectColor

		if result.IsCorrect {
			resultText = "Correct!"
			resultColor = r.correctColor
		}

		centerX := x + (r.width / 2) - 40
		centerY := y + (r.height / 2)
		ebitenutil.DebugPrintAt(screen, resultText, centerX, centerY)
	}
}

// SetColors customizes the colors used by the renderer
func (r *Renderer) SetColors(background, text, correct, incorrect color.RGBA) {
	r.backgroundColor = background
	r.textColor = text
	r.correctColor = correct
	r.incorrectColor = incorrect
}

// SetDimensions sets the dimensions of the trivia UI
func (r *Renderer) SetDimensions(width, height int) {
	r.width = width
	r.height = height
}
