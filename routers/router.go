package routers

import (
	"github.com/Eden/go-gin-example/middleware/jwt"
	"github.com/Eden/go-gin-example/pkg/setting"
	"github.com/Eden/go-gin-example/pkg/upload"
	"github.com/Eden/go-gin-example/routers/api"
	v1 "github.com/Eden/go-gin-example/routers/api/v1"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	gin.SetMode(setting.ServerSetting.RunMode)
	r.StaticFS("/upload/images", http.Dir(upload.GetImageFullPath()))
	r.GET("/auth", api.GetAuth)
	r.POST("/upload", api.UploadImage)
	apiv1 := r.Group("/api/v1")
	apiv1.Use(jwt.JWT())
	{
		apiv1.POST("/tags_list", v1.GetTags)
		apiv1.POST("/tags_add", v1.AddTag)
		apiv1.POST("/tags_edit", v1.EditTag)
		apiv1.GET("/tags_del/:id", v1.DeleteTag)
		apiv1.GET("/tags/:id", v1.ShowTag)
		apiv1.GET("/article_show/:id", v1.ShowArticle)
		apiv1.POST("/article_list", v1.GetArticles)
		apiv1.POST("/article_add", v1.AddArticle)
		apiv1.POST("/article_edit", v1.EditArticle)
		apiv1.GET("/article_del/:id", v1.DelArticle)
	}

	return r

}
