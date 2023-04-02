package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"cheatppt/controller/auth"
	msg "cheatppt/model/http"
)

func checkUserInformation(user, pass, rec *string) error {
	if *user == "" || *pass == "" {
		return errors.New("No username or password")
	} else if *rec == "" {
		return errors.New("No reCAPTCHA")
	}
	return nil
}

func UserRegister(c *gin.Context) {
	var req msg.RegisterRequest
	msg := &msg.CommonResponse{}

	if err := c.BindJSON(&req); err != nil {
		msg.Message = err.Error()
		c.AbortWithStatusJSON(http.StatusBadRequest, msg)
		return
	}

	if err := checkUserInformation(&req.Username, &req.Password, &req.Recaptcha); err != nil {
		msg.Message = err.Error()
		c.AbortWithStatusJSON(http.StatusBadRequest, msg)
		return
	}

	ip := c.ClientIP()
	auth := auth.AuthCtxCreate()
	if err := auth.UserRegister(&req, &ip); err != nil {
		msg.Message = err.Error()
		c.AbortWithStatusJSON(http.StatusBadRequest, msg)
		return
	}

	msg.Message = "Success"
	c.JSON(http.StatusOK, msg)
}

func UserLogin(c *gin.Context) {
	var req msg.LoginRequest

	if err := c.BindJSON(&req); err != nil {
		msg := &msg.CommonResponse{
			Message: err.Error(),
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, msg)
		return
	}

	if err := checkUserInformation(&req.Username, &req.Password, &req.Recaptcha); err != nil {
		msg := &msg.CommonResponse{
			Message: err.Error(),
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, msg)
		return
	}

	ip := c.ClientIP()
	auth := auth.AuthCtxCreate()
	token, err := auth.UserLogin(&req, &ip)
	if err != nil {
		msg := &msg.CommonResponse{
			Message: err.Error(),
		}
		c.AbortWithStatusJSON(http.StatusUnauthorized, msg)
		return
	}

	resp := &msg.LoginResponse{
		Token: *token,
	}
	c.JSON(http.StatusOK, &resp)
}

func UserAuthorized(c *gin.Context) {
	var req msg.AuthorizedRequest
	var msg = &msg.CommonResponse{}

	if err := c.BindJSON(&req); err != nil {
		msg.Message = err.Error()
		c.AbortWithStatusJSON(http.StatusBadRequest, msg)
		return
	}
	if req.Username == "" || req.Token == "" {
		msg.Message = "No username or token"
		c.AbortWithStatusJSON(http.StatusBadRequest, msg)
		return
	}

	auth := auth.AuthCtxCreate()
	if err := auth.UserAuthorized(&req.Token, &req.Username); err != nil {
		msg.Message = err.Error()
		c.AbortWithStatusJSON(http.StatusUnauthorized, msg)
		return
	}

	msg.Message = "Success"
	c.JSON(http.StatusOK, msg)
}

func UserLogout(c *gin.Context) {
	/*
		var req msg.LogoutRequest


				code := http.StatusBadRequest
				err := c.BindJSON(&req)
				if err == nil {
					code = http.StatusUnauthorized
					token := c.Request.Header.Get("Token")
					auth := auth.AuthCtxCreate()
					if auth.UserAuthenticate(&token) != nil {
						auth.UserLogout(&token)
						code = http.StatusOK
					}
				}


			c.JSON(code, nil)
	*/
}

func UserPasswordReset(c *gin.Context) {

}
