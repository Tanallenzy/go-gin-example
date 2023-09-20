package v1

import (
	"github.com/Eden/go-gin-example/models"
	"github.com/Eden/go-gin-example/pkg/app"
	"github.com/Eden/go-gin-example/pkg/e"
	"github.com/Eden/go-gin-example/pkg/logging"
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

	//if valid.HasErrors() {
	//	app.MarkErrors(valid.Errors)
	//	appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
	//	return
	//}
	app.MarkErrorsAndExit(valid, appG)

	articleService := article_service.Article{ID: id}
	exist, err := articleService.ExistByID()

	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_CHECK_EXIST_ARTICLE_FAIL, nil)
		return
	}

	if !exist {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_ARTICLE, nil)
		return
	}

	article, err := articleService.Get()
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

	//if valid.HasErrors() {
	//	app.MarkErrors(valid.Errors)
	//	appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
	//	return
	//}
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
	httpCode, errCode := app.BindAndValid(c, &form)
	articleService := article_service.Article{
		TagID:     com.StrTo(c.PostForm("tag_id")).MustInt(),
		Title:     c.PostForm("title"),
		Desc:      c.PostForm("desc"),
		Content:   c.PostForm("content"),
		CreatedBy: c.PostForm("created_by"),
		State:     com.StrTo(c.DefaultPostForm("state", "0")).MustInt(),
	}
	valid := validation.Validation{}
	valid.Min(articleService.TagID, 1, "tag_id").Message("标签ID必须大于0")
	valid.Required(articleService.Title, "title").Message("标题不能为空")
	valid.Required(articleService.Desc, "desc").Message("简述不能为空")
	valid.Required(articleService.Content, "content").Message("内容不能为空")
	valid.Required(articleService.CreatedBy, "created_by").Message("创建人不能为空")
	valid.Range(articleService.State, 0, 1, "state").Message("状态只允许0或1")
	app.MarkErrorsAndExit(valid, appG)

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		if models.ExistTagByID(tagId) {
			data := make(map[string]interface{})
			data["tag_id"] = tagId
			data["title"] = title
			data["desc"] = desc
			data["content"] = content
			data["created_by"] = createBy
			data["state"] = state
			models.AddArticle(data)
			code = e.SUCCESS
		} else {
			code = e.ERROR_NOT_EXIST_TAG
		}
	} else {
		for _, err := range valid.Errors {
			logging.Info("err.key: %s, err.message: %s", err.Key, err.Message)
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]string),
	})
}

func EditArticle(c *gin.Context) {
	valid := validation.Validation{}

	id := com.StrTo(c.PostForm("id")).MustInt()
	tagId := com.StrTo(c.PostForm("tag_id")).MustInt()
	title := c.PostForm("title")
	desc := c.PostForm("desc")
	content := c.PostForm("content")
	modifiedBy := c.PostForm("modified_by")

	var state int = -1
	if arg := c.PostForm("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	}

	valid.Min(id, 1, "id").Message("ID必须大于0")
	valid.MaxSize(title, 100, "title").Message("标题最长为100字符")
	valid.MaxSize(desc, 255, "desc").Message("简述最长为255字符")
	valid.MaxSize(content, 65535, "content").Message("内容最长为65535字符")
	valid.Required(modifiedBy, "modified_by").Message("修改人不能为空")
	valid.MaxSize(modifiedBy, 100, "modified_by").Message("修改人最长为100字符")

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		if models.ExistArticleByID(id) {
			if models.ExistTagByID(tagId) {
				data := make(map[string]interface{})
				if tagId > 0 {
					data["tag_id"] = tagId
				}
				if title != "" {
					data["title"] = title
				}
				if desc != "" {
					data["desc"] = desc
				}
				if content != "" {
					data["content"] = content
				}

				data["modified_by"] = modifiedBy

				models.EditArticle(id, data)
				code = e.SUCCESS
			} else {
				code = e.ERROR_NOT_EXIST_TAG
			}
		} else {
			code = e.ERROR_NOT_EXIST_ARTICLE
		}
	} else {
		for _, err := range valid.Errors {
			logging.Info("err.key: %s, err.message: %s", err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]string),
	})
}

func DelArticle(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		if models.ExistArticleByID(id) {
			models.DelArticle(id)
			code = e.SUCCESS
		} else {
			code = e.ERROR_NOT_EXIST_ARTICLE
		}
	} else {
		for _, err := range valid.Errors {
			logging.Info("err.key: %s, err.message: %s", err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]string),
	})
}
