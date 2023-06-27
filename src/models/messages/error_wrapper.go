package messages

import (
	"github.com/gin-gonic/gin"
	"go-clean/models/messages/locales"
)

type ErrorWrapper struct {
	Context    *gin.Context
	Err        error
	ErrorCode  string
	Parameters []string
	StatusCode int
}

var BasicCodes = map[string]map[string]string{
	"en": locales.BasicEN,
	"id": locales.BasicID,
}

func (wrapper *ErrorWrapper) Error() string {
	lang := wrapper.Context.Value("lang").(string)
	if wrapper.Err == nil || wrapper.ErrorCode == "" {
		return BasicCodes[lang]["BASIC-0001"]
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
