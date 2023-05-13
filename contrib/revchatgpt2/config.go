package revchatgpt2

import (
	"net/http"
)

// ClientConfig is a configuration of a client.
type ClientConfig struct {
	authToken string
	BaseURL   string

	HTTPClient *http.Client
}

func DefaultConfig(token, baseURL string) ClientConfig {
	return ClientConfig{
		authToken: token,
		BaseURL:   baseURL,

		HTTPClient: &http.Client{},
	}
}

func (c *ClientConfig) String() string {
	return "<RevChatGPT API ClientConfig>"
}
