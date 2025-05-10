// internal/game/flavor/flavor.go
package flavor

import (
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

func (m *Manager) LoadImages(assetsDir string) error {
    // Load all JPG images from the assets directory
    entries, err := os.ReadDir(assetsDir)
    if err != nil {
        return err
    }
    
    for _, entry := range entries {
        if !entry.IsDir() && (filepath.Ext(entry.Name()) == ".jpg" || filepath.Ext(entry.Name()) == ".jpeg") {
            path := filepath.Join(assetsDir, entry.Name())
            
            // Load image file
            file, err := os.Open(path)
            if err != nil {
                continue
            }
            defer file.Close()
            
            // Decode image
            img, _, err := image.Decode(file)
            if err != nil {
                continue
            }
            
            // Convert to ebiten image
            ebitenImg := ebiten.NewImageFromImage(img)
            
            // Store image with key (filename without extension)
            key := filepath.Base(entry.Name())
            key = key[:len(key)-len(filepath.Ext(key))]
            m.Images[key] = ebitenImg
            m.ImageKeys = append(m.ImageKeys, key)
        }
    }
    
    // Set initial image if available
    if len(m.ImageKeys) > 0 {
        m.CurrentImage = m.Images[m.ImageKeys[0]]
    }
    
    return nil
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