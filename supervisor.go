package openai

// SupervisorRequest represents a request structure for chat completion API.
type SupervisorRequest struct {
	Model        string                  `json:"model"`
	History      []ChatCompletionMessage `json:"history"`
	InstructTask GenericTask             `json:"instruct_task"`
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
	Token       int    `json:"token"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Value       int    `json:"value"`
}

// Map in format: {Name: Details}
type SupervisorAvailableScores map[string]SupervisorScore

type SupervisorCategory struct {
	Decription      string                    `json:"description"`
	AvailableScores SupervisorAvailableScores `json:"available_scores"`
	Threshold       int                       `json:"threshold"`
}

// Map in format: {Name: Details}
type SupervisorCategories map[string]SupervisorCategory

type SupervisorChoice struct {
	Index        int            `json:"index"`
	LogProbs     *LogProbs      `json:"logprobs,omitempty"`
	TaskName     string         `json:"task"`
	InstructTask any            `json:"instruct_task"`
	Result       TaskSupervisor `json:"result"`
}
