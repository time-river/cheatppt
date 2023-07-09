package controller

import (
	revchatgpt3 "cheatppt/controller/chat/revchatgpt"
	"cheatppt/controller/model"
	"cheatppt/controller/user"
)

func Setup() {
	model.Setup()
	revchatgpt3.Setup()
	user.Setup()
}
