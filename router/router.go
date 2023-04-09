package router

import (
	"github.com/gin-gonic/gin"

	apiv0 "cheatppt/api/v0"
	userapiv1 "cheatppt/api/v1/user"
	"cheatppt/middleware/auth"
)

const (
	userApiPrefix = "/api/v1/user"
	userReset     = "/reset"
	userCode      = "/code"
	userSignUp    = "/signup"
	userSignIn    = "/signin"
	userSignOut   = "/signout"
)

const (
	chatgptWebConfig         = "/config"
	chatgptWebChatProcess    = "/chat-process"
	chatgptWebSession        = "/session"
	chatgptWebVerify         = "/verify"
	chatgptWebRefreshSession = "/refresh-session"
)

func Initialize(router *gin.Engine) {
	//router.LoadHTMLGlob("templates/*")

	/* v2 just for no-authorization request */
	rv0 := router.Group("/api/v0")
	{
		chatgptweb := rv0.Group("/api")
		chatgptweb.POST(chatgptWebVerify, apiv0.ChatgptWebVerify)
		chatgptweb.POST(chatgptWebSession, apiv0.ChatgptWebSession)

		chatgptweb.Use(apiv0.ChatgptWebAuth)
		chatgptweb.POST(chatgptWebChatProcess, apiv0.ChatgptWebChatProcess)
		chatgptweb.POST(chatgptWebConfig, apiv0.ChatgptWebConfig)
		chatgptweb.PATCH(chatgptWebRefreshSession, apiv0.RefreshSession)
	}

	user := router.Group(userApiPrefix)
	{
		user.POST(userCode, userapiv1.UserCode)
		user.POST(userSignUp, userapiv1.SignUp)
		user.POST(userSignIn, userapiv1.SignIn)
		user.POST(userSignOut, auth.TokenVerify, userapiv1.SignOut)
		user.POST(userReset, userapiv1.Reset)
	}
}
