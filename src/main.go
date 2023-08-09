package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	server "go-clean/src/app"
	appConfig "go-clean/src/app/config"
)

func main() {
	app := gin.Default()
	server.RegisterMiddlewares(app)
	server.RegisterRoutes(app)
	err := app.Run(fmt.Sprintf(":%d", appConfig.GetConfig().App.Port))
	if err != nil {
		return
	}
}
