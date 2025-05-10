// internal/game/animation/animation.go
package animation

import (
	"time"
)

// Animation is the interface that all animations must implement
type Animation interface {
	Update(deltaTime float64) bool // Returns true when animation is complete
	Reset()                        // Reset animation to initial state
	IsActive() bool                // Check if animation is currently active
}

// AnimationManager handles all active animations
type Manager struct {
	animations map[string]Animation
	lastTime   time.Time
}

// NewManager creates a new animation manager
func NewManager() *Manager {
	return &Manager{
		animations: make(map[string]Animation),
		lastTime:   time.Now(),
	}
}

// Register adds an animation to the manager
func (m *Manager) Register(name string, anim Animation) {
	m.animations[name] = anim
}

// Start begins an animation by name
func (m *Manager) Start(name string) bool {
	anim, exists := m.animations[name]
	if !exists {
		return false
	}
	
	anim.Reset()
	return true
}

// IsAnimating checks if any animations are currently active
func (m *Manager) IsAnimating() bool {
	for _, anim := range m.animations {
		if anim.IsActive() {
			return true
		}
	}
	return false
}

// Update updates all active animations
// Returns true if any animation is still running
func (m *Manager) Update() bool {
    now := time.Now()
    deltaTime := now.Sub(m.lastTime).Seconds()
    m.lastTime = now
    
    stillAnimating := false
    
    for _, anim := range m.animations {
        if anim.IsActive() {
            // This should be !complete to indicate animation is still running
            complete := anim.Update(deltaTime)
            if !complete {
                stillAnimating = true
            }
        }
    }
    
    return stillAnimating
}