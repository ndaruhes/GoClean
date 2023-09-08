package messages

import (
	"errors"
	"github.com/gofiber/fiber/v2"
)

type ErrorWrapper struct {
	Context    *fiber.Ctx
	Err        error
	ErrorCode  string
	Parameters []string
	StatusCode int
}

func (wrapper *ErrorWrapper) Error() string {
	lang := wrapper.Context.Locals("lang").(string)
	if wrapper.Err == nil || wrapper.ErrorCode == "" {
		return ErrorCodes[lang]["ERROR-50003"]
	}
	return wrapper.Err.Error()
}

func HasError(err error) bool {
	var errorWrapper *ErrorWrapper
	switch {
	case errors.As(err, &errorWrapper):
		return err.(*ErrorWrapper).Err != nil || err.(*ErrorWrapper).ErrorCode != ""
	default:
		return err != nil
	}
}
