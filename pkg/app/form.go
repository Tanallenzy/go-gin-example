package app

import (
	"github.com/Eden/go-gin-example/pkg/e"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"net/http"
)

func BindAndValid(c *gin.Context, form interface{}) (int, int) {
	err := c.Bind(form)
	if err != nil {
		return http.StatusBadRequest, e.INVALID_PARAMS
	}

	valid := validation.Validation{}
	check, err := valid.Valid(form)
	if err != nil {
		return http.StatusInternalServerError, e.ERROR
	}
	if !check {
		MarkErrors(valid.Errors)
		return http.StatusBadRequest, e.INVALID_PARAMS
	}

	return http.StatusOK, e.SUCCESS
}

func BindAndValidAndExit(c *gin.Context, form interface{}) {
	httpCode, errCode := BindAndValid(c, form)
	if errCode != e.SUCCESS {
		appG := Gin{c}
		appG.Response(httpCode, errCode, nil)
	}
	return
}
