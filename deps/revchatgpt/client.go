package revchatgpt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/google/uuid"
	"github.com/r3labs/sse/v2"
)

func isValidUUIDv4(str string) bool {
	uuidv4Re := regexp.MustCompile("(?i)^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$")
	return len(str) > 0 && uuidv4Re.MatchString(str)
}

func uuidv4() string {
	return uuid.Must(uuid.NewRandom()).String()
}

func NewChatGPTUnofficialProxyAPI(token, url, model string) *ChatGPTUnofficialProxyAPI {
	return &ChatGPTUnofficialProxyAPI{
		apiReverseProxyUrl: url,
		model:              model,
		Headers: map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", token),
			"Accept":        "text/event-stream",
			"Content-Type":  "application/json",
		},
	}
}

func (api *ChatGPTUnofficialProxyAPI) EnableDebug(w io.Writer) {
	api.debugger = true
	api.logger = w
}

func (api *ChatGPTUnofficialProxyAPI) debug(format string, v ...any) {
	if !api.debugger {
		return
	}

	fmt.Fprintf(api.logger, format, v...)
}

func (api *ChatGPTUnofficialProxyAPI) send(body []byte, opts *SendMessageBrowserOptions) (*ChatMessage, *ChatGPTError) {
	var result = &ChatMessage{
		Role:            "assistant",
		Id:              uuidv4(),
		ParentMessageId: opts.MessageId,
		ConversationId:  "",
		Text:            "",
	}

	options := &sseOptions{
		url:       api.apiReverseProxyUrl,
		method:    "POST",
		headers:   api.Headers,
		body:      body,
		timeoutMs: opts.TimeoutMs,
		subscriber: func(msg *sse.Event) {
			api.debug("data: %s\n", msg.Data)

			if bytes.Equal(msg.Data, []byte("[DONE]")) {
				return
			}

			var convoResponseEvent ConversationResponseEvent
			if err := json.Unmarshal(msg.Data, &convoResponseEvent); err != nil {
				api.debug("%s\n", err.Error())
				return
			}

			if len(convoResponseEvent.ConversationId) > 0 {
				result.ConversationId = convoResponseEvent.ConversationId
			}

			if len(convoResponseEvent.Message.Id) > 0 {
				result.Id = convoResponseEvent.Message.Id
			}

			message := convoResponseEvent.Message

			text := message.Content.Parts[0]
			if len(text) > 0 {
				if len(text) > 0 {
					result.Text = text
				}

				if opts.OnProgress != nil {
					opts.OnProgress(*result)
				}
			}
		},
		validator: func(c *sse.Client, resp *http.Response) error {

			if resp.StatusCode != http.StatusOK {
				reason, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					reason = []byte(resp.Status)
				}

				msg := fmt.Sprintf("ChatGPT error %d: %s", resp.StatusCode, reason)
				api.debug("%s\n", msg)

				err = &ChatGPTError{
					StatusCode: resp.StatusCode,
					StatusText: msg,
				}
				return err
			}
			return nil
		},
	}

	api.debug("[%s] URL %s BODY %s\n", options.method, options.url, options.body)

	err := fetchSSE(options)
	if err != nil {
		api.debug("Error: %s\n", err.Error())
		if e, ok := err.(*ChatGPTError); ok {
			return nil, e
		} else {
			return nil, &ChatGPTError{
				StatusCode: 5,
				StatusText: "Internal Error",
			}
		}
	} else {
		return result, nil
	}
}

func (api *ChatGPTUnofficialProxyAPI) SendMessage(text string, opts SendMessageBrowserOptions) (*ChatMessage, *ChatGPTError) {
	ci := len(opts.ConversationId)
	pi := len(opts.ParentMessageId)

	if (ci == 0 && pi != 0) || (ci != 0 && pi == 0) {
		err := &ChatGPTError{
			StatusCode: 0,
			StatusText: "ChatGPTUnofficialProxyAPI.sendMessage: conversationId and parentMessageId must both be set or both be undefined",
		}
		return nil, err
	}

	if len(opts.ConversationId) > 0 && !isValidUUIDv4(opts.ConversationId) {
		err := &ChatGPTError{
			StatusCode: 1,
			StatusText: "ChatGPTUnofficialProxyAPI.sendMessage: conversationId is not a valid v4 UUID",
		}
		return nil, err
	}

	if len(opts.ParentMessageId) > 0 && !isValidUUIDv4(opts.ParentMessageId) {
		err := &ChatGPTError{
			StatusCode: 2,
			StatusText: "ChatGPTUnofficialProxyAPI.sendMessage: parentMessageId is not a valid v4 UUID",
		}
		return nil, err
	}

	if len(opts.MessageId) > 0 && !isValidUUIDv4(opts.MessageId) {
		err := &ChatGPTError{
			StatusCode: 3,
			StatusText: "ChatGPTUnofficialProxyAPI.sendMessage: messageId is not a valid v4 UUID",
		}
		return nil, err
	}

	rawBody := ConversationJSONBody{
		Action: "next",
		Messages: []Prompt{
			{
				Id: uuidv4(),
				Author: PromptAuthor{
					Role: "user",
				},
				Content: PromptContent{
					ContentType: "text",
					Parts:       []string{text},
				},
			},
		},
		Model:           api.model,
		ParentMessageId: opts.ParentMessageId,
	}

	if opts.ConversationId != "" {
		rawBody.ConverationId = opts.ConversationId
	}

	body, err := json.Marshal(rawBody)
	if err != nil {
		err := &ChatGPTError{
			StatusCode: 4,
			StatusText: "ChatGPTUnofficialProxyAPI.sendMessage: JSON body data failed",
		}
		return nil, err
	}

	return api.send(body, &opts)
}
