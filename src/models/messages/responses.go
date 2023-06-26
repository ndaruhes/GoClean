package messages

import (
	"go-clean/models/messages/locales"
	"go-clean/models/responses"
	"net/http"

	"github.com/gin-gonic/gin"
)

var SuccessCodes = map[string]map[string]string{
	"en": locales.SuccessEN,
	"id": locales.SuccessID,
}

var ErrorCodes = map[string]map[string]string{
	"en": locales.ErrorEN,
	"id": locales.ErrorID,
}

//func SetupResponse(ctx *gin.Context, basicResponse responses.BasicResponse, messageType string) (string, int) {
//	lang := ctx.Value("lang").(string)
//	var (
//		message    string
//		statusCode int
//	)
//
//	if basicResponse.StatusCode == 0 {
//		if messageType == "Success" {
//			statusCode = http.StatusOK
//		} else if messageType == "Error" {
//			statusCode = http.StatusInternalServerError
//		}
//	} else {
//		statusCode = basicResponse.StatusCode
//	}
//
//	if (basicResponse.MessageCode != "") && (MessageCodes[lang] == nil || MessageCodes[lang][basicResponse.MessageCode] == "") {
//		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
//			"success": false,
//			"error":   "Message code is not defined",
//			"status":  http.StatusText(http.StatusInternalServerError),
//		})
//		return "", 0
//	} else if basicResponse.MessageCode != "" {
//		message = MessageCodes[lang][basicResponse.MessageCode]
//	}
//
//	return message, statusCode
//}

func SendSuccessResponse(ctx *gin.Context, successResponse responses.SuccessResponse) {
	lang := ctx.Value("lang").(string)
	var message string

	if successResponse.StatusCode == 0 {
		successResponse.StatusCode = http.StatusOK
	}

	if (successResponse.SuccessCode != "") && (SuccessCodes[lang] == nil || SuccessCodes[lang][successResponse.SuccessCode] == "") {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Success code is not defined",
			"status":  http.StatusText(http.StatusInternalServerError),
		})
		return
	} else if successResponse.SuccessCode != "" {
		message = SuccessCodes[lang][successResponse.SuccessCode]
	}

	successResponseBody := gin.H{
		"success": true,
		"status":  http.StatusText(successResponse.StatusCode),
		"message": message,
	}

	if successResponse.Data != nil {
		successResponseBody["data"] = successResponse.Data
	}
	ctx.JSON(successResponse.StatusCode, successResponseBody)
}

func SendErrorResponse(ctx *gin.Context, errorResponse responses.ErrorResponse) {
	if errorResponse.StatusCode == 0 {
		errorResponse.StatusCode = http.StatusInternalServerError
	}

	if HasError(errorResponse.Error) {
		switch errorResponse.Error.(type) {
		case *ErrorWrapper:
			err := errorResponse.Error.(*ErrorWrapper)
			if err != nil {
				lang := ctx.Value("lang").(string)
				var message string

				if (err.ErrorCode != "") && (ErrorCodes[lang] == nil || ErrorCodes[lang][err.ErrorCode] == "") {
					ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
						"success": false,
						"error":   "Error code is not defined",
						"status":  http.StatusText(http.StatusInternalServerError),
					})
					return
				} else if err.ErrorCode != "" {
					message = ErrorCodes[lang][err.ErrorCode]
				}

				errorBodyResponse := gin.H{
					"success": false,
					"error":   errorResponse.Error.Error(),
					"status":  http.StatusText(err.StatusCode),
				}

				if err.ErrorCode != "" {
					errorBodyResponse["message"] = message
				}

				ctx.JSON(err.StatusCode, errorBodyResponse)
			}
		default:
			statusCode := 500
			if errorResponse.StatusCode != 500 {
				statusCode = errorResponse.StatusCode
			}
			errorBodyResponse := gin.H{
				"success": false,
				"error":   errorResponse.Error,
				"status":  http.StatusText(statusCode),
			}
			ctx.JSON(statusCode, errorBodyResponse)
		}
	}
}

//func SendErrorResponse(ctx *gin.Context, basicResponse responses.BasicResponse) {
//	//message, statusCode := SetupResponse(ctx, basicResponse, "Error")
//	lang := ctx.Value("lang").(string)
//	var message string
//
//	if basicResponse.StatusCode == 0 {
//		basicResponse.StatusCode = http.StatusInternalServerError
//	}
//
//	if HasError(basicResponse.Error) {
//		switch basicResponse.Error.(type) {
//		case *ErrorWrapper:
//			err := basicResponse.Error.(*ErrorWrapper)
//			if err != nil {
//				if (err.ErrorCode != "") && (MessageCodes[lang] == nil || MessageCodes[lang][err.ErrorCode] == "") {
//					ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
//						"success": false,
//						"error":   "Message code is not defined",
//						"status":  http.StatusText(http.StatusInternalServerError),
//					})
//					return
//				} else if err.ErrorCode != "" {
//					message = MessageCodes[lang][err.ErrorCode]
//				}
//
//				basicResponseBody := gin.H{
//					"success": false,
//					"error":   basicResponse.Error.Error(),
//					"status":  http.StatusText(err.StatusCode),
//				}
//
//				if err.ErrorCode != "" {
//					basicResponseBody["message"] = message
//				}
//
//				ctx.JSON(err.StatusCode, basicResponseBody)
//			}
//		default:
//			statusCode := 500
//			if basicResponse.StatusCode != 0 {
//				statusCode = basicResponse.StatusCode
//			}
//			ctx.AbortWithStatusJSON(statusCode, gin.H{
//				"success": false,
//				"error":   basicResponse.Error.Error(),
//				"status":  http.StatusText(statusCode),
//			})
//		}
//	}
//}
