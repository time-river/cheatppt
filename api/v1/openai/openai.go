package openaiapiv1

import (
	"context"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kr/pretty"
	openaiapi "github.com/sashabaranov/go-openai"
	log "github.com/sirupsen/logrus"

	userapiv1 "cheatppt/api/v1/user"
	"cheatppt/controller/chat/openai"
)

type ChatReq struct {
	openaiapi.ChatCompletionRequest

	SessionId string `json:"sessionId"`
}

func Chat(c *gin.Context) {
	var req ChatReq

	if err := c.BindJSON(&req); err != nil {
		rsp := &openai.ChatErrRsp{
			Error: &openai.ChatAPIErr{
				Type: "invalid_request_error",
				Code: "invalid_request_params",
			},
		}
		c.JSON(http.StatusOK, rsp)
		return
	}

	log.Debug(pretty.Sprint(req))

	ctx := context.Background()
	opts := openai.ChatOpts{
		Ctx:    ctx,
		UserId: c.GetInt(userapiv1.UserId),
		Model:  req.Model,
		Req:    req.ChatCompletionRequest,
	}

	if id, err := uuid.Parse(req.SessionId); err != nil {
		rsp := &openai.ChatErrRsp{
			Error: &openai.ChatAPIErr{
				Type: "invalid_request_error",
				Code: "invalid uuid",
			},
		}
		c.JSON(http.StatusOK, rsp)
		return
	} else {
		opts.SessionId = id
	}

	session, err := openai.NewChat(&opts)
	if err != nil {
		c.JSON(err.Error.HTTPStatusCode, err)
		return
	}
	defer session.Close()

	// for http streaming
	c.Header("Cache-Control", "no-cache")

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
