package revchatgpt2

import "fmt"

/* OpenAI request body format defination */

type PromptContent struct {
	ContentType string   `json:"content_type"`
	Parts       []string `json:"parts"`
}

type PromptAuthor struct {
	Role string `json:"role"` // 'user' | 'assistant' | 'system'
}

type Prompt struct {
	Id      string        `json:"id"`
	Author  PromptAuthor  `json:"author"`
	Content PromptContent `json:"content"`
}

type ConversationJSONBody struct {
	Action            string   `json:"action"`
	Messages          []Prompt `json:"messages"`
	ParentMessageId   string   `json:"parent_message_id"`
	Model             string   `json:"model"`
	ConverationId     *string  `json:"conversation_id,omitempty"`
	TimezoneOffsetMin int      `json:"timezone_offset_min"` // always 0 ?
}

func (b ConversationJSONBody) String() string {
	if b.ConverationId == nil {
		return fmt.Sprintf("Model: '%s' ParentMessageId: '%s'",
			b.Model, b.ParentMessageId)
	} else {
		return fmt.Sprintf("Model: '%s' ParentMessageId: '%s' ConversationId: '%s'",
			b.Model, b.ParentMessageId, *b.ConverationId)
	}
}

/* OpenAI response body format defination */

// stream response body
type MessageContent struct {
	ContentType string   `json:"content_type"`
	Parts       []string `json:"parts"`
}

type Author struct {
	Role     string      `json:"role"`
	Name     string      `json:"name"`
	Metadata interface{} `json:"metadata"`
}

type Metadata struct {
	MessageType   string       `json:"message_type"`
	ModelSlug     string       `json:"model_slug"`
	FinishDetails *interface{} `json:"finish_details"`
	Recipient     *string      `json:"recipient"`
}

type Message struct {
	Id         string         `json:"id"`
	Author     Author         `json:"author"`
	Content    MessageContent `json:"content"`
	CreateTime float32        `json:"create_time"`
	UpdateTime float32        `json:"update_time"`
	EndTurn    bool           `json:"end_turn"`
	Weight     float32        `json:"weight"`
	Recipient  string         `json:"recipient"`
	Metadata   Metadata       `json:"metadata"`
}

type ConversationResponseEvent = struct {
	Message        Message   `json:"message"`
	ConversationId string    `json:"conversation_id"`
	Error          *APIError `json:"error"`
}

// error response body
type APIError struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message"`
	Param   string `json:"param,omitempty"`
	Type    string `json:"type"`
}

func (e *APIError) String() string {
	if e.Type == "invalid_request_error" {
		return fmt.Sprintf("%s (%s)", e.Code, e.Type)
	} else {
		return fmt.Sprintf("%s: %s (%s)", e.Message, e.Code, e.Type)
	}
}
