package api

import (
	"github.com/Eden/go-gin-example/pkg/app"
	"github.com/Eden/go-gin-example/pkg/e"
	"github.com/Eden/go-gin-example/pkg/util"
	"github.com/Eden/go-gin-example/service/auth_service"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"net/http"
)

type auth struct {
	Username string `valid:"Required; MaxSize(50)"`
	Password string `valid:"Required; MaxSize(50)"`
}

func GetAuth(c *gin.Context) {
	appG := app.Gin{c}

	var authService = auth_service.Auth{
		Username: c.Query("username"),
		Password: c.Query("password"),
	}

	valid := validation.Validation{}
	a := auth{Username: authService.Username, Password: authService.Password}
	valid.Valid(&a)
	app.MarkErrorsAndExit(valid, appG)

	check, err := authService.Check()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}

	if !check {
		appG.Response(http.StatusUnauthorized, e.ERROR_AUTH, nil)
		return

	}

	token, err := util.GenerateToken(authService.Username, authService.Password)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_AUTH_TOKEN, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"token": token,
	})
}
