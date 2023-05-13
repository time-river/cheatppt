package chatgptapiv1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"cheatppt/api"
	"cheatppt/controller/chat/chatgpt"
)

type TokenReq struct {
	Token string `json:"token"`
}

func RefreshToken(c *gin.Context) {
	var req TokenReq
	var rsp = &api.Response{Status: api.FAILURE}

	if err := c.BindJSON(&req); err != nil {
		rsp.Message = "非法参数"
		c.AbortWithStatusJSON(http.StatusOK, rsp)
		return
	}

	if !chatgpt.RefreshToken(req.Token) {
		rsp.Message = "更新失败"
	} else {
		rsp.Message = "更新成功"
		rsp.Status = api.SUCCESS
	}

	c.JSON(http.StatusOK, rsp)
}
