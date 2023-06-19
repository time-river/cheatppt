package revchatgpt3

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"
)

type ConversationSession struct {
	revAccount *RevAccount
	*Session

	parentMessageId string
	isFinished      bool

	prompt  string
	message string
}

func NewConversation(opts *SendMessageBrowserOptions, session *ClientSession) (*ConversationSession, error) {
	var revAccount *RevAccount

	if err := checkOptions(opts); err != nil {
		return nil, err
	}

	if session.RevSession == nil {
		revAccount = RevAccountMgr.GetOldestRevAccount()
		if revAccount != nil {
			session.RevSession = revAccount.GetRevSession()
		}
	}

	if session.RevSession == nil {
		return nil, &APIError{
			Message: "ChatGPT不可用",
		}
	}

	session.ReqURL = url.URL{
		Scheme: "https",
		Host:   OpenAI_HOST,
		Path:   "/backend-api/conversation",
	}
	session.ReqPath = "/api/conversation"
	session.ReqParamPath = "/conversation"
	session.ReqMethod = http.MethodPost

	body := buildRequestBody(opts)
	proxyOpts := &ProxyOpts{
		ClientSession: session,
		ReqBody:       &body,
	}

	proxySession, err := proxy(proxyOpts)
	if err != nil {
		revAccount.PutRevSession()
		return nil, &APIError{
			Message: err.Error(),
		}
	}

	return &ConversationSession{
		revAccount:      revAccount,
		Session:         proxySession,
		parentMessageId: *opts.ParentMessageId,
		prompt:          opts.Prompt,
	}, nil
}

func checkOptions(opts *SendMessageBrowserOptions) error {
	if (opts.ConversationId == nil && opts.ParentMessageId != nil) || (opts.ConversationId != nil && opts.ParentMessageId == nil) {
		return &APIError{
			Message: "ChatGPTUnofficialProxyAPI.sendMessage: conversationId and parentMessageId must both be set or both be undefined",
		}
	}

	if opts.ConversationId != nil && !isValidUUIDv4(*opts.ConversationId) {
		return &APIError{
			Message: "ChatGPTUnofficialProxyAPI.sendMessage: conversationId is not a valid v4 UUID",
		}
	}

	if opts.ParentMessageId != nil && !isValidUUIDv4(*opts.ParentMessageId) {
		return &APIError{
			Message: "ChatGPTUnofficialProxyAPI.sendMessage: parentMessageId is not a valid v4 UUID",
		}
	}

	if opts.MessageId != nil && !isValidUUIDv4(*opts.MessageId) {
		return &APIError{
			Message: "ChatGPTUnofficialProxyAPI.sendMessage: messageId is not a valid v4 UUID",
		}
	}

	return nil
}

func buildRequestBody(opts *SendMessageBrowserOptions) ConversationJSONBody {
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
					Parts:       []string{opts.Prompt},
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

func (stream *ConversationSession) Recv() (response ChatMessage, err error) {
	if stream.isFinished {
		err = io.EOF
		return
	}

	var msg ConversationResponseEvent
	dataPrefix := []byte("data:")

	for {
		line, e := stream.reader.ReadBytes('\n')
		// stream message should be end with '\n', error occurs if non
		if e != nil && e == io.EOF {
			stream.isFinished = true
			err = io.EOF
			return
		} else if e != nil {
			err = fmt.Errorf("ChatGPT handle error: %s", e.Error())
			return
		} else if len(line) < len(dataPrefix) {
			// empty line, e.g. `\n`
			continue
		} else if !bytes.HasPrefix(line, dataPrefix) {
			// e.g. `event: ping\n`
			continue
		}

		log.Tracef("revchatgpt: `%s`\n", line)

		line = bytes.TrimPrefix(line, dataPrefix)
		line = bytes.TrimSpace(line)
		if bytes.Equal(line, []byte("[DONE]")) {
			stream.isFinished = true
			err = io.EOF
			return
		}

		if err = json.Unmarshal(line, &msg); err != nil {
			return
		} else if msg.Error != nil {
			// example:
			// {
			// 	 "error": {
			// 	   "code": "invalid_api_key",
			// 	   "param": null,
			// 	   "message": "Incorrect API key provided: ...",
			// 	   "type": "invalid_request_error"
			// 	 }
			// }
			response.Error = msg.Error
			err = io.EOF
			return
		}

		// Thress types of role: `system`, `user`, `assistant`
		// Only `assistant` is what we need
		if msg.Message.Author.Role != "assistant" {
			continue
		}

		// recv one message
		break
	}

	response = stream.buildChatMessage(&msg)
	return
}

func (stream *ConversationSession) buildChatMessage(msg *ConversationResponseEvent) ChatMessage {
	// not completion example:
	//
	// max token:
	// {
	//   ...
	//   "status": "finished_successfully",
	//   "end_turn": false,
	//   "weight": 1.0,
	//   "metadata": {
	//     "message_type": "next",
	//     "model_slug": "text-davinci-002-render-sha",
	//     "finish_details": {
	//       "type": "max_tokens"
	//     }
	//   },
	//   ...
	// }
	//
	// pipe broken?:
	// {
	//   ...
	//	 "status": "in_progress",
	//   "end_turn": null,
	//   "weight": 1.0,
	//   "metadata": {
	//     "message_type": "next",
	//     "model_slug": "text-davinci-002-render-sha"
	//   },
	//   ...
	// }
	//
	// completion example:
	// {
	//   "status": "finished_successfully",
	//   "end_turn": true,
	//   "weight": 1.0,
	//   "metadata": {
	//     "message_type": "next",
	//     "model_slug": "text-davinci-002-render-sha",
	//     "finish_details": {
	//   	"type": "stop",
	//   	"stop": ""
	//     }
	//   },
	// }

	var data = ChatMessage{
		MessageId:       msg.Message.Id,
		Text:            msg.Message.Content.Parts[0],
		Role:            msg.Message.Author.Role,
		Model:           msg.Message.Metadata.ModelSlug,
		Status:          ChatInProgrss,
		ParentMessageId: stream.parentMessageId,
		ConversationId:  msg.ConversationId,
	}

	if msg.Message.EndTurn != nil && *msg.Message.EndTurn {
		data.Status = ChatStop
	} else if msg.Message.EndTurn != nil && !*msg.Message.EndTurn {
		data.Status = ChatMaxTokens
	}

	stream.message = data.Text

	return data
}

func (stream *ConversationSession) Close() {
	stream.Session.Close()
	stream.revAccount.PutRevSession()
}
