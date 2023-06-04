package main

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"cheatppt/config"
	"cheatppt/log"
	"cheatppt/model/sql"
	"cheatppt/router"
)

// TODO: log

func main() {
	config.CmdlineParse()

	debug := log.Setup()

	// TODO: env check
	sql.DatabaseInit()

	engine := gin.New()
	if !debug {
		gin.SetMode(gin.ReleaseMode)
	}
	engine.Use(gin.LoggerWithWriter(log.GetWriter()))
	engine.Use(gin.Recovery())

	router.Initialize(engine)

	addr := fmt.Sprintf("%s:%d", config.GlobalCfg.Server.Addr, config.GlobalCfg.Server.Port)
	engine.Run(addr)
}
