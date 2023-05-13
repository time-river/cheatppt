package revchatgpt2

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
}

type ChatStatus int

const (
	ChatInProgrss ChatStatus = iota
	ChatStop
	ChatMaxTokens
)

type ChatMessage struct {
	MessageId string     `json:"messageId"`
	Text      string     `json:"text"`
	Role      string     `json:"role"`
	Model     string     `json:"model"`
	Status    ChatStatus `json:"status"`
	// relevant for both ChatGPTAPI and ChatGPTUnofficialProxyAPI
	ParentMessageId string `json:"parentMessageId,omitempty"`
	// only relevant for ChatGPTUnofficialProxyAPI (optional for ChatGPTAPI)
	ConversationId string `json:"conversationId,omitempty"`
}
