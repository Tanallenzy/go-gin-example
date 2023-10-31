package tag_service

import (
	"encoding/json"
	"github.com/Eden/go-gin-example/models"
	"github.com/Eden/go-gin-example/pkg/export"
	"github.com/Eden/go-gin-example/pkg/gredis"
	"github.com/Eden/go-gin-example/pkg/logging"
	"github.com/Eden/go-gin-example/service/cache_service"
	"io"
	"strconv"
)

type Tag struct {
	ID         int
	Name       string
	State      int
	CreatedBy  string
	ModifiedBy string

	PageNum  int
	PageSize int
}

func (t *Tag) Count() (int, error) {
	return models.GetTagTotal(t.getMaps())
}
func (t *Tag) ExistByID() (bool, error) {
	return models.ExistTagByID(t.ID)
}

func (t *Tag) Get(checkExist bool) (*models.Tag, error) {
	if checkExist {
		check, err := t.ExistByID()
		if err != nil {
			logging.Error(err)
		}
		if !check {
			return nil, err
		}
	}
	var cacheTag *models.Tag
	cache := cache_service.Tag{ID: t.ID}
	key := cache.GetTagKey()
	if gredis.Exists(key) {
		data, err := gredis.Get(key)
		if err != nil {
			logging.Info(err)
		} else {
			json.Unmarshal(data, &cacheTag)
			return cacheTag, nil
		}
	}
	tag, err := models.ShowTagByID(t.ID)
	if err != nil {
		return nil, err
	}
	gredis.Set(key, tag, 3600)
	return tag, nil
}

func (t *Tag) GetAll() ([]*models.Tag, error) {
	var tags []*models.Tag

	tags, err := models.GetTags(t.PageNum, t.PageSize, t.getMaps())
	if err != nil {
		return nil, err
	}

	return tags, nil
}

func (t *Tag) getMaps() map[string]interface{} {
	var maps = make(map[string]interface{})
	maps["deleted_on"] = 0
	if t.State != -1 {
		maps["state"] = t.State
	}
	return maps
}

func (t *Tag) Add() error {
	tagID, err := models.AddTag(t.Name, t.State, t.CreatedBy)
	if err != nil {
		return err
	}
	t.CacheByID(tagID)
	return nil
}

func (t *Tag) Edit() error {
	data := map[string]interface{}{
		"name":        t.Name,
		"state":       t.State,
		"modified_by": t.ModifiedBy,
	}
	err := models.EditTag(t.ID, data)
	if err != nil {
		return err
	}
	t.CacheByID(t.ID)
	return nil
}

func (t *Tag) CacheByID(tagID int) error {
	cache := cache_service.Tag{ID: tagID}
	key := cache.GetTagKey()
	tag, err := models.ShowTagByID(tagID)
	if err != nil {
		return err
	}
	err = gredis.Set(key, tag, 3600)
	if err != nil {
		logging.Error(err)
	}
	return err
}

func (t *Tag) Delete() (bool, error) {
	check, err := t.ExistByID()
	if err != nil {
		logging.Error(err)
	}
	if !check {
		return false, err
	}
	err = models.DeleteTag(t.ID)
	if err != nil {
		logging.Error(err)
		return false, err
	}
	cache := cache_service.Tag{ID: t.ID}
	key := cache.GetTagKey()
	gredis.Delete(key)
	return true, nil
}

func (t *Tag) Export() (string, error) {
	tags, err := t.GetAll()
	if err != nil {
		return "", err
	}
	titles := []string{"ID", "名称", "创建人", "创建时间", "修改人", "修改时间"}
	keys := []string{"id", "name", "created_by", "created_on", "modified_by", "modified_on"}
	var setData []map[string]string
	for _, tag := range tags {
		values := make(map[string]string)
		values["id"] = strconv.Itoa(tag.ID)
		values["name"] = tag.Name
		values["created_by"] = tag.CreatedBy
		values["created_on"] = strconv.Itoa(tag.CreatedOn)
		values["modified_by"] = tag.ModifiedBy
		values["modified_on"] = strconv.Itoa(tag.ModifiedOn)
		setData = append(setData, values)
	}
	return export.ExcelExport("tag", "标签信息", titles, keys, setData)
}

func (t *Tag) Import(r io.Reader) error {
	rows, err := export.ExcelImport(r, "标签信息")
	if err != nil {
		return err
	}
	for irow, row := range rows { //irow行数
		if irow > 0 {
			var data []string
			for _, cell := range row {
				data = append(data, cell)
			}

			models.AddTag(data[1], 1, data[2])
		}
	}

	return nil
}
