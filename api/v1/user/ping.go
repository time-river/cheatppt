package userapiv1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"cheatppt/api"
	"cheatppt/controller/model"
	"cheatppt/controller/user"
)

type PingRsp struct {
	SessionId string        `json:"sessionId"`
	Models    []model.Model `json:"models"`
}

func Ping(c *gin.Context) {
	var session string
	var rsp = &api.Response{Status: api.FAILURE}

	raw, exist := c.Get(SessionName)
	if exist {
		session = raw.(string)
	} else {
		raw, exist := c.Get(TokenName)
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

	data, err := model.ListAvailable()
	if err != nil {
		rsp.Message = err.Error()
		c.AbortWithStatusJSON(http.StatusInternalServerError, rsp)
		return
	}

	rsp.Status = api.SUCCESS
	rsp.Data = PingRsp{
		SessionId: session,
		Models:    data,
	}
	c.JSON(http.StatusOK, rsp)
}
