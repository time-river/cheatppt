package apiv0

import "cheatppt/deps/revchatgpt"

/*
 * response format:
 * {
 *    "message": message,
 *    "data": data,
 *    "status": status
 * }
 */

/*
 * data: { apiModel, reverseProxy, timeoutMs, socksProxy, httpsProxy, balance }
 * example:
 *   {
 *     "message": null,
 *     "data": {
 *       "apiModel": "ChatGPTUnofficialProxyAPI",
 *       "reverseProxy": "https://openai.vvl.me:8443/chat/api/conversation",
 *       "timeoutMs": 60000,
 *       "socksProxy": "-",
 *       "httpsProxy": "-",
 *       "balance": "-"
 *     },
 *     "status": "Success"
 *   }
 */

type ModelConfig struct {
	ApiModel     string `json:"apiModel"`
	ReversePorxy string `json:"reverseProxy"`
	TimeoutMs    uint64 `json:"timeoutMs"`
	SocksProxy   string `json:"socksProxy"`
	HttpsProxy   string `json:"httpsProxy"`
	Balance      string `json:"balance"`
}

type CommonResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
	Status  string `json:"status"`
}

type VerificationRequest struct {
	Token string `json:"token"`
}

type SessionResponse struct {
	Auth  bool   `json:"auth"`
	Model string `json:"model"`
}

type ChatContext struct {
	ConversationId  *string `json:"conversationId,omitempty"`
	ParentMessageId *string `json:"parentMessageId,omitempty"`
}

type RequestProps struct {
	Prompt        string      `json:"prompt"`
	Options       ChatContext `json:"options,omitempty"` /* algouth omit it, it's value still non-nil */
	SystemMessage string      `json:"SystemMessage"`
}

/***********************************************/
type RequestOptions struct {
	message       string
	lastContext   ChatContext
	process       func(chat revchatgpt.ChatMessage)
	systemMessage string
}

var ErrorCodeMessage = map[int]string{
	401: "[OpenAI] 提供错误的API密钥 | Incorrect API key provided",
	403: "[OpenAI] 服务器拒绝访问，请稍后再试 | Server refused to access, please try again later",
	502: "[OpenAI] 错误的网关 |  Bad Gateway",
	503: "[OpenAI] 服务器繁忙，请稍后再试 | Server is busy, please try again later",
	504: "[OpenAI] 网关超时 | Gateway Time-out",
	500: "[OpenAI] 服务器繁忙，请稍后再试 | Internal Server Error",
}
