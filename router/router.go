package router

import (
	"github.com/gin-gonic/gin"

	chatgptapiv1 "cheatppt/api/v1/chatgpt"
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
		user.POST(userSignOut, auth.TokenVerify, userapiv1.SignOut)
		user.POST(userReset, userapiv1.Reset)
	}

	chatGPT := router.Group(chatGPTApiPrefix)
	{
		chatGPT.Use(auth.TokenVerify)
		chatGPT.POST(chatGPTChat, chatgptapiv1.Chat)
		chatGPT.PATCH(chatGPTRefreshToken, chatgptapiv1.RefreshToken)
	}
}
