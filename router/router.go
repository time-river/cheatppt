package router

import (
	"github.com/gin-gonic/gin"

	chatgptapiv1 "cheatppt/api/v1/chatgpt"
	modelapiv1 "cheatppt/api/v1/model"
	openaiapiv1 "cheatppt/api/v1/openai"
	userapiv1 "cheatppt/api/v1/user"
)

const (
	userApiPrefix     = "/api/v1/user"
	userReset         = "/reset"
	userCode          = "/code"
	userSignUp        = "/signup"
	userSignIn        = "/signin"
	userSignOut       = "/signout"
	userPing          = "/ping"
	userDetail        = "/detail"
	userCDKey         = "/exgcdkey"
	userGenCDKey      = "/gencdkey"
	userListCDKey     = "/listcdkeys"
	userTopup         = "/topup"
	userBalance       = "/balance"
	userBilling       = "/billing"
	userUsage         = "/usage"
	userCreateSession = "/createsession"
)

const (
	modelApiPrefix = "/api/v1/model"
	modelAdd       = "/add"
	modelDel       = "/delete"
	modelList      = "/list"
	modelListAll   = "/listall"
)

const (
	openaiApiPrefix = "/api/v1/openai"
	openaiChat      = "conversation"
)

const (
	chatGPTApiPrefix    = "/api/v1/chatgpt"
	chatGPTChat         = "conversation"
	chatGPTAddAccount   = "addaccount"
	chatGPTDelAccount   = "delaccount"
	chatGPTListAccounts = "listaccounts"
	chatGPTLogin        = "login"
)

func Initialize(router *gin.Engine) {

	user := router.Group(userApiPrefix)
	{
		user.POST(userCode, userapiv1.UserCode)
		user.POST(userSignUp, userapiv1.SignUp)
		user.POST(userSignIn, userapiv1.SignIn)
		user.POST(userReset, userapiv1.Reset)

		user.GET(userPing, userapiv1.TokenGuard, userapiv1.Ping)
		user.GET(userDetail, userapiv1.SessionGuard, userapiv1.Detail)
		user.POST(userSignOut, userapiv1.SessionGuard, userapiv1.SignOut)
		user.POST(userCDKey, userapiv1.SessionGuard, userapiv1.ExgCDkey)

		user.POST(userGenCDKey, userapiv1.AdminGuard, userapiv1.GenCDKey)
		user.GET(userListCDKey, userapiv1.AdminGuard, userapiv1.ListCDKeys)

		user.POST(userTopup)
		user.GET(userBalance, userapiv1.SessionGuard, userapiv1.Balance)
		user.GET(userBilling, userapiv1.SessionGuard)
		user.GET(userUsage, userapiv1.SessionGuard)
		user.GET(userCreateSession, userapiv1.SessionGuard, userapiv1.CreateChatSession)
	}

	model := router.Group(modelApiPrefix)
	{
		// everyone can list models
		model.GET(modelList, userapiv1.SessionGuard, modelapiv1.ListAvailable)

		model.POST(modelAdd, userapiv1.AdminGuard, modelapiv1.Add)
		model.POST(modelDel, userapiv1.AdminGuard, modelapiv1.Del)
		model.GET(modelListAll, userapiv1.AdminGuard, modelapiv1.ListAll)
	}

	openai := router.Group(openaiApiPrefix, userapiv1.SessionGuard)
	{
		openai.POST(openaiChat, userapiv1.SessionGuard, openaiapiv1.Chat)
	}

	chatGPT := router.Group(chatGPTApiPrefix, userapiv1.SessionGuard)
	{
		chatGPT.POST(chatGPTLogin, userapiv1.SessionGuard, chatgptapiv1.Login)
		chatGPT.POST(chatGPTChat, userapiv1.SessionGuard, chatgptapiv1.Conversation)

		chatGPT.POST(chatGPTAddAccount, userapiv1.AdminGuard, chatgptapiv1.AddAccount)
		chatGPT.POST(chatGPTDelAccount, userapiv1.AdminGuard, chatgptapiv1.DelAccount)
		chatGPT.GET(chatGPTListAccounts, userapiv1.AdminGuard, chatgptapiv1.ListAccounts)
	}
}
