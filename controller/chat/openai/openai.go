package openai

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/kr/pretty"
	"github.com/sashabaranov/go-openai"
	log "github.com/sirupsen/logrus"

	"cheatppt/config"
	"cheatppt/controller/billing"
	"cheatppt/controller/model"
	"cheatppt/controller/qos"
	"cheatppt/controller/token"
)

type ChatRsp struct {
	openai.ChatCompletionStreamResponse
	Usage token.Usage `json:"usage"`
}

type ChatAPIErr openai.APIError

type ChatErrRsp struct {
	Error *ChatAPIErr `json:"error,omitempty"`
}

type ChatOpts struct {
	Ctx    context.Context
	UserId int
	Model  string
	Req    openai.ChatCompletionRequest
}

type ChatSession struct {
	stream      *openai.ChatCompletionStream
	usage       token.Usage
	UserId      int
	Model       string
	InputCoins  int
	OutputCoins int

	retry    int
	consumer *billing.Consumer
}

const (
	Provider  = "openai"
	MAX_RETRY = 3
)

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
		UserId: opts.UserId,
		Model:  opts.Model,

		retry: 0,
	}
	var err error

	// 1. does the model exist?
	model, ok := model.CacheFind(model.BuildCacheKey(Provider, opts.Model))
	if !ok || !model.Activated {
		return nil, &ChatErrRsp{
			Error: &ChatAPIErr{
				Type:    "input_error",
				Message: "模型不存在",
			},
		}
	}

	session.InputCoins = model.InputCoins
	session.OutputCoins = model.OutputCoins

	promptTokens, err := token.CountPromptToken(opts.Req.Model, opts.Req.Messages)
	if err != nil {
		return nil, &ChatErrRsp{
			Error: &ChatAPIErr{
				Type:    "internal_tiktok_error",
				Message: err.Error(),
			},
		}
	}
	session.usage.PromptTokens = promptTokens
	promptPrice := promptTokens * model.InputCoins

	// 2. try to consume
	session.consumer, err = billing.GetComsumer(session.UserId)
	if err != nil {
		return nil, &ChatErrRsp{
			Error: &ChatAPIErr{
				Type:    "internal_tiktok_error",
				Message: err.Error(),
			},
		}
	}
	session.consumer.Comsume(promptPrice)

	// 3. ratelimit
	qosMeta := qos.Meta{
		Consumer: session.consumer,
		Model:    session.Model,
		Provider: Provider,
	}
	if !qos.Allow(&qosMeta) {
		session.consumer.Rollback()

		return nil, &ChatErrRsp{
			Error: &ChatAPIErr{
				Type:    "rate_limit",
				Message: "too many request",
			},
		}
	}

	cfg := buildCfg()
	c := openai.NewClientWithConfig(cfg)

	// force stream output
	opts.Req.Stream = true

	session.stream, err = c.CreateChatCompletionStream(
		opts.Ctx, openai.ChatCompletionRequest(opts.Req))
	if err != nil {
		log.Warn(err.Error())

		session.revertRequest()

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
	if err != nil {
		price := c.OutputCoins * c.usage.CompletionTokens
		c.consumer.Comsume(price)
		c.consumer.Commit()
		if err == io.EOF {
			return nil, err
		} else {
			log.Errorf("OPENAI Recv ERROR: %s\n", err.Error())
		}
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

func (c *ChatSession) revertRequest() {
	c.consumer.Rollback()
}
