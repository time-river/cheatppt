package userapiv1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kr/pretty"
	log "github.com/sirupsen/logrus"

	"cheatppt/api"
	"cheatppt/controller/user"
	"cheatppt/utils"
)

type ResetReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Code     string `json:"code"`
}

func Reset(c *gin.Context) {
	var req ResetReq
	var rsp = &api.Response{Status: api.FAILURE}

	if err := c.BindJSON(&req); err != nil {
		rsp.Message = "非法请求"
		c.AbortWithStatusJSON(http.StatusBadRequest, rsp)
		return
	}

	log.Debug(pretty.Sprint(req))

	if !utils.UsernameCheck(req.Username) {
		rsp.Message = "用户无效"
		c.JSON(http.StatusOK, rsp)
		return
	}

	valid, err := user.ValidateResetCode(req.Username, req.Code)
	if err != nil {
		rsp.Message = err.Error()
		c.JSON(http.StatusOK, rsp)
		return
	} else if !valid {
		rsp.Message = "验证码无效"
		c.JSON(http.StatusOK, rsp)
		return
	}

	if err := user.ResetPassword(req.Username, req.Password); err != nil {
		rsp.Message = err.Error()
	} else {
		rsp.Status = api.SUCCESS
		rsp.Message = "密码重置成功"
	}

	c.JSON(http.StatusOK, rsp)
}
