package messages

import (
	"errors"
	"go-clean/src/models/messages/locales"
	"go-clean/src/models/responses"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var SuccessCodes = map[string]map[string]string{
	"en": locales.SuccessEN,
	"id": locales.SuccessID,
}

var ErrorCodes = map[string]map[string]string{
	"en": locales.ErrorEN,
	"id": locales.ErrorID,
}

func SendSuccessResponse(ctx *fiber.Ctx, successResponse responses.SuccessResponse) {
	lang := ctx.Locals("lang").(string)
	var message string

	if successResponse.StatusCode == 0 {
		successResponse.StatusCode = http.StatusOK
	}

	if (successResponse.SuccessCode != "") && (SuccessCodes[lang] == nil || SuccessCodes[lang][successResponse.SuccessCode] == "") {
		ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Success code is not defined",
			"status":  http.StatusText(http.StatusInternalServerError),
		})
		return
	} else if successResponse.SuccessCode != "" {
		message = SuccessCodes[lang][successResponse.SuccessCode]
	}

	body := fiber.Map{
		"success": true,
		"status":  http.StatusText(successResponse.StatusCode),
		"message": message,
	}

	if successResponse.Data != nil {
		body["data"] = successResponse.Data
	}
	ctx.Status(successResponse.StatusCode).JSON(body)
}

func SendErrorResponse(ctx *fiber.Ctx, errorResponse responses.ErrorResponse) {
	if HasError(errorResponse.Error) {
		var errorWrapper *ErrorWrapper
		switch {
		case errors.As(errorResponse.Error, &errorWrapper):
			var err *ErrorWrapper
			errors.As(errorResponse.Error, &err)
			if err != nil {
				lang := ctx.Locals("lang").(string)
				var message string

				if errorResponse.StatusCode == 0 {
					errorResponse.StatusCode = http.StatusInternalServerError
				}

				if err.StatusCode == 0 {
					err.StatusCode = http.StatusInternalServerError
				}

				if (err.ErrorCode != "") && (ErrorCodes[lang] == nil || ErrorCodes[lang][err.ErrorCode] == "") {
					ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
						"success": false,
						"error":   "Error code is not defined",
						"status":  http.StatusText(http.StatusInternalServerError),
					})
					return
				} else if err.ErrorCode != "" {
					message = ErrorCodes[lang][err.ErrorCode]
				}

				body := fiber.Map{
					"success": false,
					"error":   errorResponse.Error.Error(),
					"status":  http.StatusText(err.StatusCode),
				}

				if err.ErrorCode != "" {
					body["message"] = message
				}

				if len(errorResponse.FormErrors) > 0 {
					body["formErrors"] = errorResponse.FormErrors
				}

				ctx.Status(err.StatusCode).JSON(body)
			}
		default:
			statusCode := 500
			if errorResponse.StatusCode != 500 {
				statusCode = errorResponse.StatusCode
			}

			if errors.Is(errorResponse.Error, gorm.ErrRecordNotFound) {
				statusCode = http.StatusNotFound
			}

			body := fiber.Map{
				"success": false,
				"error":   errorResponse.Error.Error(),
				"status":  http.StatusText(statusCode),
			}
			ctx.Status(statusCode).JSON(body)
		}
	}
}
