package openai

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// SupervisorRequest represents a request structure for chat completion API.
type SupervisorRequest struct {
	Model        string                  `json:"model" bson:"model"`
	History      []ChatCompletionMessage `json:"history" bson:"history"`
	InstructTask GenericTask             `json:"instruct_task" bson:"instruct_task"`
	MaxTokens    int                     `json:"max_tokens" bson:"max_tokens"`
	Temperature  float64                 `json:"temperature" bson:"temperature"`
	TopP         int                     `json:"top_p" bson:"top_p"`
}

type SupervisorResponse struct {
	ID                string             `json:"id" bson:"id"`
	Object            string             `json:"object" bson:"object"`
	Created           int64              `json:"created" bson:"created"`
	Model             string             `json:"model" bson:"model"`
	Choices           []SupervisorChoice `json:"choices" bson:"choices"`
	Usage             Usage              `json:"usage" bson:"usage"`
	SystemFingerprint string             `json:"system_fingerprint" bson:"system_fingerprint"`

	httpHeader `json:"-" bson:"-"`
}

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

	type prompt struct {
		SupervisorContext supervisorContext `json:"supervisor_context" bson:"supervisor_context"`
	}

	type neolangInput struct {
		Model       string  `json:"model" bson:"model"`
		MaxTokens   int     `json:"max_tokens" bson:"max_tokens"`
		Temperature float64 `json:"temperature" bson:"temperature"`
		Prompt      string  `json:"prompt" bson:"prompt"`
	}

	messages := make([]map[string]any, len(req.History)+1)

	for idx, msg := range req.History {
		messages[idx] = map[string]any{
			"role":    msg.Role,
			"content": msg.Content,
		}
	}

	var task any
	task = req.InstructTask.Task

	if reflect.TypeOf(req.InstructTask.Task).Kind() == reflect.Ptr || reflect.TypeOf(req.InstructTask.Task).Kind() == reflect.Interface {
		task = reflect.ValueOf(req.InstructTask.Task).Elem().Interface()
	}

	messages[len(req.History)] = map[string]any{
		"role": ChatMessageRoleAssistant,
		"content": map[string]any{
			"task_guard": task,
		},
	}

	promptStr, err := json.Marshal(prompt{
		SupervisorContext: supervisorContext{
			Messages: messages,
		},
	})
	if err != nil {
		panic(fmt.Sprintf("failed to marshal prompt: %v", err))
	}

	input := neolangInput{
		Model:       req.Model,
		MaxTokens:   req.MaxTokens,
		Temperature: req.Temperature,
		Prompt:      string(promptStr),
	}
	return input
}
