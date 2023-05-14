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
	defer session.Close()

	// for http streaming
	c.Header("Cache-Control", "no-cache")
	// the response is `plain/text` because of the frontend
	// once set `text/event-stream`, the display is server-sent event,
	// therefore comment it.
	//c.Header("Content-Type", "text/event-stream")

	c.Stream(func(w io.Writer) bool {
		text, err := session.Recv()
		if err != nil && err == io.EOF {
			return false
		} else if err != nil {
			return false
		}

		w.Write([]byte(text))
		return true
	})
}
