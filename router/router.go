package router

import (
	"github.com/gin-gonic/gin"

	chatgptapiv1 "cheatppt/api/v1/chatgpt"
	modelapiv1 "cheatppt/api/v1/model"
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
	userPing      = "/ping"
	userCDKey     = "/cdkey"
	userPay       = "/pay"
)

const (
	modelApiPrefix = "/api/v1/model"
	modelAdd       = "/add"
	modelDel       = "/delete"
	modelList      = "/list"
	modelListAll   = "/listall"
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
		user.GET(userPing, userapiv1.TokenGuard, userapiv1.Ping)
		user.POST(userSignOut, userapiv1.SessionGuard, userapiv1.SignOut)
		user.POST(userReset, userapiv1.Reset)
	}

	model := router.Group(modelApiPrefix)
	{
		// everyone can list models
		model.GET(modelList, modelapiv1.ListAvailable)

		model.POST(modelAdd, userapiv1.AdminGuard, modelapiv1.Add)
		model.POST(modelDel, userapiv1.AdminGuard, modelapiv1.Del)
		model.GET(modelListAll, userapiv1.AdminGuard, modelapiv1.ListAll)
	}

	openai := router.Group(openaiApiPrefix, userapiv1.SessionGuard)
	{
		openai.POST(openaiChat, openaiapiv1.Chat)
	}

	chatGPT := router.Group(chatGPTApiPrefix, userapiv1.SessionGuard)
	{
		chatGPT.POST(chatGPTChat, chatgptapiv1.Chat)
		chatGPT.PATCH(chatGPTRefreshToken, chatgptapiv1.RefreshToken)
	}
}
