package revchatgpt2

import "fmt"

type ChatGPTError struct {
	StatusCode int
	StatusText string
	IsFinal    bool
	AccountId  string
	Err        error
}

var _ = error(&ChatGPTError{})

func (err *ChatGPTError) Error() string {
	return fmt.Sprintf("ChatGPT error %d: %s", err.StatusCode, err.StatusText)
}
