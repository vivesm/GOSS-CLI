package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Test streaming directly with LM Studio to verify the exact format
func main() {
	// Create a simple streaming request based on LM Studio docs
	reqBody := map[string]interface{}{
		"model": "openai/gpt-oss-20b",
		"messages": []map[string]string{
			{"role": "user", "content": "Say hello in one word"},
		},
		"stream":      true,
		"temperature": 0.3,
		"max_tokens":  10,
	}

	jsonBody, _ := json.Marshal(reqBody)
	
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", "http://localhost:1234/v1/chat/completions", bytes.NewReader(jsonBody))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Printf("Response Status: %s\n", resp.Status)
	fmt.Printf("Content-Type: %s\n", resp.Header.Get("Content-Type"))
	fmt.Println("--- Streaming Response ---")

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		
		// Debug: print raw line
		fmt.Printf("RAW: %q\n", line)
		
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, ":") {
			continue
		}
		
		// Check for end of stream
		if strings.TrimSpace(line) == "data: [DONE]" {
			fmt.Println("DONE: Stream completed")
			break
		}
		
		// Parse SSE data
		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")
			
			// Skip empty data
			if strings.TrimSpace(data) == "" {
				continue
			}
			
			// Try to parse JSON
			var chunk map[string]interface{}
			if err := json.Unmarshal([]byte(data), &chunk); err != nil {
				fmt.Printf("JSON Parse Error: %s, Data: %s\n", err, data)
				continue
			}
			
			// Extract content
			if choices, ok := chunk["choices"].([]interface{}); ok && len(choices) > 0 {
				if choice, ok := choices[0].(map[string]interface{}); ok {
					if delta, ok := choice["delta"].(map[string]interface{}); ok {
						// Handle reasoning (thinking)
						if reasoning, ok := delta["reasoning"].(string); ok {
							fmt.Printf("THINKING: %s", reasoning)
						}
						
						// Handle content
						if content, ok := delta["content"].(string); ok {
							fmt.Printf("CONTENT: %s", content)
						}
					}
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Scanner Error: %s\n", err)
	}
	
	fmt.Println("\n--- Test Complete ---")
}