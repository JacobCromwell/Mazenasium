package trivia

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Manager handles trivia questions and answers
type Manager struct {
	Questions    []Question
	CurrentIndex int
	Answered     bool
	Correct      bool
}

// Question represents a single trivia question
type Question struct {
	Question string
	Options  []string
	Answer   int
}

// NewManager creates a new trivia manager with default questions
func NewManager() *Manager {
	return &Manager{
		Questions:    LoadDefaultQuestions(),
		CurrentIndex: 0,
		Answered:     false,
	}
}

// LoadDefaultQuestions returns a list of default trivia questions
func LoadDefaultQuestions() []Question {
	// In a real implementation, you'd load these from a file
	return []Question{
		{
			Question: "What is the capital of France?",
			Options:  []string{"London", "Berlin", "Paris", "Madrid"},
			Answer:   2, // Paris (0-indexed)
		},
		{
			Question: "Which planet is known as the Red Planet?",
			Options:  []string{"Venus", "Mars", "Jupiter", "Saturn"},
			Answer:   1, // Mars
		},
		{
			Question: "What is the largest mammal?",
			Options:  []string{"Elephant", "Giraffe", "Blue Whale", "Hippopotamus"},
			Answer:   2, // Blue Whale
		},
		{
			Question: "What element has the chemical symbol 'O'?",
			Options:  []string{"Gold", "Oxygen", "Osmium", "Oganesson"},
			Answer:   1, // Oxygen
		},
		{
			Question: "Who wrote 'Romeo and Juliet'?",
			Options:  []string{"Charles Dickens", "William Shakespeare", "Jane Austen", "Mark Twain"},
			Answer:   1, // Shakespeare
		},
	}
}

// GetCurrentQuestion returns the current question
func (m *Manager) GetCurrentQuestion() Question {
	return m.Questions[m.CurrentIndex]
}

// SetRandomQuestion selects a random question from the available questions
func (m *Manager) SetRandomQuestion(randomFunc func(int) int) {
	m.CurrentIndex = randomFunc(len(m.Questions))
	m.Answered = false
}

// CheckAnswer checks if the provided answer index is correct
func (m *Manager) CheckAnswer(answerIndex int) bool {
	m.Answered = true
	m.Correct = (answerIndex == m.Questions[m.CurrentIndex].Answer)
	return m.Correct
}

// HandleInput processes keyboard input for trivia answering
// Returns true if an answer was selected
func (m *Manager) HandleInput() bool {
	question := m.Questions[m.CurrentIndex]
	
	// Check for answer selection
	for i := 0; i < len(question.Options); i++ {
		if inpututil.IsKeyJustPressed(ebiten.Key1 + ebiten.Key(i)) {
			m.CheckAnswer(i)
			return true
		}
	}
	
	return false
}