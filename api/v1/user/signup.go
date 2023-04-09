package userapiv1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"cheatppt/api"
	"cheatppt/controller/user"
	"cheatppt/log"
	"cheatppt/utils"
)

type SignUpReq struct {
	Username       string `json:"username"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	InvitationCdoe string `json:"invitationCode"`
	Code           string `json:"code"`
}

func SignUp(c *gin.Context) {
	var req SignUpReq
	var rsp = &api.Response{Status: api.FAILURE}

	if err := c.BindJSON(&req); err != nil {
		rsp.Message = "非法请求"
		c.AbortWithStatusJSON(http.StatusOK, rsp)
		return
	}

	// don't print password
	log.Debugf(`username: "%s", email: "%s", invitationCode: "%s" code: "%s"\n`, req.Username, req.Email, req.InvitationCdoe, req.Code)

	if !utils.UsernameCheck(req.Username) {
		rsp.Message = "非法的用户名"
		c.AbortWithStatusJSON(http.StatusOK, rsp)
		return
	} else if !utils.EmailCheck(req.Email) {
		rsp.Message = "非法的邮箱"
		c.AbortWithStatusJSON(http.StatusOK, rsp)
		return
	} else if !utils.PasswordCheck(req.Password) {
		rsp.Message = "非法的密码"
		c.AbortWithStatusJSON(http.StatusOK, rsp)
		return
	}

	if valid, err := user.ValidateSignUpCode(req.Username, req.Email, req.Code); err != nil {
		rsp.Message = err.Error()
	} else if !valid {
		rsp.Message = "验证码无效"
	}
	if err := user.SignUp(req.Username, req.Password); err != nil {
		rsp.Message = err.Error()
	} else {
		rsp.Status = api.SUCCESS
		rsp.Message = "注册成功，请登录"
	}

	c.JSON(http.StatusOK, rsp)
}
