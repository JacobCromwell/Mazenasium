// internal/game/menu/menu.go
package menu

import (
    "github.com/hajimehoshi/ebiten/v2"
    "github.com/hajimehoshi/ebiten/v2/inpututil"
)

type ItemType int

const (
    ButtonItem ItemType = iota
    SubmenuItem
)

type Item struct {
    Text     string
    Type     ItemType
    Selected bool
    Action   string
    Submenu  *Menu
}

type Menu struct {
    Title    string
    Items    []Item
    Selected int
    Parent   *Menu
}

type Manager struct {
    CurrentMenu *Menu
    RootMenu    *Menu
}

// NewManager creates a new menu manager with default menus
func NewManager() *Manager {
    // Create root menu
    rootMenu := &Menu{
        Title: "Mazenasium",
        Items: []Item{
            {Text: "Start Game", Type: ButtonItem, Selected: true, Action: "start_game"},
            {Text: "Customize", Type: SubmenuItem},
            {Text: "Quit", Type: ButtonItem, Action: "quit"},
        },
        Selected: 0,
    }
    
    // Create customize submenu
    customizeMenu := &Menu{
        Title: "Customize",
        Items: []Item{
            {Text: "Option 1", Type: ButtonItem, Selected: true, Action: "option1"},
            {Text: "Option 2", Type: ButtonItem, Action: "option2"},
            {Text: "Back", Type: ButtonItem, Action: "back"},
        },
        Selected: 0,
    }
    
    // Link menus
    rootMenu.Items[1].Submenu = customizeMenu
    customizeMenu.Parent = rootMenu
    
    return &Manager{
        CurrentMenu: rootMenu,
        RootMenu:    rootMenu,
    }
}

// MoveSelectionUp moves the selection up in the current menu
func (m *Manager) MoveSelectionUp() {
    if m.CurrentMenu == nil || len(m.CurrentMenu.Items) == 0 {
        return
    }
    
    // Deselect current item
    m.CurrentMenu.Items[m.CurrentMenu.Selected].Selected = false
    
    // Move selection up
    m.CurrentMenu.Selected--
    if m.CurrentMenu.Selected < 0 {
        m.CurrentMenu.Selected = len(m.CurrentMenu.Items) - 1
    }
    
    // Select new item
    m.CurrentMenu.Items[m.CurrentMenu.Selected].Selected = true
}

// MoveSelectionDown moves the selection down in the current menu
func (m *Manager) MoveSelectionDown() {
    if m.CurrentMenu == nil || len(m.CurrentMenu.Items) == 0 {
        return
    }
    
    // Deselect current item
    m.CurrentMenu.Items[m.CurrentMenu.Selected].Selected = false
    
    // Move selection down
    m.CurrentMenu.Selected = (m.CurrentMenu.Selected + 1) % len(m.CurrentMenu.Items)
    
    // Select new item
    m.CurrentMenu.Items[m.CurrentMenu.Selected].Selected = true
}

// SelectCurrentItem selects the current menu item
// Returns the action string if an action is selected, empty string otherwise
func (m *Manager) SelectCurrentItem() string {
    if m.CurrentMenu == nil || len(m.CurrentMenu.Items) == 0 {
        return ""
    }
    
    currentItem := m.CurrentMenu.Items[m.CurrentMenu.Selected]
    
    if currentItem.Type == SubmenuItem && currentItem.Submenu != nil {
        // Navigate to submenu
        m.CurrentMenu = currentItem.Submenu
        return ""
    } else if currentItem.Action == "back" && m.CurrentMenu.Parent != nil {
        // Navigate back to parent menu
        m.CurrentMenu = m.CurrentMenu.Parent
        return ""
    } else {
        // Return the action
        return currentItem.Action
    }
}

// HandleInput processes keyboard input for menu navigation
func (m *Manager) HandleInput() string {
    if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
        m.MoveSelectionUp()
    }
    
    if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
        m.MoveSelectionDown()
    }
    
    if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
        return m.SelectCurrentItem()
    }
    
    return ""
}