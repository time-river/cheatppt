package chatgptapiv1

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kr/pretty"

	"cheatppt/controller/chat/chatgpt"
	"cheatppt/controller/chat/model"
	"cheatppt/log"
)

type ChatReq struct {
	Model   string                      `json:"model"`
	Prompt  string                      `json:"prompt"`
	Options chatgpt.ConversationOptions `json:"options"`
	Timeout uint                        `json:"timeout,omitempty"`
}

type APIError struct {
	Code           any     `json:"code,omitempty"`
	Message        string  `json:"message"`
	Param          *string `json:"param,omitempty"`
	Type           string  `json:"type"`
	HTTPStatusCode int     `json:"-"`
}

type ChatRsp struct {
	chatgpt.ChatRsp
	Error *APIError `json:"error,omitempty"`
}

func Chat(c *gin.Context) {
	var req = ChatReq{Timeout: 0}
	var rsp = &ChatRsp{}

	if err := c.BindJSON(&req); err != nil {
		rsp.Error = &APIError{
			Message: "invalid_request_error",
			Type:    "invalid_request_params",
		}
		c.JSON(http.StatusOK, rsp)
		return
	}

	log.Debug(pretty.Sprint(req))

	if !model.Allow(req.Model) {
		rsp.Error = &APIError{
			Type: "invalid_request_error",
			Code: "invalid_model",
		}
		c.JSON(http.StatusOK, rsp)
		return
	}

	if req.Timeout == 0 {
		req.Timeout = 60
	} else if req.Timeout > 180 {
		req.Timeout = 180
	}

	parent := context.Background()
	ctx, cancel := context.WithTimeout(parent, time.Duration(req.Timeout)*time.Second)
	defer cancel()

	opts := chatgpt.ChatOpts{
		Ctx:     &ctx,
		Prompt:  req.Prompt,
		Model:   req.Model,
		Options: req.Options,
	}
	session, err := chatgpt.NewChat(&opts)
	if err != nil {
		rsp.Error = &APIError{
			Message: fmt.Sprintf("创建对话失败 | %s", err.Error()),
		}
		c.JSON(http.StatusOK, rsp)
		return
	}
	defer session.Close()

	// for http streaming
	c.Writer.Header().Set("Cache-Control", "no-cache")

	c.Stream(func(w io.Writer) bool {
		data, err := session.Recv()
		if err != nil {
			// any error will result of sending `[DONE]`
			c.SSEvent("", "[DONE]")
			return false
		}

		log.Debug(pretty.Sprint(data))

		c.SSEvent("", data)
		return true
	})
}
