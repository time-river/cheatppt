package chatgpt

import (
	"context"

	"cheatppt/config"
	"cheatppt/contrib/revchatgpt2"
)

type ChatOpts struct {
	Ctx             context.Context
	Text            string
	Model           string // empty means using the default model
	ConversationId  *string
	ParentMessageId *string
}

var action = "next"

func NewChat(opts *ChatOpts) (*revchatgpt2.ChatCompletionStream, error) {
	conf := config.GlobalCfg.ChatGPT
	config := revchatgpt2.DefaultConfig(
		conf.ChatGPTToken,
		conf.ReverseProxyUrl,
	)

	revChatGPTclient := revchatgpt2.NewClientWithConfig(config)
	revChatGPTOpts := revchatgpt2.SendMessageBrowserOptions{
		ConversationId:  opts.ConversationId,
		ParentMessageId: opts.ParentMessageId,
		Model:           opts.Model,
		Action:          &action,
	}

	return revChatGPTclient.CreateChatCompletionStream(opts.Ctx, opts.Text, revChatGPTOpts)
}
