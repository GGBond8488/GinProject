package article_service

import (
	"My-gin-Project/models"
	"My-gin-Project/pkg/gredis"
	"My-gin-Project/pkg/logging"
	"My-gin-Project/service/cache_service"
	"encoding/json"
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

func (article *Article)ExistByID()(bool,error)  {
	return models.ExistArticleByID(article.ID)
}

func (article *Article)GetAll()( []*models.Article, error)  {
	var (
		res, cacheRes []*models.Article
	)
	cache := cache_service.Article{
		TagID:    article.TagID,
		State:    article.State,
		PageNum:  article.PageNum,
		PageSize: article.PageSize,
	}
	key := cache.GetArticlesKey()
	if gredis.Exists(key){
		data,err := gredis.Get(key)
		if err != nil{
			logging.Info(err)
		}else {
			json.Unmarshal(data,&cacheRes)
			return cacheRes,nil
		}
	}
	res,err := models.GetArticles(article.PageNum,article.PageSize,article.getMaps())
	if err!=nil{
		return nil,err
	}
	gredis.Set(key,res,3600)
	return res,nil
}


func (article *Article)Get()(*models.Article, error)  {
	var cacheRes *models.Article
	cache := cache_service.Article{ID:article.ID}
	key := cache.GetArticleKey()
	if gredis.Exists(key){
		data,err := gredis.Get(key)
		if err != nil{
			logging.Error(err)
		}else {
			json.Unmarshal(data,&cacheRes)
			return cacheRes,nil
		}
	}
	res,err := models.GetArticle(article.ID)
	if err != nil{
		logging.Error(err)
		return nil,err
	}
	gredis.Set(key,res,3600)
	return res,nil
}

func (article *Article)Add()error{
		art := map[string]interface{}{
		"tag_id":          article.TagID,
		"title":           article.Title,
		"desc":            article.Desc,
		"content":         article.Content,
		"created_by":      article.CreatedBy,
		"cover_image_url": article.CoverImageUrl,
		"state":           article.State,
	}
	return models.AddArticle(art)
}

func (article *Article)Edit()error  {
	art:=models.Article{
		Model:         models.Model{ID:article.ID},
		TagID:         article.TagID,
		Title:         article.Title,
		Desc:          article.Desc,
		Content:       article.Content,
		ModifiedBy:    article.ModifiedBy,
		State:         article.State,
		CoverImageUrl: article.CoverImageUrl,
	}
	//logging.Info(art)
	return models.EditArticle(article.ID,art)
}

func (article *Article)Del()error  {
	return models.DeleteArticle(article.ID)
}

func (article *Article)Count()(int,error)  {
	return models.GetTotalAritical(article.getMaps())
}

func (a *Article) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	maps["delete_on"] = 0
	if a.State != -1 {
		maps["state"] = a.State
	}
	if a.TagID != -1 {
		maps["tag_id"] = a.TagID
	}
	if a.CreatedBy !=""{
		maps["created_by"] = a.CreatedBy
	}
	return maps
}