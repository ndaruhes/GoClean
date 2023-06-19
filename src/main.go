package main

import (
	"fmt"
	appConfig "go-clean/configs/app"

	"github.com/gin-gonic/gin"
)

func main() {
	app := gin.Default()
	RegisterMiddlewares(app)
	RegisterRoutes(app)
	app.Run(fmt.Sprintf(":%d", appConfig.GetConfig().App.Port))
}
