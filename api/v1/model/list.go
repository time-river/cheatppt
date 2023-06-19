package modelapiv1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"cheatppt/api"
	"cheatppt/controller/model"
)

func ListAvailable(c *gin.Context) {
	data, err := model.ListAvailable()
	if err != nil {
		rsp := &api.Response{
			Status:  api.FAILURE,
			Message: err.Error(),
		}

		c.JSON(http.StatusOK, rsp)
		return
	}

	rsp := &api.Response{
		Status: api.SUCCESS,
		Data:   data,
	}
	c.JSON(http.StatusOK, rsp)
}

func ListAll(c *gin.Context) {
	data, err := model.ListAll()
	if err != nil {
		rsp := &api.Response{
			Status:  api.FAILURE,
			Message: err.Error(),
		}

		c.JSON(http.StatusOK, rsp)
		return
	}

	rsp := &api.Response{
		Status: api.SUCCESS,
		Data:   data,
	}
	c.JSON(http.StatusOK, rsp)
}
