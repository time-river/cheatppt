package userapiv1

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"cheatppt/api"
	"cheatppt/controller/user"
)

const sessionName = "sessionId"
const tokenName = "tokenId"

func SessionGuard(c *gin.Context) {
	session := c.Request.Header.Get(sessionName)
	if len(session) == 0 {
		msg := &api.Response{
			Status:  api.UNAUTHORIZED,
			Message: "Error: 无访问权限 | No access rights",
		}
		c.AbortWithStatusJSON(http.StatusUnauthorized, msg)
		return
	}

	token, permit := user.ValidSession(session)
	if !permit {
		msg := &api.Response{
			Status:  api.UNAUTHORIZED,
			Message: "Error: 无访问权限 | No access rights",
		}
		c.AbortWithStatusJSON(http.StatusUnauthorized, msg)
		return
	}

	c.Set(tokenName, token)

	c.Next()
}

func TokenGuard(c *gin.Context) {
	auth := c.Request.Header.Get("Authorization")
	if len(auth) == 0 {
		msg := &api.Response{
			Status:  api.UNAUTHORIZED,
			Message: "Error: 无访问权限 | No access rights",
		}
		c.AbortWithStatusJSON(http.StatusUnauthorized, msg)
		return
	}

	token := strings.TrimSpace(strings.Replace(auth, "Bearer ", "", 1))
	if permit, err := user.ValidToken(token); err != nil {
		msg := &api.Response{
			Status:  api.UNAUTHORIZED,
			Message: "Error: 内部错误 | Internal Error",
		}
		c.AbortWithStatusJSON(http.StatusUnauthorized, msg)
		return
	} else if !permit {
		msg := &api.Response{
			Status:  api.UNAUTHORIZED,
			Message: "Error: 无访问权限 | No access rights",
		}
		c.AbortWithStatusJSON(http.StatusUnauthorized, msg)
		return
	}

	c.Set(tokenName, token)

	c.Next()
}
