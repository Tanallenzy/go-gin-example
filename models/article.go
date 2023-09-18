package models

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

func ExistArticleByID(id int) bool {
	var article Article
	db.Select("id").Where("id = ?", id).Where("deleted_on = 0").First(&article)
	if article.ID > 0 {
		return true
	}
	return false
}

func GetArticleTotal(maps interface{}) (count int) {
	db.Model(&Article{}).Where(maps).Where("deleted_on = 0").Count(&count)
	return
}

func GetArticles(pageNum int, pageSize int, maps interface{}) (articles []Article) {
	db.Preload("Tag").Where(maps).Where("deleted_on = 0").Offset(pageNum).Limit(pageSize).Find(&articles)
	return
}

func GetArticle(id int) (article Article) {
	db.Where("id = ?", id).Where("deleted_on = 0").First(&article)
	db.Model(&article).Related(&article.Tag)
	return
}

func EditArticle(id int, data interface{}) {
	db.Model(&Article{}).Where("id = ?", id).Updates(data)
	return

}

func AddArticle(data map[string]interface{}) bool {
	db.Create(&Article{
		TagID:      data["tag_id"].(int),
		Title:      data["title"].(string),
		Desc:       data["desc"].(string),
		Content:    data["content"].(string),
		CreatedBy:  data["created_by"].(string),
		State:      data["state"].(int),
		ModifiedBy: data["modified_by"].(string),
	})
	return true
}

func DelArticle(id int) bool {
	db.Where("id = ?", id).Delete(&Article{})
	return true

}
