package userapiv1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kr/pretty"

	"cheatppt/api"
	"cheatppt/controller/code"
	"cheatppt/controller/user"
	"cheatppt/log"
	"cheatppt/utils"
)

type SignInReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Code     string `json:"code"`
}

type SignInRsp struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Token    string `json:"token"`
}

func SignIn(c *gin.Context) {
	var req SignInReq
	var rsp = &api.Response{Status: api.FAILURE}

	if err := c.BindJSON(&req); err != nil {
		rsp.Message = "非法请求"
		c.AbortWithStatusJSON(http.StatusOK, rsp)
		return
	}

	log.Debug(pretty.Sprint(req))

	token := req.Code
	ip := c.Request.Header.Get("CF-Connecting-IP")
	isHuman, err := code.Confirm(ip, token)
	if err != nil {
		rsp.Message = "内部错误"
		c.AbortWithStatusJSON(http.StatusOK, rsp)
		return
	} else if !isHuman {
		rsp.Message = "非法请求"
		c.AbortWithStatusJSON(http.StatusOK, rsp)
		return
	}

	if !utils.UsernameCheck(req.Username) || !utils.PasswordCheck(req.Password) {
		log.Debug("Invalid username or password")
		rsp.Message = "用户名或密码错误"
		c.AbortWithStatusJSON(http.StatusOK, rsp)
		return
	}

	data, err := user.SignIn(req.Username, req.Password)
	if err != nil {
		rsp.Message = err.Error()
		c.AbortWithStatusJSON(http.StatusOK, rsp)
		return
	}

	rsp.Status = api.SUCCESS
	rsp.Message = "登录成功"
	rsp.Data = SignInRsp{
		Username: req.Username,
		Email:    data.Email,
		Token:    data.Token,
	}
	c.JSON(http.StatusOK, rsp)
}
