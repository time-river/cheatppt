package modelapiv1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kr/pretty"
	log "github.com/sirupsen/logrus"

	"cheatppt/api"
	"cheatppt/controller/model"
)

type DelReq struct {
	Id uint `json:"id"`
}

func Del(c *gin.Context) {
	var req DelReq
	var rsp = &api.Response{Status: api.FAILURE}

	if err := c.BindJSON(&req); err != nil {
		rsp.Message = "非法请求"
		c.AbortWithStatusJSON(http.StatusBadRequest, rsp)
		return
	}

	log.Debug(pretty.Sprint(req))

	detail := &model.DelDetail{
		Id: req.Id,
	}
	if err := model.Del(detail); err != nil {
		rsp.Message = err.Error()
		c.JSON(http.StatusOK, rsp)
		return
	}

	rsp.Status = api.SUCCESS
	c.JSON(http.StatusOK, rsp)
}
