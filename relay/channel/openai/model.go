package openai

import "encoding/json"

type VisionMessage struct {
	Role    string          `json:"role"`
	Content json.RawMessage `json:"content"`
	Name    *string         `json:"name,omitempty"`
}

type MediaMessage struct {
	Type     string `json:"type"`
	Text     string `json:"text"`
	ImageUrl any    `json:"image_url,omitempty"`
}

type Message struct {
	Role    string          `json:"role"`
	Content json.RawMessage `json:"content"`
	Name    *string         `json:"name,omitempty"`
}

type ImageURL struct {
	Url    string `json:"url,omitempty"`
	Detail string `json:"detail,omitempty"`
}

type TextContent struct {
	Type string `json:"type,omitempty"`
	Text string `json:"text,omitempty"`
}

type ImageContent struct {
	Type     string    `json:"type,omitempty"`
	ImageURL *ImageURL `json:"image_url,omitempty"`
}

type OpenAIMessageContent struct {
	Type     string    `json:"type,omitempty"`
	Text     string    `json:"text"`
	ImageURL *ImageURL `json:"image_url,omitempty"`
}

type MessageImageUrl struct {
	Url    string `json:"url"`
	Detail string `json:"detail"`
}

type MediaMessageImage struct {
	Type     string          `json:"type"`
	ImageUrl MessageImageUrl `json:"image_url"`
}

func (m Message) IsStringContent() bool {
	var content interface{}
	if err := json.Unmarshal(m.Content, &content); err != nil {
		return false
	}
	_, ok := content.(string)
	return ok
}

func (m Message) StringContent() string {
	var content interface{}
	if err := json.Unmarshal(m.Content, &content); err != nil {
		return ""
	}

	switch v := content.(type) {
	case string:
		return v
	case []interface{}:
		var contentStr string
		for _, contentItem := range v {
			contentMap, ok := contentItem.(map[string]interface{})
			if !ok {
				continue
			}
			if contentMap["type"] == ContentTypeText {
				if subStr, ok := contentMap["text"].(string); ok {
					contentStr += subStr
				}
			}
		}
		return contentStr
	default:
		return ""
	}
}

func (m Message) ParseContent() []MediaMessage {
	var contentList []MediaMessage
	var stringContent string
	if err := json.Unmarshal(m.Content, &stringContent); err == nil {
		contentList = append(contentList, MediaMessage{
			Type: ContentTypeText,
			Text: stringContent,
		})
		return contentList
	}
	var arrayContent []json.RawMessage
	if err := json.Unmarshal(m.Content, &arrayContent); err == nil {
		for _, contentItem := range arrayContent {
			var contentMap map[string]any
			if err := json.Unmarshal(contentItem, &contentMap); err != nil {
				continue
			}
			switch contentMap["type"] {
			case ContentTypeText:
				if subStr, ok := contentMap["text"].(string); ok {
					contentList = append(contentList, MediaMessage{
						Type: ContentTypeText,
						Text: subStr,
					})
				}
			case ContentTypeImageURL:
				if subObj, ok := contentMap["image_url"].(map[string]any); ok {
					detail, ok := subObj["detail"]
					if ok {
						subObj["detail"] = detail.(string)
					} else {
						subObj["detail"] = "auto"
					}
					contentList = append(contentList, MediaMessage{
						Type: ContentTypeImageURL,
						ImageUrl: MessageImageUrl{
							Url:    subObj["url"].(string),
							Detail: subObj["detail"].(string),
						},
					})
				}
			}
		}
		return contentList
	}

	return nil
}

type ResponseFormat struct {
	Type string `json:"type,omitempty"`
}

type GeneralOpenAIRequest struct {
	Model            string          `json:"model,omitempty"`
	Messages         []Message       `json:"messages,omitempty"`
	Prompt           any             `json:"prompt,omitempty"`
	Stream           bool            `json:"stream,omitempty"`
	MaxTokens        int             `json:"max_tokens,omitempty"`
	Temperature      float64         `json:"temperature,omitempty"`
	TopP             float64         `json:"top_p,omitempty"`
	N                int             `json:"n,omitempty"`
	Input            any             `json:"input,omitempty"`
	Instruction      string          `json:"instruction,omitempty"`
	Size             string          `json:"size,omitempty"`
	Functions        any             `json:"functions,omitempty"`
	FrequencyPenalty float64         `json:"frequency_penalty,omitempty"`
	PresencePenalty  float64         `json:"presence_penalty,omitempty"`
	ResponseFormat   *ResponseFormat `json:"response_format,omitempty"`
	Seed             float64         `json:"seed,omitempty"`
	Tools            any             `json:"tools,omitempty"`
	ToolChoice       any             `json:"tool_choice,omitempty"`
	User             string          `json:"user,omitempty"`
}

func (r GeneralOpenAIRequest) ParseInput() []string {
	if r.Input == nil {
		return nil
	}
	var input []string
	switch r.Input.(type) {
	case string:
		input = []string{r.Input.(string)}
	case []any:
		input = make([]string, 0, len(r.Input.([]any)))
		for _, item := range r.Input.([]any) {
			if str, ok := item.(string); ok {
				input = append(input, str)
			}
		}
	}
	return input
}

type ChatRequest struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	MaxTokens int       `json:"max_tokens"`
}

type TextRequest struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	Prompt    string    `json:"prompt"`
	MaxTokens int       `json:"max_tokens"`
	//Stream   bool      `json:"stream"`
}

// ImageRequest docs: https://platform.openai.com/docs/api-reference/images/create
type ImageRequest struct {
	Model          string `json:"model"`
	Prompt         string `json:"prompt" binding:"required"`
	N              int    `json:"n,omitempty"`
	Size           string `json:"size,omitempty"`
	Quality        string `json:"quality,omitempty"`
	ResponseFormat string `json:"response_format,omitempty"`
	Style          string `json:"style,omitempty"`
	User           string `json:"user,omitempty"`
}

type WhisperJSONResponse struct {
	Text string `json:"text,omitempty"`
}

type WhisperVerboseJSONResponse struct {
	Task     string    `json:"task,omitempty"`
	Language string    `json:"language,omitempty"`
	Duration float64   `json:"duration,omitempty"`
	Text     string    `json:"text,omitempty"`
	Segments []Segment `json:"segments,omitempty"`
}

type Segment struct {
	Id               int     `json:"id"`
	Seek             int     `json:"seek"`
	Start            float64 `json:"start"`
	End              float64 `json:"end"`
	Text             string  `json:"text"`
	Tokens           []int   `json:"tokens"`
	Temperature      float64 `json:"temperature"`
	AvgLogprob       float64 `json:"avg_logprob"`
	CompressionRatio float64 `json:"compression_ratio"`
	NoSpeechProb     float64 `json:"no_speech_prob"`
}

type TextToSpeechRequest struct {
	Model string `json:"model" binding:"required"`
	Input string `json:"input" binding:"required"`
	Voice string `json:"voice" binding:"required"`
}

type AudioResponse struct {
	Text string `json:"text,omitempty"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type Error struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Param   string `json:"param"`
	Code    any    `json:"code"`
}

type ErrorWithStatusCode struct {
	Error
	StatusCode int `json:"status_code"`
}

type SlimTextResponse struct {
	Choices []TextResponseChoice `json:"choices"`
	Usage   `json:"usage"`
	Error   Error `json:"error"`
}

type TextResponseChoice struct {
	Index        int `json:"index"`
	Message      `json:"message"`
	FinishReason string `json:"finish_reason"`
}

type TextResponse struct {
	Id      string               `json:"id"`
	Model   string               `json:"model,omitempty"`
	Object  string               `json:"object"`
	Created int64                `json:"created"`
	Choices []TextResponseChoice `json:"choices"`
	Usage   `json:"usage"`
}
type OpenAIError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Param   string `json:"param"`
	Code    any    `json:"code"`
}
type EmbeddingResponseItem struct {
	Object    string    `json:"object"`
	Index     int       `json:"index"`
	Embedding []float64 `json:"embedding"`
}

type EmbeddingResponse struct {
	Object string                  `json:"object"`
	Data   []EmbeddingResponseItem `json:"data"`
	Model  string                  `json:"model"`
	Usage  `json:"usage"`
}

type ImageResponse struct {
	Created int `json:"created"`
	Data    []struct {
		Url string `json:"url"`
	}
}

type ChatCompletionsStreamResponseChoice struct {
	Delta struct {
		Content string `json:"content"`
	} `json:"delta"`
	FinishReason *string `json:"finish_reason,omitempty"`
}

type ChatCompletionsStreamResponse struct {
	Id      string                                `json:"id"`
	Object  string                                `json:"object"`
	Created int64                                 `json:"created"`
	Model   string                                `json:"model"`
	Choices []ChatCompletionsStreamResponseChoice `json:"choices"`
}

type CompletionsStreamResponse struct {
	Choices []struct {
		Text         string `json:"text"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
}
