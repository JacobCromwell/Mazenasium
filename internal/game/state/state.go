package state

import (
	"fmt"
	"math/rand"
	"image/color"

	"github.com/JacobCromwell/Mazenasium/internal/game/maze"
	"github.com/JacobCromwell/Mazenasium/internal/game/npc"
	"github.com/JacobCromwell/Mazenasium/internal/game/player"
	"github.com/JacobCromwell/Mazenasium/internal/game/trivia"
	"github.com/JacobCromwell/Mazenasium/internal/game/turn"
	"github.com/JacobCromwell/Mazenasium/internal/game/ui"
)

// GameState represents the current state of the game
type GameState int

const (
	Playing GameState = iota
	AnsweringTrivia
	GameOver
)

// Manager handles all game state logic
type Manager struct {
	CurrentState GameState
	TurnManager  *turn.Manager
	Player       *player.Player
	NPCManager   *npc.Manager
	Maze         *maze.Maze
	TriviaMgr    *trivia.Manager
	UIRenderer   *ui.Renderer
	InputHandler *ui.InputHandler
	Winner       string
}

// New creates a new game state manager
func New(screenWidth, screenHeight int) *Manager {
	mazeWidth := 10
	mazeHeight := 10

	manager := &Manager{
		CurrentState: Playing,
		TurnManager:  turn.NewManager(),
		Player:       player.New(1, 1, maze.TileSize),
		NPCManager:   npc.NewManager(),
		Maze:         maze.New(mazeWidth, mazeHeight, screenWidth-maze.Radius-20, screenHeight-maze.Radius-20),
		TriviaMgr:    trivia.NewManager(),
		UIRenderer:   ui.NewRenderer(),
		InputHandler: ui.NewInputHandler(),
		Winner:       "",
	}

	// Create NPCs
	npc1 := npc.New(0, 3, 3, maze.TileSize, color.RGBA{255, 0, 0, 255})
	npc2 := npc.New(1, 5, 5, maze.TileSize, color.RGBA{0, 255, 0, 255})
	
	// Add NPCs to manager
	manager.NPCManager.AddNPC(npc1)
	manager.NPCManager.AddNPC(npc2)

	return manager
}

// Update updates the game state
func (m *Manager) Update() {
	switch m.CurrentState {
	case Playing:
		m.updatePlaying()
	case AnsweringTrivia:
		m.updateTrivia()
	case GameOver:
		if m.InputHandler.CheckRestartKey() {
			// Reset game
			*m = *New(ui.ScreenWidth, ui.ScreenHeight)
		}
	}

	// Update action message timer in the UI renderer
	m.UIRenderer.UpdateActionTimer()
}

// Update while playing
func (m *Manager) updatePlaying() {
	// Update positions for smooth movement
	m.updatePositions()

	// Process based on turn state
	switch m.TurnManager.CurrentState {
	case turn.WaitingForMove:
		if m.TurnManager.IsPlayerTurn() {
			m.handlePlayerMovement()
		} else {
			m.processNPCTurn()
		}

	case turn.WaitingForAction:
		if m.InputHandler.CheckActionKey() {
			m.UIRenderer.SetActionMessage("Action used!", 60)
			m.TurnManager.NextState(turn.WaitingForEndTurn)
		}

		// Allow skipping the action
		if m.InputHandler.CheckSkipActionKey() {
			m.TurnManager.NextState(turn.WaitingForEndTurn)
		}

	case turn.WaitingForEndTurn:
		if m.InputHandler.CheckEndTurnKey() {
			// End turn and switch to next actor
			m.TurnManager.EndTurn()
			// Reset NPC movement tracking for the new turn if switching to NPC turn
			if m.TurnManager.CurrentOwner == turn.NPCTurn {
				m.NPCManager.ResetMovedStatus()
			}
		}

	case turn.ProcessingNPCTurn:
		// If no NPCs are still moving, process their next action
		if !m.NPCManager.AnyMoving() {
			m.processNPCTurn()
		}
	}
}

// Update positions for smooth movement
func (m *Manager) updatePositions() {
	// Update player position with smooth movement
	playerGridX, playerGridY := m.Player.GetGridPosition()
	
	// Update player, and check if they've arrived at destination
	if arrived := m.Player.Update(5.0); arrived {
		// Check if player reached the goal
		if m.Maze.IsGoal(playerGridX, playerGridY) {
			m.Winner = "Player"
			m.CurrentState = GameOver
			return
		}

		// If this was a player move, show trivia
		if m.TurnManager.IsPlayerTurn() && m.TurnManager.CurrentState == turn.WaitingForMove {
			m.CurrentState = AnsweringTrivia
			m.TurnManager.NextState(turn.WaitingForTrivia)
			m.TriviaMgr.Answered = false
			m.TriviaMgr.SetRandomQuestion(rand.Intn)
		}
	}

	// Update NPCs positions using the manager
	arrivedNPCs := m.NPCManager.UpdatePositions(5.0)
	
	// Check if any NPCs reached the goal
	for _, arrivedNPC := range arrivedNPCs {
		if m.Maze.IsGoal(arrivedNPC.GridX, arrivedNPC.GridY) {
			m.Winner = fmt.Sprintf("NPC %d", arrivedNPC.ID+1)
			m.CurrentState = GameOver
			return
		}
	}
}

// Handle player movement
func (m *Manager) handlePlayerMovement() {
	if m.Player.IsMoving() {
		return // Already moving
	}

	playerGridX, playerGridY := m.Player.GetGridPosition()
	dx, dy := m.InputHandler.CheckPlayerMovement()
	
	if dx == 0 && dy == 0 {
		return // No movement input
	}
	
	newGridX, newGridY := playerGridX + dx, playerGridY + dy

	// Check if movement is valid (not a wall and within bounds)
	if m.Maze.IsValidMove(newGridX, newGridY) {
		// Set destination for smooth movement
		m.Player.SetDestination(newGridX, newGridY, maze.TileSize)
	}
}

// Process NPC turn using the NPC manager
func (m *Manager) processNPCTurn() {
	// Check if all NPCs have moved
	if m.NPCManager.AllMoved() {
		m.TurnManager.EndTurn() // Switch back to player's turn
		return
	}

	// Process one NPC's turn using a callback to check valid moves
	validMoveFn := func(x, y int) bool {
		return m.Maze.IsValidMove(x, y)
	}
	
	m.NPCManager.ProcessTurn(validMoveFn)
}

// Update trivia state
func (m *Manager) updateTrivia() {
	// Get input from the input handler
	answer := m.InputHandler.CheckTriviaInput()
	
	if answer > 0 {
		// Process the answer
		correct := m.TriviaMgr.CheckAnswer(answer - 1) // Convert from 1-based to 0-based
		m.TriviaMgr.Answered = true
		m.TriviaMgr.Correct = correct
		
		// Return to game after brief delay
		m.CurrentState = Playing
		m.TurnManager.NextState(turn.WaitingForAction)
	}
}