package openai

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/kr/pretty"
	"github.com/sashabaranov/go-openai"

	"cheatppt/config"
	"cheatppt/controller/token"
	"cheatppt/log"
)

type ChatRsp struct {
	openai.ChatCompletionStreamResponse
	Usage token.Usage `json:"usage"`
}

type ChatAPIErr openai.APIError

type ChatErrRsp struct {
	Error *ChatAPIErr `json:"error,omitempty"`
}

type ChatReq openai.ChatCompletionRequest

type ChatOpts struct {
	Ctx context.Context
	Req ChatReq
}

type ChatSession struct {
	stream *openai.ChatCompletionStream
	usage  token.Usage
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
	var session = ChatSession{
		usage: token.Usage{
			PromptTokens:     0,
			CompletionTokens: 0,
		},
	}
	var err error

	session.usage.PromptTokens, err = token.CountPromptToken(opts.Req.Model, opts.Req.Messages)
	if err != nil {
		return nil, &ChatErrRsp{
			Error: &ChatAPIErr{
				Type:    "internal_tiktok_error",
				Message: err.Error(),
			},
		}
	}

	cfg := buildCfg()
	c := openai.NewClientWithConfig(cfg)

	// force stream output
	opts.Req.Stream = true

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

func (c *ChatSession) Recv() (*ChatRsp, error) {
	var rsp ChatRsp

	data, err := c.stream.Recv()
	if err != nil && err == io.EOF {
		return nil, err
	} else if err != nil {
		log.Warnf("Recv Msg ERROR: %s\n", err.Error())
		return nil, err
	}

	log.Trace(pretty.Sprint(data))

	// each response in stream regards as one token
	c.usage.CompletionTokens += 1
	rsp = ChatRsp{
		ChatCompletionStreamResponse: data,
		Usage:                        c.usage,
	}

	return &rsp, nil
}

func (c *ChatSession) Close() {
	c.stream.Close()
}
