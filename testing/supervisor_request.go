package main

import (
	"encoding/json"

	"github.com/neospace-ai/go-openai"
)

func main() {
	request := openai.SupervisorRequest{
		Model: "gpt-4",
		History: []openai.ChatCompletionMessage{
			{
				Role:    "user",
				Content: "What is the capital of France?",
			},
			{
				Role:    "assistant",
				Content: "The capital of France is Paris.",
			},
		},
		SearchQuery: []string{
			"France capital city",
			"capital of France",
		},
		SearchAnswer: "Paris",
		Category:     "Geography",
		Description:  "A question about the capital city of France.",
		AvailableScores: map[string]string{
			"accuracy":  "High",
			"relevance": "Moderate",
		},
		MaxTokens:   100,
		Temperature: 0.7,
		TopP:        1,
	}

	x, err := json.Marshal(request)
	if err != nil {
		panic(err)
	}
	print(string(x))
}
