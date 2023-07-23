package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	appConfig "go-clean/configs/app"
)

func main() {
	app := gin.Default()
	RegisterMiddlewares(app)
	RegisterRoutes(app)
	err := app.Run(fmt.Sprintf(":%d", appConfig.GetConfig().App.Port))
	if err != nil {
		return
	}
}
