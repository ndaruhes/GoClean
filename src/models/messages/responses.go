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

func SendBasicResponse(fiberCtx *fiber.Ctx, basicResponse responses.BasicResponse) {
	lang := fiberCtx.Locals("lang").(string)
	var message string

	if basicResponse.StatusCode == 0 {
		basicResponse.StatusCode = http.StatusOK
	}

	if (basicResponse.SuccessCode != "") && (SuccessCodes[lang] == nil || SuccessCodes[lang][basicResponse.SuccessCode] == "") {
		fiberCtx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Success code is not defined",
			"status":  http.StatusText(http.StatusInternalServerError),
		})
		return
	} else if basicResponse.SuccessCode != "" {
		message = SuccessCodes[lang][basicResponse.SuccessCode]
	}

	body := fiber.Map{
		"success": true,
		"status":  http.StatusText(basicResponse.StatusCode),
		"message": message,
	}

	if basicResponse.Data != nil {
		body["data"] = basicResponse.Data
	}
	fiberCtx.Status(basicResponse.StatusCode).JSON(body)
}

func SendErrorResponse(fiberCtx *fiber.Ctx, errorResponse responses.ErrorResponse) {
	if HasError(errorResponse.Error) {
		var errorWrapper *ErrorWrapper
		switch {
		case errors.As(errorResponse.Error, &errorWrapper):
			var err *ErrorWrapper
			errors.As(errorResponse.Error, &err)
			if err != nil {
				lang := fiberCtx.Locals("lang").(string)
				var message string

				if errorResponse.StatusCode == 0 {
					errorResponse.StatusCode = http.StatusInternalServerError
				}

				if err.StatusCode == 0 {
					err.StatusCode = http.StatusInternalServerError
				}

				if (err.ErrorCode != "") && (ErrorCodes[lang] == nil || ErrorCodes[lang][err.ErrorCode] == "") {
					fiberCtx.Status(http.StatusInternalServerError).JSON(fiber.Map{
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

				fiberCtx.Status(err.StatusCode).JSON(body)
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
			fiberCtx.Status(statusCode).JSON(body)
		}
	}
}
