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

	"github.com/JacobCromwell/mazenasium/cmd/game/model/varbar.go"
)

const (
	screenWidth  = 800
	screenHeight = 600
	mazeRadius   = 120 // Radius of the circular maze
	tileSize     = 40  // Size of each tile in the maze
	playerSize   = 38  // Player size (2 pixels smaller than tile)
)

// Game implements ebiten.Game interface.
type Game struct {
	gameState   GameState
	turnState   TurnState
	currentTurn TurnOwner
	player      Player
	npcs        []NPC
	maze        Maze
	triviaMgr   TriviaManager
	actionMsg   string
	actionTimer int
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

// NPC represents a non-player character
type NPC struct {
	id           int
	gridX, gridY int
	x, y         float64 // Actual position for smooth movement
	destX, destY float64 // Destination for smooth movement
	moving       bool
	size         float64
	color        color.RGBA
}

// Maze represents the maze grid
type Maze struct {
	grid          [][]MazeTile
	width, height int
	centerX       float64
	centerY       float64
	rotationAngle float64
}

// MazeTile represents a single tile in the maze
type MazeTile struct {
	isWall   bool
	visited  bool
	hasItem  bool
	itemType int
}

// TriviaManager handles trivia questions and answers
type TriviaManager struct {
	questions    []TriviaQuestion
	currentIndex int
	answered     bool
	correct      bool
}

// TriviaQuestion represents a single trivia question
type TriviaQuestion struct {
	question string
	options  []string
	answer   int
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
		npcs: []NPC{
			{
				id:    0,
				gridX: 3,
				gridY: 3,
				size:  playerSize,
				color: color.RGBA{255, 0, 0, 255},
			},
			{
				id:    1,
				gridX: 5,
				gridY: 5,
				size:  playerSize,
				color: color.RGBA{0, 255, 0, 255},
			},
		},
		maze: Maze{
			width:         mazeWidth,
			height:        mazeHeight,
			centerX:       screenWidth - mazeRadius - 20,
			centerY:       screenHeight - mazeRadius - 20,
			rotationAngle: 0,
		},
		triviaMgr: TriviaManager{
			questions:    loadTriviaQuestions(),
			currentIndex: 0,
			answered:     false,
		},
		actionMsg:   "",
		actionTimer: 0,
	}

	// Generate maze grid
	game.maze.grid = createMazeGrid(mazeWidth, mazeHeight)

	// Set initial positions for player and NPCs
	game.player.x = float64(game.player.gridX) * tileSize
	game.player.y = float64(game.player.gridY) * tileSize
	game.player.destX = game.player.x
	game.player.destY = game.player.y

	for i := range game.npcs {
		game.npcs[i].x = float64(game.npcs[i].gridX) * tileSize
		game.npcs[i].y = float64(game.npcs[i].gridY) * tileSize
		game.npcs[i].destX = game.npcs[i].x
		game.npcs[i].destY = game.npcs[i].y
	}

	return game
}

// Create a simple maze grid
func createMazeGrid(width, height int) [][]MazeTile {
	grid := make([][]MazeTile, height)
	for y := range grid {
		grid[y] = make([]MazeTile, width)
		for x := range grid[y] {
			// Create walls around the edges and some random walls
			if x == 0 || y == 0 || x == width-1 || y == height-1 || (rand.Intn(100) < 20 && x > 1 && y > 1) {
				grid[y][x].isWall = true
			}
		}
	}

	// Ensure the starting positions are not walls
	grid[1][1].isWall = false // Player start
	grid[3][3].isWall = false // NPC1 start
	grid[5][5].isWall = false // NPC2 start

	return grid
}

// Load trivia questions
func loadTriviaQuestions() []TriviaQuestion {
	// In a real implementation, you'd load these from a file
	return []TriviaQuestion{
		{
			question: "What is the capital of France?",
			options:  []string{"London", "Berlin", "Paris", "Madrid"},
			answer:   2, // Paris (0-indexed)
		},
		{
			question: "Which planet is known as the Red Planet?",
			options:  []string{"Venus", "Mars", "Jupiter", "Saturn"},
			answer:   1, // Mars
		},
		{
			question: "What is the largest mammal?",
			options:  []string{"Elephant", "Giraffe", "Blue Whale", "Hippopotamus"},
			answer:   2, // Blue Whale
		},
		{
			question: "What element has the chemical symbol 'O'?",
			options:  []string{"Gold", "Oxygen", "Osmium", "Oganesson"},
			answer:   1, // Oxygen
		},
		{
			question: "Who wrote 'Romeo and Juliet'?",
			options:  []string{"Charles Dickens", "William Shakespeare", "Jane Austen", "Mark Twain"},
			answer:   1, // Shakespeare
		},
	}
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
			} else {
				g.currentTurn = PlayerTurn
				g.turnState = WaitingForMove
			}
		}

	case ProcessingNPCTurn:
		// If no NPCs are still moving, process their next action
		if !g.anyNPCMoving() {
			g.processNPCTurn()
		}
	}

	// Maze rotation (can be done anytime during player's turn)
	if g.currentTurn == PlayerTurn {
		if ebiten.IsKeyPressed(ebiten.KeyQ) {
			g.maze.rotationAngle -= 0.05
		}
		if ebiten.IsKeyPressed(ebiten.KeyE) {
			g.maze.rotationAngle += 0.05
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

			// If this was a player move, show trivia
			if g.currentTurn == PlayerTurn && g.turnState == WaitingForMove {
				g.gameState = AnsweringTrivia
				g.turnState = WaitingForTrivia
				g.triviaMgr.answered = false
				g.triviaMgr.currentIndex = rand.Intn(len(g.triviaMgr.questions))
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

	// Update NPCs position with smooth movement
	for i := range g.npcs {
		if g.npcs[i].moving {
			moveSpeed := 5.0
			dx := g.npcs[i].destX - g.npcs[i].x
			dy := g.npcs[i].destY - g.npcs[i].y

			if math.Abs(dx) < moveSpeed && math.Abs(dy) < moveSpeed {
				// Arrived at destination
				g.npcs[i].x = g.npcs[i].destX
				g.npcs[i].y = g.npcs[i].destY
				g.npcs[i].moving = false
			} else {
				// Move toward destination
				if dx != 0 {
					g.npcs[i].x += math.Copysign(moveSpeed, dx)
				}
				if dy != 0 {
					g.npcs[i].y += math.Copysign(moveSpeed, dy)
				}
			}
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
	if newGridX >= 0 && newGridX < g.maze.width &&
		newGridY >= 0 && newGridY < g.maze.height &&
		!g.maze.grid[newGridY][newGridX].isWall {
		// Update grid position
		g.player.gridX = newGridX
		g.player.gridY = newGridY

		// Set destination for smooth movement
		g.player.destX = float64(newGridX) * tileSize
		g.player.destY = float64(newGridY) * tileSize
		g.player.moving = true
	}
}

// Process NPC turn
func (g *Game) processNPCTurn() {
	if g.anyNPCMoving() {
		return // Wait for movement to complete
	}

	// Process all NPCs
	for i := range g.npcs {
		if g.npcs[i].moving {
			continue // Skip if already moving
		}

		// Simple NPC movement - move in a random valid direction
		directions := []struct{ dx, dy int }{
			{-1, 0}, {1, 0}, {0, -1}, {0, 1}, // Left, right, up, down
		}

		// Shuffle directions
		rand.Shuffle(len(directions), func(i, j int) {
			directions[i], directions[j] = directions[j], directions[i]
		})

		moved := false
		for _, dir := range directions {
			newGridX := g.npcs[i].gridX + dir.dx
			newGridY := g.npcs[i].gridY + dir.dy

			// Check if movement is valid
			if newGridX >= 0 && newGridX < g.maze.width &&
				newGridY >= 0 && newGridY < g.maze.height &&
				!g.maze.grid[newGridY][newGridX].isWall {
				// Update grid position
				g.npcs[i].gridX = newGridX
				g.npcs[i].gridY = newGridY

				// Set destination for smooth movement
				g.npcs[i].destX = float64(newGridX) * tileSize
				g.npcs[i].destY = float64(newGridY) * tileSize
				g.npcs[i].moving = true
				moved = true
				break
			}
		}

		if !moved {
			// If NPC can't move, just skip its turn
			continue
		}
	}

	// After processing all NPCs, return to player's turn
	if !g.anyNPCMoving() {
		g.currentTurn = PlayerTurn
		g.turnState = WaitingForMove
	}
}

// Check if any NPC is currently moving
func (g *Game) anyNPCMoving() bool {
	for i := range g.npcs {
		if g.npcs[i].moving {
			return true
		}
	}
	return false
}

// Update trivia state
func (g *Game) updateTrivia() {
	question := g.triviaMgr.questions[g.triviaMgr.currentIndex]

	// Check for answer selection
	for i := 0; i < len(question.options); i++ {
		if inpututil.IsKeyJustPressed(ebiten.Key1 + ebiten.Key(i)) {
			g.triviaMgr.answered = true
			g.triviaMgr.correct = (i == question.answer)

			// Return to game after brief delay
			go func() {
				// Note: In a real game, you'd want to use a more robust timer or state system
				// This is just a simple demonstration
				g.gameState = Playing
				g.turnState = WaitingForAction
			}()
		}
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
		ebitenutil.DebugPrint(screen, "Game Over! Press SPACE to restart")
	}
}

// Draw the playing state
func (g *Game) drawPlaying(screen *ebiten.Image) {
	// Draw the maze grid
	g.drawMaze(screen)

	// Draw NPCs
	for _, npc := range g.npcs {
		ebitenutil.DrawRect(screen, npc.x+1, npc.y+1, npc.size, npc.size, npc.color)
	}

	// Draw player
	ebitenutil.DrawRect(screen, g.player.x+1, g.player.y+1, g.player.size, g.player.size, color.RGBA{0, 0, 255, 255})

	// Draw circular maze in the corner
	g.drawCircularMaze(screen)

	// Draw UI info
	g.drawUI(screen)

	// Draw action message if active
	if g.actionMsg != "" {
		ebitenutil.DebugPrintAt(screen, g.actionMsg, screenWidth/2-50, screenHeight/2)
	}
}

// Draw the maze grid
func (g *Game) drawMaze(screen *ebiten.Image) {
	// Draw grid lines and tiles
	for y := 0; y < g.maze.height; y++ {
		for x := 0; x < g.maze.width; x++ {
			// Calculate tile position
			tileX := float64(x) * tileSize
			tileY := float64(y) * tileSize

			// Draw tile border
			borderColor := color.RGBA{100, 100, 100, 255}
			ebitenutil.DrawLine(screen, tileX, tileY, tileX+tileSize, tileY, borderColor)
			ebitenutil.DrawLine(screen, tileX, tileY, tileX, tileY+tileSize, borderColor)
			ebitenutil.DrawLine(screen, tileX+tileSize, tileY, tileX+tileSize, tileY+tileSize, borderColor)
			ebitenutil.DrawLine(screen, tileX, tileY+tileSize, tileX+tileSize, tileY+tileSize, borderColor)

			// Draw wall or floor
			if g.maze.grid[y][x].isWall {
				ebitenutil.DrawRect(screen, tileX, tileY, tileSize, tileSize, color.RGBA{70, 70, 70, 255})
			} else {
				ebitenutil.DrawRect(screen, tileX, tileY, tileSize, tileSize, color.RGBA{200, 200, 200, 100})
			}
		}
	}
}

// Draw the circular maze in the corner
func (g *Game) drawCircularMaze(screen *ebiten.Image) {
	// Draw outer circle
	ebitenutil.DrawCircle(screen, g.maze.centerX, g.maze.centerY, mazeRadius, color.RGBA{200, 200, 200, 100})

	// Draw a simplified representation of the maze in the circle
	// This is just a placeholder - in a real game, you'd want to create a proper radial maze
	cellAngle := 2 * math.Pi / float64(g.maze.width)
	cellRadius := mazeRadius / float64(g.maze.height)

	for y := 0; y < g.maze.height; y++ {
		radius := float64(y+1) * cellRadius

		for x := 0; x < g.maze.width; x++ {
			angle := g.maze.rotationAngle + float64(x)*cellAngle

			// Calculate position
			cellX := g.maze.centerX + math.Cos(angle)*radius
			cellY := g.maze.centerY + math.Sin(angle)*radius

			// Draw cell
			if g.maze.grid[y][x].isWall {
				ebitenutil.DrawCircle(screen, cellX, cellY, cellRadius/2, color.RGBA{70, 70, 70, 255})
			}
		}
	}

	// Draw player position in the minimap
	playerAngle := g.maze.rotationAngle + float64(g.player.gridX)*cellAngle
	playerRadius := float64(g.player.gridY+1) * cellRadius
	playerMiniX := g.maze.centerX + math.Cos(playerAngle)*playerRadius
	playerMiniY := g.maze.centerY + math.Sin(playerAngle)*playerRadius
	ebitenutil.DrawCircle(screen, playerMiniX, playerMiniY, cellRadius/2, color.RGBA{0, 0, 255, 255})

	// Draw rotation controls
	ebitenutil.DebugPrintAt(screen, "Q/E: Rotate", int(g.maze.centerX)-40, int(g.maze.centerY)+mazeRadius+10)
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
}

// Draw the trivia screen
func (g *Game) drawTrivia(screen *ebiten.Image) {
	question := g.triviaMgr.questions[g.triviaMgr.currentIndex]

	// Draw question background
	ebitenutil.DrawRect(screen, 50, 50, screenWidth-100, screenHeight-100, color.RGBA{50, 50, 80, 240})

	// Draw question
	ebitenutil.DebugPrintAt(screen, question.question, 70, 70)

	// Draw options
	for i, option := range question.options {
		optionYpadding := 10 * i
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d: %s", i+1, option), 70, (140 + optionYpadding))
	}

	// Draw instructions
	ebitenutil.DebugPrintAt(screen, "Press 1-4 to answer", 70, screenHeight-100)

	// If answered, show result
	if g.triviaMgr.answered {
		resultText := "Incorrect!"
		//resultColor := color.RGBA{255, 0, 0, 255}

		if g.triviaMgr.correct {
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
