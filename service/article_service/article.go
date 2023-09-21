package article_service

import (
	"encoding/json"
	"github.com/Eden/go-gin-example/models"
	"github.com/Eden/go-gin-example/pkg/gredis"
	"github.com/Eden/go-gin-example/pkg/logging"
	"github.com/Eden/go-gin-example/service/cache_service"
)

type Article struct {
	ID            int
	TagID         int
	Title         string
	Desc          string
	Content       string
	CoverImageUrl string
	State         int
	CreatedBy     string
	ModifiedBy    string

	PageNum  int
	PageSize int
}

func (a *Article) Count() (int, error) {
	return models.GetArticleTotal(a.getMaps())
}
func (a *Article) ExistByID() (bool, error) {
	return models.ExistArticleByID(a.ID)
}
func (a *Article) Get(checkExist bool) (*models.Article, error) {
	if checkExist {
		check, err := a.ExistByID()
		if err != nil {
			logging.Error(err)
		}
		if !check {
			return nil, err
		}
	}
	var cacheArticle *models.Article
	cache := cache_service.Article{ID: a.ID}
	key := cache.GetArticleKey()
	if gredis.Exists(key) {
		data, err := gredis.Get(key)
		if err != nil {
			logging.Info(err)
		} else {
			json.Unmarshal(data, &cacheArticle)
			return cacheArticle, nil
		}
	}
	article, err := models.GetArticle(a.ID)
	if err != nil {
		return nil, err
	}
	gredis.Set(key, article, 3600)
	return article, nil
}

func (a *Article) GetAll() ([]*models.Article, error) {
	var articles []*models.Article

	articles, err := models.GetArticles(a.PageNum, a.PageSize, a.getMaps())
	if err != nil {
		return nil, err
	}

	return articles, nil
}

func (a *Article) getMaps() map[string]interface{} {
	var maps = make(map[string]interface{})
	maps["deleted_on"] = 0
	if a.State != -1 {
		maps["state"] = a.State
	}
	if a.TagID != -1 {
		maps["tag_id"] = a.TagID
	}
	return maps
}

func (a *Article) Add() error {
	data := map[string]interface{}{
		"tag_id":          a.TagID,
		"title":           a.Title,
		"desc":            a.Desc,
		"content":         a.Content,
		"cover_image_url": a.CoverImageUrl,
		"created_by":      a.CreatedBy,
		"state":           a.State,
		"modified_by":     a.ModifiedBy,
	}
	articleID, err := models.AddArticle(data)
	if err != nil {
		return err
	}
	a.CacheByID(articleID)
	return nil
}

func (a *Article) Edit() error {
	data := map[string]interface{}{
		"tag_id":          a.TagID,
		"title":           a.Title,
		"desc":            a.Desc,
		"content":         a.Content,
		"cover_image_url": a.CoverImageUrl,
		"created_by":      a.CreatedBy,
		"state":           a.State,
		"modified_by":     a.ModifiedBy,
	}
	err := models.EditArticle(a.ID, data)
	if err != nil {
		return err
	}
	a.CacheByID(a.ID)
	return nil
}

func (a *Article) CacheByID(articleID int) error {
	a.ID = articleID
	cache := cache_service.Article{ID: a.ID}
	key := cache.GetArticleKey()
	article, err := models.GetArticle(a.ID)
	if err != nil {
		return err
	}
	err = gredis.Set(key, article, 3600)
	if err != nil {
		logging.Error(err)
	}
	return err
}

func (a *Article) Delete() (bool, error) {
	check, err := a.ExistByID()
	if err != nil {
		logging.Error(err)
	}
	if !check {
		return false, err
	}
	err = models.DelArticle(a.ID)
	if err != nil {
		logging.Error(err)
		return false, err
	}
	cache := cache_service.Article{ID: a.ID}
	key := cache.GetArticleKey()
	gredis.Delete(key)
	return true, nil
}
