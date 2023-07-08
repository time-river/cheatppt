package userapiv1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"cheatppt/api"
	"cheatppt/controller/billing"
)

type UsageRange struct {
	From   int `form:"from"`
	Length int `form:"length"`
}

func Usage(c *gin.Context) {
	var req UsageRange
	var rsp = &api.Response{Status: api.FAILURE}

	if err := c.BindQuery(&req); err != nil {
		rsp.Message = "非法请求"
		c.AbortWithStatusJSON(http.StatusBadRequest, rsp)
		return
	}

	rng := &billing.Range{
		UserId: c.GetInt(UserId),
		From:   req.From,
		Length: req.Length,
	}
	data, err := billing.GetUsages(rng)
	if err != nil {
		rsp.Message = err.Error()
	} else {
		rsp.Data = data
		rsp.Status = api.SUCCESS
	}

	c.JSON(http.StatusOK, rsp)
}
