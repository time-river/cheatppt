package chatgptapiv1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kr/pretty"
	log "github.com/sirupsen/logrus"

	"cheatppt/api"
	revchatgpt3 "cheatppt/controller/chat/revchatgpt"
)

func Login(c *gin.Context) {
	var req revchatgpt3.LoginOpts
	var rsp = &api.Response{Status: api.FAILURE}

	if err := c.BindJSON(&req); err != nil {
		rsp.Message = "非法请求"
		c.AbortWithStatusJSON(http.StatusBadRequest, rsp)
		return
	}

	log.Debug(pretty.Sprint(req))

	data, err := revchatgpt3.Login(&req)
	if err != nil {
		rsp.Message = err.Error()
	} else {
		rsp = &api.Response{
			Data:   data,
			Status: api.SUCCESS,
		}
	}

	c.JSON(http.StatusOK, rsp)
}
