package userapiv1

import "github.com/gin-gonic/gin"

type TopupReq struct {
	AppId    string `form:"appId"`
	AuthCode string `form:"authCode"`
}

func Topup(c *gin.Context) {

}
