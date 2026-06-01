package intent

import (
	"regexp"
	"strings"
)

// Result represents a parsed intent.
type Result struct {
	Action     string
	Target     *string
	Query      *string
	Confidence float64
	ExecuteURL string
}

// Parser is the rule-based intent parser (server-side mirror of the extension parser).
type Parser struct {
	siteMap   map[string]string
	searchMap map[string]string // site → search URL template (%s = query)
	aliases   map[string]string
}

func NewParser() *Parser {
	return &Parser{
		siteMap: map[string]string{
			"youtube":       "https://www.youtube.com",
			"google":        "https://www.google.com",
			"github":        "https://www.github.com",
			"twitter":       "https://www.x.com",
			"x":             "https://www.x.com",
			"reddit":        "https://www.reddit.com",
			"spotify":       "https://open.spotify.com",
			"discord":       "https://discord.com",
			"gmail":         "https://mail.google.com",
			"netflix":       "https://www.netflix.com",
			"twitch":        "https://www.twitch.tv",
			"linkedin":      "https://www.linkedin.com",
			"stackoverflow": "https://stackoverflow.com",
		},
		searchMap: map[string]string{
			"youtube": "https://www.youtube.com/results?search_query=%s",
			"google":  "https://www.google.com/search?q=%s",
			"github":  "https://github.com/search?q=%s",
			"spotify": "https://open.spotify.com/search/%s",
			"reddit":  "https://www.reddit.com/search/?q=%s",
		},
		aliases: map[string]string{
			"tube": "youtube",
			"yt":   "youtube",
			"gh":   "github",
			"tw":   "twitter",
			"so":   "stackoverflow",
		},
	}
}

var (
	playPattern       = regexp.MustCompile(`^open\s+(\w+)\s+play\s+(.+)$`)
	openSearchPattern = regexp.MustCompile(`^open\s+(\w+)\s+(.+)$`)
	openPattern       = regexp.MustCompile(`^(?:open|go to|goto)\s+(.+)$`)
	searchPattern     = regexp.MustCompile(`^(?:search for|search|look up|find)\s+(.+)$`)
)

// Parse takes a transcript and returns the parsed intent.
func (p *Parser) Parse(raw string) Result {
	text := strings.ToLower(strings.TrimSpace(raw))

	// Browser control
	if r := p.parseBrowserControl(text); r != nil {
		return *r
	}

	// "open [site] play [query]"
	if m := playPattern.FindStringSubmatch(text); m != nil {
		site := p.findSite(m[1])
		if site != "" {
			query := m[2]
			return Result{
				Action:     "open-and-play",
				Target:     &site,
				Query:      &query,
				Confidence: 0.92,
				ExecuteURL: p.buildSearchURL(site, query),
			}
		}
	}

	// "open [site] [query]"
	if m := openSearchPattern.FindStringSubmatch(text); m != nil {
		site := p.findSite(m[1])
		if site != "" {
			query := m[2]
			return Result{
				Action:     "open-and-search",
				Target:     &site,
				Query:      &query,
				Confidence: 0.9,
				ExecuteURL: p.buildSearchURL(site, query),
			}
		}
	}

	// "open [site]" or "go to [site]"
	if m := openPattern.FindStringSubmatch(text); m != nil {
		input := strings.TrimSpace(m[1])
		// Check for URL (has dot)
		if strings.Contains(input, ".") {
			url := input
			if !strings.HasPrefix(url, "http") {
				url = "https://" + url
			}
			return Result{
				Action:     "navigate",
				Target:     &url,
				Confidence: 0.85,
				ExecuteURL: url,
			}
		}
		site := p.findSite(input)
		if site != "" {
			url := p.siteMap[site]
			return Result{
				Action:     "navigate",
				Target:     &site,
				Confidence: 0.92,
				ExecuteURL: url,
			}
		}
	}

	// "search [query]"
	if m := searchPattern.FindStringSubmatch(text); m != nil {
		query := m[1]
		url := strings.ReplaceAll(p.searchMap["google"], "%s", strings.ReplaceAll(query, " ", "+"))
		return Result{
			Action:     "search",
			Query:      &query,
			Confidence: 0.9,
			ExecuteURL: url,
		}
	}

	// Fallback: treat as search
	query := text
	url := strings.ReplaceAll(p.searchMap["google"], "%s", strings.ReplaceAll(query, " ", "+"))
	return Result{
		Action:     "search",
		Query:      &query,
		Confidence: 0.6,
		ExecuteURL: url,
	}
}

func (p *Parser) parseBrowserControl(text string) *Result {
	controls := map[string]string{
		"go back":        "back",
		"back":           "back",
		"go forward":     "forward",
		"forward":        "forward",
		"close tab":      "close-tab",
		"close this tab": "close-tab",
		"new tab":        "new-tab",
		"open new tab":   "new-tab",
	}

	for pattern, action := range controls {
		if text == pattern {
			return &Result{
				Action:     "browser-control",
				Query:      &action,
				Confidence: 0.95,
			}
		}
	}

	if strings.HasPrefix(text, "scroll down") {
		q := "scroll-down"
		return &Result{Action: "browser-control", Query: &q, Confidence: 0.9}
	}
	if strings.HasPrefix(text, "scroll up") {
		q := "scroll-up"
		return &Result{Action: "browser-control", Query: &q, Confidence: 0.9}
	}

	return nil
}

func (p *Parser) findSite(input string) string {
	clean := strings.ToLower(regexp.MustCompile(`[^a-z0-9]`).ReplaceAllString(input, ""))

	// Direct match
	if _, ok := p.siteMap[clean]; ok {
		return clean
	}

	// Alias match
	if alias, ok := p.aliases[clean]; ok {
		return alias
	}

	// Partial match
	for key := range p.siteMap {
		if strings.Contains(key, clean) || strings.Contains(clean, key) {
			return key
		}
	}

	return ""
}

func (p *Parser) buildSearchURL(site, query string) string {
	template, ok := p.searchMap[site]
	if !ok {
		// Fall back to site base URL
		if url, ok := p.siteMap[site]; ok {
			return url
		}
		return ""
	}
	return strings.ReplaceAll(template, "%s", strings.ReplaceAll(query, " ", "+"))
}
