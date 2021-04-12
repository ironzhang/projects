package spam

import "strings"

func NewEngine(patterns []string) *Engine {
	return &Engine{patterns: patterns}
}

type Engine struct {
	patterns []string
}

func (e *Engine) IsSpam(message string) bool {
	for _, pattern := range e.patterns {
		if strings.Contains(message, pattern) {
			return true
		}
	}
	return false
}
