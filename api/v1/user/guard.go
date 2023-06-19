package userapiv1

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kr/pretty"
	log "github.com/sirupsen/logrus"

	"cheatppt/api"
	"cheatppt/controller/user"
)

const SessionName = "X-SessionId"
const TokenName = "X-TokenId"
const UserId = "X-UserId"
const UserLevel = "X-UserLevel"

func sessionParse(c *gin.Context) *user.Claims {
	session := c.Request.Header.Get(SessionName)
	if len(session) == 0 {
		return nil
	}

	token, permit := user.ValidSession(session)
	if !permit {
		return nil
	}

	c.Set(TokenName, token)
	return user.TokenParse(token)
}

func SessionGuard(c *gin.Context) {
	claims := sessionParse(c)
	if claims == nil {
		msg := &api.Response{
			Status:  api.UNAUTHORIZED,
			Message: "Error: 无访问权限 | No access rights",
		}
		c.AbortWithStatusJSON(http.StatusUnauthorized, msg)
		return
	}

	log.Trace(pretty.Sprint(claims))

	c.Set(UserId, claims.UserID)
	c.Set(UserLevel, claims.UserLevel)
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

	c.Set(TokenName, token)

	c.Next()
}

func AdminGuard(c *gin.Context) {
	claims := sessionParse(c)
	if claims == nil || claims.UserLevel != 0 {
		msg := &api.Response{
			Status:  api.UNAUTHORIZED,
			Message: "Error: 无访问权限 | No access rights",
		}
		c.AbortWithStatusJSON(http.StatusUnauthorized, msg)
		return
	}

	log.Trace(pretty.Sprint(claims))

	c.Set(UserId, claims.UserID)
	c.Set(UserLevel, claims.UserLevel)

	c.Next()
}
