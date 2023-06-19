package revchatgpt3

import (
	"crypto/rand"
	"errors"
	"io"
	"sync"
	"time"

	"github.com/kr/pretty"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"cheatppt/model/sql"
)

type Account struct {
	Email     string
	Password  string
	Activated bool
}

// only one account can be operated at the same time
var mu sync.RWMutex

func activate(account *sql.ChatGPTAccount, loginReq *LoginOpts) (*RevSession, error) {
	revSession, err := Login(loginReq)
	if err != nil {
		return nil, err
	}
	account.AccessTokenRfAt = time.Now()

	_, ciphertext, err := encrypt(revSession.Token, account.IV)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	account.AccessToken = ciphertext
	if len(revSession.Puid) != 0 {
		_, ciphertext, err := encrypt(revSession.Puid, account.IV)
		if err != nil {
			log.Error(err.Error())
			return nil, err
		}

		account.Puid = ciphertext
		account.PuidRfAt = account.AccessTokenRfAt
	}

	return revSession, nil
}

func fixupPassword(account *sql.ChatGPTAccount, detail *Account) (bool, error) {
	plaintext := detail.Password

	if account.Password != nil {
		plaintextPtr, err := decrypt(account.IV, account.Password)
		if err != nil {
			return false, err
		} else if plaintextPtr == plaintext {
			return false, nil
		}
	}

	IV, ciphertext, err := encrypt(plaintext, account.IV)
	if err != nil {
		return false, err
	}
	account.Password = ciphertext
	account.IV = IV

	return true, nil
}

func needLogin(fixed bool, detail *Account) bool {
	if !detail.Activated {
		// request don't allow login
		return false
	} else if fixed {
		// password change
		return true
	} else if !fixed && !RevAccountMgr.FindRevAccountByEmail(detail.Email) {
		// password don't change but not register
		return true
	} else {
		return false
	}
}

func AddAccount(detail *Account) error {
	mu.Lock()
	defer mu.Unlock()

	var err error

	account := sql.ChatGPTAccount{Email: detail.Email}

	db := sql.NewSQLClient()
	err = db.Model(&sql.ChatGPTAccount{}).Where("email = ?", detail.Email).First(&account).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		// sql error
		log.Errorf("revchatgpt3 AddAccount First ERROR: %s\n", err.Error())
		return err
	}

	fixed, err := fixupPassword(&account, detail)
	if err != nil {
		return err
	}

	if needLogin(fixed, detail) {
		loginReq := LoginOpts{
			Email:    detail.Email,
			Password: detail.Password,
		}

		revSession, err := activate(&account, &loginReq)
		if err != nil {
			return err
		}

		revAccount, _ := RevAccountMgr.GetRevAccountByEmail(detail.Email)
		if revAccount != nil {
			revAccount.UpdateRevSession(revSession)
		} else {
			RevAccountMgr.RegisterRevAccount(&account, revSession, detail.Password)
		}
	}

	if !detail.Activated {
		RevAccountMgr.DelRevAccount(detail.Email)
	}

	account.Activated = detail.Activated

	if err := db.Save(&account).Error; err != nil {
		// DON'T FREE CACHE
		log.Errorf(err.Error())
		return err
	}

	log.Trace(pretty.Sprint(account))

	return nil
}

func DelAccount(email *string) {
	mu.Lock()
	defer mu.Unlock()

	account := sql.ChatGPTAccount{Email: *email}
	db := sql.NewSQLClient()
	db.Where("email = ?", *email).Delete(&account)

	RevAccountMgr.DelRevAccount(account.Email)
}

func ListAccounts() ([]Account, error) {
	mu.RLock()
	defer mu.RUnlock()

	var accounts []sql.ChatGPTAccount
	db := sql.NewSQLClient()

	if err := db.Find(&accounts).Error; err != nil {
		log.Errorf("revchatgpt3 ListAccounts db.Find Error: %s\n", err.Error())
		return nil, err
	}

	var data = make([]Account, len(accounts))
	for i, account := range accounts {
		plaintext, err := decrypt(account.IV, account.Password)
		if err != nil {
			return nil, err
		}

		data[i] = Account{
			Email:     account.Email,
			Password:  plaintext,
			Activated: account.Activated,
		}
	}

	return data, nil
}

func newIV(IV []byte, size int) ([]byte, error) {
	var iv []byte

	if IV == nil || len(IV) != size {
		iv = make([]byte, size)

		if _, err := io.ReadFull(rand.Reader, iv); err != nil {
			return nil, err
		}
	} else {
		zero := true
		iv = IV

		for _, val := range iv {
			if val != 0 {
				zero = false
				break
			}
		}

		if zero {
			if _, err := io.ReadFull(rand.Reader, iv); err != nil {
				return nil, err
			}
		}
	}

	return iv, nil
}
