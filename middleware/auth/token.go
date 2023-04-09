package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"cheatppt/api"
	"cheatppt/model/redis"
)

func authorize(tokenString string) bool {
	rds := redis.NewRedisCient()
	_, err := rds.TokenVerify(tokenString)
	return err == nil
}

func TokenVerify(c *gin.Context) {
	rawToken := c.Request.Header.Get("Authorization")
	if len(rawToken) == 0 {
		msg := &api.Response{
			Message: "Error: 无访问权限 | No access rights",
		}
		c.AbortWithStatusJSON(http.StatusUnauthorized, msg)
		return
	}

	token := strings.TrimSpace(strings.Replace(rawToken, "Bearer ", "", 1))
	if !authorize(token) {
		msg := &api.Response{
			Message: "Error: 无访问权限 | No access rights",
		}
		c.AbortWithStatusJSON(http.StatusUnauthorized, msg)
		return
	}

	c.Next()
}
