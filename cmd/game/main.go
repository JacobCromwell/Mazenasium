package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/JacobCromwell/Mazenasium/internal/game/maze"
	"github.com/JacobCromwell/Mazenasium/internal/game/npc"
	"github.com/JacobCromwell/Mazenasium/internal/game/player"
	"github.com/JacobCromwell/Mazenasium/internal/game/trivia"
	"github.com/JacobCromwell/Mazenasium/internal/game/turn"
	"github.com/JacobCromwell/Mazenasium/internal/game/ui"
)

// Game implements ebiten.Game interface.
type Game struct {
	gameState    GameState
	turnManager  *turn.Manager
	player       *player.Player
	npcManager   *npc.Manager
	maze         *maze.Maze
	triviaMgr    *trivia.Manager
	uiRenderer   *ui.Renderer   // UI renderer
	inputHandler *ui.InputHandler // Input handler
	winner       string       // Track who wins the game
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
		gameState:    Playing,
		turnManager:  turn.NewManager(),
		player:       player.New(1, 1, maze.TileSize),
		npcManager:   npc.NewManager(),
		maze:         maze.New(mazeWidth, mazeHeight, ui.ScreenWidth-maze.Radius-20, ui.ScreenHeight-maze.Radius-20),
		triviaMgr:    trivia.NewManager(),
		uiRenderer:   ui.NewRenderer(),     // Initialize UI renderer
		inputHandler: ui.NewInputHandler(), // Initialize input handler
		winner:       "",
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
		if g.inputHandler.CheckRestartKey() {
			// Reset game
			*g = *NewGame()
		}
	}

	// Update action message timer in the UI renderer
	g.uiRenderer.UpdateActionTimer()

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
		if g.inputHandler.CheckActionKey() {
			g.uiRenderer.SetActionMessage("Action used!", 60) // Use UI renderer to set message
			g.turnManager.NextState(turn.WaitingForEndTurn)
		}

		// Allow skipping the action
		if g.inputHandler.CheckSkipActionKey() {
			g.turnManager.NextState(turn.WaitingForEndTurn)
		}

	case turn.WaitingForEndTurn:
		if g.inputHandler.CheckEndTurnKey() {
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
	dx, dy := g.inputHandler.CheckPlayerMovement()
	
	if dx == 0 && dy == 0 {
		return // No movement input
	}
	
	newGridX, newGridY := playerGridX + dx, playerGridY + dy

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
	// Get input from the input handler
	answer := g.inputHandler.CheckTriviaInput()
	
	if answer > 0 {
		// Process the answer
		correct := g.triviaMgr.CheckAnswer(answer - 1) // Convert from 1-based to 0-based
		g.triviaMgr.Answered = true
		g.triviaMgr.Correct = correct
		
		// Return to game after brief delay
		go func() {
			// Note: In a real game, you'd want to use a more robust timer or state system
			// This is just a simple demonstration
			g.gameState = Playing
			g.turnManager.NextState(turn.WaitingForAction)
		}()
	}
}

// Draw the game - delegates to UI renderer
func (g *Game) Draw(screen *ebiten.Image) {
	g.uiRenderer.Draw(
		screen,
		int(g.gameState),
		g.maze,
		g.player,
		g.npcManager,
		g.turnManager,
		g.triviaMgr,
		g.winner,
	)
}

// Layout implements ebiten.Game's Layout
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ui.ScreenWidth, ui.ScreenHeight
}

func main() {
	ebiten.SetWindowSize(ui.ScreenWidth, ui.ScreenHeight)
	ebiten.SetWindowTitle("Mazenasium")

	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}