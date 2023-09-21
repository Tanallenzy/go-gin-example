package v1

import (
	"github.com/Eden/go-gin-example/pkg/app"
	"github.com/Eden/go-gin-example/pkg/e"
	"github.com/Eden/go-gin-example/pkg/setting"
	"github.com/Eden/go-gin-example/pkg/util"
	"github.com/Eden/go-gin-example/service/article_service"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"net/http"
)

func ShowArticle(c *gin.Context) {
	appG := app.Gin{c}
	id := com.StrTo(c.Param("id")).MustInt()
	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	app.MarkErrorsAndExit(valid, appG)

	articleService := article_service.Article{ID: id}

	article, err := articleService.Get(true)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_ARTICLE_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, article)

}

func GetArticles(c *gin.Context) {
	var appG = app.Gin{c}
	var articleService = article_service.Article{
		PageSize: setting.AppSetting.PageSize,
		PageNum:  util.GetPage(c),
		State:    -1,
		TagID:    -1,
	}
	valid := validation.Validation{}
	if arg := c.PostForm("state"); arg != "" {
		articleService.State = com.StrTo(arg).MustInt()
		valid.Range(articleService.State, 0, 1, "state").Message("状态只允许0或1")
	}
	if arg := c.PostForm("tag_id"); arg != "" {
		articleService.TagID = com.StrTo(arg).MustInt()
		valid.Min(articleService.TagID, 1, "tag_id").Message("标签ID必须大于0")
	}

	app.MarkErrorsAndExit(valid, appG)

	total, err := articleService.Count()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_COUNT_ARTICLE_FAIL, nil)
		return
	}

	articles, err := articleService.GetAll()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_ARTICLES_FAIL, nil)
		return
	}
	data := make(map[string]interface{})
	data["lists"] = articles
	data["total"] = total
	appG.Response(http.StatusOK, e.SUCCESS, data)
}

type AddArticleForm struct {
	TagID         int    `form:"tag_id" valid:"Required;Min(1)"`
	Title         string `form:"title" valid:"Required;MaxSize(100)"`
	Desc          string `form:"desc" valid:"Required;MaxSize(255)"`
	Content       string `form:"content" valid:"Required;MaxSize(65535)"`
	CreatedBy     string `form:"created_by" valid:"Required;MaxSize(100)"`
	CoverImageUrl string `form:"cover_image_url" valid:"Required;MaxSize(255)"`
	State         int    `form:"state" valid:"Range(0,1)"`
}

func AddArticle(c *gin.Context) {
	appG := app.Gin{c}
	var form AddArticleForm
	app.BindAndValidAndExit(c, &form)

	articleService := article_service.Article{
		TagID:         form.TagID,
		Title:         form.Title,
		Desc:          form.Desc,
		Content:       form.Content,
		CreatedBy:     form.CreatedBy,
		CoverImageUrl: form.CoverImageUrl,
		State:         form.State,
	}
	err := articleService.Add()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_ARTICLE_FAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)

}

type EditArticleForm struct {
	ID            int    `form:"id" valid:"Required;Min(1)"`
	TagID         int    `form:"tag_id" valid:"Required;Min(1)"`
	Title         string `form:"title" valid:"Required;MaxSize(100)"`
	Desc          string `form:"desc" valid:"Required;MaxSize(255)"`
	Content       string `form:"content" valid:"Required;MaxSize(65535)"`
	ModifiedBy    string `form:"modified_by" valid:"Required;MaxSize(100)"`
	CoverImageUrl string `form:"cover_image_url" valid:"Required;MaxSize(255)"`
	State         int    `form:"state" valid:"Range(0,1)"`
}

func EditArticle(c *gin.Context) {
	var appG = app.Gin{c}
	var form EditArticleForm

	app.BindAndValidAndExit(c, &form)

	articleService := article_service.Article{
		ID:            form.ID,
		TagID:         form.TagID,
		Title:         form.Title,
		Desc:          form.Desc,
		Content:       form.Content,
		ModifiedBy:    form.ModifiedBy,
		CoverImageUrl: form.CoverImageUrl,
		State:         form.State,
	}

	err := articleService.Edit()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EDIT_ARTICLE_FAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)

}

func DelArticle(c *gin.Context) {
	appG := app.Gin{c}
	articleService := article_service.Article{
		ID: com.StrTo(c.Param("id")).MustInt(),
	}
	valid := validation.Validation{}
	valid.Min(articleService.ID, 1, "id").Message("ID必须大于0")
	app.MarkErrorsAndExit(valid, appG)
	check, err := articleService.Delete()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_DELETE_ARTICLE_FAIL, nil)
	}
	if !check {
		appG.Response(http.StatusOK, e.ERROR_DELETE_ARTICLE_FAIL, nil)
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
