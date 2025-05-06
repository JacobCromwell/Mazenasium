// internal/game/ui/input.go
package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// InputHandler manages input processing for the game
type InputHandler struct{}

// NewInputHandler creates a new input handler
func NewInputHandler() *InputHandler {
	return &InputHandler{}
}

// IsKeyJustPressed checks if a specific key was just pressed
func (i *InputHandler) IsKeyJustPressed(key ebiten.Key) bool {
	return inpututil.IsKeyJustPressed(key)
}

// IsKeyPressed checks if a specific key is being held down
func (i *InputHandler) IsKeyPressed(key ebiten.Key) bool {
	return ebiten.IsKeyPressed(key)
}

// CheckPlayerMovement checks for player movement input and returns the direction
// Returns dx, dy indicating the direction of movement
func (i *InputHandler) CheckPlayerMovement() (int, int) {
	dx, dy := 0, 0
	
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		dy = -1
	} else if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		dy = 1
	} else if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		dx = -1
	} else if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		dx = 1
	}
	
	return dx, dy
}

// CheckMazeRotation checks for maze rotation input
// Returns: -1 for left rotation, 1 for right rotation, 0 for no rotation
func (i *InputHandler) CheckMazeRotation() int {
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		return -1 // Rotate left
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		return 1 // Rotate right
	}
	return 0 // No rotation
}

// CheckActionKey checks if the action key was pressed
func (i *InputHandler) CheckActionKey() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyA)
}

// CheckSkipActionKey checks if the skip action key was pressed
func (i *InputHandler) CheckSkipActionKey() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyTab)
}

// CheckEndTurnKey checks if the end turn key was pressed
func (i *InputHandler) CheckEndTurnKey() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeySpace)
}

// CheckTriviaInput checks for trivia answer input (1-4)
// Returns: 0 for no input, 1-4 for answers
func (i *InputHandler) CheckTriviaInput() int {
	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		return 1
	}
	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		return 2
	}
	if inpututil.IsKeyJustPressed(ebiten.Key3) {
		return 3
	}
	if inpututil.IsKeyJustPressed(ebiten.Key4) {
		return 4
	}
	return 0
}

// CheckRestartKey checks if the restart key was pressed
func (i *InputHandler) CheckRestartKey() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeySpace)
}

// CheckXRotateLeftKey checks if the X-rotate left key was pressed
func (ih *InputHandler) CheckXRotateLeftKey() bool {
    return inpututil.IsKeyJustPressed(ebiten.KeyF)
}

// CheckXRotateRightKey checks if the X-rotate right key was pressed
func (ih *InputHandler) CheckXRotateRightKey() bool {
    return inpututil.IsKeyJustPressed(ebiten.KeyR)
}

// CheckConfirmKey checks if the confirm key was pressed
func (ih *InputHandler) CheckConfirmKey() bool {
    return inpututil.IsKeyJustPressed(ebiten.KeyEnter)
}

// CheckCancelKey checks if the cancel key was pressed
func (ih *InputHandler) CheckCancelKey() bool {
    return inpututil.IsKeyJustPressed(ebiten.KeyEscape)
}

// CheckActionSelectionInput checks for action selection input (1-9)
// Returns: 0 for no input, 1-9 for action selection
func (i *InputHandler) CheckActionSelectionInput() int {
    if inpututil.IsKeyJustPressed(ebiten.Key1) {
        return 1
    }
    if inpututil.IsKeyJustPressed(ebiten.Key2) {
        return 2
    }
    if inpututil.IsKeyJustPressed(ebiten.Key3) {
        return 3
    }
    if inpututil.IsKeyJustPressed(ebiten.Key4) {
        return 4
    }
    if inpututil.IsKeyJustPressed(ebiten.Key5) {
        return 5
    }
    if inpututil.IsKeyJustPressed(ebiten.Key6) {
        return 6
    }
    if inpututil.IsKeyJustPressed(ebiten.Key7) {
        return 7
    }
    if inpututil.IsKeyJustPressed(ebiten.Key8) {
        return 8
    }
    if inpututil.IsKeyJustPressed(ebiten.Key9) {
        return 9
    }
    return 0
}