package auth

import (
	"bytes"
	"errors"
	"sync"

	"cheatppt/config"
	msg "cheatppt/model/http"
	"cheatppt/model/redis"
	"cheatppt/model/sql"
	"cheatppt/model/sql/db"
)

type Auth struct {
	digest    *Digest
	token     *Token
	reCAPTCHA *ReCAPTCHA
}

var onceConf sync.Once
var auth *Auth

func AuthCtxCreate() *Auth {
	if auth == nil {
		onceConf.Do(func() {
			serverConf := config.GlobalCfg.Server
			captchaConf := config.GlobalCfg.ReCAPTCHA
			auth = &Auth{
				digest: &Digest{
					salt: serverConf.Secret,
				},
				token: &Token{
					secret: config.GlobalKey[:],
				},
				reCAPTCHA: &ReCAPTCHA{
					secret: captchaConf.Secret,
					server: captchaConf.Server,
				},
			}
		})
	}
	return auth
}

func (l *Auth) reCAPTCHAVerify(clientIP *string, response *string) error {

	r, err := l.reCAPTCHA.check(clientIP, response)
	if err != nil {
		return errors.New("Internal error")
	}

	// TODO
	return nil

	if !r.Success {
		return errors.New(r.ErrorCodes[0])
	} else if r.Score < 0.8 { // TODO: score configurize?
		return errors.New("Are you robots?")
	}

	return nil
}

func (l *Auth) UserRegister(req *msg.RegisterRequest, clientIP *string) error {
	if err := l.reCAPTCHAVerify(clientIP, &req.Recaptcha); err != nil {
		return err
	}

	sql := sql.SQLCtxCreate()
	user := db.User{
		Username:      req.Username,
		Email:         req.Email,
		Password:      l.digest.digest(&req.Password),
		Level:         100,
		EmailVerified: false,
	}
	if err := sql.UserCreate(&user); err != nil {
		return errors.New("User has been registered")
	}

	return nil
}

func (l *Auth) passwordVerify(username *string, password *string) error {
	digest := l.digest.digest(password)

	sql := sql.SQLCtxCreate()
	found, err := sql.PasswdLookup(username)
	if err != nil {
		return errors.New("Bad username or password")
	}

	if !bytes.Equal(digest[:], found[:]) {
		return errors.New("Bad username or password")
	} else {
		return nil
	}
}

func (l *Auth) UserLogin(req *msg.LoginRequest, clientIP *string) (*string, error) {
	if err := l.reCAPTCHAVerify(clientIP, &req.Recaptcha); err != nil {
		return nil, err
	}

	if err := l.passwordVerify(&req.Username, &req.Password); err != nil {
		return nil, err
	}

	token, err := l.token.generate(&req.Username)
	if err != nil {
		return nil, errors.New("Internal error")
	}

	rds := redis.RedisCtxCreate()
	if err := rds.TokenLease(*token, req.Username); err != nil {
		return nil, errors.New("Internal error")
	}

	return token, nil
}

func (l *Auth) UserAuthorized(tokenString *string, username *string) error {
	rds := redis.RedisCtxCreate()
	result, err := rds.TokenVerify(*tokenString)
	if err != nil || *username != *result {
		return errors.New("No authorization")
	}

	return nil
}

func (l *Auth) UserLogout(token *string) {
	rds := redis.RedisCtxCreate()
	rds.TokenRevoke(*token)
}
