package userapiv1

import (
	"cheatppt/api"
	"cheatppt/controller/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ListCDKeys(c *gin.Context) {
	rsp := &api.Response{Status: api.SUCCESS}

	data, err := user.ListCDKeys()
	if err != nil {
		rsp = &api.Response{
			Status:  api.FAILURE,
			Message: err.Error(),
		}
	}
	rsp.Data = data

	c.JSON(http.StatusOK, rsp)
}
