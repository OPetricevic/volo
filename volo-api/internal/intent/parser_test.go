package intent

import (
	"testing"
)

func TestParser_Search(t *testing.T) {
	p := NewParser()

	tests := []struct {
		name       string
		input      string
		wantAction string
		wantQuery  string
		minConf    float64
	}{
		{"search keyword", "search react hooks", "search", "react hooks", 0.9},
		{"search for keyword", "search for best restaurants", "search", "best restaurants", 0.9},
		{"look up keyword", "look up weather tomorrow", "search", "weather tomorrow", 0.9},
		{"find keyword", "find golang patterns", "search", "golang patterns", 0.9},
		{"fallback to search", "where does burek come from", "search", "where does burek come from", 0.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := p.Parse(tt.input)
			if result.Action != tt.wantAction {
				t.Errorf("action = %q, want %q", result.Action, tt.wantAction)
			}
			if result.Query == nil || *result.Query != tt.wantQuery {
				got := ""
				if result.Query != nil {
					got = *result.Query
				}
				t.Errorf("query = %q, want %q", got, tt.wantQuery)
			}
			if result.Confidence < tt.minConf {
				t.Errorf("confidence = %f, want >= %f", result.Confidence, tt.minConf)
			}
		})
	}
}

func TestParser_Navigate(t *testing.T) {
	p := NewParser()

	tests := []struct {
		name       string
		input      string
		wantTarget string
	}{
		{"open youtube", "open youtube", "youtube"},
		{"go to github", "go to github", "github"},
		{"goto reddit", "goto reddit", "reddit"},
		{"fuzzy tube", "open tube", "youtube"},
		{"fuzzy yt", "open yt", "youtube"},
		{"fuzzy gh", "open gh", "github"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := p.Parse(tt.input)
			if result.Action != "navigate" {
				t.Errorf("action = %q, want %q", result.Action, "navigate")
			}
			if result.Target == nil || *result.Target != tt.wantTarget {
				got := ""
				if result.Target != nil {
					got = *result.Target
				}
				t.Errorf("target = %q, want %q", got, tt.wantTarget)
			}
		})
	}
}

func TestParser_NavigateRawURL(t *testing.T) {
	p := NewParser()

	result := p.Parse("open example.com")
	if result.Action != "navigate" {
		t.Errorf("action = %q, want navigate", result.Action)
	}
	if result.Target == nil || *result.Target != "https://example.com" {
		got := ""
		if result.Target != nil {
			got = *result.Target
		}
		t.Errorf("target = %q, want https://example.com", got)
	}
}

func TestParser_OpenAndSearch(t *testing.T) {
	p := NewParser()

	tests := []struct {
		name       string
		input      string
		wantTarget string
		wantQuery  string
	}{
		{"youtube search", "open youtube lofi beats", "youtube", "lofi beats"},
		{"spotify search", "open spotify jazz", "spotify", "jazz"},
		{"github search", "open github react", "github", "react"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := p.Parse(tt.input)
			if result.Action != "open-and-search" {
				t.Errorf("action = %q, want open-and-search", result.Action)
			}
			if result.Target == nil || *result.Target != tt.wantTarget {
				got := ""
				if result.Target != nil {
					got = *result.Target
				}
				t.Errorf("target = %q, want %q", got, tt.wantTarget)
			}
			if result.Query == nil || *result.Query != tt.wantQuery {
				got := ""
				if result.Query != nil {
					got = *result.Query
				}
				t.Errorf("query = %q, want %q", got, tt.wantQuery)
			}
		})
	}
}

func TestParser_OpenAndPlay(t *testing.T) {
	p := NewParser()

	tests := []struct {
		name       string
		input      string
		wantTarget string
		wantQuery  string
	}{
		{"youtube play", "open youtube play lofi hip hop", "youtube", "lofi hip hop"},
		{"spotify play", "open spotify play jazz", "spotify", "jazz"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := p.Parse(tt.input)
			if result.Action != "open-and-play" {
				t.Errorf("action = %q, want open-and-play", result.Action)
			}
			if result.Target == nil || *result.Target != tt.wantTarget {
				got := ""
				if result.Target != nil {
					got = *result.Target
				}
				t.Errorf("target = %q, want %q", got, tt.wantTarget)
			}
			if result.Query == nil || *result.Query != tt.wantQuery {
				got := ""
				if result.Query != nil {
					got = *result.Query
				}
				t.Errorf("query = %q, want %q", got, tt.wantQuery)
			}
		})
	}
}

func TestParser_BrowserControl(t *testing.T) {
	p := NewParser()

	tests := []struct {
		name      string
		input     string
		wantQuery string
	}{
		{"go back", "go back", "back"},
		{"back", "back", "back"},
		{"go forward", "go forward", "forward"},
		{"forward", "forward", "forward"},
		{"close tab", "close tab", "close-tab"},
		{"close this tab", "close this tab", "close-tab"},
		{"new tab", "new tab", "new-tab"},
		{"open new tab", "open new tab", "new-tab"},
		{"scroll down", "scroll down", "scroll-down"},
		{"scroll up", "scroll up", "scroll-up"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := p.Parse(tt.input)
			if result.Action != "browser-control" {
				t.Errorf("action = %q, want browser-control", result.Action)
			}
			if result.Query == nil || *result.Query != tt.wantQuery {
				got := ""
				if result.Query != nil {
					got = *result.Query
				}
				t.Errorf("query = %q, want %q", got, tt.wantQuery)
			}
			if result.Confidence < 0.9 {
				t.Errorf("confidence = %f, want >= 0.9", result.Confidence)
			}
		})
	}
}

func TestParser_CaseInsensitive(t *testing.T) {
	p := NewParser()

	result := p.Parse("OPEN YOUTUBE")
	if result.Action != "navigate" {
		t.Errorf("action = %q, want navigate", result.Action)
	}
	if result.Target == nil || *result.Target != "youtube" {
		t.Errorf("target = %v, want youtube", result.Target)
	}

	result2 := p.Parse("Search For React Hooks")
	if result2.Action != "search" {
		t.Errorf("action = %q, want search", result2.Action)
	}
	if result2.Query == nil || *result2.Query != "react hooks" {
		t.Errorf("query = %v, want 'react hooks'", result2.Query)
	}
}

func TestParser_ExecuteURL(t *testing.T) {
	p := NewParser()

	tests := []struct {
		name    string
		input   string
		wantURL string
	}{
		{"search builds google url", "search react hooks", "https://www.google.com/search?q=react+hooks"},
		{"navigate youtube", "open youtube", "https://www.youtube.com"},
		{"open-and-search youtube", "open youtube lofi", "https://www.youtube.com/results?search_query=lofi"},
		{"raw url", "open example.com", "https://example.com"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := p.Parse(tt.input)
			if result.ExecuteURL != tt.wantURL {
				t.Errorf("executeURL = %q, want %q", result.ExecuteURL, tt.wantURL)
			}
		})
	}
}
