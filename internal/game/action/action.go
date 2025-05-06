// internal/game/action/action.go
package action

import (
	"fmt"
)

// ActionType represents different types of actions
type ActionType int

const (
	XRotateLeft ActionType = iota
	XRotateRight
	// Future actions can be added here
)

// Action represents a player action
type Action struct {
	Type        ActionType
	Name        string
	Description string
	Cooldown    int // Cooldown in frames
}

// Manager handles player actions
type Manager struct {
	Actions       []Action
	Cooldowns     map[ActionType]int // Current cooldown for each action
	SelectedIndex int                // Currently selected action in the popup
}

// NewManager creates a new action manager
func NewManager() *Manager {
	// Initialize with default actions
	actions := []Action{
		{
			Type:        XRotateLeft,
			Name:        "Rotate Row Left",
			Description: "Rotate the current row to the left",
			Cooldown:    120, // 2 seconds at 60 FPS
		},
		{
			Type:        XRotateRight,
			Name:        "Rotate Row Right",
			Description: "Rotate the current row to the right",
			Cooldown:    120,
		},
	}

	cooldowns := make(map[ActionType]int)
	for _, action := range actions {
		cooldowns[action.Type] = 0
	}

	return &Manager{
		Actions:       actions,
		Cooldowns:     cooldowns,
		SelectedIndex: -1, // No action selected by default
	}
}

// UpdateCooldowns decreases all cooldowns by 1 (called each frame)
func (m *Manager) UpdateCooldowns() {
	for actionType, cooldown := range m.Cooldowns {
		if cooldown > 0 {
			m.Cooldowns[actionType]--
		}
	}
}

// IsActionAvailable checks if an action is available (not on cooldown)
func (m *Manager) IsActionAvailable(actionType ActionType) bool {
	return m.Cooldowns[actionType] == 0
}

// UseAction puts an action on cooldown
func (m *Manager) UseAction(actionType ActionType) {
	for _, action := range m.Actions {
		if action.Type == actionType {
			m.Cooldowns[actionType] = action.Cooldown
			break
		}
	}
}

// GetAvailableActions returns a list of currently available actions
func (m *Manager) GetAvailableActions() []Action {
	available := []Action{}
	for _, action := range m.Actions {
		if m.IsActionAvailable(action.Type) {
			available = append(available, action)
		}
	}
	return available
}

// GetActionByNumber returns an action by its number in the list (1-based)
// Returns nil if the number is invalid
func (m *Manager) GetActionByNumber(number int) *Action {
	availableActions := m.GetAvailableActions()
	if number < 1 || number > len(availableActions) {
		return nil
	}
	
	return &availableActions[number-1]
}

// FormatActionsList returns a formatted string of available actions
func (m *Manager) FormatActionsList() string {
	availableActions := m.GetAvailableActions()
	if len(availableActions) == 0 {
		return "No actions available"
	}

	result := "Available Actions:\n"
	for i, action := range availableActions {
		result += fmt.Sprintf("%d: %s - %s\n", i+1, action.Name, action.Description)
	}
	
	return result
}