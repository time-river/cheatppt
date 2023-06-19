package userapiv1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"cheatppt/api"
	"cheatppt/controller/user"
)

type ExgCDKeyReq struct {
	CDKey string `json:"cdkey"`
}

// exchange cd-key
func ExgCDkey(c *gin.Context) {
	var req ExgCDKeyReq
	var rsp = &api.Response{Status: api.FAILURE}

	if err := c.BindJSON(&req); err != nil {
		rsp.Message = "非法请求"
		c.AbortWithStatusJSON(http.StatusBadRequest, rsp)
		return
	}

	err := user.ExgCDkey(req.CDKey, c.GetInt(UserId))
	if err != nil {
		rsp.Message = err.Error()
		c.JSON(http.StatusOK, rsp)
		return
	}

	rsp.Status = api.SUCCESS
	c.JSON(http.StatusOK, rsp)
}
