package revchatgpt2

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type Client struct {
	config ClientConfig

	requestBuilder requestBuilder
}

// NewClientWithConfig creates new OpenAI API client for specified config.
func NewClientWithConfig(config ClientConfig) *Client {
	return &Client{
		config:         config,
		requestBuilder: newRequestBuilder(),
	}
}

// ChatCompletionStream
// Note: Perhaps it is more elegant to abstract Stream using generics.
type ChatCompletionStream struct {
	*streamReader
}

// CreateChatCompletionStream — API call to create a chat completion w/ streaming
// support. It sets whether to stream back partial progress. If set, tokens will be
// sent as data-only server-sent events as they become available, with the
// stream terminated by a data: [DONE] message.
func (c *Client) CreateChatCompletionStream(
	ctx context.Context,
	prompt string,
	opts SendMessageBrowserOptions,
) (*ChatCompletionStream, error) {
	if err := checkOptions(&opts); err != nil {
		return nil, err
	}

	body := buildRequestBody(prompt, &opts)

	req, err := c.newStreamRequest(ctx, "POST", c.config.BaseURL, body)
	if err != nil {
		return nil, err
	}

	resp, err := c.config.HTTPClient.Do(req) //nolint:bodyclose // body is closed in stream.Close()
	if err != nil {
		err = &ChatGPTError{
			StatusCode: 10,
			StatusText: err.Error(),
			Err:        err,
		}
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		// example:
		// - `{"detail": "Conversation not found"}`
		// - StatusCode: 403, reason: `error code: 1020`
		var reason string
		rawReason, errResp := ioutil.ReadAll(resp.Body)
		if errResp != nil {
			reason = resp.Status
		} else {
			reason = string(rawReason)
		}

		err = &ChatGPTError{
			StatusCode: resp.StatusCode,
			StatusText: reason,
			Err:        err,
		}
		return nil, err
	}

	stream := &ChatCompletionStream{
		streamReader: &streamReader{
			parentMessageId: *opts.ParentMessageId,
			reader:          bufio.NewReader(resp.Body),
			response:        resp,
			responseBuilder: newEvewntResponseBuilder(),
		},
	}

	return stream, nil
}

func (c *Client) newStreamRequest(
	ctx context.Context,
	method string,
	url string,
	body any) (*http.Request, error) {

	log.Tracef("[%s] to %s [BODY] %s\n", method, url, body)

	req, err := c.requestBuilder.build(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.config.authToken))

	req.Header.Set("Accept", "text/event-stream; charset=utf-8")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	return req, nil
}

func checkOptions(opts *SendMessageBrowserOptions) error {
	if (opts.ConversationId == nil && opts.ParentMessageId != nil) || (opts.ConversationId != nil && opts.ParentMessageId == nil) {
		return &ChatGPTError{
			StatusCode: 0,
			StatusText: "ChatGPTUnofficialProxyAPI.sendMessage: conversationId and parentMessageId must both be set or both be undefined",
		}
	}

	if opts.ConversationId != nil && !isValidUUIDv4(*opts.ConversationId) {
		return &ChatGPTError{
			StatusCode: 1,
			StatusText: "ChatGPTUnofficialProxyAPI.sendMessage: conversationId is not a valid v4 UUID",
		}
	}

	if opts.ParentMessageId != nil && !isValidUUIDv4(*opts.ParentMessageId) {
		return &ChatGPTError{
			StatusCode: 2,
			StatusText: "ChatGPTUnofficialProxyAPI.sendMessage: parentMessageId is not a valid v4 UUID",
		}
	}

	if opts.MessageId != nil && !isValidUUIDv4(*opts.MessageId) {
		return &ChatGPTError{
			StatusCode: 3,
			StatusText: "ChatGPTUnofficialProxyAPI.sendMessage: messageId is not a valid v4 UUID",
		}
	}

	return nil
}

func buildRequestBody(text string, opts *SendMessageBrowserOptions) ConversationJSONBody {
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

	if len(opts.Model) == 0 {
		opts.Model = "text-davinci-002-render-sha"
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

	body := ConversationJSONBody{
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
		Model:           opts.Model,
		ParentMessageId: *opts.ParentMessageId,
	}

	if opts.ConversationId != nil {
		// create new one to prevent user action
		new := *opts.ConversationId
		body.ConverationId = &new
	}

	return body
}
