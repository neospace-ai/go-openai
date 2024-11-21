package openai

// SupervisorRequest represents a request structure for chat completion API.
type SupervisorRequest struct {
	Model        string                  `json:"model"`
	History      []ChatCompletionMessage `json:"history"`
	InstructTask TaskResultCollection    `json:"instruct_task"`
	MaxTokens    int                     `json:"max_tokens"`
	Temperature  float64                 `json:"temperature"`
	TopP         int                     `json:"top_p"`
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

type SupervisorScore struct {
	Index        int    `json:"index"`
	SpecialToken int    `json:"special_token"`
	Name         string `json:"name"`
	Description  string `json:"description"`
}

type SupervisorCategory struct {
	Category        string            `json:"category"`
	Decription      string            `json:"description"`
	AvailableScores []SupervisorScore `json:"available_scores"`
}

type SupervisorChoice struct {
	Index        int                  `json:"index"`
	LogProbs     *LogProbs            `json:"logprobs,omitempty"`
	RawResponse  string               `json:"raw_response"`
	Task         string               `json:"task"`
	InstructTask any                  `json:"instruct_task"`
	Categories   []SupervisorCategory `json:"categories"`
	Reasoning    string               `json:"reasoning"`
	Score        []string             `json:"score"`
	Feedback     string               `json:"feedback"`
}
