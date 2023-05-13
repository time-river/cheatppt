package chatgptapiv1

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"cheatppt/controller/chat/chatgpt"
	"cheatppt/controller/chat/model"
	"cheatppt/log"
)

type ConversationRequest struct {
	ConversationId  *string `json:"conversationId,omitempty"`
	ParentMessageId *string `json:"parentMessageId,omitempty"`
}

type ChatReq struct {
	Model   string              `json:"model"`
	Prompt  string              `json:"prompt"`
	Options ConversationRequest `json:"options"`
	Timeout uint                `json:"timeout,omitempty"`
}

type ChatRsp struct {
	Error   string              `json:"error,omitempty"`
	Text    string              `json:"text,omitempty"`
	Role    string              `json:"role,omitempty"`
	Options ConversationRequest `json:"options"`
}

func Chat(c *gin.Context) {
	var req = ChatReq{Timeout: 0}
	var rsp = &ChatRsp{}

	if err := c.BindJSON(&req); err != nil {
		rsp.Error = "非法参数"
		c.AbortWithStatusJSON(http.StatusOK, rsp)
		return
	}

	if !model.Allow(req.Model) {
		rsp.Error = "模型不存在"
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
		Ctx:             ctx,
		Text:            req.Prompt,
		Model:           req.Model,
		ConversationId:  req.Options.ConversationId,
		ParentMessageId: req.Options.ParentMessageId,
	}
	stream, err := chatgpt.NewChat(&opts)
	if err != nil {
		rsp.Error = fmt.Sprintf("创建对话失败 | %s", err.Error())
		c.JSON(http.StatusOK, rsp)
	}

	// for http streaming
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Content-Type", "application/octet-stream; charset=utf-8")

	for {
		rawData, err := stream.Recv()
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			log.Errorf("Stream.Recv ERROR: %s\n", err.Error())
			break
		}

		var data = &ChatRsp{
			Text: rawData.Text,
			Role: rawData.Role,
			Options: ConversationRequest{
				ConversationId:  &rawData.ConversationId,
				ParentMessageId: &rawData.ParentMessageId,
			},
		}

		msg, err := json.Marshal(data)
		if err != nil {
			log.Warnf("json.Marshal ERROR: %s\n", err.Error())
			break
		}

		msg = append(msg, '\n')
		c.JSON(http.StatusOK, msg)
	}
}
