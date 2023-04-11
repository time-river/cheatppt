package revchatgpt2

import (
	"fmt"
	"net/http"
	"time"
)

// ClientConfig is a configuration of a client.
type ClientConfig struct {
	authToken string
	BaseURL   string
	Model     string
	Timeout   time.Duration

	HTTPClient *http.Client
}

func DefaultConfig(token, baseURL, model string) ClientConfig {
	return ClientConfig{
		authToken: token,
		BaseURL:   baseURL,
		Model:     model,

		HTTPClient: &http.Client{},
	}
}

func (c *ClientConfig) String() string {
	return fmt.Sprintf("<RevChatGPT API ClientConfig>, model: %s", c.Model)
}
