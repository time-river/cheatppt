package apiv0

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"

	"cheatppt/config"
	"cheatppt/deps/revchatgpt"
)

func ChatgptWebAuth(c *gin.Context) {
	secretKey := config.GlobalCfg.ChatgptWeb.AuthSecretKey
	if len(secretKey) == 0 {
		c.Next()
		return
	}

	rawAuth := c.Request.Header.Get("Authorization")
	if len(rawAuth) == 0 {
		msg := &CommonResponse{
			Status:  "Unauthorized",
			Message: "Error: 无访问权限 | No access rights",
			Data:    "",
		}
		c.AbortWithStatusJSON(http.StatusUnauthorized, msg)
		return
	}

	auth := strings.Trim(strings.Replace(rawAuth, "Bearer ", "", 1), " \t\n")
	if len(auth) == 0 || auth != secretKey {
		msg := &CommonResponse{
			Status:  "Unauthorized",
			Message: "Please authenticate.",
			Data:    "",
		}
		c.AbortWithStatusJSON(http.StatusUnauthorized, msg)
		return
	}

	c.Next()
}

func fetchBalance(apiBase, apiToken string) string {
	if len(apiBase) == 0 || len(apiToken) == 0 {
		return "-"
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/dashboard/billing/credit_grants", apiBase), nil)
	if err != nil {
		return "-"
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiToken))

	resp, err := client.Do(req)
	if err != nil {
		return "-"
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "-"
	}

	var data struct {
		TotalAvailable float64 `json:"total_available"`
	}
	if err := json.Unmarshal(body, &data); err != nil {
		return "-"
	}

	return fmt.Sprintf("%.3f", data.TotalAvailable)
}

func ChatgptWebConfig(c *gin.Context) {
	var msg = &CommonResponse{
		Status: "Success",
	}

	conf := config.GlobalCfg.ChatgptWeb
	balance := fetchBalance(conf.OpenaiApiBase, conf.OpenaiApiToken)

	msg.Data = ModelConfig{
		ApiModel:     conf.ApiModel,
		ReversePorxy: conf.ReverseProxyUrl,
		TimeoutMs:    conf.TimeoutMs,
		SocksProxy:   conf.SocksProxy,
		HttpsProxy:   conf.HttpsProxy,
		Balance:      balance,
	}

	c.JSON(http.StatusOK, msg)
}

func ChatgptWebVerify(c *gin.Context) {
	var req VerificationRequest
	var msg = &CommonResponse{}

	if err := c.BindJSON(&req); err != nil {
		msg.Message = "Bad request"
		c.AbortWithStatusJSON(http.StatusBadRequest, msg)
		return
	}

	conf := config.GlobalCfg.ChatgptWeb
	if len(conf.AuthSecretKey) == 0 {
		msg.Status = "Fail"
		msg.Message = "Secret key is empty"
	} else if req.Token != conf.AuthSecretKey {
		msg.Status = "Fail"
		msg.Message = "密钥无效 | Secret key is invalid"
	} else {
		msg.Status = "Success"
		msg.Message = "Verify successfully"
	}

	c.JSON(http.StatusOK, msg)
}

func ChatgptWebSession(c *gin.Context) {
	var msg = &CommonResponse{
		Status: "Success",
	}

	conf := config.GlobalCfg.ChatgptWeb

	msg.Data = SessionResponse{
		Auth:  (len(conf.AuthSecretKey) != 0),
		Model: conf.ApiModel,
	}

	c.JSON(http.StatusOK, msg)
}

var onceConf sync.Once
var api *revchatgpt.ChatGPTUnofficialProxyAPI

func chatReplyProcess(params *RequestOptions) *CommonResponse {
	conf := config.GlobalCfg.ChatgptWeb

	if api == nil {
		onceConf.Do(func() {
			api = revchatgpt.NewChatGPTUnofficialProxyAPI(
				conf.OpenaiApiToken,
				conf.ReverseProxyUrl,
				conf.ApiModel,
			)
			api.EnableDebug(os.Stdout)
		})
	}

	/* only support unofficial chatgpt api currently */
	// TODO
	opts := revchatgpt.SendMessageBrowserOptions{
		ConversationId:  params.lastContext.ConversationId,
		ParentMessageId: params.lastContext.ParentMessageId,
		OnProgress:      params.process,
		TimeoutMs:       conf.TimeoutMs,
	}

	if _, err := api.SendMessage(params.message, opts); err != nil {
		var msg = &CommonResponse{
			Status: "Fail",
		}

		code := err.StatusCode
		if value, ok := ErrorCodeMessage[code]; ok {
			msg.Message = value
		} else if len(err.StatusText) > 0 {
			msg.Message = err.StatusText
		} else {
			msg.Message = "Please check the back-end console"
		}
		return msg
	}

	return nil
}

func ChatgptWebChatProcess(c *gin.Context) {
	var req RequestProps

	c.Header("Content-Type", "application/octet-stream")

	if err := c.BindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, nil)
		return
	}

	//firstChunk := true
	params := &RequestOptions{
		message:     req.Prompt,
		lastContext: req.Options,
		process: func(chat revchatgpt.ChatMessage) {
			c.JSON(http.StatusOK, &chat)
			/*
				if firstChunk {
					c.JSON(http.StatusOK, &chat)
				} else if data, err := json.Marshal(chat); err == nil {
					//c.String(http.StatusOK, fmt.Sprintf("\n%s", data))
					c.JSON(http.StatusOK, &chat)
				} // ignore others
			*/
			//firstChunk = false
		},
		systemMessage: req.SystemMessage,
	}

	msg := chatReplyProcess(params)
	c.JSON(http.StatusOK, msg)
}
