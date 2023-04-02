package router

import (
	"github.com/gin-gonic/gin"

	"cheatppt/api/v1"
)

const (
	register    = "/register"
	login       = "/login"
	authorized  = "/authorized"
	reset       = "/reset"
	chatProcess = "/chat-process"
	listModels  = "/list-models"
)

func Initialize(router *gin.Engine) {

	router.LoadHTMLGlob("templates/*")

	apiv1 := router.Group("/api/v1")
	{
		apiv1.POST(register, api.UserRegister)
		apiv1.POST(login, api.UserLogin)
		apiv1.POST(authorized, api.UserAuthorized)
		apiv1.POST(reset, api.UserPasswordReset)
		apiv1.POST(chatProcess, api.ChatProcess)
		apiv1.GET(listModels, api.ListModels)
	}
}
