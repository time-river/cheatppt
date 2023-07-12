package userapiv1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"cheatppt/api"
	user "cheatppt/controller/user/session"
)

func CreateChatSession(c *gin.Context) {
	var rsp = &api.Response{Status: api.FAILURE}

	userId := c.GetInt(UserId)
	data, err := user.CreateSession(userId)
	if err != nil {
		rsp.Message = err.Error()
	} else {
		rsp = &api.Response{
			Status: api.SUCCESS,
			Data:   data,
		}
	}

	c.JSON(http.StatusOK, rsp)
}
