package chatgptapiv1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kr/pretty"
	log "github.com/sirupsen/logrus"

	"cheatppt/api"
	revchatgpt3 "cheatppt/controller/chat/revchatgpt"
)

type AddAccountReq struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	Activated bool   `json:"activated"`
}

func AddAccount(c *gin.Context) {
	var req AddAccountReq
	var rsp = &api.Response{Status: api.FAILURE}

	if err := c.BindJSON(&req); err != nil {
		rsp.Message = "非法请求"
		c.AbortWithStatusJSON(http.StatusBadRequest, rsp)
		return
	}

	// TODO: password leak
	log.Debug(pretty.Sprint(req))

	account := revchatgpt3.Account{
		Email:     req.Email,
		Password:  req.Password,
		Activated: req.Activated,
	}
	if err := revchatgpt3.AddAccount(&account); err != nil {
		rsp.Message = err.Error()
	} else {
		rsp.Status = api.SUCCESS
	}

	c.JSON(http.StatusOK, rsp)
}
