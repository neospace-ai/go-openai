package openai

import (
	"context"
	"net/http"
)

const (
	SPECIAL_TOKEN_BEGIN_TEXT     = "<|begin_of_text|>"
	SPECIAL_TOKEN_START_HEADER   = "<|start_header_id|>"
	SPECIAL_TOKEN_END_HEADER     = "<|end_header_id|>"
	SPECIAL_TOKEN_EOT            = "<|eot_id|>"
	SPECIAL_TOKEN_START_ANALYSIS = "<|start_analysis|>"
	SPECIAL_TOKEN_END_ANALYSIS   = "<|end_analysis|>"
	SPECIAL_TOKEN_TOOLS_REQUEST  = "<|tools_request|>"
	SPECIAL_TOKEN_GUARD_RAIL     = "<|guard|>"
)

type ChatCompletionStreamChoiceDelta struct {
	Content      string                   `json:"content,omitempty"`
	Role         string                   `json:"role,omitempty"`
	FunctionCall *FunctionCall            `json:"function_call,omitempty"`
	ToolCalls    []ToolCall               `json:"tool_calls,omitempty"`
	Analysis     *ChatCompletionAnalysis  `json:"Analysis,omitempty"`
	GuardRail    *ChatCompletionGuardRail `json:"guard,omitempty"`
}

type ChatCompletionStreamChoice struct {
	Index                int                             `json:"index"`
	Delta                ChatCompletionStreamChoiceDelta `json:"delta"`
	FinishReason         FinishReason                    `json:"finish_reason"`
	ContentFilterResults ContentFilterResults            `json:"content_filter_results,omitempty"`
}

type PromptFilterResult struct {
	Index                int                  `json:"index"`
	ContentFilterResults ContentFilterResults `json:"content_filter_results,omitempty"`
}

type ChatCompletionStreamResponse struct {
	ID                  string                       `json:"id"`
	Object              string                       `json:"object"`
	Created             int64                        `json:"created"`
	Model               string                       `json:"model"`
	Choices             []ChatCompletionStreamChoice `json:"choices"`
	SystemFingerprint   string                       `json:"system_fingerprint"`
	PromptAnnotations   []PromptAnnotation           `json:"prompt_annotations,omitempty"`
	PromptFilterResults []PromptFilterResult         `json:"prompt_filter_results,omitempty"`
	// An optional field that will only be present when you set stream_options: {"include_usage": true} in your request.
	// When present, it contains a null value except for the last chunk which contains the token usage statistics
	// for the entire request.
	Usage *Usage `json:"usage,omitempty"`
}

// ChatCompletionStream
// Note: Perhaps it is more elegant to abstract Stream using generics.
type ChatCompletionStream struct {
	*streamReader[ChatCompletionStreamResponse]
}

// CreateChatCompletionStream — API call to create a chat completion w/ streaming
// support. It sets whether to stream back partial progress. If set, tokens will be
// sent as data-only server-sent events as they become available, with the
// stream terminated by a data: [DONE] message.
func (c *Client) CreateChatCompletionStream(
	ctx context.Context,
	request ChatCompletionRequest,
) (stream *ChatCompletionStream, err error) {
	urlSuffix := chatCompletionsSuffix
	if !checkEndpointSupportsModel(urlSuffix, request.Model) {
		err = ErrChatCompletionInvalidModel
		return
	}

	request.Stream = true
	req, err := c.newRequest(ctx, http.MethodPost, c.fullURL(urlSuffix, request.Model), withBody(request))
	if err != nil {
		return nil, err
	}

	resp, err := sendRequestStream[ChatCompletionStreamResponse](c, req)
	if err != nil {
		return
	}
	stream = &ChatCompletionStream{
		streamReader: resp,
	}
	return
}
