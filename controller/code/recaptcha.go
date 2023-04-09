// Package recaptcha handles reCaptcha (http://www.google.com/recaptcha) form submissions
//
// This package is designed to be called from within an HTTP server or web framework
// which offers reCaptcha form inputs and requires them to be evaluated for correctness
//
// Edit the recaptchaPrivateKey constant before building and using
package code

import (
	"cheatppt/config"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
	"time"

	"cheatppt/log"

	"github.com/kr/pretty"
)

const urlSuffix = "/recaptcha/api/siteverify"

type reCAPTCHA struct {
	server string
	secret string
}

type recaptchaResponse struct {
	Success     bool      `json:"success"`
	Score       float64   `json:"score"`
	Action      string    `json:"action"`
	ChallengeTS time.Time `json:"challenge_ts"`
	Hostname    string    `json:"hostname"`
	ErrorCodes  []string  `json:"error-codes"`
}

// check uses the client ip address, the challenge code from the reCaptcha form,
// and the client's response input to that challenge to determine whether or not
// the client answered the reCaptcha input question correctly.
// It returns a boolean value indicating whether or not the client answered correctly.
func (rc *reCAPTCHA) check(remoteIp, response string) (r recaptchaResponse, err error) {
	resp, err := http.PostForm(rc.server,
		url.Values{"secret": {rc.secret}, "remoteip": {remoteIp}, "response": {response}})
	if err != nil {
		log.Errorf("Post error: %s\n", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("Read error: could not read body: %s", err)
		return
	}
	err = json.Unmarshal(body, &r)
	if err != nil {
		log.Errorf("Read error: got invalid JSON: %s", err)
		return
	}
	return
}

var rc *reCAPTCHA
var onceConf sync.Once

// ReCaptcha v2
func Confirm(remotIp, response string) (bool, error) {
	if len(response) == 0 {
		return false, fmt.Errorf("验证码无效")
	}

	if rc == nil {
		onceConf.Do(func() {
			url := ""
			if config.Code.Host == "" {
				url = fmt.Sprintf("https://%s%s", "www.google.com", urlSuffix)
			} else {
				url = fmt.Sprintf("https://%s%s", config.Code.Host, urlSuffix)
			}

			rc = &reCAPTCHA{
				server: url,
				secret: config.Code.Secret,
			}
		})
	}

	rsp, err := rc.check(remotIp, response)
	if err != nil {
		return false, fmt.Errorf("内部错误")
	} else {
		log.Debug(pretty.Sprint(rsp))
		return rsp.Success, nil
	}
}
