package router

import (
	"github.com/gin-gonic/gin"

	chatgptapiv1 "cheatppt/api/v1/chatgpt"
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

const (
	chatGPTApiPrefix    = "/api/v1/chatgpt"
	chatGPTChat         = "/chat"
	chatGPTRefreshToken = "/refresh"
)

func Initialize(router *gin.Engine) {

	user := router.Group(userApiPrefix)
	{
		user.POST(userCode, userapiv1.UserCode)
		user.POST(userSignUp, userapiv1.SignUp)
		user.POST(userSignIn, userapiv1.SignIn)
		user.POST(userSignOut, userapiv1.Guard, userapiv1.SignOut)
		user.POST(userReset, userapiv1.Reset)
	}

	openai := router.Group(openaiApiPrefix)
	{
		openai.Use(userapiv1.Guard)
		openai.POST(openaiChat, openaiapiv1.Chat)
	}

	chatGPT := router.Group(chatGPTApiPrefix)
	{
		openai.Use(userapiv1.Guard)
		chatGPT.POST(chatGPTChat, chatgptapiv1.Chat)
		chatGPT.PATCH(chatGPTRefreshToken, chatgptapiv1.RefreshToken)
	}
}
