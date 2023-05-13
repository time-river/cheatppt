package chatgpt

import (
	"context"
	"sync"
	"time"

	"cheatppt/config"
	"cheatppt/log"
)

var mutex sync.Mutex

const timeout = 30 // unit: second

func RefreshToken(token string) bool {
	mutex.Lock()
	defer mutex.Unlock()

	conf := config.ChatGPT
	conf.ChatGPTToken = token

	parent := context.Background()
	ctx, cancel := context.WithTimeout(parent, timeout*time.Second)
	defer cancel()

	opts := ChatOpts{
		Ctx:    &ctx,
		Prompt: "Hi",
	}
	session, err := NewChat(&opts)
	if err != nil {
		log.Errorf("Refresh Token CreateChat ERROR: %s\n", err.Error())
		return false
	}
	defer session.Close()

	return false
}
