package chatgptapiv1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kr/pretty"
	log "github.com/sirupsen/logrus"

	"cheatppt/api"
	revchatgpt3 "cheatppt/controller/chat/revchatgpt"
)

type DelAccountReq struct {
	Email string `json:"email"`
}

func DelAccount(c *gin.Context) {
	var req DelAccountReq
	var rsp = &api.Response{Status: api.FAILURE}

	if err := c.BindJSON(&req); err != nil {
		rsp.Message = "非法请求"
		c.AbortWithStatusJSON(http.StatusBadRequest, rsp)
		return
	}

	log.Debug(pretty.Sprint(req))

	revchatgpt3.DelAccount(&req.Email)
	rsp.Status = api.SUCCESS

	c.JSON(http.StatusOK, rsp)
}
