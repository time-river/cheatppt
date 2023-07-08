package revchatgpt3

import (
	"context"
	"fmt"
	"sync"
	"time"

	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/maps"
	"gorm.io/gorm"

	rds "cheatppt/model/redis"
	"cheatppt/model/sql"
	"cheatppt/utils"
)

type RevAccountManager struct {
	// <conversationId, email>
	redis *redis.Client
	db    *gorm.DB

	// <email, RevAccount>
	size        int
	revAccounts map[string]*RevAccount
	lru         *lru.Cache[string, *RevAccount]

	mu *sync.RWMutex
}

const (
	PERIOD_WEEK = 60 * 60 * 24 * 7 // seconds of one week
	PERIOD_DAY  = 60 * 60 * 24     // seconds of one day
)

var RevAccountMgr *RevAccountManager

func (m *RevAccountManager) GetOldestRevAccount() *RevAccount {
	if _, revAccount, ok := m.lru.GetOldest(); ok {
		return revAccount
	} else {
		return nil
	}
}

func (m *RevAccountManager) SetRevAccount(conversationId, email string) {
	ctx := context.Background()
	key := conversationId
	val := email

	m.redis.Set(ctx, key, val, PERIOD_DAY*time.Second)

	account := sql.ChatGPTConversationMapping{
		ConversationId: key,
		AccountEmail:   val,
	}
	m.db.Save(&account)
}

func (m *RevAccountManager) GetRevAccountByEmail(email string) (*RevAccount, error) {
	if revAccount, ok := m.revAccounts[email]; ok {
		return revAccount, nil
	} else {
		return nil, fmt.Errorf("ChatGPT账号已消失")
	}
}

func (m *RevAccountManager) FindRevAccountByEmail(email string) bool {
	_, ok := m.revAccounts[email]

	log.Tracef("ChatGPT Account %s found: %v\n", email, ok)
	return ok
}

func (m *RevAccountManager) GetRevAccountByConversationId(conversationId string) (*RevAccount, error) {
	ctx := context.Background()
	key := conversationId

	email, err := m.redis.Get(ctx, key).Result()
	if err != nil && err != redis.Nil {
		// redis error
		return nil, err
	} else if err == redis.Nil {
		// not found record in redis, try SQL
		var mapping sql.ChatGPTConversationMapping

		if err := m.db.First(&mapping, key).Error; err != nil {
			// sql error
			return nil, err
		}

		email = mapping.AccountEmail
		m.redis.Set(ctx, key, email, PERIOD_DAY*time.Second)
	}

	// found record, try get revAccount
	return m.GetRevAccountByEmail(email)
}

func (m *RevAccountManager) DelRevAccount(email string) {
	key := email

	log.Infof("delete ChatGPT account `%s`\n", key)

	m.lru.Remove(key)
	delete(m.revAccounts, key)
}

func (m *RevAccountManager) AddRevAccount(email string, account *RevAccount) {
	key := email
	val := account

	m.mu.Lock()
	if len(m.revAccounts) == m.size {
		size := utils.Find2power(m.size)
		m.lru.Resize(size)
	}
	m.mu.Unlock()

	m.revAccounts[key] = val
	m.lru.Add(key, val)

	log.Infof("add chatGPT account `%s`\n", key)
}

func (m *RevAccountManager) RegisterRevAccount(account *sql.ChatGPTAccount, revSession *RevSession, plaintext string) {

	revAccount := RevAccount{
		ChatGPTAccount: account,
		RevSession:     *revSession,
		mu:             sync.Mutex{},

		password:  plaintext,
		refreshAt: account.AccessTokenRfAt,
	}

	m.AddRevAccount(account.Email, &revAccount)
}

func revAccountsSetup() {
	db := sql.NewSQLClient()

	var accounts []sql.ChatGPTAccount

	if err := db.Find(&accounts).Error; err != nil {
		panic(err)
	}

	size := len(accounts)
	if size == 0 {
		size = utils.Find2power(size)
	}

	lru, err := lru.New[string, *RevAccount](size)
	if err != nil {
		panic(err)
	}

	RevAccountMgr = &RevAccountManager{
		redis: rds.NewRedisCient().GetClient(),
		db:    db,

		size:        size,
		revAccounts: make(map[string]*RevAccount),
		lru:         lru,

		mu: &sync.RWMutex{},
	}

	for _, account := range accounts {
		if !account.Activated {
			log.Infof("ignore account(`%s`) signin\n", account.Email)
			continue
		}

		log.Infof("account (`%s`) signin\n", account.Email)

		plaintext, err := decrypt(account.IV, account.Password)
		if err != nil {
			panic(err)
		}

		now := time.Now().Unix()
		ref := account.AccessTokenRfAt.Unix()
		var revSession *RevSession

		if (now - ref) > PERIOD_WEEK {
			loginReq := LoginOpts{
				Email:    account.Email,
				Password: plaintext,
			}

			revSession, err = activate(&account, &loginReq)
			if err != nil {
				panic(err)
			}

			db.Save(&account)
		} else {
			var puid string

			accessToken, err := decrypt(account.IV, account.AccessToken)
			if err != nil {
				panic(err)
			}

			if len(account.Puid) != 0 {
				puid, err = decrypt(account.IV, account.Puid)
				if err != nil {
					panic(err)
				}
			}

			revSession = &RevSession{
				Token: accessToken,
				Puid:  puid,
			}
		}

		RevAccountMgr.RegisterRevAccount(&account, revSession, plaintext)
	}

	go RevAccountMgr.refreshRoutine()
}

func (m *RevAccountManager) refreshRoutine() {
	for range time.Tick(time.Hour) {
		db := sql.NewSQLClient()

		keys := maps.Keys(m.revAccounts)
		for _, key := range keys {
			now := time.Now().Unix()
			if revAccount, exist := m.revAccounts[key]; exist {
				if (now - revAccount.refreshAt.Unix()) > PERIOD_WEEK {
					account := revAccount.ChatGPTAccount
					plaintext := revAccount.password

					loginReq := LoginOpts{
						Email:    account.Email,
						Password: plaintext,
					}

					revSession, err := activate(account, &loginReq)
					if err != nil {
						log.Error(err.Error())
						continue
					}

					revAccount.UpdateRevSession(revSession)

					db.Save(account)
				}
			}
		}
	}
}
