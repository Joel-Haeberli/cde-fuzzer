package core

import (
	"fmt"
	"strings"

	"github.com/Joel-Haeberli/cde-fuzzer/internal/result"
)

// MockLLMClient is a simple mock implementation for testing
type MockLLMClient struct{}

func (m *MockLLMClient) Complete(prompt, text string) (string, error) {
	// Simple mock logic: extract information based on the prompt
	// In a real implementation, this would call an actual LLM API

	// Check if the prompt asks for something specific
	promptLower := strings.ToLower(prompt)
	textLower := strings.ToLower(text)

	// Example: if prompt asks for age and text contains age information
	if strings.Contains(promptLower, "age") && strings.Contains(promptLower, "extract") {
		if strings.Contains(textLower, "years old") || strings.Contains(textLower, "year-old") {
			// Extract age number
			for i := 0; i < len(text); i++ {
				if text[i] >= '0' && text[i] <= '9' {
					start := i
					for i < len(text) && text[i] >= '0' && text[i] <= '9' {
						i++
					}
					ageStr := text[start:i]
					// Check if this is followed by age-related terms
					suffix := strings.ToLower(text[i:min(i+20, len(text))])
					if strings.Contains(suffix, "year") || strings.Contains(suffix, "old") {
						return ageStr, nil
					}
				}
			}
		}
		return "no data", nil
	}

	// Example: if prompt asks for diagnosis (be specific to avoid false positives)
	if (strings.Contains(promptLower, "diagnosis") && strings.Contains(promptLower, "extract")) ||
		(strings.Contains(promptLower, "carcinoma") && strings.Contains(promptLower, "extract") && !strings.Contains(promptLower, "er")) {
		if strings.Contains(textLower, "carcinoma") {
			// Find and return the carcinoma mention
			start := strings.Index(textLower, "carcinoma")
			if start >= 0 {
				// Find the start of the word (go back to find space or beginning)
				wordStart := start
				for wordStart > 0 && text[wordStart-1] != ' ' && text[wordStart-1] != '\t' && text[wordStart-1] != '\n' {
					wordStart--
				}
				// Find the end of the word
				wordEnd := start + len("carcinoma")
				for wordEnd < len(text) && text[wordEnd] != ' ' && text[wordEnd] != ',' && text[wordEnd] != '.' {
					wordEnd++
				}
				return text[wordStart:wordEnd], nil
			}
		}
		return "no data", nil
	}

	// Example: if prompt asks for ER status
	if strings.Contains(promptLower, "er status") || strings.Contains(promptLower, "estrogen receptor") || strings.Contains(promptLower, "extract the er") {
		if strings.Contains(textLower, "er positive") {
			// Find the exact "ER positive" text
			start := strings.Index(textLower, "er positive")
			if start >= 0 {
				// Find the start of the phrase
				wordStart := start
				for wordStart > 0 && text[wordStart-1] != ' ' && text[wordStart-1] != '\t' && text[wordStart-1] != '\n' {
					wordStart--
				}
				// Find the end of the phrase (up to space, comma, or parenthesis)
				wordEnd := start + len("er positive")
				for wordEnd < len(text) && text[wordEnd] != ' ' && text[wordEnd] != ',' && text[wordEnd] != ')' {
					wordEnd++
				}
				return text[wordStart:wordEnd], nil
			}
		}
		if strings.Contains(textLower, "er negative") {
			start := strings.Index(textLower, "er negative")
			if start >= 0 {
				wordStart := start
				for wordStart > 0 && text[wordStart-1] != ' ' && text[wordStart-1] != '\t' && text[wordStart-1] != '\n' {
					wordStart--
				}
				wordEnd := start + len("er negative")
				for wordEnd < len(text) && text[wordEnd] != ' ' && text[wordEnd] != ',' && text[wordEnd] != ')' {
					wordEnd++
				}
				return text[wordStart:wordEnd], nil
			}
		}
		if strings.Contains(textLower, "er+") {
			return "ER+", nil
		}
		if strings.Contains(textLower, "er-") {
			return "ER-", nil
		}
		return "no data", nil
	}

	// Default: return "no data" if we can't determine what to extract
	return "no data", nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// LLMRule implements Rule using a language model prompt for extraction.
// It sends the text to an LLM with a specific prompt and extracts the answer.
type LLMRule struct {
	name      string
	prompt    string
	accuracy  float64
	llmClient LLMClient
}

// LLMClient is an interface for interacting with language models
type LLMClient interface {
	Complete(prompt, text string) (string, error)
}

// NewLLMRule creates a rule that uses an LLM for extraction
func NewLLMRule(name, prompt string, accuracy float64, client LLMClient) *LLMRule {
	return &LLMRule{
		name:      name,
		prompt:    prompt,
		accuracy:  accuracy,
		llmClient: client,
	}
}

func (r *LLMRule) Name() string { return r.name }

func (r *LLMRule) Match(text string) bool {
	// For LLM rules, we always attempt to match since the LLM determines relevance
	return true
}

func (r *LLMRule) Apply(text string) ([]result.Match, error) {
	// Construct the full prompt by combining the rule prompt with the text
	fullPrompt := fmt.Sprintf("%s\n\nText: %s\n\nAnswer:", r.prompt, text)

	// Get the LLM completion
	answer, err := r.llmClient.Complete(fullPrompt, text)
	if err != nil {
		return nil, fmt.Errorf("llm completion failed: %v", err)
	}

	// Check if the answer indicates no data was found
	if strings.ToLower(strings.TrimSpace(answer)) == "no data" {
		return []result.Match{}, nil
	}

	// Find the answer in the original text to get proper offsets
	start := strings.Index(text, answer)
	if start >= 0 {
		return []result.Match{
			{
				Value: answer,
				Start: start,
				End:   start + len(answer),
			},
		}, nil
	}

	// If we can't find the exact answer in the text, return it anyway with offset 0
	return []result.Match{
		{
			Value: answer,
			Start: 0,
			End:   len(answer),
		},
	}, nil
}

func (r *LLMRule) Accuracy() float64 { return r.accuracy }
