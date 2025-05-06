// internal/game/turn/turn.go
package turn

// State represents the current state within a turn
type State int

const (
	WaitingForMove State = iota
	WaitingForTrivia
	WaitingForAction // Now optional
	WaitingForEndTurn
	ProcessingNPCTurn
	SelectingAction // New state for action selection popup
)

// Owner represents who currently has a turn
type Owner int

const (
	PlayerTurn Owner = iota
	NPCTurn
)

// Manager handles the turn-based logic of the game
type Manager struct {
	CurrentState State
	CurrentOwner Owner
}

// NewManager creates a new turn manager
func NewManager() *Manager {
	return &Manager{
		CurrentState: WaitingForMove,
		CurrentOwner: PlayerTurn,
	}
}

// NextState advances to the next state based on the current state and owner
func (m *Manager) NextState(newState State) {
	m.CurrentState = newState
}

// EndTurn ends the current turn and switches to the next actor
func (m *Manager) EndTurn() {
	if m.CurrentOwner == PlayerTurn {
		m.CurrentOwner = NPCTurn
		m.CurrentState = ProcessingNPCTurn
	} else {
		m.CurrentOwner = PlayerTurn
		m.CurrentState = WaitingForMove
	}
}

// IsPlayerTurn checks if it's currently the player's turn
func (m *Manager) IsPlayerTurn() bool {
	return m.CurrentOwner == PlayerTurn
}

// StateText returns descriptive text for the current state
func (m *Manager) StateText() string {
	switch m.CurrentState {
	case WaitingForMove:
		return "Arrow Keys: Move"
	case WaitingForAction:
		return "A: Show Actions, Space: End Turn"
	case SelectingAction:
		return "Enter 1-9 to select action"
	case WaitingForEndTurn:
		return "Space: End Turn"
	case ProcessingNPCTurn:
		return "NPCs are moving..."
	default:
		return ""
	}
}

// OwnerText returns descriptive text for the current turn owner
func (m *Manager) OwnerText() string {
	if m.CurrentOwner == PlayerTurn {
		return "Player's Turn"
	}
	return "NPC's Turn"
}