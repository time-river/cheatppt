package router

import (
	apiv0 "cheatppt/api/v0"
	"cheatppt/api/v1"

	"github.com/gin-gonic/gin"
)

const (
	register    = "/register"
	login       = "/login"
	authorized  = "/authorized"
	emailVerfiy = "/email-verify"
	reset       = "/reset"
	logout      = "/logout"
	chatProcess = "/chat-process"
	listModels  = "/list-models"
)

const (
	chatgptWebConfig      = "/config"
	chatgptWebChatProcess = "/chat-process"
	chatgptWebSession     = "/session"
	chatgptWebVerify      = "/verify"
)

func Initialize(router *gin.Engine) {
	router.LoadHTMLGlob("templates/*")

	/* v2 just for no-authorization request */
	rv0 := router.Group("/api/v0")
	{
		rv0.POST(chatgptWebVerify, apiv0.ChatgptWebVerify)
		rv0.POST(chatgptWebSession, apiv0.ChatgptWebSession)

		rv0.Use(apiv0.ChatgptWebAuth)
		rv0.POST(chatgptWebChatProcess, apiv0.ChatgptWebChatProcess)
		rv0.POST(chatgptWebConfig, apiv0.ChatgptWebConfig)
	}

	apiv1 := router.Group("/api/v1")
	{
		apiv1.POST(register, api.UserRegister)
		apiv1.POST(login, api.UserLogin)
		apiv1.GET(authorized, api.UserAuthorized)
		apiv1.POST(emailVerfiy, api.EmailVerfiy)
		apiv1.POST(reset, api.UserPasswordReset)
		apiv1.GET(logout, api.UserLogout)
		apiv1.POST(chatProcess, api.ChatProcess)
		apiv1.GET(listModels, api.ListModels)
	}
}
