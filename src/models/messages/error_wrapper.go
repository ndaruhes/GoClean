package messages

import (
	"context"
	"errors"
)

type ErrorWrapper struct {
	Context    context.Context
	Err        error
	ErrorCode  string
	Parameters []string
	StatusCode int
}

func (wrapper *ErrorWrapper) Error() string {
	lang := wrapper.Context.Value("lang").(string)
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
