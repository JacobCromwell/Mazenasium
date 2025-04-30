// pkg/trivia/models.go
package trivia

// Question represents a single trivia question with options and the correct answer
type Question struct {
	Text    string   // The question text
	Options []string // Available answer options
	Answer  int      // Index of the correct answer (0-based)
}

// Result represents the outcome of answering a question
type Result struct {
	Question   *Question // The question that was answered
	Selected   int       // The option that was selected
	IsCorrect  bool      // Whether the selection was correct
	PointValue int       // How many points this question was worth
}

// QuestionSet represents a collection of questions that can be used in the game
type QuestionSet struct {
	Questions []Question
	Name      string
	Category  string
}
