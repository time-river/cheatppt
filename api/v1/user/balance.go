package userapiv1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"cheatppt/api"
	"cheatppt/controller/user"
)

type BalanceRsp struct {
	Balance float32 `json:"balance"`
}

func Balance(c *gin.Context) {
	var rsp = &api.Response{Status: api.FAILURE}

	userId := c.GetInt(UserId)
	data, err := user.Detail(userId)
	if err != nil {
		rsp.Message = err.Error()
	} else {
		rsp.Data = BalanceRsp{
			// convert to yuan from cent
			Balance: float32(data.Credit) / 100.0,
		}
		rsp.Status = api.SUCCESS
	}

	c.JSON(http.StatusOK, rsp)
}
