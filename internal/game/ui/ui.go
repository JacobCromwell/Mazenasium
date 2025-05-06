// internal/game/ui/ui.go
package ui

import (
	"fmt"
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"github.com/JacobCromwell/Mazenasium/internal/game/action"
	"github.com/JacobCromwell/Mazenasium/internal/game/maze"
	"github.com/JacobCromwell/Mazenasium/internal/game/npc"
	"github.com/JacobCromwell/Mazenasium/internal/game/player"
	"github.com/JacobCromwell/Mazenasium/internal/game/trivia"
	"github.com/JacobCromwell/Mazenasium/internal/game/turn"
)

const (
	ScreenWidth  = 800
	ScreenHeight = 600
)

// Renderer handles all UI rendering for the game
type Renderer struct {
	actionMsg   string
	actionTimer int
}

// NewRenderer creates a new UI renderer
func NewRenderer() *Renderer {
	return &Renderer{
		actionMsg:   "",
		actionTimer: 0,
	}
}

// SetActionMessage sets a temporary action message to display
func (r *Renderer) SetActionMessage(msg string, duration int) {
	r.actionMsg = msg
	r.actionTimer = duration
}

// UpdateActionTimer updates the action message timer
func (r *Renderer) UpdateActionTimer() {
	if r.actionTimer > 0 {
		r.actionTimer--
		if r.actionTimer == 0 {
			r.actionMsg = ""
		}
	}
}

// Draw renders the entire game based on state
func (r *Renderer) Draw(
	screen *ebiten.Image,
	gameState int, // GameState
	mazeObj *maze.Maze,
	playerObj *player.Player,
	npcManager *npc.Manager,
	turnManager *turn.Manager,
	triviaManager *trivia.Manager,
	actionManager *action.Manager, // Added action manager
	winner string,
) {
	// Draw background
	screen.Fill(color.RGBA{40, 45, 55, 255})

	switch gameState {
	case 0: // Playing
		r.drawPlaying(screen, mazeObj, playerObj, npcManager, turnManager, actionManager)
	case 1: // AnsweringTrivia
		r.drawTrivia(screen, triviaManager)
	case 2: // GameOver
		r.drawGameOver(screen, winner)
	}
}

// Draw the game over screen
func (r *Renderer) drawGameOver(screen *ebiten.Image, winner string) {
	// Draw message background
	ebitenutil.DrawRect(screen, 100, 200, ScreenWidth-200, 100, color.RGBA{50, 50, 80, 240})
	
	// Draw winner message
	winMessage := fmt.Sprintf("%s reached the goal first and won!", winner)
	ebitenutil.DebugPrintAt(screen, winMessage, ScreenWidth/2-120, ScreenHeight/2-10)
	ebitenutil.DebugPrintAt(screen, "Press SPACE to restart", ScreenWidth/2-100, ScreenHeight/2+20)
}

// Draw the playing state
func (r *Renderer) drawPlaying(
	screen *ebiten.Image,
	mazeObj *maze.Maze,
	playerObj *player.Player,
	npcManager *npc.Manager,
	turnManager *turn.Manager,
	actionManager *action.Manager, // Added action manager
) {
	// Draw the maze grid
	mazeObj.Draw(screen)

	// Draw NPCs
	for _, npc := range npcManager.NPCs {
		ebitenutil.DrawRect(screen, npc.X+1, npc.Y+1, npc.Size, npc.Size, npc.Color)
	}

	// Draw player
	playerX, playerY := playerObj.GetPosition()
	ebitenutil.DrawRect(screen, playerX+1, playerY+1, playerObj.Size, playerObj.Size, color.RGBA{0, 0, 255, 255})

	// Draw UI info
	r.drawUI(screen, turnManager)

	// Draw action selection popup if in SelectingAction state
	if turnManager.CurrentState == turn.SelectingAction {
		r.drawActionPopup(screen, actionManager)
	}

	// Draw action message if active
	if r.actionMsg != "" {
		ebitenutil.DebugPrintAt(screen, r.actionMsg, ScreenWidth/2-150, ScreenHeight-50)
	}
}

// Draw the UI
func (r *Renderer) drawUI(screen *ebiten.Image, turnManager *turn.Manager) {
	// Draw turn info using the turn manager
	ebitenutil.DebugPrintAt(screen, turnManager.OwnerText(), 10, 10)

	// Draw turn state info using the turn manager
	ebitenutil.DebugPrintAt(screen, turnManager.StateText(), 10, 30)

	// Draw goal info
	ebitenutil.DebugPrintAt(screen, "Reach the purple goal to win!", 10, 50)
}

// Draw the action selection popup
func (r *Renderer) drawActionPopup(screen *ebiten.Image, actionManager *action.Manager) {
	// Get formatted list of available actions
	actionText := actionManager.FormatActionsList()
	lines := strings.Split(actionText, "\n")
	
	// Calculate popup dimensions based on content
	width := 300
	height := 40 + (len(lines) * 20)
	x := (ScreenWidth - width) / 2
	y := (ScreenHeight - height) / 2
	
	// Draw popup background
	ebitenutil.DrawRect(screen, float64(x), float64(y), float64(width), float64(height), color.RGBA{70, 70, 100, 240})
	ebitenutil.DrawRect(screen, float64(x+2), float64(y+2), float64(width-4), float64(height-4), color.RGBA{40, 40, 70, 240})
	
	// Draw action list
	for i, line := range lines {
		ebitenutil.DebugPrintAt(screen, line, x+10, y+20+(i*20))
	}
	
	// Draw instructions at the bottom
	ebitenutil.DebugPrintAt(screen, "Press number to select, ESC to cancel", x+10, y+height-20)
}

// Draw the trivia screen
func (r *Renderer) drawTrivia(screen *ebiten.Image, triviaManager *trivia.Manager) {
	currentQuestion := triviaManager.GetCurrentQuestion()

	// Draw question background
	ebitenutil.DrawRect(screen, 50, 50, ScreenWidth-100, ScreenHeight-100, color.RGBA{50, 50, 80, 240})

	// Draw question
	ebitenutil.DebugPrintAt(screen, currentQuestion.Question, 70, 70)

	// Draw options
	for i, option := range currentQuestion.Options {
		optionYpadding := 30 * i
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d: %s", i+1, option), 70, (140 + optionYpadding))
	}

	// Draw instructions
	ebitenutil.DebugPrintAt(screen, "Press 1-4 to answer", 70, ScreenHeight-100)

	// If answered, show result
	if triviaManager.Answered {
		resultText := "Incorrect!"
		//resultColor := color.RGBA{255, 0, 0, 255}

		if triviaManager.Correct {
			resultText = "Correct!"
			//resultColor = color.RGBA{0, 255, 0, 255}
		}

		ebitenutil.DebugPrintAt(screen, resultText, ScreenWidth/2-40, ScreenHeight/2)
	}
}