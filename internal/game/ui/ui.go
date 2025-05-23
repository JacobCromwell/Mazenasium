package ui

import (
	"fmt"
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"github.com/JacobCromwell/Mazenasium/internal/game/action"
	"github.com/JacobCromwell/Mazenasium/internal/game/maze"
	"github.com/JacobCromwell/Mazenasium/internal/game/npc"
	"github.com/JacobCromwell/Mazenasium/internal/game/player"
	"github.com/JacobCromwell/Mazenasium/internal/game/trivia"
	"github.com/JacobCromwell/Mazenasium/internal/game/turn"
	"github.com/JacobCromwell/Mazenasium/internal/game/flavor"
	"github.com/JacobCromwell/Mazenasium/internal/game/menu"
)

const (
	ScreenWidth  = 1200
	ScreenHeight = 1000
)

// Renderer handles all UI rendering for the game
type Renderer struct {
	actionMsg   string
	actionTimer int
}

// NewRenderer creates a new UI renderer
func NewRenderer() *Renderer {
	return &Renderer{
		actionMsg:   "",
		actionTimer: 0,
	}
}

// SetActionMessage sets a temporary action message to display
func (r *Renderer) SetActionMessage(msg string, duration int) {
	r.actionMsg = msg
	r.actionTimer = duration
}

// UpdateActionTimer updates the action message timer
func (r *Renderer) UpdateActionTimer() {
	if r.actionTimer > 0 {
		r.actionTimer--
		if r.actionTimer == 0 {
			r.actionMsg = ""
		}
	}
}

// Add this method to the Renderer struct
func (r *Renderer) drawMenu(screen *ebiten.Image, menuManager *menu.Manager) {
    if menuManager == nil || menuManager.CurrentMenu == nil {
        return
    }
    
    currentMenu := menuManager.CurrentMenu
    
    // Draw menu background
    ebitenutil.DrawRect(screen, 100, 100, ScreenWidth-200, ScreenHeight-200, color.RGBA{50, 50, 80, 240})
    
    // Draw menu title
    titleX := ScreenWidth/2 - len(currentMenu.Title)*4
    DrawText(screen, currentMenu.Title, titleX, 120)
    
    // Draw menu items
    for i, item := range currentMenu.Items {
        itemY := 160 + (i * 40)
        itemText := item.Text
        
        // Add indicator for submenu
        if item.Type == menu.SubmenuItem {
            itemText += " ►"
        }
        
        // Draw selection indicator for selected item
        if item.Selected {
            DrawText(screen, "> " + itemText, ScreenWidth/2 - 100, itemY)
        } else {
            DrawText(screen, "  " + itemText, ScreenWidth/2 - 100, itemY)
        }
    }
    
    // Draw instructions
    DrawText(screen, "↑/↓: Navigate, Enter: Select", ScreenWidth/2 - 120, ScreenHeight - 150)
}

// Update the Draw method to include the menu state
func (r *Renderer) Draw(
    screen *ebiten.Image,
    gameState int, // GameState
    mazeObj *maze.Maze,
    playerObj *player.Player,
    npcManager *npc.Manager,
    turnManager *turn.Manager,
    triviaManager *trivia.Manager,
    actionManager *action.Manager,
    menuManager *menu.Manager,
    flavorManager *flavor.Manager, // Add flavor manager
    winner string,
) {
    // Draw background
    screen.Fill(color.RGBA{40, 45, 55, 255})

    switch gameState {
    case 0: // Menu
        r.drawMenu(screen, menuManager)
    case 1: // Playing
        r.drawPlayingSplitScreen(screen, mazeObj, playerObj, npcManager, turnManager, actionManager, flavorManager)
    case 2: // AnsweringTrivia
        r.drawTrivia(screen, triviaManager)
    case 3: // GameOver
        r.drawGameOver(screen, winner)
    }
}

// Add a new method for split-screen rendering
func (r *Renderer) drawPlayingSplitScreen(
    screen *ebiten.Image,
    mazeObj *maze.Maze,
    playerObj *player.Player,
    npcManager *npc.Manager,
    turnManager *turn.Manager,
    actionManager *action.Manager,
    flavorManager *flavor.Manager,
) {
    // Create a layout manager
    layout := NewLayoutManager(ScreenWidth, ScreenHeight)
    
    // Get the maze section
    mazeSection := layout.GetSection(MazeSection)
    
    // Draw maze section border
    if mazeSection.Border {
        // Draw section border
        borderColor := color.RGBA{70, 70, 100, 255}
        ebitenutil.DrawRect(
            screen,
            float64(mazeSection.Rect.X),
            float64(mazeSection.Rect.Y),
            float64(mazeSection.Rect.Width),
            float64(mazeSection.Rect.Height),
            borderColor,
        )
    }
    
    // Draw section title
    if mazeSection.Title != "" {
        DrawText(screen, mazeSection.Title, mazeSection.Rect.X + 10, mazeSection.Rect.Y + 20)
    }
    
    // Draw the maze with proper offset to center it in the section
    // Account for title height
    mazeOffsetY := float64(mazeSection.Rect.Y) + 40 // Add space for title
    
    // Calculate maze width and height in pixels
    mazeWidthPixels := float64(mazeObj.State.Width * int(maze.TileSize))
    //mazeHeightPixels := float64(mazeObj.State.Height * int(maze.TileSize))
    
    // Center the maze in the section
    mazeOffsetX := float64(mazeSection.Rect.X) + (float64(mazeSection.Rect.Width) - mazeWidthPixels) / 2
    
    // Draw the maze
    DrawMaze(screen, mazeObj, mazeOffsetX, mazeOffsetY)
    
    // Draw NPCs
    for _, npc := range npcManager.NPCs {
        ebitenutil.DrawRect(
            screen, 
            mazeOffsetX + npc.X + 1, 
            mazeOffsetY + npc.Y + 1, 
            npc.Size, 
            npc.Size, 
            npc.Color,
        )
    }
    
    // Draw player
    playerX, playerY := playerObj.GetPosition()
    ebitenutil.DrawRect(
        screen, 
        mazeOffsetX + playerX + 1, 
        mazeOffsetY + playerY + 1, 
        playerObj.Size, 
        playerObj.Size, 
        color.RGBA{0, 0, 255, 255},
    )
    
    // Get the flavor section
    flavorSection := layout.GetSection(FlavorSection)
    
    // Draw flavor section border
    if flavorSection.Border {
        // Draw section border
        borderColor := color.RGBA{70, 70, 100, 255}
        ebitenutil.DrawRect(
            screen,
            float64(flavorSection.Rect.X),
            float64(flavorSection.Rect.Y),
            float64(flavorSection.Rect.Width),
            float64(flavorSection.Rect.Height),
            borderColor,
        )
    }
    
    // Draw section title
    if flavorSection.Title != "" {
        DrawText(screen, flavorSection.Title, flavorSection.Rect.X + 10, flavorSection.Rect.Y + 20)
    }
    
    // Draw flavor image if available
    if flavorManager != nil && flavorManager.CurrentImage != nil {
        // Draw the flavor image in its section
        flavorManager.Draw(
            screen,
            flavorSection.Rect.X,
            flavorSection.Rect.Y + 40, // Add space for title
            flavorSection.Rect.Width,
            flavorSection.Rect.Height - 40,
        )
    } else {
        // Draw a placeholder message
        DrawText(
            screen,
            "No flavor image available",
            flavorSection.Rect.X + 50,
            flavorSection.Rect.Y + 100,
        )
    }
    
    // Draw UI info in the maze section
    // Display near the top of the maze section
    DrawText(screen, turnManager.OwnerText(), mazeSection.Rect.X + 10, mazeSection.Rect.Y + 60)
    DrawText(screen, turnManager.StateText(), mazeSection.Rect.X + 10, mazeSection.Rect.Y + 80)
    
    // Draw action selection popup if in SelectingAction state
    if turnManager.CurrentState == turn.SelectingAction {
        r.drawActionPopup(screen, actionManager)
    }
    
    // Draw action message if active - overlay at the bottom of the screen
    if r.actionMsg != "" {
        // Calculate message width for centering
        msgWidth := len(r.actionMsg) * 14 // Approximate width based on character count
        
        // Draw a background rectangle for the message
        msgBgX := ScreenWidth/2 - msgWidth/2 - 10
        msgBgWidth := msgWidth + 20
        
        ebitenutil.DrawRect(screen, float64(msgBgX), ScreenHeight-60, float64(msgBgWidth), 30, color.RGBA{0, 0, 0, 180})
        DrawText(screen, r.actionMsg, ScreenWidth/2-msgWidth/2, ScreenHeight-50)
    }
}

// Draw the game over screen
func (r *Renderer) drawGameOver(screen *ebiten.Image, winner string) {
	// Draw message background
	ebitenutil.DrawRect(screen, 100, 200, ScreenWidth-200, 100, color.RGBA{50, 50, 80, 240})
	
	// Draw winner message
	winMessage := fmt.Sprintf("%s reached the goal first and won!", winner)
	DrawText(screen, winMessage, ScreenWidth/2-120, ScreenHeight/2-10)
	DrawText(screen, "Press SPACE to restart", ScreenWidth/2-100, ScreenHeight/2+20)
}

// Draw the playing state
func (r *Renderer) drawPlaying(
	screen *ebiten.Image,
	mazeObj *maze.Maze,
	playerObj *player.Player,
	npcManager *npc.Manager,
	turnManager *turn.Manager,
	actionManager *action.Manager,
) {
	// Draw the maze grid using our new function
	DrawMaze(screen, mazeObj, 0, 0) // Use 0, 0 as the offset (or adjust as needed)

	// Draw NPCs
	for _, npc := range npcManager.NPCs {
		ebitenutil.DrawRect(screen, npc.X+1, npc.Y+1, npc.Size, npc.Size, npc.Color)
	}

	// Draw player
	playerX, playerY := playerObj.GetPosition()
	ebitenutil.DrawRect(screen, playerX+1, playerY+1, playerObj.Size, playerObj.Size, color.RGBA{0, 0, 255, 255})

	// Draw UI info
	r.drawUI(screen, turnManager)

	// Draw action selection popup if in SelectingAction state
	if turnManager.CurrentState == turn.SelectingAction {
		r.drawActionPopup(screen, actionManager)
	}

	// Draw action message if active
	if r.actionMsg != "" {
		// Calculate message width for centering
		msgWidth := len(r.actionMsg) * 14 // Approximate width based on character count
		
		// Draw a background rectangle for the message
		msgBgX := ScreenWidth/2 - msgWidth/2 - 10
		msgBgWidth := msgWidth + 20
		
		ebitenutil.DrawRect(screen, float64(msgBgX), ScreenHeight-60, float64(msgBgWidth), 30, color.RGBA{0, 0, 0, 180})
		DrawText(screen, r.actionMsg, ScreenWidth/2-msgWidth/2, ScreenHeight-50)
	}
}

// Draw the UI
func (r *Renderer) drawUI(screen *ebiten.Image, turnManager *turn.Manager) {
	// Draw turn info using the turn manager
	DrawText(screen, turnManager.OwnerText(), 10, 10)

	// Draw turn state info using the turn manager
	DrawText(screen, turnManager.StateText(), 10, 30)

	// Draw goal info
	DrawText(screen, "Reach the purple goal to win!", 10, 50)
}

// Draw the action selection popup
func (r *Renderer) drawActionPopup(screen *ebiten.Image, actionManager *action.Manager) {
	// Get formatted list of available actions
	actionText := actionManager.FormatActionsList()
	lines := strings.Split(actionText, "\n")
	
	// Calculate popup dimensions based on content
	// Find the longest line to determine width
	maxLineLength := 0
	for _, line := range lines {
		if len(line) > maxLineLength {
			maxLineLength = len(line)
		}
	}
	
	// Calculate width and height with padding
	width := maxLineLength*7 + 40 // Approximate width based on character count plus padding
	if width < 300 {
		width = 300 // Minimum width
	}
	
	height := 40 + (len(lines) * 20) // Height based on number of lines plus padding
	if height < 100 {
		height = 100 // Minimum height
	}
	
	x := (ScreenWidth - width) / 2
	y := (ScreenHeight - height) / 2
	
	// Draw popup background
	ebitenutil.DrawRect(screen, float64(x), float64(y), float64(width), float64(height), color.RGBA{70, 70, 100, 240})
	ebitenutil.DrawRect(screen, float64(x+2), float64(y+2), float64(width-4), float64(height-4), color.RGBA{40, 40, 70, 240})
	
	// Draw action list
	for i, line := range lines {
		DrawText(screen, line, x+10, y+20+(i*20))
	}
	
	// Draw instructions at the bottom
	DrawText(screen, "Press number to select, ESC to cancel", x+10, y+height-20)
}

// Draw the trivia screen
func (r *Renderer) drawTrivia(screen *ebiten.Image, triviaManager *trivia.Manager) {
	currentQuestion := triviaManager.GetCurrentQuestion()

	// Draw question background
	ebitenutil.DrawRect(screen, 50, 50, ScreenWidth-100, ScreenHeight-100, color.RGBA{50, 50, 80, 240})

	// Draw question
	DrawText(screen, currentQuestion.Question, 70, 70)

	// Draw options
	for i, option := range currentQuestion.Options {
		optionYpadding := 60 * i
		DrawText(screen, fmt.Sprintf("%d: %s", i+1, option), 70, (140 + optionYpadding))
	}

	// Draw instructions
	DrawText(screen, "Press 1-4 to answer", 70, ScreenHeight-100)

	// If answered, show result
	if triviaManager.Answered {
		resultText := "Incorrect!"
		//resultColor := color.RGBA{255, 0, 0, 255}

		if triviaManager.Correct {
			resultText = "Correct!"
			//resultColor = color.RGBA{0, 255, 0, 255}
		}

		// Calculate message width for centering
		msgWidth := len(resultText) * 9 // Approximate width based on character count
		
		DrawText(screen, resultText, ScreenWidth/2-msgWidth/2, ScreenHeight/2)
	}
}

