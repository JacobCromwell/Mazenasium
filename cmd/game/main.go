package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/JacobCromwell/Mazenasium/internal/game/maze"
	"github.com/JacobCromwell/Mazenasium/internal/game/npc"
	"github.com/JacobCromwell/Mazenasium/internal/game/player"
	"github.com/JacobCromwell/Mazenasium/internal/game/trivia"
	"github.com/JacobCromwell/Mazenasium/internal/game/turn"
)

const (
	screenWidth  = 800
	screenHeight = 600
)

// Game implements ebiten.Game interface.
type Game struct {
	gameState   GameState
	turnManager *turn.Manager // Changed from individual turn states
	player      *player.Player
	npcManager  *npc.Manager
	maze        *maze.Maze
	triviaMgr   *trivia.Manager
	actionMsg   string
	actionTimer int
	winner      string // Track who wins the game
}

// GameState represents the current state of the game
type GameState int

const (
	Playing GameState = iota
	AnsweringTrivia
	GameOver
)

// Initialize new game
func NewGame() *Game {
	mazeWidth := 10
	mazeHeight := 10

	game := &Game{
		gameState:   Playing,
		turnManager: turn.NewManager(), // Initialize turn manager
		player:      player.New(1, 1, maze.TileSize),
		npcManager:  npc.NewManager(),
		maze:        maze.New(mazeWidth, mazeHeight, screenWidth-maze.Radius-20, screenHeight-maze.Radius-20),
		triviaMgr:   trivia.NewManager(),
		actionMsg:   "",
		actionTimer: 0,
		winner:      "",
	}

	// Create NPCs
	npc1 := npc.New(0, 3, 3, maze.TileSize, color.RGBA{255, 0, 0, 255})
	npc2 := npc.New(1, 5, 5, maze.TileSize, color.RGBA{0, 255, 0, 255})
	
	// Add NPCs to manager
	game.npcManager.AddNPC(npc1)
	game.npcManager.AddNPC(npc2)

	return game
}

// Update game state
func (g *Game) Update() error {
	switch g.gameState {
	case Playing:
		g.updatePlaying()
	case AnsweringTrivia:
		g.updateTrivia()
	case GameOver:
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			// Reset game
			*g = *NewGame()
		}
	}

	// Update action message timer
	if g.actionTimer > 0 {
		g.actionTimer--
		if g.actionTimer == 0 {
			g.actionMsg = ""
		}
	}

	return nil
}

// Update while playing
func (g *Game) updatePlaying() {
	// Update positions for smooth movement
	g.updatePositions()

	// Process based on turn state
	switch g.turnManager.CurrentState {
	case turn.WaitingForMove:
		if g.turnManager.IsPlayerTurn() {
			g.handlePlayerMovement()
		} else {
			g.processNPCTurn()
		}

	case turn.WaitingForAction:
		if inpututil.IsKeyJustPressed(ebiten.KeyA) {
			g.actionMsg = "Action used!"
			g.actionTimer = 60 // Show message for 60 frames
			g.turnManager.NextState(turn.WaitingForEndTurn)
		}

		// Allow skipping the action
		if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
			g.turnManager.NextState(turn.WaitingForEndTurn)
		}

	case turn.WaitingForEndTurn:
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			// End turn and switch to next actor
			g.turnManager.EndTurn()
			// Reset NPC movement tracking for the new turn if switching to NPC turn
			if g.turnManager.CurrentOwner == turn.NPCTurn {
				g.npcManager.ResetMovedStatus()
			}
		}

	case turn.ProcessingNPCTurn:
		// If no NPCs are still moving, process their next action
		if !g.npcManager.AnyMoving() {
			g.processNPCTurn()
		}
	}
}

// Update positions for smooth movement
func (g *Game) updatePositions() {
	// Update player position with smooth movement
	playerGridX, playerGridY := g.player.GetGridPosition()
	
	// Update player, and check if they've arrived at destination
	if arrived := g.player.Update(5.0); arrived {
		// Check if player reached the goal
		if g.maze.IsGoal(playerGridX, playerGridY) {
			g.winner = "Player"
			g.gameState = GameOver
			return
		}

		// If this was a player move, show trivia
		if g.turnManager.IsPlayerTurn() && g.turnManager.CurrentState == turn.WaitingForMove {
			g.gameState = AnsweringTrivia
			g.turnManager.NextState(turn.WaitingForTrivia)
			g.triviaMgr.Answered = false
			g.triviaMgr.SetRandomQuestion(rand.Intn)
		}
	}

	// Update NPCs positions using the manager
	arrivedNPCs := g.npcManager.UpdatePositions(5.0)
	
	// Check if any NPCs reached the goal
	for _, arrivedNPC := range arrivedNPCs {
		if g.maze.IsGoal(arrivedNPC.GridX, arrivedNPC.GridY) {
			g.winner = fmt.Sprintf("NPC %d", arrivedNPC.ID+1)
			g.gameState = GameOver
			return
		}
	}
}

// Handle player movement
func (g *Game) handlePlayerMovement() {
	if g.player.IsMoving() {
		return // Already moving
	}

	playerGridX, playerGridY := g.player.GetGridPosition()
	newGridX, newGridY := playerGridX, playerGridY

	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		newGridY--
	} else if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		newGridY++
	} else if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		newGridX--
	} else if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		newGridX++
	} else {
		return // No movement key pressed
	}

	// Check if movement is valid (not a wall and within bounds)
	if g.maze.IsValidMove(newGridX, newGridY) {
		// Set destination for smooth movement
		g.player.SetDestination(newGridX, newGridY, maze.TileSize)
	}
}

// Process NPC turn using the NPC manager
func (g *Game) processNPCTurn() {
	// Check if all NPCs have moved
	if g.npcManager.AllMoved() {
		g.turnManager.EndTurn() // Switch back to player's turn
		return
	}

	// Process one NPC's turn using a callback to check valid moves
	validMoveFn := func(x, y int) bool {
		return g.maze.IsValidMove(x, y)
	}
	
	g.npcManager.ProcessTurn(validMoveFn)
}

// Update trivia state
func (g *Game) updateTrivia() {
	// Use the trivia package's HandleInput method
	if g.triviaMgr.HandleInput() {
		// Return to game after brief delay
		go func() {
			// Note: In a real game, you'd want to use a more robust timer or state system
			// This is just a simple demonstration
			g.gameState = Playing
			g.turnManager.NextState(turn.WaitingForAction)
		}()
	}
}

// Draw the game
func (g *Game) Draw(screen *ebiten.Image) {
	// Draw background
	screen.Fill(color.RGBA{40, 45, 55, 255})

	switch g.gameState {
	case Playing:
		g.drawPlaying(screen)
	case AnsweringTrivia:
		g.drawTrivia(screen)
	case GameOver:
		g.drawGameOver(screen)
	}
}

// Draw the game over screen
func (g *Game) drawGameOver(screen *ebiten.Image) {
	// Draw message background
	ebitenutil.DrawRect(screen, 100, 200, screenWidth-200, 100, color.RGBA{50, 50, 80, 240})
	
	// Draw winner message
	winMessage := fmt.Sprintf("%s reached the goal first and won!", g.winner)
	ebitenutil.DebugPrintAt(screen, winMessage, screenWidth/2-120, screenHeight/2-10)
	ebitenutil.DebugPrintAt(screen, "Press SPACE to restart", screenWidth/2-100, screenHeight/2+20)
}

// Draw the playing state
func (g *Game) drawPlaying(screen *ebiten.Image) {
	// Draw the maze grid
	g.maze.Draw(screen)

	// Draw NPCs
	for _, npc := range g.npcManager.NPCs {
		ebitenutil.DrawRect(screen, npc.X+1, npc.Y+1, npc.Size, npc.Size, npc.Color)
	}

	// Draw player
	playerX, playerY := g.player.GetPosition()
	ebitenutil.DrawRect(screen, playerX+1, playerY+1, g.player.Size, g.player.Size, color.RGBA{0, 0, 255, 255})

	// Draw UI info
	g.drawUI(screen)

	// Draw action message if active
	if g.actionMsg != "" {
		ebitenutil.DebugPrintAt(screen, g.actionMsg, screenWidth/2-50, screenHeight/2)
	}
}

// Draw the UI
func (g *Game) drawUI(screen *ebiten.Image) {
	// Draw turn info using the turn manager
	ebitenutil.DebugPrintAt(screen, g.turnManager.OwnerText(), 10, 10)

	// Draw turn state info using the turn manager
	ebitenutil.DebugPrintAt(screen, g.turnManager.StateText(), 10, 30)

	// Draw goal info
	ebitenutil.DebugPrintAt(screen, "Reach the purple goal to win!", 10, 50)
}

// Draw the trivia screen
func (g *Game) drawTrivia(screen *ebiten.Image) {
	currentQuestion := g.triviaMgr.GetCurrentQuestion()

	// Draw question background
	ebitenutil.DrawRect(screen, 50, 50, screenWidth-100, screenHeight-100, color.RGBA{50, 50, 80, 240})

	// Draw question
	ebitenutil.DebugPrintAt(screen, currentQuestion.Question, 70, 70)

	// Draw options
	for i, option := range currentQuestion.Options {
		optionYpadding := 30 * i
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d: %s", i+1, option), 70, (140 + optionYpadding))
	}

	// Draw instructions
	ebitenutil.DebugPrintAt(screen, "Press 1-4 to answer", 70, screenHeight-100)

	// If answered, show result
	if g.triviaMgr.Answered {
		resultText := "Incorrect!"
		//resultColor := color.RGBA{255, 0, 0, 255}

		if g.triviaMgr.Correct {
			resultText = "Correct!"
			//resultColor = color.RGBA{0, 255, 0, 255}
		}

		ebitenutil.DebugPrintAt(screen, resultText, screenWidth/2-40, screenHeight/2)
	}
}

// Layout implements ebiten.Game's Layout
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Mazenasium")

	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}