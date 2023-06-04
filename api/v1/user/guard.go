package userapiv1

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"cheatppt/api"
	"cheatppt/controller/user"
)

func Guard(c *gin.Context) {
	auth := c.Request.Header.Get("Authorization")
	if len(auth) == 0 {
		msg := &api.Response{
			Status:  api.UNAUTHORIZED,
			Message: "Error: 无访问权限 | No access rights",
			Data:    "",
		}
		c.AbortWithStatusJSON(http.StatusUnauthorized, msg)
		return
	}

	token := strings.Trim(strings.Replace(auth, "Bearer ", "", 1), " \t\n")
	if !user.ValidToken(token) {
		msg := &api.Response{
			Status:  api.UNAUTHORIZED,
			Message: "Error: 无访问权限 | No access rights",
			Data:    "",
		}
		c.AbortWithStatusJSON(http.StatusUnauthorized, msg)
		return
	}

	c.Next()
}
