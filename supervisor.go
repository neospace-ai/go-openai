package openai

import (
	"encoding/json"
	"fmt"
)

// SupervisorRequest represents a request structure for chat completion API.
type SupervisorRequest struct {
	Model           string                  `json:"model"`
	History         []ChatCompletionMessage `json:"history"`
	SearchQuery     []string                `json:"search_query"`
	SearchAnswer    string                  `json:"answer"`
	Category        string                  `json:"category"`
	Description     string                  `json:"description"`
	AvailableScores map[string]string       `json:"available_scores"`
	MaxTokens       int                     `json:"max_tokens"`
	Temperature     float64                 `json:"temperature"`
	TopP            int                     `json:"top_p"`
}

func (m SupervisorRequest) MarshalJSON() ([]byte, error) {
	type searchMessageContentQuery struct {
		Search []string `json:"search"`
		Answer string   `json:"answer"`
	}

	type searchMessageContent struct {
		TaskSelectExpertise searchMessageContentQuery `json:"task_select_expertise"`
	}

	type searchMessage struct {
		Role    string               `json:"role"`
		Content searchMessageContent `json:"content"`
	}

	type supervisorContext struct {
		Messages []any `json:"messages"`
	}

	type supervisorAnalysis struct {
		Category        string            `json:"category"`
		Description     string            `json:"description"`
		AvailableScores map[string]string `json:"available_scores"`
	}

	type promptFieldObject struct {
		Context  supervisorContext  `json:"supervisor_context"`
		Analysis supervisorAnalysis `json:"supervisor_analysis"`
	}

	type outObject struct {
		Model       string  `json:"model"`
		Prompt      string  `json:"prompt"`
		LogProbs    bool    `json:"logprobs"`
		TopLogProbs int     `json:"top_logprobs"`
		MaxTokens   int     `json:"max_tokens"`
		Temperature float64 `json:"temperature"`
		TopP        int     `json:"top_p"`
	}

	messages := make([]any, len(m.History)+1)
	for i, msg := range m.History {
		messages[i] = msg
	}

	messages[len(m.History)] = searchMessage{
		Role: "system",
		Content: searchMessageContent{
			TaskSelectExpertise: searchMessageContentQuery{
				Search: m.SearchQuery,
				Answer: m.SearchAnswer,
			},
		},
	}

	promptObj := promptFieldObject{
		Context: supervisorContext{
			Messages: messages,
		},
		Analysis: supervisorAnalysis{
			Category:        m.Category,
			Description:     m.Description,
			AvailableScores: m.AvailableScores,
		},
	}

	prompt, err := json.Marshal(promptObj)
	if err != nil {
		return nil, err
	}

	out := outObject{
		Model:       m.Model,
		Prompt:      string(prompt),
		MaxTokens:   m.MaxTokens,
		Temperature: m.Temperature,
		TopP:        m.TopP,
	}

	return json.Marshal(out)

}

func (m *SupervisorRequest) UnmarshalJSON(data []byte) error {
	type outObject struct {
		Model       string  `json:"model"`
		Prompt      string  `json:"prompt"`
		LogProbs    bool    `json:"logprobs"`
		TopLogProbs int     `json:"top_logprobs"`
		MaxTokens   int     `json:"max_tokens"`
		Temperature float64 `json:"temperature"`
		TopP        int     `json:"top_p"`
	}

	var out outObject
	if err := json.Unmarshal(data, &out); err != nil {
		return err
	}

	type searchMessageContentQuery struct {
		Search []string `json:"search"`
		Answer string   `json:"answer"`
	}

	type searchMessageContent struct {
		TaskSelectExpertise searchMessageContentQuery `json:"task_select_expertise"`
	}

	type searchMessage struct {
		Role    string               `json:"role"`
		Content searchMessageContent `json:"content"`
	}

	type supervisorContext struct {
		Messages []json.RawMessage `json:"messages"`
	}

	type supervisorAnalysis struct {
		Category        string            `json:"category"`
		Description     string            `json:"description"`
		AvailableScores map[string]string `json:"available_scores"`
	}

	type promptFieldObject struct {
		Context  supervisorContext  `json:"supervisor_context"`
		Analysis supervisorAnalysis `json:"supervisor_analysis"`
	}

	var promptObj promptFieldObject
	if err := json.Unmarshal([]byte(out.Prompt), &promptObj); err != nil {
		return err
	}

	messages := promptObj.Context.Messages

	if len(messages) == 0 {
		return fmt.Errorf("no messages found in supervisor_context")
	}

	decodedHistory := make([]ChatCompletionMessage, 0, len(messages)-1)
	for i := 0; i < len(messages)-1; i++ {
		var chatMsg ChatCompletionMessage
		if err := json.Unmarshal(messages[i], &chatMsg); err != nil {
			return fmt.Errorf("error unmarshalling history message: %w", err)
		}
		decodedHistory = append(decodedHistory, chatMsg)
	}

	var searchMsg searchMessage
	if err := json.Unmarshal(messages[len(messages)-1], &searchMsg); err != nil {
		return fmt.Errorf("error unmarshalling search message: %w", err)
	}

	m.Model = out.Model
	m.History = decodedHistory
	m.SearchQuery = searchMsg.Content.TaskSelectExpertise.Search
	m.SearchAnswer = searchMsg.Content.TaskSelectExpertise.Answer
	m.Category = promptObj.Analysis.Category
	m.Description = promptObj.Analysis.Description
	m.AvailableScores = promptObj.Analysis.AvailableScores
	m.MaxTokens = out.MaxTokens
	m.Temperature = out.Temperature
	m.TopP = out.TopP

	return nil
}

type SupervisorChoice struct {
	Index      int                 `json:"index"`
	LogProbs   *LogProbs           `json:"logprobs,omitempty"`
	TaskResult QuerySupervisorTask `json:"task_result"`
}

type SupervisorResponse struct {
	ID                string             `json:"id"`
	Object            string             `json:"object"`
	Created           int64              `json:"created"`
	Model             string             `json:"model"`
	Choices           []SupervisorChoice `json:"choices"`
	Usage             Usage              `json:"usage"`
	SystemFingerprint string             `json:"system_fingerprint"`

	httpHeader
}
