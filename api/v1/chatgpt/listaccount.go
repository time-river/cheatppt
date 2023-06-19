package chatgptapiv1

import (
	"cheatppt/api"
	revchatgpt3 "cheatppt/controller/chat/revchatgpt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ListAccounts(c *gin.Context) {
	var rsp *api.Response

	data, err := revchatgpt3.ListAccounts()
	if err != nil {
		rsp = &api.Response{
			Status:  api.FAILURE,
			Message: err.Error(),
		}
	} else {
		rsp = &api.Response{
			Status: api.SUCCESS,
			Data:   data,
		}
	}

	c.JSON(http.StatusOK, rsp)
}
