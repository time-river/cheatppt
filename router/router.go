package router

import (
	"github.com/gin-gonic/gin"

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

func Initialize(router *gin.Engine) {

	user := router.Group(userApiPrefix)
	{
		user.POST(userCode, userapiv1.UserCode)
		user.POST(userSignUp, userapiv1.SignUp)
		user.POST(userSignIn, userapiv1.SignIn)
		user.POST(userSignOut, userapiv1.SignOut)
		user.POST(userReset, userapiv1.Reset)
	}
}
