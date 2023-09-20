package models

import (
	"github.com/jinzhu/gorm"
)

type Article struct {
	Model
	TagID      int    `json:"tag_id"`
	Tag        Tag    `json:"tag"`
	Title      string `json:"title"`
	Desc       string `json:"desc"`
	State      int    `json:"state"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	Content    string `json:"content"`
}

//func (article *Article) BeforeCreate(scope *gorm.Scope) error {
//	scope.SetColumn("CreatedOn", time.Now().Unix())
//	return nil
//
//}
//
//func (article *Article) BeforeUpdate(scope *gorm.Scope) error {
//	scope.SetColumn("ModifiedOn", time.Now().Unix())
//	return nil
//}

func ExistArticleByID(id int) (bool, error) {
	var article Article
	err := db.Select("id").Where("id = ?", id).Where("deleted_on = 0").First(&article).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if article.ID > 0 {
		return true, nil
	}
	return false, nil
}

func GetArticleTotal(maps interface{}) (count int, err error) {
	err = db.Model(&Article{}).Where(maps).Where("deleted_on = 0").Count(&count).Error
	if err != nil {
		return 0, err
	}
	return
}

func GetArticles(pageNum int, pageSize int, maps interface{}) (articles []*Article, err error) {
	err = db.Preload("Tag").Where(maps).Offset(pageNum).Limit(pageSize).Find(&articles).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return articles, nil
}

//func GetArticle(id int) (article Article) {
//	db.Where("id = ?", id).Where("deleted_on = 0").First(&article)
//	db.Model(&article).Related(&article.Tag)
//	return
//}

func GetArticle(id int) (*Article, error) {
	var article Article
	err := db.Where("id = ? AND deleted_on = ? ", id, 0).First(&article).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	err = db.Model(&article).Related(&article.Tag).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &article, nil
}

func EditArticle(id int, data interface{}) {
	db.Model(&Article{}).Where("id = ?", id).Updates(data)
	return

}

func AddArticle(data map[string]interface{}) (bool, error) {
	err := db.Create(&Article{
		TagID:      data["tag_id"].(int),
		Title:      data["title"].(string),
		Desc:       data["desc"].(string),
		Content:    data["content"].(string),
		CreatedBy:  data["created_by"].(string),
		State:      data["state"].(int),
		ModifiedBy: data["modified_by"].(string),
	}).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

func DelArticle(id int) bool {
	db.Where("id = ?", id).Delete(&Article{})
	return true

}
