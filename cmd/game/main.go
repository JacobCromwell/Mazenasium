package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/JacobCromwell/Mazenasium/internal/game/npc"
	"github.com/JacobCromwell/Mazenasium/internal/game/trivia"
	"github.com/JacobCromwell/Mazenasium/internal/game/maze"
)

const (
	screenWidth  = 800
	screenHeight = 600
	playerSize   = 38  // Player size (2 pixels smaller than tile)
)

// Game implements ebiten.Game interface.
type Game struct {
	gameState   GameState
	turnState   TurnState
	currentTurn TurnOwner
	player      Player
	npcManager  *npc.Manager
	maze        *maze.Maze  // Now refers to maze.Maze type from our package
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

// TurnState represents the current state within a turn
type TurnState int

const (
	WaitingForMove TurnState = iota
	WaitingForTrivia
	WaitingForAction
	WaitingForEndTurn
	ProcessingNPCTurn
)

// TurnOwner represents who currently has a turn
type TurnOwner int

const (
	PlayerTurn TurnOwner = iota
	NPCTurn
)

// Player represents the player character
type Player struct {
	gridX, gridY int
	x, y         float64 // Actual position for smooth movement
	destX, destY float64 // Destination for smooth movement
	moving       bool
	size         float64
}

// Initialize new game
func NewGame() *Game {
	mazeWidth := 10
	mazeHeight := 10

	game := &Game{
		gameState:   Playing,
		turnState:   WaitingForMove,
		currentTurn: PlayerTurn,
		player: Player{
			gridX: 1,
			gridY: 1,
			size:  playerSize,
		},
		npcManager: npc.NewManager(),
		maze:       maze.New(mazeWidth, mazeHeight, screenWidth-maze.Radius-20, screenHeight-maze.Radius-20),
		triviaMgr:  trivia.NewManager(),
		actionMsg:  "",
		actionTimer: 0,
		winner:     "",
	}

	// Set initial position for player
	game.player.x = float64(game.player.gridX) * maze.TileSize
	game.player.y = float64(game.player.gridY) * maze.TileSize
	game.player.destX = game.player.x 
	game.player.destY = game.player.y

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
	switch g.turnState {
	case WaitingForMove:
		if g.currentTurn == PlayerTurn {
			g.handlePlayerMovement()
		} else {
			g.processNPCTurn()
		}

	case WaitingForAction:
		if inpututil.IsKeyJustPressed(ebiten.KeyA) {
			g.actionMsg = "Action used!"
			g.actionTimer = 60 // Show message for 60 frames
			g.turnState = WaitingForEndTurn
		}

		// Allow skipping the action
		if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
			g.turnState = WaitingForEndTurn
		}

	case WaitingForEndTurn:
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			// End turn and switch to next actor
			if g.currentTurn == PlayerTurn {
				g.currentTurn = NPCTurn
				g.turnState = ProcessingNPCTurn
				// Reset NPC movement tracking for the new turn
				g.npcManager.ResetMovedStatus()
			} else {
				g.currentTurn = PlayerTurn
				g.turnState = WaitingForMove
			}
		}

	case ProcessingNPCTurn:
		// If no NPCs are still moving, process their next action
		if !g.npcManager.AnyMoving() {
			g.processNPCTurn()
		}
	}

	// Maze rotation (can be done anytime during player's turn)
	if g.currentTurn == PlayerTurn {
		if ebiten.IsKeyPressed(ebiten.KeyQ) {
			g.maze.RotateLeft()
		}
		if ebiten.IsKeyPressed(ebiten.KeyE) {
			g.maze.RotateRight()
		}
	}
}

// Update positions for smooth movement
func (g *Game) updatePositions() {
	// Update player position with smooth movement
	if g.player.moving {
		moveSpeed := 5.0
		dx := g.player.destX - g.player.x
		dy := g.player.destY - g.player.y

		if math.Abs(dx) < moveSpeed && math.Abs(dy) < moveSpeed {
			// Arrived at destination
			g.player.x = g.player.destX
			g.player.y = g.player.destY
			g.player.moving = false

			// Check if player reached the goal
			if g.maze.IsGoal(g.player.gridX, g.player.gridY) {
				g.winner = "Player"
				g.gameState = GameOver
				return
			}

			// If this was a player move, show trivia
			if g.currentTurn == PlayerTurn && g.turnState == WaitingForMove {
				g.gameState = AnsweringTrivia
				g.turnState = WaitingForTrivia
				g.triviaMgr.Answered = false
				g.triviaMgr.SetRandomQuestion(rand.Intn)
			}
		} else {
			// Move toward destination
			if dx != 0 {
				g.player.x += math.Copysign(moveSpeed, dx)
			}
			if dy != 0 {
				g.player.y += math.Copysign(moveSpeed, dy)
			}
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
	if g.player.moving {
		return // Already moving
	}

	newGridX, newGridY := g.player.gridX, g.player.gridY

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
		// Update grid position
		g.player.gridX = newGridX
		g.player.gridY = newGridY

		// Set destination for smooth movement
		g.player.destX = float64(newGridX) * maze.TileSize
		g.player.destY = float64(newGridY) * maze.TileSize
		g.player.moving = true
	}
}

// Process NPC turn using the NPC manager
func (g *Game) processNPCTurn() {
	// Check if all NPCs have moved
	if g.npcManager.AllMoved() {
		g.currentTurn = PlayerTurn
		g.turnState = WaitingForMove
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
			g.turnState = WaitingForAction
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
	ebitenutil.DrawRect(screen, g.player.x+1, g.player.y+1, g.player.size, g.player.size, color.RGBA{0, 0, 255, 255})

	// Draw circular maze in the corner
	g.maze.DrawCircular(screen, g.player.gridX, g.player.gridY)

	// Draw UI info
	g.drawUI(screen)

	// Draw action message if active
	if g.actionMsg != "" {
		ebitenutil.DebugPrintAt(screen, g.actionMsg, screenWidth/2-50, screenHeight/2)
	}
}

// Draw the UI
func (g *Game) drawUI(screen *ebiten.Image) {
	// Draw turn info
	turnText := "Player's Turn"
	if g.currentTurn == NPCTurn {
		turnText = "NPC's Turn"
	}
	ebitenutil.DebugPrintAt(screen, turnText, 10, 10)

	// Draw turn state info
	stateText := ""
	switch g.turnState {
	case WaitingForMove:
		stateText = "Arrow Keys: Move"
	case WaitingForAction:
		stateText = "A: Use Action, Tab: Skip"
	case WaitingForEndTurn:
		stateText = "Space: End Turn"
	case ProcessingNPCTurn:
		stateText = "NPCs are moving..."
	}
	ebitenutil.DebugPrintAt(screen, stateText, 10, 30)

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