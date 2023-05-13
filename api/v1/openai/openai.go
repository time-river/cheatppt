package openaiapiv1

import (
	"context"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kr/pretty"

	"cheatppt/controller/chat/model"
	"cheatppt/controller/chat/openai"
	"cheatppt/log"
)

func Chat(c *gin.Context) {
	var req openai.ChatReq

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

	if !model.Allow(req.Model) {
		rsp := &openai.ChatErrRsp{
			Error: &openai.ChatAPIErr{
				Type: "invalid_request_error",
				Code: "invalid_model",
			},
		}
		c.JSON(http.StatusOK, rsp)
		return
	}

	ctx := context.Background()
	opts := openai.ChatOpts{
		Ctx: ctx,
		Req: openai.ChatReq(req),
	}

	session, err := openai.NewChat(&opts)
	if err != nil {
		c.JSON(err.Error.HTTPStatusCode, err)
		return
	}

	// for http streaming
	c.Header("Cache-Control", "no-cache")
	// once set `text/event-stream`, the display isn't stream,
	// therefore comment it
	//c.Header("Content-Type", "text/event-stream")

	c.Stream(func(w io.Writer) bool {
		data, err := session.Recv()
		if err != nil && err == io.EOF {
			return false
		} else if err != nil {
			log.Warnf("Recv Msg ERROR: %s\n", err.Error())
			return false
		}

		log.Trace(pretty.Sprint(data))

		w.Write([]byte(data.Choices[0].Delta.Content))
		return true
	})
}
