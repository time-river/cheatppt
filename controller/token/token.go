package token

import (
	"reflect"
	"strings"

	"github.com/pkoukk/tiktoken-go"
	"github.com/sashabaranov/go-openai"
	log "github.com/sirupsen/logrus"
)

type Usage struct {
	PromptTokens     int `json:"promptTokens"`
	CompletionTokens int `json:"completionTokens"`
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
		log.Warnf("tiktok EncodingForModel ERROR: %s\n", err.Error())
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

func CountStringToken(prompt string) (int, error) {
	promptTokens := 0

	// use `gpt-3.5-turbo` model
	tiktok, err := initTikToken("gpt-3.5-turbo")
	if err != nil {
		return 0, err
	}

	promptTokens += tiktok.tokensPerMessage
	promptTokens += len(tiktok.tkm.Encode(prompt, nil, nil))
	promptTokens += tiktok.tokensPerName

	return promptTokens, nil
}
