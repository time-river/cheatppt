package revchatgpt

import (
	"io"
)

type ChatGPTUnofficialProxyAPI struct {
	apiReverseProxyUrl string
	Headers            map[string]string
	model              string
	logger             io.Writer
	debugger           bool
}

type ChatGPTError struct {
	StatusCode int
	StatusText string
	IsFinal    bool
	AccountId  string
}

var _ = error.Error(&ChatGPTError{})

/* Called by `SendMessage()` */

type PartialResponser func(partialResponse ChatMessage)

type SendMessageBrowserOptions struct {
	ConversationId  string
	ParentMessageId string
	MessageId       string
	Action          string /* 'next' | 'variant' */
	TimeoutMs       uint64
	OnProgress      PartialResponser
}

type ChatMessage struct {
	Id     string      `json:"id"`
	Text   string      `json:"text"`
	Role   string      `json:"role"`
	Name   string      `json:"name,omitempty"`
	Delta  string      `json:"delta,omitempty"`
	Detail interface{} `json:"detail,omitempty"`
	// relevant for both ChatGPTAPI and ChatGPTUnofficialProxyAPI
	ParentMessageId string `json:"parentMessageId,omitempty"`
	// only relevant for ChatGPTUnofficialProxyAPI (optional for ChatGPTAPI)
	ConversationId string `json:"conversationId,omitempty"`
}

/* OpenAI request body format defination */

type PromptContent struct {
	ContentType string   `json:"content_type"`
	Parts       []string `json:"parts"`
}

type PromptAuthor struct {
	Role string `json:"role"`
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
	ConverationId     string   `json:"conversation_id"`
	TimezoneOffsetMin int      `json:"timezone_offset_min"` // always 0 ?
}

/* OpenAI response *body format defination */
/* Don't imit any member because of parsing */
/* 429: {"detail": "" } */

// TODO: interface{} ?

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
	MessageType string `json:"message_type"`
	ModelSlug   string `json:"model_slug"`
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
	Message        Message     `json:"message"`
	ConversationId string      `json:"conversation_id"`
	Error          interface{} `json:"error"`
}
