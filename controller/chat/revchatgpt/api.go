package revchatgpt3

import "net/url"

type SendMessageBrowserOptions struct {
	/* `*Id`` can do not exist:
	 * - exist: not nil
	 * - not exist: nil
	 */
	ConversationId  *string
	ParentMessageId *string
	MessageId       *string
	Action          *string /* 'next' | 'variant' */
	Model           string
	Prompt          string
}

type ClientOpts struct {
	UserId       int
	ReqURL       url.URL
	ReqPath      string
	ReqParamPath string
	ReqMethod    string
}

type ClientSession struct {
	ClientOpts
	*RevSession
}

type ChatStatus int

const (
	ChatInProgrss ChatStatus = iota
	ChatStop
	ChatMaxTokens
)

type ChatMessage struct {
	MessageId string     `json:"messageId,omitempty"`
	Text      string     `json:"text,omitempty"`
	Role      string     `json:"role,omitempty"`
	Model     string     `json:"model,omitempty"`
	Status    ChatStatus `json:"status,omitempty"`
	// relevant for both ChatGPTAPI and ChatGPTUnofficialProxyAPI
	ParentMessageId string `json:"parentMessageId,omitempty"`
	// only relevant for ChatGPTUnofficialProxyAPI (optional for ChatGPTAPI)
	ConversationId string    `json:"conversationId,omitempty"`
	Error          *APIError `json:"error,omitempty"`
}
