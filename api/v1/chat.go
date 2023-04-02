package api

import (
	"cheatppt/controller/auth"
	"cheatppt/controller/chat"
	model "cheatppt/model/http"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ListModels(c *gin.Context) {
	var resp model.ListModelsResponse

	code := http.StatusUnauthorized
	token := c.Request.Header.Get("Token")
	auth := auth.AuthCtxCreate()
	if err := auth.UserAuthorized(&token, nil); err == nil {
		code = http.StatusOK
		resp.Models = chat.ListModels()
	}

	c.JSON(code, &resp)
}

func ChatProcess(c *gin.Context) {

}
