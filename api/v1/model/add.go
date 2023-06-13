package modelapiv1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"cheatppt/api"
	"cheatppt/controller/model"
)

const (
	Add_REQ_CREATE = "create"
	ADD_REQ_MODIFY = "modify"
)

type AddReq struct {
	Type        string `json:"type"`
	Id          uint   `json:"Id,omitempty"`
	DisplayName string `json:"displayName"`
	ModelName   string `json:"modelName"`
	Provider    string `json:"provider"`
	LeastCoins  int    `json:"leastCoins"`
	Activated   bool   `json:"activated,omitempty"`
}

func Add(c *gin.Context) {
	var req AddReq
	var rsp = &api.Response{Status: api.FAILURE}

	if err := c.BindJSON(&req); err != nil {
		rsp.Message = "非法请求"
		c.AbortWithStatusJSON(http.StatusOK, rsp)
		return
	}

	detail := model.AddDetail{
		DisplayName: req.DisplayName,
		ModelName:   req.ModelName,
		Provider:    req.Provider,
		LeastCoins:  req.LeastCoins,
		Activated:   req.Activated,
	}

	if req.Type == Add_REQ_CREATE {
		if err := model.Add(&detail, true); err != nil {
			rsp.Message = err.Error()
		} else {
			rsp.Status = api.SUCCESS
		}
	} else if req.Type == ADD_REQ_MODIFY {
		detail.Id = req.Id

		if err := model.Add(&detail, false); err != nil {
			rsp.Message = err.Error()
		} else {
			rsp.Status = api.SUCCESS
		}
	} else {
		rsp.Message = "非法请求"
	}

	c.JSON(http.StatusOK, rsp)
}
