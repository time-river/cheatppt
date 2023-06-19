package revchatgpt3

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	http "github.com/bogdanfinn/fhttp"
	"github.com/kr/pretty"
	log "github.com/sirupsen/logrus"
)

type ProxyOpts struct {
	*ClientSession
	ReqBody *ConversationJSONBody
}

type Session struct {
	response *http.Response
	reader   *bufio.Reader

	StatusCode int
	Header     *http.Header
}

// TODO: timeout
func proxy(opts *ProxyOpts) (*Session, error) {
	client := getReqClient()

	// Remove _cfuvid cookie from session
	client.jar.SetCookies(&opts.ReqURL, []*http.Cookie{})

	var url string
	var err error
	var request_method string
	var request *http.Request
	var response *http.Response

	// https://chat.openai.com/backend-api/conversation
	url = "https://" + OpenAI_HOST + "/backend-api" + opts.ReqParamPath

	if opts.ReqURL.RawQuery != "" {
		url += "?" + opts.ReqURL.RawQuery
	}

	request_method = opts.ReqMethod

	if opts.ReqPath == "/api/conversation" {
		if strings.HasPrefix(opts.ReqBody.Model, "gpt-4") {
			token, err := get_arkose_token()
			if err == nil {
				arkose_token = token
			}
			opts.ReqBody.ArkoseToken = &arkose_token
		}
	}

	body_json, err := json.Marshal(opts.ReqBody)
	if err != nil {
		return nil, err
	}
	original_body := bytes.NewReader(body_json)
	request, _ = http.NewRequest(request_method, url, original_body)

	request.Header.Set("Host", ""+OpenAI_HOST+"")
	request.Header.Set("Origin", "https://"+OpenAI_HOST+"/chat")
	request.Header.Set("Connection", "keep-alive")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Keep-Alive", "timeout=360")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", opts.Token))
	request.Header.Set("sec-ch-ua", "\"Chromium\";v=\"112\", \"Brave\";v=\"112\", \"Not:A-Brand\";v=\"99\"")
	request.Header.Set("sec-ch-ua-mobile", "?0")
	request.Header.Set("sec-ch-ua-platform", "\"Linux\"")
	request.Header.Set("sec-fetch-dest", "empty")
	request.Header.Set("sec-fetch-mode", "cors")
	request.Header.Set("sec-fetch-site", "same-origin")
	request.Header.Set("sec-gpc", "1")
	request.Header.Set("user-agent", user_agent)

	if opts.RevSession.Puid != "" {
		request.Header.Set("cookie", "_puid="+opts.Puid+";")
	}

	log.Trace(pretty.Sprint(request.Header))

	response, err = client.client.Do(request)
	if err != nil {
		return nil, err
	}

	log.Trace(response.StatusCode)

	return &Session{
		response:   response,
		reader:     bufio.NewReader(response.Body),
		StatusCode: response.StatusCode,
		Header:     &response.Header,
	}, nil
}

func (s *Session) Read(buf []byte) (int, error) {
	return s.reader.Read(buf)
}

func (s *Session) Close() {
	s.response.Body.Close()
}
