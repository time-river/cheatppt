package revchatgpt

import (
	"bytes"
	"net/http"
	"time"

	"github.com/r3labs/sse/v2"
)

type sseOptions struct {
	url        string
	method     string
	headers    map[string]string
	body       []byte
	timeoutMs  uint64
	subscriber func(msg *sse.Event)
	validator  func(c *sse.Client, resp *http.Response) error
}

func fetchSSE(opts *sseOptions) error {
	client := sse.NewClient(opts.url, func(c *sse.Client) {
		for key, value := range opts.headers {
			c.Headers[key] = value
		}
		c.Method = "POST"
		c.Connection.Timeout = time.Duration(opts.timeoutMs) * time.Millisecond

		c.ResponseValidator = opts.validator
		c.Body = bytes.NewReader(opts.body)
	})

	return client.SubscribeRaw(opts.subscriber)
}
