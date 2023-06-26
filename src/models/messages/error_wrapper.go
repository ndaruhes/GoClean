package messages

import (
	"github.com/gin-gonic/gin"
)

type ErrorWrapper struct {
	Context    *gin.Context
	Err        error
	ErrorCode  string
	Parameters []string
	StatusCode int
}

func (wrapper *ErrorWrapper) Error() string {
	if wrapper.Err == nil || wrapper.ErrorCode == "" {
		return "Something went wrong"
	}
	return wrapper.Err.Error()
}

func HasError(err error) bool {
	switch err.(type) {
	case *ErrorWrapper:
		return err.(*ErrorWrapper).Err != nil || err.(*ErrorWrapper).ErrorCode != ""
	default:
		return err != nil
	}
}
