package router

import (
	"github.com/gin-gonic/gin"

	openaiapiv1 "cheatppt/api/v1/openai"
	userapiv1 "cheatppt/api/v1/user"
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
	openaiApiPrefix = "/api/openai/v1"
	openaiChat      = "/chat/completions"
)

func Initialize(router *gin.Engine) {

	user := router.Group(userApiPrefix)
	{
		user.POST(userCode, userapiv1.UserCode)
		user.POST(userSignUp, userapiv1.SignUp)
		user.POST(userSignIn, userapiv1.SignIn)
		user.POST(userSignOut, userapiv1.SignOut)
		user.POST(userReset, userapiv1.Reset)
	}

	openai := router.Group(openaiApiPrefix)
	{
		openai.POST(openaiChat, openaiapiv1.Chat)
	}
}
