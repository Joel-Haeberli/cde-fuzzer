package cde

// CDE represents a Common Data Element — a question with a set of possible answers.
type CDE struct {
	ID       string
	Question string
	Answers  []string
}
