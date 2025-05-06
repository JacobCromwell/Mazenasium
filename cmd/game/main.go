package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/JacobCromwell/Mazenasium/internal/game/state"
	"github.com/JacobCromwell/Mazenasium/internal/game/ui"
)

// Game implements ebiten.Game interface.
type Game struct {
	stateManager *state.Manager
}

// Initialize new game
func NewGame() *Game {
	return &Game{
		stateManager: state.New(ui.ScreenWidth, ui.ScreenHeight),
	}
}

// Update game state
func (g *Game) Update() error {
	g.stateManager.Update()
	return nil
}

// Draw the game - delegates to UI renderer
func (g *Game) Draw(screen *ebiten.Image) {
	g.stateManager.UIRenderer.Draw(
		screen,
		int(g.stateManager.CurrentState),
		g.stateManager.Maze,
		g.stateManager.Player,
		g.stateManager.NPCManager,
		g.stateManager.TurnManager,
		g.stateManager.TriviaMgr,
		g.stateManager.ActionMgr, // Add ActionMgr parameter
		g.stateManager.Winner,
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