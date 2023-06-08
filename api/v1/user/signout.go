package userapiv1

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"cheatppt/api"
	"cheatppt/controller/user"
)

func SignOut(c *gin.Context) {
	var msg = &api.Response{Status: api.FAILURE}

	auth := c.Request.Header.Get("Authorization")
	token := strings.TrimSpace(strings.Replace(auth, "Bearer ", "", 1))
	user.SignOut(token)

	msg.Status = api.SUCCESS
	msg.Message = "退出成功"
	c.JSON(http.StatusOK, msg)
}
