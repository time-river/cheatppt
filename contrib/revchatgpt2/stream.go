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
	prefix := "data: "

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
		} else if len(line) < len(prefix) {
			// empty line, e.g. `\n`
			continue
		}

		log.Debugf("revchatgpt: `%s`\n", line)

		line = bytes.TrimPrefix(line, []byte(prefix))
		if bytes.Equal(line, []byte("[DONE]\n")) {
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
	var result = ChatMessage{
		ParentMessageId: stream.parentMessageId,
	}

	result.ConversationId = msg.ConversationId

	if len(msg.Message.Id) > 0 {
		result.Id = msg.Message.Id
	}

	message := msg.Message
	text := message.Content.Parts[0]

	if len(text) > 0 {
		result.Text = text
	}

	return result
}
