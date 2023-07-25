package middlewares

import (
	"go-clean/src/setup/app"
	"go-clean/src/shared/utils"

	"github.com/gin-gonic/gin"
)

func LangMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Header.Get("lang") == "" {
			c.Request.Header.Set("lang", app.GetConfig().App.Locale)
		}
		c.Set("lang", c.GetHeader("lang"))
		utils.OverrideGinRequest(c, "lang", c.GetHeader("lang"))
		c.Next()
	}
}
