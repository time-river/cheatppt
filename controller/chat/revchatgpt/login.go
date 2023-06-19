package revchatgpt3

import (
	"fmt"

	"github.com/acheong08/OpenAIAuth/auth"
	log "github.com/sirupsen/logrus"
)

type LoginOpts struct {
	Email    string `json:"openaiEmail"`
	Password string `json:"openaiPassword"`
}

func Login(opts *LoginOpts) (*RevSession, error) {
	var rsp RevSession

	authenticator := auth.NewAuthenticator(opts.Email, opts.Password, getHttpProxy())
	err := authenticator.Begin()
	if err != nil {
		return nil, fmt.Errorf("StatusCode: %d Details: %s", err.StatusCode, err.Details)
	}

	rsp.Token = authenticator.GetAccessToken()

	puid, err := authenticator.GetPUID()
	if err == nil {
		rsp.Puid = puid
	}

	log.Debug("login success")

	return &rsp, nil
}
