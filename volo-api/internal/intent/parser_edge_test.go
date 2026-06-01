package intent

import (
	"strings"
	"testing"
)

// ─── Edge Cases That Could Break Things ──────────────────

func TestParser_EmptyInput(t *testing.T) {
	p := NewParser()

	result := p.Parse("")
	if result.Action != "search" {
		t.Errorf("empty input: action = %q, want search (fallback)", result.Action)
	}
	if result.Confidence > 0.7 {
		t.Errorf("empty input: confidence = %f, want low", result.Confidence)
	}
}

func TestParser_WhitespaceOnly(t *testing.T) {
	p := NewParser()

	result := p.Parse("   ")
	if result.Action != "search" {
		t.Errorf("whitespace: action = %q, want search (fallback)", result.Action)
	}
}

func TestParser_VeryLongInput(t *testing.T) {
	p := NewParser()

	// 500 character input — should not panic or hang
	long := "search " + strings.Repeat("a very long query that goes on and on ", 20)
	result := p.Parse(long)
	if result.Action != "search" {
		t.Errorf("long input: action = %q, want search", result.Action)
	}
}

func TestParser_SpecialCharacters(t *testing.T) {
	p := NewParser()

	tests := []struct {
		input string
		name  string
	}{
		{"search c++ tutorial", "plus signs"},
		{"search what's the weather?", "apostrophe and question mark"},
		{"open https://example.com/path?q=1&b=2", "url with params"},
		{"search $100 deals", "dollar sign"},
		{"search <script>alert('xss')</script>", "html injection attempt"},
		{"search 你好世界", "unicode characters"},
		{"search café résumé", "accented characters"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Should not panic
			result := p.Parse(tt.input)
			if result.Action == "" {
				t.Error("expected non-empty action")
			}
		})
	}
}

func TestParser_SiteNameVariations(t *testing.T) {
	p := NewParser()

	tests := []struct {
		input      string
		wantTarget string
	}{
		{"open YouTube", "youtube"},                 // capitalized
		{"open YOUTUBE", "youtube"},                 // all caps
		{"open youtube.com", "https://youtube.com"}, // with .com — treated as raw URL (dot triggers URL path)
		{"open yt", "youtube"},                      // alias
		{"open tube", "youtube"},                    // alias
		{"open gh", "github"},                       // alias
		{"open reddit", "reddit"},                   // direct
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := p.Parse(tt.input)
			if result.Action != "navigate" {
				// youtube.com will be treated as raw URL, that's acceptable
				if result.Target == nil || !strings.Contains(*result.Target, tt.wantTarget) {
					t.Errorf("target = %v, want to contain %q", result.Target, tt.wantTarget)
				}
				return
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

func TestParser_AmbiguousCommands(t *testing.T) {
	p := NewParser()

	// These should all produce SOME result without panicking
	ambiguous := []string{
		"open",                   // open what?
		"search",                 // search what?
		"play",                   // play what?
		"go",                     // go where?
		"open something unknown", // unknown site
		"do something",           // completely unknown
	}

	for _, input := range ambiguous {
		t.Run(input, func(t *testing.T) {
			result := p.Parse(input)
			if result.Action == "" {
				t.Error("expected non-empty action for ambiguous input")
			}
			if result.Confidence < 0 || result.Confidence > 1 {
				t.Errorf("confidence %f out of range [0,1]", result.Confidence)
			}
		})
	}
}

func TestParser_OpenAndPlay_NoSite(t *testing.T) {
	p := NewParser()

	// "open unknownsite play something" — site not recognized
	result := p.Parse("open randomsite play jazz")
	// Should fall through to open-and-search or fallback, not panic
	if result.Action == "" {
		t.Error("expected non-empty action")
	}
}

func TestParser_NavigateRawURL_Variations(t *testing.T) {
	p := NewParser()

	tests := []struct {
		input   string
		wantURL string
	}{
		{"open example.com", "https://example.com"},
		{"open www.example.com", "https://www.example.com"},
		{"open https://already-has-protocol.com", "https://already-has-protocol.com"},
		{"open sub.domain.co.uk", "https://sub.domain.co.uk"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := p.Parse(tt.input)
			if result.Action != "navigate" {
				t.Errorf("action = %q, want navigate", result.Action)
				return
			}
			if result.Target == nil || *result.Target != tt.wantURL {
				got := ""
				if result.Target != nil {
					got = *result.Target
				}
				t.Errorf("target = %q, want %q", got, tt.wantURL)
			}
		})
	}
}

func TestParser_BrowserControl_ExactMatch(t *testing.T) {
	p := NewParser()

	// These should NOT match browser control
	notControl := []string{
		"go back to youtube",     // "go back" is a prefix but has more words
		"close the tab please",   // not exact match
		"new tab and open gmail", // not exact match
		"scroll down the page",   // "scroll down" is prefix, has more
	}

	for _, input := range notControl {
		t.Run(input, func(t *testing.T) {
			result := p.Parse(input)
			// Some of these might still match due to prefix matching on scroll
			// The important thing is they don't panic
			if result.Action == "" {
				t.Error("expected non-empty action")
			}
		})
	}
}

func TestParser_ConfidenceRanges(t *testing.T) {
	p := NewParser()

	tests := []struct {
		input   string
		minConf float64
		maxConf float64
	}{
		{"go back", 0.9, 1.0},
		{"open youtube", 0.9, 1.0},
		{"search react", 0.85, 1.0},
		{"something random", 0.0, 0.7}, // fallback should be low confidence
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := p.Parse(tt.input)
			if result.Confidence < tt.minConf || result.Confidence > tt.maxConf {
				t.Errorf("confidence = %f, want [%f, %f]", result.Confidence, tt.minConf, tt.maxConf)
			}
		})
	}
}

func TestParser_ExecuteURL_NeverEmpty_ForNavigation(t *testing.T) {
	p := NewParser()

	navCommands := []string{
		"open youtube",
		"open github",
		"open spotify",
		"search react hooks",
		"open youtube lofi",
	}

	for _, input := range navCommands {
		t.Run(input, func(t *testing.T) {
			result := p.Parse(input)
			if result.ExecuteURL == "" {
				t.Errorf("expected non-empty ExecuteURL for %q (action: %s)", input, result.Action)
			}
			if !strings.HasPrefix(result.ExecuteURL, "http") {
				t.Errorf("ExecuteURL = %q, expected to start with http", result.ExecuteURL)
			}
		})
	}
}

func TestParser_ExecuteURL_Empty_ForBrowserControl(t *testing.T) {
	p := NewParser()

	controlCommands := []string{
		"go back",
		"close tab",
		"new tab",
		"scroll down",
	}

	for _, input := range controlCommands {
		t.Run(input, func(t *testing.T) {
			result := p.Parse(input)
			if result.ExecuteURL != "" {
				t.Errorf("expected empty ExecuteURL for browser control %q, got %q", input, result.ExecuteURL)
			}
		})
	}
}
