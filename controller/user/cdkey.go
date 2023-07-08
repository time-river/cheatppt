package user

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"cheatppt/model/redis"
	"cheatppt/model/sql"
)

type CDKeyMeta struct {
	Nr      int
	Comment string
	Credit  float32
	Expire  int
}

type CDKeyClaims struct {
	Credit    float32 `json:"credit"`
	Comment   string  `json:"comment"`
	NotBefore int     `json:"notBefore"` // unix time
}

type CDKeyDetail struct {
	CDKeyClaims
	Id string `json:"id"`
}

const KeyLength = 16
const cdkeyPrefix = "ck"

func getCode(length int) string {

	const lettersAndDigits = `abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789`
	code := make([]byte, length)
	for i := range code {
		code[i] = lettersAndDigits[rand.Intn(len(lettersAndDigits))]
	}

	return string(code)
}

func GenCDKeys(meta *CDKeyMeta) ([]string, error) {
	ctx := context.Background()
	rds := redis.NewRedisCient().GetClient()
	keys := make([]string, 0)

	if time.Now().Unix() > int64(meta.Expire) {
		return nil, fmt.Errorf("无效日期")
	}

	claims := &CDKeyClaims{
		Credit:    meta.Credit,
		Comment:   meta.Comment,
		NotBefore: meta.Expire,
	}
	cdkeyJSON, err := json.Marshal(claims)
	if err != nil {
		log.Errorf("CDKEY json.Marshal ERROR: %s\n", err.Error())
		return nil, fmt.Errorf("内部错误")
	}

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < meta.Nr; i++ {
		id := getCode(KeyLength)
		key := fmt.Sprintf("%s-%s", cdkeyPrefix, id)

		expire := time.Duration(meta.Expire-int(time.Now().Unix())) * time.Second
		if expire <= 0 {
			goto out
		}
		if err := rds.Set(ctx, key, cdkeyJSON, expire).Err(); err != nil {
			log.Errorf("CDKEY rds SET ERROR: %s\n", err.Error())
			goto out
		}

		keys = append(keys, key)
	}

out:
	return keys, nil
}

func ExgCDkey(cdkey string, userId int) error {
	ctx := context.Background()
	rds := redis.NewRedisCient().GetClient()
	var claims CDKeyClaims

	tx := rds.TxPipeline()
	get := tx.Get(ctx, cdkey)
	tx.Del(ctx, cdkey)
	if _, err := tx.Exec(ctx); err != nil {
		log.Errorf("CDKEY ExgCDKey ERROR: %s\n", err.Error())
		return fmt.Errorf("KEY不存在")
	}
	if err := json.Unmarshal([]byte(get.Val()), &claims); err != nil {
		return fmt.Errorf("内部错误")
	}
	if claims.NotBefore < int(time.Now().Unix()) {
		return fmt.Errorf("CDKey无效")
	}

	db := sql.NewSQLClient()
	db.Transaction(func(tx *gorm.DB) error {
		var user sql.User
		if err := tx.First(&user, userId).Error; err != nil {
			return err
		}

		user.Coins += sql.RMB2Coins(claims.Credit)
		tx.Save(&user)
		return nil
	})

	return nil
}

func ListCDKeys() ([]CDKeyDetail, error) {
	ctx := context.Background()
	rds := redis.NewRedisCient().GetClient()
	cdkeys := make([]CDKeyDetail, 0)

	pattern := fmt.Sprintf("%s-*", cdkeyPrefix)
	keys, err := rds.Keys(ctx, pattern).Result()
	if err != nil {
		log.Errorf("CDKEY ListCDKEY ERROR: %s\n", err.Error())
		return nil, fmt.Errorf("内部错误")
	}

	for _, key := range keys {
		var claims CDKeyClaims

		result := rds.Get(ctx, key)
		if err := json.Unmarshal([]byte(result.Val()), &claims); err != nil {
			return nil, fmt.Errorf("未知错误")
		}

		detail := CDKeyDetail{
			CDKeyClaims: claims,
			Id:          key,
		}
		cdkeys = append(cdkeys, detail)
	}
	return cdkeys, nil
}
