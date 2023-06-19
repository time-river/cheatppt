package userapiv1

import (
	"cheatppt/api"
	"cheatppt/controller/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Detail(c *gin.Context) {
	rsp := &api.Response{Status: api.FAILURE}

	data, err := user.Detail(c.GetInt(UserId))
	if err != nil {
		rsp.Message = err.Error()
	} else {
		rsp.Data = data
		rsp.Status = api.SUCCESS
	}

	c.JSON(http.StatusOK, rsp)
}
