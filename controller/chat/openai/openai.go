package openai

import (
	"context"
	"errors"
	"net/http"

	"github.com/sashabaranov/go-openai"

	"cheatppt/config"
	"cheatppt/log"
)

type ChatReq openai.ChatCompletionRequest
type ChatMsgRsp openai.ChatCompletionResponse
type ChatAPIErr openai.APIError
type ChatErrRsp struct {
	Error *ChatAPIErr `json:"error,omitempty"`
}

type ChatOpts struct {
	Ctx context.Context
	Req ChatReq
}

type ChatSession struct {
	stream *openai.ChatCompletionStream
}

func buildCfg() openai.ClientConfig {
	conf := config.OpenAI

	cfg := openai.DefaultConfig(conf.Token)
	if conf.BaseURL != "" {
		cfg.BaseURL = conf.BaseURL
	}
	if conf.OrdID != "" {
		cfg.OrgID = conf.OrdID
	}

	return cfg
}

func NewChat(opts *ChatOpts) (*ChatSession, *ChatErrRsp) {
	var session ChatSession
	var err error

	cfg := buildCfg()
	c := openai.NewClientWithConfig(cfg)

	session.stream, err = c.CreateChatCompletionStream(opts.Ctx, openai.ChatCompletionRequest(opts.Req))
	if err != nil {
		log.Errorf("ChatCompletionStream ERROR: %v\n", err)

		// request error
		reqErr := &openai.RequestError{}
		if errors.As(err, &reqErr) {
			e := &openai.APIError{
				Type:           "invalid_request_error",
				Message:        err.Error(),
				HTTPStatusCode: reqErr.HTTPStatusCode,
			}
			return nil, &ChatErrRsp{
				Error: (*ChatAPIErr)(e),
			}
		}

		// response error
		resErr := &openai.APIError{}
		if errors.As(err, &resErr) {
			return nil, &ChatErrRsp{
				Error: (*ChatAPIErr)(resErr),
			}
		}

		// others
		return nil, &ChatErrRsp{
			Error: &ChatAPIErr{
				Type:           "invalid_request_error",
				Message:        err.Error(),
				HTTPStatusCode: http.StatusOK,
			},
		}
	}

	return &session, nil
}

func (c *ChatSession) Recv() (openai.ChatCompletionStreamResponse, error) {
	return c.stream.Recv()
}
