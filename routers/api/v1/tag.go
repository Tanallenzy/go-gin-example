package v1

import (
	"github.com/Eden/go-gin-example/pkg/app"
	"github.com/Eden/go-gin-example/pkg/e"
	"github.com/Eden/go-gin-example/pkg/setting"
	"github.com/Eden/go-gin-example/pkg/util"
	"github.com/Eden/go-gin-example/service/tag_service"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"net/http"
)

// 获取多个文章标签
func GetTags(c *gin.Context) {
	appG := app.Gin{c}
	var tagService = tag_service.Tag{
		PageSize: setting.AppSetting.PageSize,
		PageNum:  util.GetPage(c),
		State:    -1,
	}
	valid := validation.Validation{}
	if arg := c.PostForm("state"); arg != "" {
		tagService.State = com.StrTo(arg).MustInt()
		valid.Range(tagService.State, 0, 1, "state").Message("状态只允许0或1")
	}
	if arg := c.PostForm("name"); arg != "" {
		tagService.Name = arg
		valid.MaxSize(tagService.Name, 100, "name").Message("名称最长为100字符")
	}

	app.MarkErrorsAndExit(valid, appG)

	total, err := tagService.Count()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_COUNT_TAG_FAIL, nil)
		return
	}

	tags, err := tagService.GetAll()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_TAGS_FAIL, nil)
		return
	}
	data := make(map[string]interface{})
	data["lists"] = tags
	data["total"] = total
	appG.Response(http.StatusOK, e.SUCCESS, data)

}

// 新增文章标签
type AddTagForm struct {
	Name      string `form:"name" valid:"Required;MaxSize(100)"`
	CreatedBy string `form:"created_by" valid:"Required;MaxSize(100)"`
	State     int    `form:"state" valid:"Range(0,1)"`
}

func AddTag(c *gin.Context) {
	appG := app.Gin{c}
	var form AddTagForm
	app.BindAndValidAndExit(c, &form)

	tagService := tag_service.Tag{
		Name:      form.Name,
		CreatedBy: form.CreatedBy,
		State:     form.State,
	}
	err := tagService.Add()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_TAG_FAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)

}

// 修改文章标签
type EditTagForm struct {
	ID         int    `form:"id" valid:"Required;Min(1)"`
	Name       string `form:"name" valid:"Required;MaxSize(100)"`
	ModifiedBy string `form:"modified_by" valid:"Required;MaxSize(100)"`
	State      int    `form:"state" valid:"Range(0,1)"`
}

func EditTag(c *gin.Context) {
	var appG = app.Gin{c}
	var form EditTagForm

	app.BindAndValidAndExit(c, &form)

	tagService := tag_service.Tag{
		ID:         form.ID,
		Name:       form.Name,
		ModifiedBy: form.ModifiedBy,
		State:      form.State,
	}

	err := tagService.Edit()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EDIT_TAG_FAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)

}

// 删除文章标签
func DeleteTag(c *gin.Context) {
	appG := app.Gin{c}
	tagService := tag_service.Tag{
		ID: com.StrTo(c.Param("id")).MustInt(),
	}
	valid := validation.Validation{}
	valid.Min(tagService.ID, 1, "id").Message("ID必须大于0")
	app.MarkErrorsAndExit(valid, appG)
	check, err := tagService.Delete()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_DELETE_TAG_FAIL, nil)
	}
	if !check {
		appG.Response(http.StatusOK, e.ERROR_DELETE_TAG_FAIL, nil)
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

func ShowTag(c *gin.Context) {
	appG := app.Gin{c}
	id := com.StrTo(c.Param("id")).MustInt()
	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	app.MarkErrorsAndExit(valid, appG)

	tagService := tag_service.Tag{ID: id}

	tag, err := tagService.Get(true)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_ARTICLE_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, tag)
}
