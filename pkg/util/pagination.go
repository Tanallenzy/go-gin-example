package util

import (
	"github.com/Eden/go-gin-example/pkg/setting"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
)

func GetPage(c *gin.Context) int {
	result := 0
	if page, _ := com.StrTo(c.DefaultPostForm("page", "1")).Int(); page > 0 {
		result = (page - 1) * setting.AppSetting.PageSize
	}
	return result
}
