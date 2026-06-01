package intent

import (
	"sync"
	"testing"
)

// Verify the parser is safe for concurrent use (no shared mutable state).
func TestParser_ConcurrentSafety(t *testing.T) {
	p := NewParser()

	inputs := []string{
		"search react hooks",
		"open youtube play lofi",
		"go back",
		"open github",
		"close tab",
		"search for golang patterns",
		"open spotify jazz",
		"scroll down",
		"open example.com",
		"what is the meaning of life",
	}

	var wg sync.WaitGroup
	errors := make(chan string, 100)

	// Run 100 concurrent parses
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			input := inputs[idx%len(inputs)]
			result := p.Parse(input)

			// Basic sanity checks
			if result.Action == "" {
				errors <- "empty action for: " + input
			}
			if result.Confidence < 0 || result.Confidence > 1 {
				errors <- "invalid confidence for: " + input
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		t.Error(err)
	}
}

// Verify multiple parser instances don't interfere with each other.
func TestParser_MultipleInstances(t *testing.T) {
	p1 := NewParser()
	p2 := NewParser()

	r1 := p1.Parse("open youtube")
	r2 := p2.Parse("search react")

	if r1.Action != "navigate" {
		t.Errorf("p1: action = %q, want navigate", r1.Action)
	}
	if r2.Action != "search" {
		t.Errorf("p2: action = %q, want search", r2.Action)
	}
}
