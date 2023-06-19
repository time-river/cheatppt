package chatgptapiv1

import (
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kr/pretty"
	log "github.com/sirupsen/logrus"

	userapiv1 "cheatppt/api/v1/user"
	revchatgpt3 "cheatppt/controller/chat/revchatgpt"
	"cheatppt/controller/token"
)

type ConversationReq struct {
	Model      string                          `json:"model"`
	Prompt     string                          `json:"prompt"`
	Options    revchatgpt3.ConversationOptions `json:"options"`
	Timeout    uint                            `json:"timeout,omitempty"`
	RevSession *revchatgpt3.RevSession         `json:"revSession,omitempty"`
}

type Message struct {
	PartContent string                          `json:"partContent"`
	Options     revchatgpt3.ConversationOptions `json:"options"`
	EndTurn     string                          `json:"endTurn"`
	Usage       token.Usage                     `json:"usage"`
}

type ConversationRsp struct {
	Message
	Error error `json:"error,omitempty"`
}

func Conversation(c *gin.Context) {
	var req = ConversationReq{Timeout: 0}
	var rsp = &ConversationRsp{}

	if err := c.BindJSON(&req); err != nil {
		rsp.Error = &revchatgpt3.APIError{
			Message: "invalid_request_error",
			Type:    "invalid_request_params",
		}
		c.JSON(http.StatusOK, rsp)
		return
	}

	log.Debug(pretty.Sprint(req))

	opts := &revchatgpt3.SendMessageBrowserOptions{
		ConversationId:  req.Options.ConversationId,
		ParentMessageId: req.Options.ParentMessageId,
		MessageId:       req.Options.MessageId,
		Model:           req.Model,
		Prompt:          req.Prompt,
	}
	clientSession := &revchatgpt3.ClientSession{
		ClientOpts: revchatgpt3.ClientOpts{
			UserId: c.GetInt(userapiv1.UserId),
		},
		RevSession: req.RevSession,
	}
	session, err := revchatgpt3.NewConversation(opts, clientSession)
	if err != nil {
		rsp.Error = err
		c.JSON(http.StatusOK, rsp)
		return
	}
	defer session.Close()

	// for http streaming
	c.Writer.Header().Set("Cache-Control", "no-cache")

	for k, v := range *session.Header {
		if strings.ToLower(k) == "content-encoding" {
			continue
		}
		c.Header(k, v[0])
	}
	c.Status(session.StatusCode)

	if session.StatusCode == http.StatusOK {
		c.Stream(func(w io.Writer) bool {
			data, err := session.Recv()
			if err != nil {
				// any error will result of sending `[DONE]`
				c.SSEvent("", "[DONE]")
				return false
			}

			rsp = &ConversationRsp{
				Message: Message{
					PartContent: data.Text,
					Options: revchatgpt3.ConversationOptions{
						MessageId:       &data.MessageId,
						ParentMessageId: &data.ParentMessageId,
						ConversationId:  &data.ConversationId,
					},
				},
			}

			log.Debug(pretty.Sprint(rsp))

			c.SSEvent("", rsp)
			return true
		})
	} else {
		buf := make([]byte, 4096)

		for {
			n, err := session.Read(buf)
			if n > 0 {
				_, writeErr := c.Writer.Write(buf[:n])
				if writeErr != nil {
					break
				}
				c.Writer.Flush() // flush buffer to make sure the data is sent to client in time.
			}
			if err == io.EOF {
				break
			} else if err != nil {
				break
			}
		}
	}
}
