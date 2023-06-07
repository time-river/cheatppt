package userapiv1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kr/pretty"

	"cheatppt/api"
	"cheatppt/config"
	"cheatppt/controller/code"
	"cheatppt/controller/user"
	"cheatppt/log"
	"cheatppt/utils"
)

const (
	CODE_REQ_SIGNUP = "signup" // need username and email
	CODE_REQ_RESET  = "reset"  // only need username
)

type CodeReq struct {
	Type     string `json:"type"` // 'signup' | 'reset'
	Username string `json:"username"`
	Email    string `json:"email,omitempty"`
	Code     string `json:"code"`
}

type CodeRsp struct {
	Code string `json:"code"`
}

func UserCode(c *gin.Context) {
	var rsp = &api.Response{Status: api.FAILURE}
	var req CodeReq

	if err := c.BindJSON((&req)); err != nil {
		rsp.Message = "非法请求"
		c.AbortWithStatusJSON(http.StatusOK, rsp)
		return
	}

	log.Debug(pretty.Sprint(req))

	if req.Type == CODE_REQ_SIGNUP && !config.Server.EnableRegister {
		rsp.Message = "管理员禁止注册"
		c.JSON(http.StatusOK, rsp)
		return
	}

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

	if req.Type == CODE_REQ_SIGNUP && utils.UsernameCheck(req.Username) && utils.EmailCheck(req.Email) {
		if err := user.GenerateSignUpCode(req.Username, req.Email); err != nil {
			rsp.Message = err.Error()
		} else {
			rsp.Status = api.SUCCESS
			rsp.Message = "发送成功，请查收"
		}
	} else if req.Type == CODE_REQ_RESET && utils.UsernameCheck(req.Username) {
		if err := user.GenerateResetCode(req.Username); err != nil {
			rsp.Message = err.Error()
		} else {
			rsp.Status = api.SUCCESS
			rsp.Message = "发送成功，请查收"
		}
	} else {
		rsp.Message = "非法请求"
	}

	c.JSON(http.StatusOK, rsp)
}
