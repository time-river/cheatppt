package main

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"cheatppt/config"
	"cheatppt/model/sql"
	"cheatppt/router"
)

// TODO: log

func main() {
	config.CmdlineParse()

	// TODO: env check
	sql.DatabaseInit()

	if !config.GlobalCfg.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.Default()
	engine.Use(gin.Logger())
	router.Initialize(engine)

	addr := fmt.Sprintf("%s:%d", config.GlobalCfg.Server.Addr, config.GlobalCfg.Server.Port)
	engine.Run(addr)
}
