package revchatgpt3

import (
	"sync"
	"time"

	"cheatppt/model/sql"
)

type RevAccount struct {
	*sql.ChatGPTAccount
	RevSession

	mu sync.Mutex

	password  string
	refreshAt time.Time
}

func (m *RevAccount) Acquire() {
	m.mu.Lock()
}

func (m *RevAccount) Release() {
	m.mu.Unlock()
}

func (m *RevAccount) GetRevSession() *RevSession {
	m.Acquire()

	return &m.RevSession
}

func (m *RevAccount) PutRevSession() {
	m.Release()
}

func (m *RevAccount) UpdateRevSession(session *RevSession) {
	m.Acquire()
	defer m.Release()

	m.RevSession = *session
}
