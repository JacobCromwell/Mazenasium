package maze

import (
	"fmt"
    "math/rand"
)

// Generator handles maze generation algorithms
type Generator struct {
    // Any configuration options for generation
    RandomSeed int64
}

// NewGenerator creates a new maze generator
func NewGenerator(seed int64) *Generator {
    return &Generator{
        RandomSeed: seed,
    }
}

// Generate creates a new maze with the given dimensions
func (g *Generator) Generate(width, height int) *State {
    // Create a new empty state
    state := NewState(width, height)
    
    // Use a local random source to ensure deterministic generation with the same seed
    r := rand.New(rand.NewSource(g.RandomSeed))
    
    // Generate the maze using a depth-first search algorithm
    g.generatePathways(state, 1, 1, r)
    
    // Add some random additional paths
    g.addRandomPaths(state, r)
    
    // Choose a goal position in the bottom-right quarter
    goalX, goalY := g.chooseGoalPosition(state, r)
    state.SetTileType(goalX, goalY, Goal)
    state.GoalX = goalX
    state.GoalY = goalY
    
    // Ensure there's a path to the goal
    g.ensurePathToGoal(state, 1, 1, goalX, goalY)
    
    // Ensure the starting positions for player and NPCs are clear
    state.SetTileType(1, 1, Floor) // Player start
    state.SetTileType(3, 3, Floor) // NPC1 start
    state.SetTileType(5, 5, Floor) // NPC2 start
    
    // Set flavor images for tiles
    g.setFlavorImages(state)
    
    return state
}

// chooseGoalPosition selects a position for the goal
func (g *Generator) chooseGoalPosition(state *State, r *rand.Rand) (int, int) {
    width, height := state.Width, state.Height
    
    // Choose a goal in the bottom-right quarter
    goalX, goalY := 0, 0
    for {
        goalX = width - 2 - r.Intn(width/4)
        goalY = height - 2 - r.Intn(height/4)
        
        // Ensure the goal isn't too close to the start
        if abs(goalX-1) + abs(goalY-1) >= (width + height)/3 {
            break
        }
    }
    
    return goalX, goalY
}

// generatePathways creates the initial maze structure
func (g *Generator) generatePathways(state *State, startX, startY int, r *rand.Rand) {
    // Initialize visited grid
    visited := make([][]bool, state.Height)
    for y := range visited {
        visited[y] = make([]bool, state.Width)
    }

    // Stack for backtracking
    stack := []Position{{X: startX, Y: startY}}
    
    // Mark the starting position as a floor and visited
    state.SetTileType(startX, startY, Floor)
    visited[startY][startX] = true

    // Directions: North, East, South, West
    dx := []int{0, 1, 0, -1}
    dy := []int{-1, 0, 1, 0}

    // Keep going until the stack is empty
    for len(stack) > 0 {
        // Get the current position
        current := stack[len(stack)-1]
        
        // Find unvisited neighbors
        neighbors := []int{}
        for d := 0; d < 4; d++ {
            nx, ny := current.X + dx[d]*2, current.Y + dy[d]*2
            
            // Check if the neighbor is valid and unvisited
            if nx >= 1 && nx < state.Width-1 && ny >= 1 && ny < state.Height-1 && !visited[ny][nx] {
                neighbors = append(neighbors, d)
            }
        }
        
        if len(neighbors) > 0 {
            // Randomly choose a neighbor
            d := neighbors[r.Intn(len(neighbors))]
            nx, ny := current.X + dx[d]*2, current.Y + dy[d]*2
            
            // Carve a passage by setting both the neighbor and the wall in between to floor
            betweenX, betweenY := current.X + dx[d], current.Y + dy[d]
            state.SetTileType(betweenX, betweenY, Floor)
            state.SetTileType(nx, ny, Floor)
            
            // Mark as visited and push to stack
            visited[ny][nx] = true
            stack = append(stack, Position{X: nx, Y: ny})
        } else {
            // No unvisited neighbors, backtrack
            stack = stack[:len(stack)-1]
        }
    }
}

// addRandomPaths adds some random paths to make the maze more interesting
func (g *Generator) addRandomPaths(state *State, r *rand.Rand) {
    // Number of random paths to add (adjustable)
    extraPaths := (state.Width + state.Height) / 3
    
    for i := 0; i < extraPaths; i++ {
        // Pick a random wall that's not on the border
        x, y := 0, 0
        for {
            x = r.Intn(state.Width-2) + 1
            y = r.Intn(state.Height-2) + 1
            
            if state.GetTile(x, y).Type == Wall && 
               x > 0 && x < state.Width-1 && y > 0 && y < state.Height-1 {
                break
            }
        }
        
        // Count adjacent floor tiles
        floorCount := 0
        if state.GetTile(y-1, x) != nil && state.GetTile(y-1, x).Type == Floor { floorCount++ }
        if state.GetTile(y+1, x) != nil && state.GetTile(y+1, x).Type == Floor { floorCount++ }
        if state.GetTile(y, x-1) != nil && state.GetTile(y, x-1).Type == Floor { floorCount++ }
        if state.GetTile(y, x+1) != nil && state.GetTile(y, x+1).Type == Floor { floorCount++ }
        
        // Only remove walls that connect two different passages
        // This creates loops in the maze
        if floorCount >= 2 {
            state.SetTileType(x, y, Floor)
        }
    }
}

// ensurePathToGoal makes sure there's a path from start to goal
func (g *Generator) ensurePathToGoal(state *State, startX, startY, goalX, goalY int) {
    // Use breadth-first search to check if there's a path
    if g.hasPath(state, startX, startY, goalX, goalY) {
        return
    }
    
    // If no path exists, create one
    currentX, currentY := startX, startY
    
    // Move toward the goal with a slight randomness
    for currentX != goalX || currentY != goalY {
        // Decide whether to move in X or Y direction
        moveX := rand.Intn(2) == 0
        
        if moveX && currentX != goalX {
            // Move in X direction
            dx := 1
            if currentX > goalX {
                dx = -1
            }
            state.SetTileType(currentX + dx, currentY, Floor)
            currentX += dx
        } else if currentY != goalY {
            // Move in Y direction
            dy := 1
            if currentY > goalY {
                dy = -1
            }
            state.SetTileType(currentX, currentY + dy, Floor)
            currentY += dy
        }
        
        // Set current position to floor
        state.SetTileType(currentX, currentY, Floor)
    }
}

// hasPath checks if there's a path from start to goal
func (g *Generator) hasPath(state *State, startX, startY, goalX, goalY int) bool {
    // Initialize visited grid
    visited := make([][]bool, state.Height)
    for y := range visited {
        visited[y] = make([]bool, state.Width)
    }
    
    // Queue for BFS
    queue := []Position{{X: startX, Y: startY}}
    visited[startY][startX] = true
    
    // Directions: North, East, South, West
    dx := []int{0, 1, 0, -1}
    dy := []int{-1, 0, 1, 0}
    
    // BFS
    for len(queue) > 0 {
        // Get the current position
        current := queue[0]
        queue = queue[1:]
        
        // Check if we reached the goal
        if current.X == goalX && current.Y == goalY {
            return true
        }
        
        // Check all four directions
        for d := 0; d < 4; d++ {
            nx, ny := current.X + dx[d], current.Y + dy[d]
            
            // Check if the neighbor is valid, a floor, and unvisited
            if nx >= 0 && nx < state.Width && ny >= 0 && ny < state.Height {
                tile := state.GetTile(nx, ny)
                if tile != nil && tile.Type != Wall && !visited[ny][nx] {
                    visited[ny][nx] = true
                    queue = append(queue, Position{X: nx, Y: ny})
                }
            }
        }
    }
    
    return false
}

// setFlavorImages assigns flavor images to tiles
func (g *Generator) setFlavorImages(state *State) {
    // Assign flavor images to non-wall tiles based on their ID
    for y := 0; y < state.Height; y++ {
        for x := 0; x < state.Width; x++ {
            tile := state.GetTile(x, y)
            if tile == nil {
                continue
            }
            
            // Skip walls - they don't get flavor images
            if tile.Type == Wall {
                continue
            }
            
            // Assign flavor image based on tile ID
            // Format: assets/hallway/{id}.jpg
            //imagePath := fmt.Sprintf("assets/hallway/%d.jpg", tile.ID)

			if((tile.ID % 2) == 0){
				imagePath := fmt.Sprintf("assets/hallway/1.jpg")	
				tile.SetFlavorImage(imagePath)
			} else {
				imagePath := fmt.Sprintf("assets/hallway/2.jpg")
				tile.SetFlavorImage(imagePath)
			}

            //tile.SetFlavorImage(imagePath)
        }
    }
}

// Helper functions...
func abs(x int) int {
    if x < 0 {
        return -x
    }
    return x
}
