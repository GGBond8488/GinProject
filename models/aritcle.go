package models

import "github.com/jinzhu/gorm"

type Article struct {
	Model
	//作为主键
	TagID int `json:"tag_id" gorm:"index"`
	//Tag字段，实际是一个嵌套的struct，
	//它利用TagID与Tag模型相互关联，在执行查询的时候，
	//能够达到Article、Tag关联查询的功能
	Tag   Tag `json:"tag"`

	Title         string `json:"title"`
	Desc          string `json:"desc"`
	Content       string `json:"content"`
	CreatedBy     string `json:"created_by"`
	ModifiedBy    string `json:"modified_by"`
	State         int    `json:"state"`
	CoverImageUrl string `json:"cover_image_url"`
}

func ExistArticleByID(id int)(bool,error) {
	var article Article
	err := db.Select("id").Where("id=? AND delete_on=?",id,0).First(&article).Error

	if err!=nil&&err != gorm.ErrRecordNotFound {
		return false,err
	}

	if article.ID>0{
		return true,nil
	}
	return false,nil
}

func GetTotalAritical(maps interface{}) (count int, err error) {
	err = db.Model(&Article{}).Where(maps).Count(&count).Error
	if err!=nil&&err != gorm.ErrRecordNotFound {
		return count,err
	}
	return count,nil
}
//Preload就是一个预加载器，它会执行两条SQL，分别是SELECT * FROM blog_articles;和SELECT * FROM blog_tag WHERE id IN (1,2,3,4);，
//那么在查询出结构后，gorm内部处理对应的映射逻辑，将其填充到Article的Tag中，会特别方便，并且避免了循环查询
func GetArticles(pageNum int,pageSize int,maps interface{})(articles []*Article,err error)  {
	err = db.Preload("Tag").Where(maps).Offset(pageNum).Limit(pageSize).Find(articles).Error
	if err!=nil&&err != gorm.ErrRecordNotFound {
		return articles,err
	}
	return articles,nil
}
//Article有一个结构体成员是TagID，就是外键。gorm会通过类名+ID的方式去找到这两个类之间的关联关系
//Article有一个结构体成员是Tag，就是我们嵌套在Article里的Tag结构体，我们可以通过Related进行关联查询

func GetArticle(id int)(article *Article,err error)  {
	err = db.Where("id = ?",id).First(article).Error
	if err!=nil&&err != gorm.ErrRecordNotFound {
		return nil,err
	}
	err = db.Model(article).Related(article.Tag).Error
	if err!=nil&&err != gorm.ErrRecordNotFound {
		return nil,err
	}
	//db.Model(&user).Related(&profile)
	////// SELECT * FROM profiles WHERE id = 111; // 111是user的外键ProfileID
	return
}

func EditArticle(id int,data interface{})error {
	return db.Model(&Article{}).Where("id = ?",id).Updates(data).Error

}

func DeleteArticle(id int)error {
	return db.Where("id =?", id).Delete(Article{}).Error
}

func AddArticle(data map[string]interface {}) error {
	return  db.Create(&Article {
		TagID : data["tag_id"].(int),
		Title : data["title"].(string),
		Desc : data["desc"].(string),
		Content : data["content"].(string),
		CreatedBy : data["created_by"].(string),
		State : data["state"].(int),
		CoverImageUrl:data["cover_image_url"].(string),
	}).Error

}


func DeleteAllAriticle() error {
	return db.Unscoped().Where("delete_on != ?",0).Delete(&Tag{}).Error

}


//func (article *Article) BeforeCreate(scope *gorm.Scope) error {
//	scope.SetColumn("CreatedOn", time.Now().Unix())
//
//	return nil
//}
//
//func (article *Article) BeforeUpdate(scope *gorm.Scope) error {
//	scope.SetColumn("ModifiedOn", time.Now().Unix())
//
//	return nil
//}
