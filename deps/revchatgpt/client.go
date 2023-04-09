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
	"gopkg.in/cenkalti/backoff.v1"
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
		ParentMessageId: *opts.MessageId,
		ConversationId:  opts.ConversationId,
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
				api.debug("json.Unmarshal error: %s, body: %s\n", err.Error(), msg.Data)
				return
			}

			if convoResponseEvent.ConversationId != nil {
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
				// library won't close the socket manually, therefore
				// we do that.
				defer resp.Body.Close()

				reason, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					reason = []byte(resp.Status)
				}

				msg := fmt.Sprintf("ChatGPT error %d: %s", resp.StatusCode, reason)
				api.debug("%s\n", msg)

				// returns a *PermanentError to prevent retry according the library comments.
				// otherwise we should define `sse.Client.ReconnectStrategy`.
				err = &ChatGPTError{
					StatusCode: resp.StatusCode,
					StatusText: msg,
				}
				return backoff.Permanent(err)
			}
			return nil
		},
	}

	api.debug("POST %s\n  body: %s\n  headers: %s\n", options.url, options.body, options.headers)

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

	if (opts.ConversationId == nil && opts.ParentMessageId != nil) || (opts.ConversationId != nil && opts.ParentMessageId == nil) {
		err := &ChatGPTError{
			StatusCode: 0,
			StatusText: "ChatGPTUnofficialProxyAPI.sendMessage: conversationId and parentMessageId must both be set or both be undefined",
		}
		return nil, err
	}

	if opts.ConversationId != nil && !isValidUUIDv4(*opts.ConversationId) {
		err := &ChatGPTError{
			StatusCode: 1,
			StatusText: "ChatGPTUnofficialProxyAPI.sendMessage: conversationId is not a valid v4 UUID",
		}
		return nil, err
	}

	if opts.ParentMessageId != nil && !isValidUUIDv4(*opts.ParentMessageId) {
		err := &ChatGPTError{
			StatusCode: 2,
			StatusText: "ChatGPTUnofficialProxyAPI.sendMessage: parentMessageId is not a valid v4 UUID",
		}
		return nil, err
	}

	if opts.MessageId != nil && !isValidUUIDv4(*opts.MessageId) {
		err := &ChatGPTError{
			StatusCode: 3,
			StatusText: "ChatGPTUnofficialProxyAPI.sendMessage: messageId is not a valid v4 UUID",
		}
		return nil, err
	}

	/* both of conversationId and parentMessageId are undefined */
	if opts.ParentMessageId == nil {
		uuid := uuidv4()
		opts.ParentMessageId = &uuid
	}

	if opts.MessageId == nil {
		uuid := uuidv4()
		opts.MessageId = &uuid
	}

	if opts.Action == nil {
		action := "next"
		opts.Action = &action
	}

	// Body example:
	//
	//   first chat:
	//     body: {
	//    	  action: 'next',
	//    	  messages: [ [Object] ],
	//    	  model: 'gpt-4',
	//    	  parent_message_id: '4436c401-1da3-4375-88fa-b2a5f6efbb3d'
	//      },
	//
	//   the following:
	//     body: {
	//	      action: 'next',
	//	      messages: [ [Object] ],
	//	      model: 'gpt-4',
	//	      parent_message_id: 'a6638807-f58a-43e3-9aa5-70dec9623f81',
	//	      conversation_id: 'b08ea722-9c4f-4ecf-b4c5-4a5da0f0ff7e'
	//      },
	//     body: {
	//    	  action: 'next',
	//    	  messages: [ [Object] ],
	//    	  model: 'gpt-4',
	//    	  parent_message_id: '8f0c2874-a471-471d-90f9-524677945a47',
	//    	  conversation_id: 'b08ea722-9c4f-4ecf-b4c5-4a5da0f0ff7e'
	//      },
	//     body: {
	//    	  action: 'next',
	//    	  messages: [ [Object] ],
	//    	  model: 'gpt-4',
	//    	  parent_message_id: 'f1a6a62b-2e93-47b3-8b34-6774a88dd00b',
	//    	  conversation_id: 'b08ea722-9c4f-4ecf-b4c5-4a5da0f0ff7e'
	//      },

	rawBody := ConversationJSONBody{
		Action: *opts.Action,
		Messages: []Prompt{
			{
				Id: *opts.MessageId,
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
		ParentMessageId: *opts.ParentMessageId,
	}

	if opts.ConversationId != nil {
		// create new one to prevent user action
		new := *opts.ConversationId
		rawBody.ConverationId = &new
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
