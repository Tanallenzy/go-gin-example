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

func (a *Article) ExistByID() (bool, error) {
	return models.ExistArticleByID(a.ID)
}

func (a *Article) Get() (*models.Article, error) {
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