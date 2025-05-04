package ui

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

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
	winner string,
) {
	// Draw background
	screen.Fill(color.RGBA{40, 45, 55, 255})

	switch gameState {
	case 0: // Playing
		r.drawPlaying(screen, mazeObj, playerObj, npcManager, turnManager)
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

	// Draw action message if active
	if r.actionMsg != "" {
		ebitenutil.DebugPrintAt(screen, r.actionMsg, ScreenWidth/2-50, ScreenHeight/2)
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