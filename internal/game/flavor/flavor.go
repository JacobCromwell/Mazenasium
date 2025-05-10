// internal/game/flavor/flavor.go
package flavor

import (
	"fmt"
    "image"
    _ "image/jpeg" // Register JPEG decoder
    "os"
    "path/filepath"

    "github.com/hajimehoshi/ebiten/v2"
)

type Manager struct {
    Images       map[string]*ebiten.Image
    CurrentImage *ebiten.Image
    ImageKeys    []string // To allow cycling through images
    CurrentIndex int
}

func NewManager() *Manager {
    return &Manager{
        Images:       make(map[string]*ebiten.Image),
        ImageKeys:    make([]string, 0),
        CurrentIndex: 0,
    }
}

// Update internal/game/flavor/flavor.go
func (m *Manager) LoadImages(assetsDir string) error {
    // Safety check
    if m == nil {
        return fmt.Errorf("flavor manager is nil")
    }

    // Ensure the assets directory exists
    if _, err := os.Stat(assetsDir); os.IsNotExist(err) {
        return fmt.Errorf("assets directory does not exist: %v", err)
    }

    // Look for images in the hallway subdirectory
    hallwayDir := filepath.Join(assetsDir, "hallway")
    
    // Create the directory if it doesn't exist
    if _, err := os.Stat(hallwayDir); os.IsNotExist(err) {
        fmt.Println("Creating hallway directory:", hallwayDir)
        if err := os.MkdirAll(hallwayDir, 0755); err != nil {
            return fmt.Errorf("failed to create hallway directory: %v", err)
        }
        
        // Return early since the directory is empty
        fmt.Println("Hallway directory is empty, no images to load")
        return nil
    }
    
    // Initialize maps if they haven't been already
    if m.Images == nil {
        m.Images = make(map[string]*ebiten.Image)
    }
    
    if m.ImageKeys == nil {
        m.ImageKeys = make([]string, 0)
    }
    
	return nil
    // Rest of the method...
}

func (m *Manager) UpdateImage(playerX, playerY int) {
    // Simple algorithm to change image based on player position
    // Each tile position maps to a different image
    index := (playerX + playerY) % len(m.ImageKeys)
    if index < 0 {
        index += len(m.ImageKeys)
    }
    
    if len(m.ImageKeys) > 0 {
        m.CurrentIndex = index
        m.CurrentImage = m.Images[m.ImageKeys[index]]
    }
}

func (m *Manager) Draw(screen *ebiten.Image, x, y, width, height int) {
    if m.CurrentImage == nil {
        return
    }
    
    // Draw the current image, scaled to fit the section
    op := &ebiten.DrawImageOptions{}
    
    // Scale image to fit the section while maintaining aspect ratio
    imgWidth, imgHeight := m.CurrentImage.Size()
    scaleX := float64(width) / float64(imgWidth)
    scaleY := float64(height) / float64(imgHeight)
    
    // Use the smaller scale to avoid stretching
    scale := scaleX
    if scaleY < scaleX {
        scale = scaleY
    }
    
    op.GeoM.Scale(scale, scale)
    
    // Center the image in the section
    centeredX := x + (width - int(float64(imgWidth)*scale))/2
    centeredY := y + (height - int(float64(imgHeight)*scale))/2
    
    op.GeoM.Translate(float64(centeredX), float64(centeredY))
    
    screen.DrawImage(m.CurrentImage, op)
}

// Update in internal/game/flavor/flavor.go
func (m *Manager) SetImageByPath(path string) {
    // Safety check
    if m == nil || m.Images == nil {
        fmt.Println("Warning: Flavor manager or images map is nil")
        return
    }
    
    // Check if the image is already loaded
    img, exists := m.Images[path]
    if exists {
        m.CurrentImage = img
        return
    }
    
    // If the image isn't loaded yet, try to load it
    file, err := os.Open(path)
    if err != nil {
        fmt.Printf("Warning: Could not open image %s: %v\n", path, err)
        return
    }
    defer file.Close()
    
    // Decode image
    decodedImg, _, err := image.Decode(file)
    if err != nil {
        fmt.Printf("Warning: Could not decode image %s: %v\n", path, err)
        return
    }
    
    // Convert to ebiten image
    ebitenImg := ebiten.NewImageFromImage(decodedImg)
    
    // Store image and set as current
    m.Images[path] = ebitenImg
    m.ImageKeys = append(m.ImageKeys, path)
    m.CurrentImage = ebitenImg
}