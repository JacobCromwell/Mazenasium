// internal/game/state/state.go
package state

import (
	"fmt"
	"image/color"
	//"math/rand" // skipping trivia for now

	"github.com/JacobCromwell/Mazenasium/internal/game/action"
	"github.com/JacobCromwell/Mazenasium/internal/game/flavor"
	"github.com/JacobCromwell/Mazenasium/internal/game/maze"
	"github.com/JacobCromwell/Mazenasium/internal/game/menu"
	"github.com/JacobCromwell/Mazenasium/internal/game/npc"
	"github.com/JacobCromwell/Mazenasium/internal/game/player"
	"github.com/JacobCromwell/Mazenasium/internal/game/trivia"
	"github.com/JacobCromwell/Mazenasium/internal/game/turn"
	"github.com/JacobCromwell/Mazenasium/internal/game/ui"
)

// GameState represents the current state of the game
type GameState int

const (
	Menu GameState = iota
	Playing
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
	ActionMgr    *action.Manager
	MenuMgr      *menu.Manager
	UIRenderer   *ui.Renderer
	InputHandler *ui.InputHandler
	Flavor       *flavor.Manager
	Winner       string

	// fields for xRotateAction
	xRotateActive    bool // Whether X-rotate mode is active
	xRotateDirection int  // 1 for right, -1 for left
}

// In internal/game/state/state.go
// Update the New function to ensure proper initialization of the Flavor manager

func New(screenWidth, screenHeight int) *Manager {
    // Increased base size for the maze - will be doubled in maze.New
    mazeWidth := 10
    mazeHeight := 10

    // Create and initialize the flavor manager first
    flavorMgr := flavor.NewManager()
    
    manager := &Manager{
        CurrentState:     Menu, // Start with Menu state
        TurnManager:      turn.NewManager(),
        Player:           player.New(1, 1, maze.TileSize),
        NPCManager:       npc.NewManager(),
        Maze:             maze.New(mazeWidth, mazeHeight, 0, 0),
        TriviaMgr:        trivia.NewManager(),
        ActionMgr:        action.NewManager(),
        MenuMgr:          menu.NewManager(), // Initialize menu manager
        UIRenderer:       ui.NewRenderer(),
        InputHandler:     ui.NewInputHandler(),
        Flavor:           flavorMgr, // Make sure this is set
        Winner:           "",
        xRotateActive:    false,
        xRotateDirection: 0,
    }

    // Create NPCs
    npc1 := npc.New(0, 3, 3, maze.TileSize, color.RGBA{255, 0, 0, 255})
    npc2 := npc.New(1, 5, 5, maze.TileSize, color.RGBA{0, 255, 0, 255})

    // Add NPCs to manager
    manager.NPCManager.AddNPC(npc1)
    manager.NPCManager.AddNPC(npc2)

    // Try to load flavor images after initializing the manager
    if flavorMgr != nil {
        // Use a try/catch pattern to prevent crash if image loading fails
        func() {
            defer func() {
                if r := recover(); r != nil {
                    fmt.Println("Warning: Error while loading flavor images:", r)
                }
            }()
            
            err := flavorMgr.LoadImages("assets")
            if err != nil {
                fmt.Println("Warning: Failed to load flavor images:", err)
            }
        }()
    }

    return manager
}

// Update the Update method to handle menu state
func (m *Manager) Update() {
	switch m.CurrentState {
	case Menu:
		m.updateMenu()
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

	// Update action cooldowns
	m.ActionMgr.UpdateCooldowns()
}

// Add the updateMenu method
func (m *Manager) updateMenu() {
	action := m.MenuMgr.HandleInput()

	if action == "start_game" {
		// Start the game
		m.CurrentState = Playing
	} else if action == "quit" {
		// Quit the game
		// In a real implementation, you'd handle this differently
		// For now, we'll just switch to game over state
		m.Winner = "Quit"
		m.CurrentState = GameOver
	}
}

// Update while playing
func (m *Manager) updatePlaying() {
	// Update positions for smooth movement
	m.updatePositions()

	// If X-rotate is active, handle confirmation or cancellation
	if m.xRotateActive {
		m.handleXRotateConfirmation()
		return
	}

	// Process based on turn state
	switch m.TurnManager.CurrentState {
	case turn.WaitingForMove:
		if m.TurnManager.IsPlayerTurn() {
			m.handlePlayerMovement()
		} else {
			m.processNPCTurn()
		}

	case turn.WaitingForAction:
		// Player can now either show the action menu or end their turn directly
		if m.InputHandler.CheckActionKey() {
			// Show action menu
			m.TurnManager.NextState(turn.SelectingAction)
		} else if m.InputHandler.CheckEndTurnKey() {
			// Skip action and end turn
			m.TurnManager.EndTurn()
			// Reset NPC movement tracking for the new turn if switching to NPC turn
			if m.TurnManager.CurrentOwner == turn.NPCTurn {
				m.NPCManager.ResetMovedStatus()
			}
		}

	case turn.SelectingAction:
		// Handle action selection or cancellation
		if m.InputHandler.CheckCancelKey() {
			// Return to the WaitingForAction state
			m.TurnManager.NextState(turn.WaitingForAction)
			m.UIRenderer.SetActionMessage("Action selection cancelled", 60)
		} else {
			// Check for action number input
			actionNum := m.InputHandler.CheckActionSelectionInput()
			if actionNum > 0 {
				// Get the selected action
				selectedAction := m.ActionMgr.GetActionByNumber(actionNum)
				if selectedAction != nil {
					// Process the selected action
					m.handleActionSelection(*selectedAction)
				}
			}
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

// Add this method to the Manager struct to collect entity positions
func (m *Manager) collectEntityPositions() []maze.Position {
    positions := []maze.Position{}
    
    // Add player position
    playerGridX, playerGridY := m.Player.GetGridPosition()
    positions = append(positions, maze.Position{X: playerGridX, Y: playerGridY})
    
    // Add NPC positions
    for _, npc := range m.NPCManager.NPCs {
        positions = append(positions, maze.Position{X: npc.GridX, Y: npc.GridY})
    }
    
    return positions
}

// Modify the handleXRotateConfirmation method to check for collisions
func (m *Manager) handleXRotateConfirmation() {
    // Check for confirmation
    if m.InputHandler.CheckConfirmKey() {
        playerGridX, playerGridY := m.Player.GetGridPosition()
        
        // Collect all entity positions
        entityPositions := m.collectEntityPositions()
        
        // Check for collisions
        hasCollision := m.Maze.CheckXRotateCollisions(
            playerGridX,
            playerGridY,
            m.xRotateDirection,
            entityPositions,
        )
        
        if hasCollision {
            // Cancel the action due to collision
            m.Maze.ClearHighlights()
            m.xRotateActive = false
            m.UIRenderer.SetActionMessage("Cannot move wall segments on top of players or NPCs", 120)
            m.TurnManager.NextState(turn.WaitingForAction)
            return
        }

		// No collision, perform the rotation
		m.Maze.PerformXRotate(playerGridX, playerGridY, m.xRotateDirection)

		// Mark the action as used
		if m.xRotateDirection > 0 {
			m.ActionMgr.UseAction(action.XRotateRight)
			m.UIRenderer.SetActionMessage("X-Rotate Right Used!", 60)
		} else {
			m.ActionMgr.UseAction(action.XRotateLeft)
			m.UIRenderer.SetActionMessage("X-Rotate Left Used!", 60)
		}

		// Clear state and move to end turn
		m.xRotateActive = false
		m.TurnManager.NextState(turn.WaitingForEndTurn)
	}

	// Check for cancellation
	if m.InputHandler.CheckCancelKey() {
		// Clear highlights and exit X-rotate mode
		m.Maze.ClearHighlights()
		m.xRotateActive = false
		m.UIRenderer.SetActionMessage("X-Rotate Cancelled", 60)
		m.TurnManager.NextState(turn.WaitingForAction)
	}
}

// Handle the selected action
func (m *Manager) handleActionSelection(selectedAction action.Action) {
	switch selectedAction.Type {
	case action.XRotateLeft:
		playerGridX, playerGridY := m.Player.GetGridPosition()
		m.Maze.HighlightXRotation(playerGridX, playerGridY)
		m.xRotateActive = true
		m.xRotateDirection = -1
		m.UIRenderer.SetActionMessage("X-Rotate Left? (Confirm: Enter, Cancel: Esc)", 0) // 0 for no timeout

	case action.XRotateRight:
		playerGridX, playerGridY := m.Player.GetGridPosition()
		m.Maze.HighlightXRotation(playerGridX, playerGridY)
		m.xRotateActive = true
		m.xRotateDirection = 1
		m.UIRenderer.SetActionMessage("X-Rotate Right? (Confirm: Enter, Cancel: Esc)", 0)

	// Add more cases for future actions

	default:
		m.UIRenderer.SetActionMessage("Unknown action selected", 60)
		m.TurnManager.NextState(turn.WaitingForAction)
	}
}

// Update positions for smooth movement
func (m *Manager) updatePositions() {
	// Update player position with smooth movement
	playerGridX, playerGridY := m.Player.GetGridPosition()

	// Update player, and check if they've arrived at destination
	if arrived := m.Player.Update(5.0); arrived {

		if m.Flavor != nil {
			playerGridX, playerGridY := m.Player.GetGridPosition()
			tile := m.Maze.State.GetTile(playerGridX, playerGridY)
			
			if tile != nil && tile.Type != maze.Wall {
				// Get the flavor image path from the tile
				imagePath := tile.GetFlavorImage()
				if imagePath != "" {
					// Update the flavor image based on the tile's assigned image
					m.Flavor.SetImageByPath(imagePath)
				}
			}
		}
        

		// Check if player reached the goal
		if m.Maze.IsGoal(playerGridX, playerGridY) {
			m.Winner = "Player"
			m.CurrentState = GameOver
			return
		}

		// DEBUGGING: Skip trivia and go directly to action phase
		if m.TurnManager.IsPlayerTurn() && m.TurnManager.CurrentState == turn.WaitingForMove {
			// Comment out the trivia part
			// m.CurrentState = AnsweringTrivia
			// m.TurnManager.NextState(turn.WaitingForTrivia)
			// m.TriviaMgr.Answered = false
			// m.TriviaMgr.SetRandomQuestion(rand.Intn)

			// Instead, go directly to waiting for action
			m.TurnManager.NextState(turn.WaitingForAction)
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

	newGridX, newGridY := playerGridX+dx, playerGridY+dy

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
