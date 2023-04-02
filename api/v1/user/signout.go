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

	raw := c.Request.Header.Get("Token")
	token := strings.TrimSpace(strings.Replace(raw, "Bearer ", "", 1))
	user.Logout(token)

	msg.Status = api.SUCCESS
	msg.Message = "退出成功"
	c.JSON(http.StatusOK, msg)
}
