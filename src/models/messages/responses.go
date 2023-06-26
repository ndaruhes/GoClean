package messages

import (
	"github.com/gin-gonic/gin"
	"go-clean/models/responses"
	"net/http"
)

var MessageCodes = map[string]map[string]string{
	"en": EN,
	"id": ID,
}

func SendSuccessResponse(ctx *gin.Context, basicResponse responses.BasicResponse) {
	lang := ctx.Value("lang").(string)
	var message string
	if (basicResponse.MessageCode != "") && (MessageCodes[lang] == nil || MessageCodes[lang][basicResponse.MessageCode] == "") {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Message code is not defined",
			"status":  http.StatusText(http.StatusInternalServerError),
		})
		return
	} else if basicResponse.MessageCode != "" {
		message = MessageCodes[lang][basicResponse.MessageCode]
	}

	if basicResponse.StatusCode == 0 {
		basicResponse.StatusCode = http.StatusOK
	}
	basicResponseBody := gin.H{
		"success": true,
		"status":  http.StatusText(basicResponse.StatusCode),
		"message": message,
	}
	if basicResponse.Data != nil {
		basicResponseBody["data"] = basicResponse.Data
	}
	ctx.JSON(basicResponse.StatusCode, basicResponseBody)
}

func SendErrorResponse(ctx *gin.Context, basicResponse responses.BasicResponse) {
	if HasError(basicResponse.Error) {
		if basicResponse.StatusCode == 0 {
			basicResponse.StatusCode = http.StatusInternalServerError
		}
		switch basicResponse.Error.(type) {
		case *ErrorWrapper:
			err := basicResponse.Error.(*ErrorWrapper)
			if err != nil {
				lang := ctx.Value("lang").(string)
				var message string
				if (basicResponse.MessageCode != "") && (MessageCodes[lang] == nil || MessageCodes[lang][basicResponse.MessageCode] == "") {
					ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
						"success": false,
						"error":   "Message code is not defined",
						"status":  http.StatusText(http.StatusInternalServerError),
					})
					return
				} else if basicResponse.MessageCode != "" {
					message = MessageCodes[lang][basicResponse.MessageCode]
				}

				basicResponseBody := gin.H{
					"success": false,
					"status":  http.StatusText(basicResponse.StatusCode),
					"message": message,
				}
				if basicResponse.Data != nil {
					basicResponseBody["data"] = basicResponse.Data
				}
				ctx.JSON(basicResponse.StatusCode, basicResponseBody)
			}
		default:
			ctx.AbortWithStatusJSON(basicResponse.StatusCode, gin.H{
				"success": false,
				"error":   basicResponse.Error.Error(),
				"status":  http.StatusText(basicResponse.StatusCode),
			})
		}
	}
}
