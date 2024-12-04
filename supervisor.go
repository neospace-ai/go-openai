package openai

import "strconv"

// SupervisorRequest represents a request structure for chat completion API.
type SupervisorRequest struct {
	Model        string                  `json:"model" bson:"model"`
	History      []ChatCompletionMessage `json:"history" bson:"history"`
	InstructTask GenericTask             `json:"instruct_task" bson:"instruct_task"`
	MaxTokens    int                     `json:"max_tokens" bson:"max_tokens"`
	Temperature  float64                 `json:"temperature" bson:"temperature"`
	TopP         int                     `json:"top_p" bson:"top_p"`
	Categories   SupervisorCategories    `json:"categories" bson:"categories"`
}

type SupervisorResponse struct {
	ID                string             `json:"id" bson:"id"`
	Object            string             `json:"object" bson:"object"`
	Created           int64              `json:"created" bson:"created"`
	Model             string             `json:"model" bson:"model"`
	Choices           []SupervisorChoice `json:"choices" bson:"choices"`
	Usage             Usage              `json:"usage" bson:"usage"`
	SystemFingerprint string             `json:"system_fingerprint" bson:"system_fingerprint"`

	httpHeader
}

type SupervisorScore struct {
	Token       int    `json:"token" bson:"token"`
	Description string `json:"description" bson:"description"`
	Value       int    `json:"value" bson:"value"`
}

// Map in format: {Name: Details}
type SupervisorAvailableScores map[string]SupervisorScore

type SupervisorCategory struct {
	Decription      string                    `json:"description" bson:"description"`
	AvailableScores SupervisorAvailableScores `json:"available_scores" bson:"available_scores"`
	Threshold       int                       `json:"threshold" bson:"threshold"`
}

// Map in format: {Name: Details}
type SupervisorCategories map[string]SupervisorCategory

type SupervisorChoice struct {
	Index        int            `json:"index" bson:"index"`
	LogProbs     *LogProbs      `json:"logprobs,omitempty" bson:"logprobs,omitempty"`
	TaskName     string         `json:"task" bson:"task"`
	InstructTask any            `json:"instruct_task" bson:"instruct_task"`
	Result       TaskSupervisor `json:"result" bson:"result"`
}

func (req SupervisorRequest) ToNeolangInput() any {
	type supervisorContext struct {
		Messages []map[string]any `json:"messages" bson:"messages"`
	}

	type component struct {
		Category        string            `json:"category" bson:"category"`
		Description     string            `json:"description" bson:"description"`
		AvailableScores map[string]string `json:"available_scores" bson:"available_scores"`
	}

	type supervisorMechanics struct {
		Task       TaskDefinition `json:"task" bson:"task"`
		Components []component    `json:"components" bson:"components"`
	}

	type prompt struct {
		SupervisorContext   supervisorContext   `json:"supervisor_context" bson:"supervisor_context"`
		SupervisorMechanics supervisorMechanics `json:"supervisor_mechanics" bson:"supervisor_mechanics"`
	}

	type neolangInput struct {
		Model       string  `json:"model" bson:"model"`
		MaxTokens   int     `json:"max_tokens" bson:"max_tokens"`
		Temperature float64 `json:"temperature" bson:"temperature"`
		Prompt      prompt  `json:"prompt" bson:"prompt"`
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
