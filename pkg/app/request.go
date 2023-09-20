package app

import (
	"github.com/Eden/go-gin-example/pkg/e"
	"github.com/Eden/go-gin-example/pkg/logging"
	"github.com/astaxie/beego/validation"
	"net/http"
)

func MarkErrors(errors []*validation.Error) {
	for _, err := range errors {
		logging.Info(err.Key, err.Message)
	}

	return
}

func MarkErrorsAndExit(valid validation.Validation, appG Gin) {
	if valid.HasErrors() {
		MarkErrors(valid.Errors)
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
	}
	return
}
