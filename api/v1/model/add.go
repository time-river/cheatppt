package modelapiv1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kr/pretty"
	log "github.com/sirupsen/logrus"

	"cheatppt/api"
	"cheatppt/controller/model"
)

const (
	Add_REQ_CREATE = "create"
	ADD_REQ_MODIFY = "modify"
)

type AddReq struct {
	Type        string `json:"type"`
	Id          uint   `json:"id,omitempty"`
	DisplayName string `json:"displayName"`
	ModelName   string `json:"modelName"`
	Provider    string `json:"provider"`
	InputCoins  int    `json:"inputCoins"`
	OutputCoins int    `json:"outputCoins"`
	Activated   bool   `json:"activated,omitempty"`
}

func Add(c *gin.Context) {
	var req AddReq
	var rsp = &api.Response{Status: api.FAILURE}

	if err := c.BindJSON(&req); err != nil {
		rsp.Message = "非法请求"
		c.AbortWithStatusJSON(http.StatusBadRequest, rsp)
		return
	}

	log.Debug(pretty.Sprint(req))

	detail := model.AddDetail{
		DisplayName: req.DisplayName,
		ModelName:   req.ModelName,
		Provider:    req.Provider,
		InputCoins:  req.InputCoins,
		OutputCoins: req.OutputCoins,
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
