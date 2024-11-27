package openai

import "strconv"

// SupervisorRequest represents a request structure for chat completion API.
type SupervisorRequest struct {
	Model        string                  `json:"model"`
	History      []ChatCompletionMessage `json:"history"`
	InstructTask GenericTask             `json:"instruct_task"`
	MaxTokens    int                     `json:"max_tokens"`
	Temperature  float64                 `json:"temperature"`
	TopP         int                     `json:"top_p"`
	Categories   SupervisorCategories    `json:"categories"`
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

func (req SupervisorRequest) ToNeolangInput() any {
	type supervisorContext struct {
		Messages []map[string]any `json:"messages"`
	}

	type component struct {
		Category        string            `json:"category"`
		Description     string            `json:"description"`
		AvailableScores map[string]string `json:"available_scores"`
	}

	type supervisorMechanics struct {
		Task       TaskDefinition `json:"task"`
		Components []component    `json:"components"`
	}

	type prompt struct {
		SupervisorContext   supervisorContext   `json:"supervisor_context"`
		SupervisorMechanics supervisorMechanics `json:"supervisor_mechanics"`
	}

	type neolangInput struct {
		Model       string  `json:"model"`
		MaxTokens   int     `json:"max_tokens"`
		Temperature float64 `json:"temperature"`
		Prompt      prompt  `json:"prompt"`
	}

	messages := make([]map[string]any, len(req.History)+1)

	for idx, msg := range req.History {
		messages[idx] = map[string]any{
			"role":    msg.Role,
			"content": msg.Content,
		}
	}

	messages[len(req.History)] = map[string]any{
		"role": ChatMessageRoleAssistant,
		"content": map[string]any{
			"task_guard": req.InstructTask.Task,
		},
	}

	components := make([]component, len(req.Categories))
	for catName, catDetails := range req.Categories {
		availableScores := make(map[string]string, len(catDetails.AvailableScores))
		for _, scoreDetails := range catDetails.AvailableScores {
			availableScores[strconv.Itoa(scoreDetails.Token)] = scoreDetails.Description
		}

		components = append(components, component{
			Category:        catName,
			Description:     catDetails.Decription,
			AvailableScores: availableScores,
		})
	}

	input := neolangInput{
		Model:       req.Model,
		MaxTokens:   req.MaxTokens,
		Temperature: req.Temperature,
		Prompt: prompt{
			SupervisorContext: supervisorContext{
				Messages: messages,
			},
			SupervisorMechanics: supervisorMechanics{
				Task:       GUARD_TASK_DEFINITION,
				Components: components,
			},
		},
	}
	return input
}
