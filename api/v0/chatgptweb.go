package apiv0

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"cheatppt/config"
	"cheatppt/contrib/revchatgpt2"
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

var clientConfig *revchatgpt2.ClientConfig
var revchatgptClient *revchatgpt2.Client
var onceConf sync.Once

func initConfig() {
	conf := config.GlobalCfg.ChatgptWeb

	if revchatgptClient == nil {
		onceConf.Do(func() {
			config := revchatgpt2.DefaultConfig(
				conf.OpenaiApiToken,
				conf.ReverseProxyUrl,
				conf.ApiModel,
			)
			clientConfig = &config
			revchatgptClient = revchatgpt2.NewClientWithConfig(config)
		})
	}
}

func parseErrorMsg(err error) string {
	fmt.Printf("Error: %s\n", err.Error())

	if e, ok := err.(*revchatgpt2.ChatGPTError); ok {
		if value, ok := ErrorCodeMessage[e.StatusCode]; ok {
			return value
		} else {
			return "ChatGPT unknow error"
		}
	} else {
		return "Internal Error"
	}
}

func ChatgptWebChatProcess(c *gin.Context) {
	var req RequestProps

	if err := c.BindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, nil)
		return
	}

	initConfig()

	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Content-Type", "application/octet-stream; charset=utf-8")

	/* only support unofficial chatgpt api currently */
	opts := revchatgpt2.SendMessageBrowserOptions{
		ConversationId:  req.Options.ConversationId,
		ParentMessageId: req.Options.ParentMessageId,
		Model:           clientConfig.Model,
	}
	text := req.Prompt

	messages := make(chan []byte)
	var msg *CommonResponse

	go func() {
		parent := context.Background()
		ctx, cancel := context.WithTimeout(parent, 60*time.Second)
		defer cancel()

		stream, err := revchatgptClient.CreateChatCompletionStream(ctx, text, opts)
		if err != nil && err.Error() == context.DeadlineExceeded.Error() {
			msg = &CommonResponse{
				Status:  "Fail",
				Message: "ChatGPT Server Request Timeout",
			}
		} else if err != nil {
			msg = &CommonResponse{
				Status:  "Fail",
				Message: parseErrorMsg(err),
			}
		} else {
			defer stream.Close()

			for {
				data, err := stream.Recv()
				if err != nil && err == io.EOF {
					// don't send anything
					break
				} else if err != nil {
					msg = &CommonResponse{
						Status:  "Fail",
						Message: parseErrorMsg(err),
					}
					break
				} else if chunk, err := data.Marshal(); err == nil {
					messages <- chunk
				}
			}
		}

		close(messages)
	}()

	c.Stream(func(w io.Writer) bool {
		keep := false

		if message, ok := <-messages; ok {
			message = append(message, '\n')
			w.Write(message)
			keep = true
		}

		return keep
	})

	if msg != nil {
		// send error message
		c.JSON(http.StatusOK, msg)
	}
}
