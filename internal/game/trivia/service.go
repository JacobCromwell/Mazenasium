// pkg/trivia/service.go
package trivia

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
)

// Manager handles loading, selecting, and checking answers for trivia questions
type Manager struct {
	questionSets []QuestionSet
	currentSet   *QuestionSet
	currentIndex int
	answered     bool
	lastResult   *Result
}

// NewManager creates a new trivia manager
func NewManager() *Manager {
	return &Manager{
		questionSets: []QuestionSet{},
		currentIndex: 0,
		answered:     false,
	}
}

// LoadQuestionSet loads trivia questions from a file
func (m *Manager) LoadQuestionSet(filename string) error {
	// For initial implementation, we'll just load built-in questions
	// In the future, this would load from a JSON file
	if filename != "" {
		data, err := os.ReadFile(filename)
		if err != nil {
			return fmt.Errorf("failed to load questions from %s: %w", filename, err)
		}

		var questionSet QuestionSet
		if err := json.Unmarshal(data, &questionSet); err != nil {
			return fmt.Errorf("failed to parse questions from %s: %w", filename, err)
		}

		m.questionSets = append(m.questionSets, questionSet)
		if m.currentSet == nil {
			m.currentSet = &m.questionSets[0]
		}
		return nil
	}

	// If no file specified, load default questions
	defaultSet := m.getDefaultQuestions()
	m.questionSets = append(m.questionSets, defaultSet)
	m.currentSet = &m.questionSets[0]
	return nil
}

// GetCurrentQuestion returns the current trivia question
func (m *Manager) GetCurrentQuestion() *Question {
	if m.currentSet == nil || len(m.currentSet.Questions) == 0 {
		return nil
	}
	return &m.currentSet.Questions[m.currentIndex]
}

// NextQuestion advances to the next trivia question
func (m *Manager) NextQuestion() *Question {
	if m.currentSet == nil || len(m.currentSet.Questions) == 0 {
		return nil
	}

	m.answered = false
	m.lastResult = nil
	m.currentIndex = (m.currentIndex + 1) % len(m.currentSet.Questions)
	return m.GetCurrentQuestion()
}

// RandomQuestion selects a random question
func (m *Manager) RandomQuestion() *Question {
	if m.currentSet == nil || len(m.currentSet.Questions) == 0 {
		return nil
	}

	m.answered = false
	m.lastResult = nil
	m.currentIndex = rand.Intn(len(m.currentSet.Questions))
	return m.GetCurrentQuestion()
}

// Answer checks if the provided answer index is correct
func (m *Manager) Answer(selectedIndex int) *Result {
	question := m.GetCurrentQuestion()
	if question == nil {
		return nil
	}

	isCorrect := selectedIndex == question.Answer
	result := &Result{
		Question:   question,
		Selected:   selectedIndex,
		IsCorrect:  isCorrect,
		PointValue: 1, // Default point value
	}

	m.answered = true
	m.lastResult = result
	return result
}

// IsAnswered returns whether the current question has been answered
func (m *Manager) IsAnswered() bool {
	return m.answered
}

// GetLastResult returns the result of the last answered question
func (m *Manager) GetLastResult() *Result {
	return m.lastResult
}

// GetDefaultQuestions returns a default set of trivia questions
func (m *Manager) getDefaultQuestions() QuestionSet {
	return QuestionSet{
		Name:     "General Knowledge",
		Category: "Mixed",
		Questions: []Question{
			{
				Text:    "What is the capital of France?",
				Options: []string{"London", "Berlin", "Paris", "Madrid"},
				Answer:  2, // Paris (0-indexed)
			},
			{
				Text:    "Which planet is known as the Red Planet?",
				Options: []string{"Venus", "Mars", "Jupiter", "Saturn"},
				Answer:  1, // Mars
			},
			{
				Text:    "What is the largest mammal?",
				Options: []string{"Elephant", "Giraffe", "Blue Whale", "Hippopotamus"},
				Answer:  2, // Blue Whale
			},
			{
				Text:    "What element has the chemical symbol 'O'?",
				Options: []string{"Gold", "Oxygen", "Osmium", "Oganesson"},
				Answer:  1, // Oxygen
			},
			{
				Text:    "Who wrote 'Romeo and Juliet'?",
				Options: []string{"Charles Dickens", "William Shakespeare", "Jane Austen", "Mark Twain"},
				Answer:  1, // Shakespeare
			},
			{
				Text:    "Which of these is not a primary color of light?",
				Options: []string{"Red", "Green", "Yellow", "Blue"},
				Answer:  2, // Yellow
			},
			{
				Text:    "What is the capital of Japan?",
				Options: []string{"Seoul", "Beijing", "Tokyo", "Bangkok"},
				Answer:  2, // Tokyo
			},
			{
				Text:    "Which language is most widely spoken in Brazil?",
				Options: []string{"Spanish", "Portuguese", "English", "French"},
				Answer:  1, // Portuguese
			},
			{
				Text:    "What year did the first person land on the moon?",
				Options: []string{"1965", "1969", "1973", "1980"},
				Answer:  1, // 1969
			},
			{
				Text:    "Which of these is not one of the Great Lakes?",
				Options: []string{"Michigan", "Ontario", "Huron", "Okeechobee"},
				Answer:  3, // Okeechobee
			},
		},
	}
}
