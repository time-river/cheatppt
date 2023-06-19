package userapiv1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kr/pretty"
	log "github.com/sirupsen/logrus"

	"cheatppt/api"
	"cheatppt/controller/user"
)

type GenCDKeyReq struct {
	Number  int    `json:"number"`
	Comment string `json:"comment"`
	Coins   int    `json:"coins"`
	Expire  int    `json:"expire"`
}

// generate cd-key
func GenCDKey(c *gin.Context) {
	var req GenCDKeyReq
	var rsp = &api.Response{Status: api.FAILURE}

	if err := c.BindJSON(&req); err != nil || req.Number < 1 {
		rsp.Message = "非法请求"
		c.AbortWithStatusJSON(http.StatusBadRequest, rsp)
		return
	}

	log.Debug(pretty.Sprint(req))

	meta := user.CDKeyMeta{
		Nr:      req.Number,
		Comment: req.Comment,
		Coins:   req.Coins,
		Expire:  req.Expire,
	}
	data, err := user.GenCDKeys(&meta)
	if err != nil {
		rsp.Message = err.Error()
		c.JSON(http.StatusOK, rsp)
		return
	}

	rsp.Data = data
	rsp.Status = api.SUCCESS
	c.JSON(http.StatusOK, rsp)
}
