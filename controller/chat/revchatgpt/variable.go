package revchatgpt3

import "cheatppt/config"

const (
	user_agent  = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36"
	OpenAI_HOST = "chat.openai.com"
)

var (
	arkose_token string
)

func getHttpProxy() string {
	return config.ChatGPT.HttpProxy
}
