package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/vivesm/GOSS-CLI/agentic-cli/openai"
)

// RateLimiter implements a simple token bucket rate limiter for web searches
type RateLimiter struct {
	tokens    int
	maxTokens int
	refillRate time.Duration
	lastRefill time.Time
	mu         sync.Mutex
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(maxTokens int, refillRate time.Duration) *RateLimiter {
	return &RateLimiter{
		tokens:     maxTokens,
		maxTokens:  maxTokens,
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
}

// Allow checks if a request can proceed
func (rl *RateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(rl.lastRefill)
	
	// Refill tokens based on elapsed time
	if elapsed >= rl.refillRate {
		tokensToAdd := int(elapsed / rl.refillRate)
		rl.tokens = min(rl.maxTokens, rl.tokens+tokensToAdd)
		rl.lastRefill = now
	}

	if rl.tokens > 0 {
		rl.tokens--
		return true
	}

	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Global rate limiter for web searches (5 requests per minute)
var webSearchRateLimiter = NewRateLimiter(5, time.Minute/5)

// CreateWebSearchTools returns web search MCP tools
func CreateWebSearchTools() []openai.Tool {
	return []openai.Tool{
		{
			Type: "function",
			Function: openai.ToolFunction{
				Name:        "web_search",
				Description: "Search the web using Brave Search API",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"query": map[string]interface{}{
							"type":        "string",
							"description": "Search query string",
						},
						"count": map[string]interface{}{
							"type":        "integer",
							"description": "Number of search results to return (default: 5, max: 20)",
							"minimum":     1,
							"maximum":     20,
						},
					},
					"required": []string{"query"},
				},
				Handler: webSearchHandler,
			},
		},
	}
}

// BraveSearchResponse represents the Brave Search API response
type BraveSearchResponse struct {
	Web struct {
		Type    string `json:"type"`
		Results []struct {
			Title       string `json:"title"`
			URL         string `json:"url"`
			Description string `json:"description"`
			Age         string `json:"age,omitempty"`
		} `json:"results"`
	} `json:"web"`
	Query struct {
		Original string      `json:"original"`
		Show     interface{} `json:"show_strict_warning"`
		Altered  string      `json:"altered,omitempty"`
	} `json:"query"`
}

func webSearchHandler(ctx context.Context, args map[string]interface{}) (string, error) {
	// Rate limiting check
	if !webSearchRateLimiter.Allow() {
		return "", fmt.Errorf("rate limit exceeded: maximum 5 web searches per minute allowed")
	}

	query, ok := args["query"].(string)
	if !ok {
		return "", fmt.Errorf("query must be a string")
	}

	// Input validation
	if len(strings.TrimSpace(query)) == 0 {
		return "", fmt.Errorf("search query cannot be empty")
	}

	if len(query) > 1000 {
		return "", fmt.Errorf("search query too long (max 1000 characters)")
	}

	count := 5 // default
	if c, ok := args["count"].(float64); ok {
		count = int(c)
		if count > 20 {
			count = 20
		}
		if count < 1 {
			count = 1
		}
	}

	// For now, we'll use a simple HTTP search instead of Brave API
	// In production, you'd want to use proper Brave Search API with API key
	results, err := performWebSearch(ctx, query, count)
	if err != nil {
		return "", fmt.Errorf("web search failed: %w", err)
	}

	return formatSearchResults(query, results), nil
}

func performWebSearch(ctx context.Context, query string, count int) ([]SearchResult, error) {
	// Try to load Brave API key from .env.brave.api file
	apiKey := loadBraveAPIKey()

	client := &http.Client{Timeout: 10 * time.Second}

	if apiKey != "" {
		// Use Brave Search API if we have a key
		apiURL := fmt.Sprintf("https://api.search.brave.com/res/v1/web/search?q=%s&count=%d",
			url.QueryEscape(query), count)

		req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
		if err != nil {
			return nil, err
		}

		req.Header.Set("Accept", "application/json")
		req.Header.Set("X-Subscription-Token", apiKey)
		req.Header.Set("User-Agent", "goss-cli/1.0")

		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			// Fallback to DuckDuckGo if Brave API fails
			return performDuckDuckGoSearch(ctx, query, count)
		}

		var braveResp BraveSearchResponse
		if err := json.NewDecoder(resp.Body).Decode(&braveResp); err != nil {
			return nil, err
		}

		var results []SearchResult
		for i, r := range braveResp.Web.Results {
			if i >= count {
				break
			}
			results = append(results, SearchResult{
				Title:       r.Title,
				URL:         r.URL,
				Description: r.Description,
			})
		}
		return results, nil
	}

	// Fallback to DuckDuckGo if no API key
	return performDuckDuckGoSearch(ctx, query, count)
}

func performDuckDuckGoSearch(ctx context.Context, query string, count int) ([]SearchResult, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	// Use DuckDuckGo API as fallback
	apiURL := fmt.Sprintf("https://api.duckduckgo.com/?q=%s&format=json&no_html=1&skip_disambig=1",
		url.QueryEscape(query))

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "goss-cli/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ddgResponse DuckDuckGoResponse
	if err := json.NewDecoder(resp.Body).Decode(&ddgResponse); err != nil {
		return nil, err
	}

	var results []SearchResult

	// Add instant answer if available
	if ddgResponse.AbstractText != "" {
		results = append(results, SearchResult{
			Title:       ddgResponse.Heading,
			URL:         ddgResponse.AbstractURL,
			Description: ddgResponse.AbstractText,
		})
	}

	// Add related topics
	for i, topic := range ddgResponse.RelatedTopics {
		if i >= count-1 { // Save space for instant answer
			break
		}
		if topic.Text != "" {
			results = append(results, SearchResult{
				Title:       extractTitle(topic.Text),
				URL:         topic.FirstURL,
				Description: topic.Text,
			})
		}
	}

	// If we don't have enough results, add some mock results
	// In production, you'd implement proper search API integration
	if len(results) == 0 {
		results = []SearchResult{
			{
				Title:       "Search Results",
				URL:         "https://duckduckgo.com/?q=" + url.QueryEscape(query),
				Description: fmt.Sprintf("No instant results found for '%s'. Try searching directly on the web.", query),
			},
		}
	}

	return results, nil
}

type SearchResult struct {
	Title       string
	URL         string
	Description string
}

type DuckDuckGoResponse struct {
	AbstractText  string         `json:"AbstractText"`
	AbstractURL   string         `json:"AbstractURL"`
	Heading       string         `json:"Heading"`
	RelatedTopics []RelatedTopic `json:"RelatedTopics"`
}

type RelatedTopic struct {
	FirstURL string `json:"FirstURL"`
	Text     string `json:"Text"`
}

func extractTitle(text string) string {
	parts := strings.SplitN(text, " - ", 2)
	if len(parts) > 1 {
		return parts[0]
	}
	// Take first 60 characters as title
	if len(text) > 60 {
		return text[:57] + "..."
	}
	return text
}

func loadBraveAPIKey() string {
	// First try environment variable
	if key := os.Getenv("BRAVE_API_KEY"); key != "" {
		return strings.TrimSpace(key)
	}

	// Then try .env.brave.api file
	data, err := os.ReadFile(".env.brave.api")
	if err == nil {
		return strings.TrimSpace(string(data))
	}

	// Try from home directory
	homeDir, err := os.UserHomeDir()
	if err == nil {
		data, err = os.ReadFile(homeDir + "/.env.brave.api")
		if err == nil {
			return strings.TrimSpace(string(data))
		}
	}

	return ""
}

func formatSearchResults(query string, results []SearchResult) string {
	var output strings.Builder
	output.WriteString(fmt.Sprintf("Web search results for: %s\n\n", query))

	if len(results) == 0 {
		output.WriteString("No results found. Try a different search query.\n")
		return output.String()
	}

	for i, result := range results {
		output.WriteString(fmt.Sprintf("%d. %s\n", i+1, result.Title))
		output.WriteString(fmt.Sprintf("   URL: %s\n", result.URL))
		output.WriteString(fmt.Sprintf("   %s\n\n", result.Description))
	}

	return output.String()
}
