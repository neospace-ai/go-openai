package openai

import (
	"context"
	"fmt"
	"math/rand/v2"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type MockInstructOptions string

const (
	TEST_STANDART     MockInstructOptions = "STANDART"
	TEST_GUARD_UNSAFE MockInstructOptions = "TEST_GUARD_UNSAFE"
	TEST_GUARD_SAFE   MockInstructOptions = "TEST_GUARD_SAFE"
)

// Returns a prederermined response for ChatCompletion API.
func (c *Client) MockChatCompletion(
	ctx context.Context,
	request ChatCompletionRequest,
	mockOption MockInstructOptions,
) (response ChatCompletionResponse, err error) {
	if request.Stream {
		err = ErrChatCompletionStreamNotSupported
		return
	}

	urlSuffix := chatCompletionsSuffix
	if !checkEndpointSupportsModel(urlSuffix, request.Model) {
		err = ErrChatCompletionInvalidModel
		return
	}

	req, err := c.newRequest(ctx, http.MethodPost, c.fullURL(urlSuffix, request.Model), withBody(request))
	if err != nil {
		return
	}
	print(fmt.Sprintf("Request sent to openai via mock: %+v", req))
	return getChatCompletionResponseMock(request, mockOption), nil
}

func getChatCompletionResponseMock(request ChatCompletionRequest, mockOption MockInstructOptions) ChatCompletionResponse {
	var guard *TaskGuard
	switch mockOption {
	case TEST_GUARD_UNSAFE:
		guard = &TaskGuard{
			Safe:            false,
			GuardReasoning:  "User is trying to test a guard activation implementation function, i should help him",
			GuardCategories: []string{"Racism", "Bullying", "Other"},
		}
	case TEST_GUARD_SAFE:
		guard = &TaskGuard{
			Safe:           true,
			GuardReasoning: "User is trying to test a guard safe implementation function, i should help him",
		}
	}

	return ChatCompletionResponse{
		ID:      uuid.NewString(),
		Object:  "text_completion",
		Created: time.Now().Unix(),
		Model:   request.Model,
		Choices: []ChatCompletionChoice{
			{
				Index:        0,
				FinishReason: "stop",
				Message: ChatCompletionMessage{
					Role:      "assistant",
					Content:   "This is a mocked response. Hope your tests are going well! :)",
					Reasoning: "The user is testing the implementation of the system. I should help him by providing a mocked response.",
				},
				TaskResults: TaskResultCollection{
					RawResponse: "<|mocked_response|> Lorem Ipsum <|eot|>",
					TaskGuard:   guard,
				},
			},
		},
		Usage: Usage{
			PromptTokens:     123456,
			CompletionTokens: 654321,
			TotalTokens:      123456 + 654321,
		},
		SystemFingerprint: "MockedSystemFingerprint",
	}
}

type MockSupervisorOptions string

const (
	TEST_SUPERVISOR_SELECT_BEST   MockSupervisorOptions = "TEST_SUPERVISOR_SELECT_BEST"
	TEST_SUPERVISOR_SELECT_WORST  MockSupervisorOptions = "TEST_SUPERVISOR_SELECT_WORST"
	TEST_SUPERVISOR_SELECT_RANDOM MockSupervisorOptions = "TEST_SUPERVISOR_SELECT_RANDOM"
)

func (c *Client) MockSupervisorCompletion(
	ctx context.Context,
	request SupervisorRequest,
	mockOptions MockSupervisorOptions,
) (response SupervisorResponse, err error) {
	urlSuffix := supervisorSuffix
	if !checkEndpointSupportsModel(urlSuffix, request.Model) {
		err = ErrChatCompletionInvalidModel
		return
	}

	req, err := c.newRequest(ctx, http.MethodPost, c.fullURL(urlSuffix, request.Model), withBody(request.ToNeolangInput()))
	if err != nil {
		return
	}

	print(fmt.Sprintf("Request sent to openai via supervisor mock: %+v", req))

	return getSupervisorResponseMock(request, mockOptions), nil
}

func getSupervisorResponseMock(request SupervisorRequest, mockOption MockSupervisorOptions) SupervisorResponse {
	var taskCat []SupervisorTaskCategory
	var supervisorReasoning string
	var supervisorFeedback string
	switch mockOption {
	case TEST_SUPERVISOR_SELECT_BEST:
		taskCat = getBestFromCategories(request.Categories)
		supervisorReasoning = "I am a mock that thinks the mock did a nice job"
	case TEST_SUPERVISOR_SELECT_WORST:
		taskCat = getWorstFromCategories(request.Categories)
		supervisorReasoning = "I am a mock that thinks the mock did a very bad job"
		supervisorFeedback = "The instruct should be better"
	case TEST_SUPERVISOR_SELECT_RANDOM:
		taskCat = getRandomFromCategories(request.Categories)
		supervisorReasoning = "I am mock and i dont know what i think"
		supervisorFeedback = "I cant give feedback, i dont know what i am doing"
	}
	return SupervisorResponse{
		ID:      uuid.NewString(),
		Object:  "text_completion",
		Created: time.Now().Unix(),
		Model:   request.Model,
		Choices: []SupervisorChoice{
			{
				Index:        0,
				TaskName:     request.InstructTask.TaskType,
				InstructTask: request.InstructTask.Task,
				Result: TaskSupervisor{
					RawResponse: "<|mocked_response|> supervisor mock <|eot|>",
					Reasoning:   supervisorReasoning,
					Feedback:    supervisorFeedback,
					Categories:  taskCat,
				},
			},
		},
		Usage: Usage{
			PromptTokens:     123456,
			CompletionTokens: 654321,
			TotalTokens:      123456 + 654321,
		},
		SystemFingerprint: "MockedSystemFingerprint",
	}
}

func getBestFromCategories(categories SupervisorCategories) []SupervisorTaskCategory {
	result := make([]SupervisorTaskCategory, len(categories))
	idx := 0
	for name, cat := range categories {
		scores := make([]SupervisorTaskScore, 0)
		maxScoreName := ""
		for scoreName, score := range cat.AvailableScores {
			scores = append(scores, SupervisorTaskScore{
				Token:       score.Token,
				Description: score.Description,
			})
			if maxScoreName == "" || score.Value > cat.AvailableScores[maxScoreName].Value {
				maxScoreName = scoreName
			}
		}
		chosenTok := cat.AvailableScores[maxScoreName].Token
		result[idx] = SupervisorTaskCategory{
			Name:            name,
			Decription:      cat.Decription,
			AvailableScores: scores,
			Chosen:          &chosenTok,
		}
		idx++
	}
	return result
}

func getWorstFromCategories(categories SupervisorCategories) []SupervisorTaskCategory {
	result := make([]SupervisorTaskCategory, len(categories))
	idx := 0
	for name, cat := range categories {
		scores := make([]SupervisorTaskScore, 0)
		maxScoreName := ""
		for scoreName, score := range cat.AvailableScores {
			scores = append(scores, SupervisorTaskScore{
				Token:       score.Token,
				Description: score.Description,
			})
			if maxScoreName == "" || score.Value < cat.AvailableScores[maxScoreName].Value {
				maxScoreName = scoreName
			}
		}
		chosenTok := cat.AvailableScores[maxScoreName].Token
		result[idx] = SupervisorTaskCategory{
			Name:            name,
			Decription:      cat.Decription,
			AvailableScores: scores,
			Chosen:          &chosenTok,
		}
		idx++
	}
	return result
}

func getRandomFromCategories(categories SupervisorCategories) []SupervisorTaskCategory {
	result := make([]SupervisorTaskCategory, len(categories))
	idx := 0
	for name, cat := range categories {
		scores := make([]SupervisorTaskScore, 0)
		for _, score := range cat.AvailableScores {
			scores = append(scores, SupervisorTaskScore{
				Token:       score.Token,
				Description: score.Description,
			})
		}
		chosenTok := scores[rand.Int()%len(scores)].Token
		result[idx] = SupervisorTaskCategory{
			Name:            name,
			Decription:      cat.Decription,
			AvailableScores: scores,
			Chosen:          &chosenTok,
		}
		idx++
	}
	return result
}