package revchatgpt2

import "encoding/json"

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

func (msg ChatMessage) Marshal() ([]byte, error) {
	return json.Marshal(msg)
}
