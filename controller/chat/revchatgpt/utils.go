package revchatgpt3

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"sync"

	log "github.com/sirupsen/logrus"

	"cheatppt/config"
)

var onceConf sync.Once
var aesEncryptKey []byte // password encryption

func Setup() {
	revAccountsSetup()
}

func getAesEncryptKey() []byte {
	if aesEncryptKey == nil {
		onceConf.Do(func() {
			key := []byte(config.Server.Secret)
			hash := md5.Sum(key)
			aesEncryptKey = hash[:]
		})
	}

	return aesEncryptKey
}

func encrypt(plaintext string, IV []byte) ([]byte, []byte, error) {
	key := getAesEncryptKey()

	if len(plaintext) == 0 {
		log.Warn("encrypt nothing!")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Error(err)
		return nil, nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}

	iv, err := newIV(IV, aesgcm.NonceSize())
	if err != nil {
		return nil, nil, err
	}

	ciphertext := aesgcm.Seal(nil, iv, []byte(plaintext), nil)
	return iv, ciphertext, nil
}

func decrypt(iv, ciphertext []byte) (string, error) {
	key := getAesEncryptKey()

	if len(ciphertext) == 0 {
		log.Warn("encrypt nothing!")
	} else if len(iv) == 0 {
		panic("iv is empty")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	decrypted, err := aesgcm.Open(nil, iv, ciphertext, nil)
	if err != nil {
		return "", err
	}

	plaintext := string(decrypted)
	return plaintext, nil
}
