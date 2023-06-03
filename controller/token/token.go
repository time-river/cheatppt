package token

import (
	"reflect"
	"strings"

	"github.com/pkoukk/tiktoken-go"
	"github.com/sashabaranov/go-openai"

	"cheatppt/log"
)

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
}

type tiktok struct {
	tkm              *tiktoken.Tiktoken
	tokensPerMessage int
	tokensPerName    int
}

func initTikToken(model string) (*tiktok, error) {
	var tiktok tiktok

	tkm, err := tiktoken.EncodingForModel(model)
	if err != nil {
		return nil, err
	}

	tiktok.tkm = tkm

	if strings.HasPrefix(model, "gpt-3.5-turbo") {
		tiktok.tokensPerMessage = 4
		tiktok.tokensPerName = -1
	} else if strings.HasPrefix(model, "gpt-4") {
		tiktok.tokensPerMessage = 3
		tiktok.tokensPerName = 1
	} else {
		log.Warn("Warning: model not found. Using cl100k_base encoding.")
		tiktok.tokensPerMessage = 3
		tiktok.tokensPerName = 1
	}

	return &tiktok, nil
}

func CountPromptToken(model string, prompts []openai.ChatCompletionMessage) (int, error) {
	promptTokens := 0

	tiktok, err := initTikToken(model)
	if err != nil {
		return 0, err
	}

	for _, prompt := range prompts {
		promptTokens += tiktok.tokensPerMessage

		t := reflect.TypeOf(prompt)
		v := reflect.ValueOf(prompt)
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			value := v.Field(i)
			promptTokens += len(tiktok.tkm.Encode(value.String(), nil, nil))
			if field.Name == "name" {
				promptTokens += tiktok.tokensPerName
			}
		}
	}

	return promptTokens, nil
}
