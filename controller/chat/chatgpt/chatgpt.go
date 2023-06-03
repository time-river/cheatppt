package chatgpt

import (
	"context"
	"fmt"
	"io"

	"cheatppt/config"
	"cheatppt/contrib/revchatgpt2"
	"cheatppt/controller/token"
	"cheatppt/log"
)

type ConversationOptions struct {
	MessageId       *string `json:"messageId,omitempty"`
	ConversationId  *string `json:"conversationId,omitempty"`
	ParentMessageId *string `json:"parentMessageId,omitempty"`
}

type ChatRsp struct {
	PartContent string              `json:"partContent"`
	Options     ConversationOptions `json:"options"`
	Usage       token.Usage         `json:"usage"`
}

type ChatOpts struct {
	Ctx     *context.Context
	Prompt  string
	Model   string // empty means using the default model
	Options ConversationOptions
}

type ChatSession struct {
	stream []*revchatgpt2.ChatCompletionStream
	usage  token.Usage

	messageIdx int

	// used for continue session
	ctx            *context.Context
	client         *revchatgpt2.Client
	model          string
	conversationId string
	messageId      string
}

var action = "next"

const (
//networkErr = "\nIf a 'network error' happens, please revert and resume your answer.\n"
//networkErr = "如果网络出现问题，请回退之前的回答并重新输出。"
)

func NewChat(opts *ChatOpts) (*ChatSession, error) {
	var session = ChatSession{
		ctx:    opts.Ctx,
		stream: make([]*revchatgpt2.ChatCompletionStream, 0),
		usage: token.Usage{
			PromptTokens:     0,
			CompletionTokens: 0,
		},
	}
	var err error

	conf := config.GlobalCfg.ChatGPT
	config := revchatgpt2.DefaultConfig(
		conf.ChatGPTToken,
		conf.ReverseProxyUrl,
	)

	session.client = revchatgpt2.NewClientWithConfig(config)
	revChatGPTOpts := revchatgpt2.SendMessageBrowserOptions{
		MessageId:       opts.Options.MessageId,
		ConversationId:  opts.Options.ConversationId,
		ParentMessageId: opts.Options.ParentMessageId,
		Model:           opts.Model,
		Action:          &action,
	}

	prompt := opts.Prompt
	session.usage.PromptTokens, err = token.CountStringToken(prompt)
	if err != nil {
		return nil, fmt.Errorf("内部出错了，换用其他模型试试吧")
	}

	stream, err := session.client.CreateChatCompletionStream(*opts.Ctx, prompt, revChatGPTOpts)
	if err != nil {
		return nil, fmt.Errorf("ChatGPT出错了，换用其他模型吧")
	}

	session.stream = append(session.stream, stream)

	return &session, nil
}

func (c *ChatSession) continueChat(msg string) {
	var prompt = "Please continue to answer without adding additional characters."

	c.messageIdx = 0

	revChatGPTOpts := revchatgpt2.SendMessageBrowserOptions{
		ConversationId:  &c.conversationId,
		ParentMessageId: &c.messageId,
		Model:           c.model,
		Action:          &action,
	}

	stream, err := c.client.CreateChatCompletionStream(*c.ctx, prompt, revChatGPTOpts)
	if err != nil {
		log.Warnf("ChatGPT continueReq Error: %s\n", err.Error())
	} else {
		c.stream = append(c.stream, stream)
	}
}

func (c *ChatSession) Recv() (*ChatRsp, error) {
	var msg ChatRsp
	var contRecv = true

	for contRecv {
		data, err := c.stream[len(c.stream)-1].Recv()
		if err != nil && err == io.EOF {
			return nil, err
		} else if err != nil {
			return nil, err
		}

		// first response
		if c.messageIdx == 0 {
			c.model = data.Model
			c.conversationId = data.ConversationId
			c.messageId = data.MessageId
		}

		message := data.Text
		// there maybe serval responses have the same content, don't ignore them because of the end flag (`EndTurn`)
		msg = ChatRsp{
			PartContent: message[c.messageIdx:],
			Options: ConversationOptions{
				MessageId:       &data.MessageId,
				ConversationId:  &data.ConversationId,
				ParentMessageId: &data.ParentMessageId,
			},
			Usage: c.usage,
		}
		c.messageIdx = len(message)
		c.usage.CompletionTokens += 1

		if data.Status == revchatgpt2.ChatMaxTokens {
			c.continueChat(message)
		}
		contRecv = false
	}

	return &msg, nil
}

func (c *ChatSession) Close() {
	for _, stream := range c.stream {
		stream.Close()
	}
}
