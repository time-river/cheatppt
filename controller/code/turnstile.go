package code

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/kr/pretty"
	log "github.com/sirupsen/logrus"

	"cheatppt/config"
)

type response struct {
	// Success indicates if the challenge was passed
	Success bool `json:"success"`
	// ChallengeTs is the timestamp of the captcha
	ChallengeTs string `json:"challenge_ts"`
	// Hostname is the hostname of the passed captcha
	Hostname string `json:"hostname"`
	// ErrorCodes contains error codes returned by hCaptcha (optional)
	ErrorCodes []string `json:"error-codes"`
	// Action  is the customer widget identifier passed to the widget on the client side
	Action string `json:"action"`
	// CData is the customer data passed to the widget on the client side
	CData string `json:"cdata"`
}

const server = "https://challenges.cloudflare.com/turnstile/v0/siteverify"

func check(ip, token string) (*response, error) {
	secret := config.Code.Secret

	data := url.Values{
		"secret":   {secret},
		"response": {token},
		"remoteip": {ip},
	}

	resp, err := http.PostForm(server, data)
	if err != nil {
		log.Warnf("Turnstile POST ERROR: %s\n", err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Warnf("Turnstile read ERROR: %s\n", err.Error())
		return nil, err
	}

	var r response
	err = json.Unmarshal(body, &r)
	if err != nil {
		log.Warnf("Turnstile ERROR: get invalid JSON(`%s`): %s\n", pretty.Sprint(body), err.Error())
		return nil, err
	}
	return &r, nil
}

func Confirm(ip, token string) (bool, error) {
	log.Tracef("Turnstile ip: %s token: %s\n", ip, token)

	if len(token) == 0 {
		return false, fmt.Errorf("验证码无效")
	}

	rsp, err := check(ip, token)
	if err != nil {
		return false, fmt.Errorf("内部错误")
	} else {
		log.Trace(pretty.Sprint(rsp))
		return rsp.Success, nil
	}
}
