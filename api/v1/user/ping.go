package userapiv1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"cheatppt/api"
	"cheatppt/controller/user"
)

type PingRsp struct {
	SessionId string `json:"sessionId"`
}

func Ping(c *gin.Context) {
	var session string
	var rsp = &api.Response{Status: api.FAILURE}

	raw, exist := c.Get(sessionName)
	if exist {
		session = raw.(string)
	} else {
		raw, exist := c.Get(tokenName)
		if !exist {
			rsp.Message = "内部错误"
			c.AbortWithStatusJSON(http.StatusInternalServerError, rsp)
			return
		}

		token := raw.(string)
		val := user.NewSession(token)
		if val == nil {
			rsp.Message = "内部错误"
			c.AbortWithStatusJSON(http.StatusInternalServerError, rsp)
			return
		}

		session = *val
	}

	rsp.Status = api.SUCCESS
	rsp.Data = PingRsp{
		SessionId: session,
	}
	c.JSON(http.StatusOK, rsp)
}
