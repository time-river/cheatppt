package revchatgpt2

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"

	"cheatppt/log"
)

type streamReader struct {
	isFinished      bool
	parentMessageId string

	reader   *bufio.Reader
	response *http.Response

	responseBuilder responseBuilder
}

func (stream *streamReader) Close() {
	stream.response.Body.Close()
}

func (stream *streamReader) Recv() (response ChatMessage, err error) {
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
			text := fmt.Sprintf("ChatGPT handle error: %s", e.Error())
			err = &ChatGPTError{
				StatusCode: 20,
				StatusText: text,
				Err:        e,
			}
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

		e = stream.responseBuilder.build(line, &msg)
		if e != nil {
			err = &ChatGPTError{
				StatusCode: 21,
				StatusText: e.Error(),
				Err:        e,
			}
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
			err = &ChatGPTError{
				StatusCode: 22,
				StatusText: msg.Error.String(),
			}

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

func (stream *streamReader) buildChatMessage(msg *ConversationResponseEvent) ChatMessage {
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

	return data
}
