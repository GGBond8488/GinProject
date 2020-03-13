package version1

import (
	"My-gin-Project/pkg/app"
	"My-gin-Project/pkg/e"
	"My-gin-Project/pkg/setting"
	"My-gin-Project/pkg/util"
	"My-gin-Project/service/article_service"
	"My-gin-Project/service/tag_service"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"

	"net/http"
)

// @Summary Get a single article
// @Produce  json
// @Param id path int true "ID"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/articles/{id} [get]
func GetArticle(c *gin.Context) {
	appG := app.Gin{c}
	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")
	if valid.HasErrors(){
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusBadRequest,e.INVALID_PARAMS,nil)
		return
	}
	articleService := article_service.Article{ID:id}
	exists,err := articleService.ExistByID()
	if err!=nil{
		appG.Response(http.StatusOK,e.ERROR_CHECK_EXIST_ARTICLE_FAIL,nil)
		return
	}
	if !exists{
		appG.Response(http.StatusOK,e.ERROR_NOT_EXIST_ARTICLE,nil)
		return
	}
	article ,err := articleService.Get()
	if err != nil{
		appG.Response(http.StatusInternalServerError,e.ERROR_GET_ARTICLE_FAIL,nil)
	}
	appG.Response(http.StatusOK,e.SUCCESS,article)

	//if !valid.HasErrors() {
	//	if models.ExistArticleByID(id) {
	//		data,ok = models.GetArticle(id)
	//		if ok {
	//			code = e.SUCCESS
	//		}else {
	//			code = e.ERROR_DATABASE_EXCEPTION
	//		}
	//	} else {
	//		code = e.ERROR_NOT_EXIST_ARTICLE
	//	}
	//} else {
	//	for _, err := range valid.Errors {
	//		logging.Info(err.Key+"  "+err.Message)
	//	}
	//}
	//
	//c.JSON(http.StatusOK, gin.H{
	//	"code": code,
	//	"msg":  e.GetMsg(code),
	//	"data": data,
	//})
}


// @Summary Get multiple articles
// @Produce  json
// @Param tag_id body int false "TagID"
// @Param state body int false "State"
// @Param created_by body int false "CreatedBy"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/articles [get]
func GetArticles(c *gin.Context) {
	appG := app.Gin{c}
	valid := validation.Validation{}

	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	}
	var tagId int = -1
	if arg := c.Query("tag_id"); arg != "" {
		tagId = com.StrTo(arg).MustInt()
		valid.Min(tagId, 1, "tag_id").Message("标签ID必须大于0")
	}
	var CreatedBys string
	if CreatedBys := c.Query("created_by");CreatedBys!=""{
		valid.MaxSize(CreatedBys,100,"created_by")
	}
	if valid.HasErrors(){
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusBadRequest,e.INVALID_PARAMS,nil)
		return
	}
	articleService := article_service.Article{
		State: state,
		TagID:	tagId,
		PageSize:setting.AppSetting.PageSize,
		PageNum:util.GetPage(c),
		CreatedBy: CreatedBys,
	}
	total,err := articleService.Count()
	if err!= nil{
		appG.Response(http.StatusInternalServerError,e.ERROR_GET_ARTICLES_FAIL,nil)
		return
	}
	articles,err := articleService.GetAll()
	if err!= nil{
		appG.Response(http.StatusInternalServerError,e.ERROR_GET_ARTICLES_FAIL,nil)
		return
	}
	data := make(map[string]interface{})
	data["lists"] = articles
	data["total"] = total
	appG.Response(http.StatusOK,e.SUCCESS,data)
	//code := e.INVALID_PARAMS
	//if !valid.HasErrors() {
	//	code = e.ERROR_DATABASE_EXCEPTION
	//	if ok {
	//		data["lists"], ok = models.GetArticles(util.GetPage(c), setting.AppSetting.PageSize, maps)
	//
	//	}
	//	if ok{
	//		data["total"],ok = models.GetTotalAritical(maps)
	//		code = e.SUCCESS
	//	}
	//
	//} else {
	//	for _, err := range valid.Errors {
	//		logging.Info(err.Key+"  "+err.Message)
	//	}
	//}
	//
	//c.JSON(http.StatusOK, gin.H{
	//	"code": code,
	//	"msg":  e.GetMsg(code),
	//	"data": data,
	//})
}

// @Summary Add article
// @Produce  json
// @Param tag_id body int true "TagID"
// @Param title body string true "Title"
// @Param desc body string true "Desc"
// @Param content body string true "Content"
// @Param created_by body string true "CreatedBy"
// @Param state body int true "State"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/articles [post]
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
	form := AddArticleForm{}
	httpCode,errCode := app.BindAndValid(c,&form)
	if errCode != e.SUCCESS{
		appG.Response(httpCode,errCode,nil)
		return
	}
	tagService :=tag_service.Tag{
		ID: form.TagID,
	}
	exists, err := tagService.ExistByID()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}

	if !exists {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}
	articleService := article_service.Article{
		TagID:         form.TagID,
		Title:         form.Title,
		Desc:          form.Desc,
		Content:       form.Content,
		CoverImageUrl: form.CoverImageUrl,
		State:         form.State,
	}
	if err := articleService.Add(); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_ARTICLE_FAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
	//data := make(map[string]interface{})
	//tagId := com.StrTo(c.Query("tag_id")).MustInt()
	//title := c.Query("title")
	//desc := c.Query("desc")
	//content := c.Query("content")
	//createdBy := c.Query("created_by")
	//state := com.StrTo(c.DefaultQuery("state", "0")).MustInt()
	//image := c.Query("image")
	//file,image,err := c.Request.FormFile("image")
	//if err != nil {
	//	logging.Warn(err)
	//	code := e.ERROR
	//	c.JSON(http.StatusOK,gin.H{
	//		"code":code,
	//		"msg":e.GetMsg(code),
	//		"data":data,
	//	})
	//}
	//if image == nil {
	//	code = e.INVALID_PARAMS
	//}else {
	//	imageName := upload.GetImageName(image.Filename)
	//	fullPath := upload.GetImageFullPath()
	//	savePath := upload.GetImagePath()
	//
	//	src := fullPath+imageName
	//
	//	if !upload.CheckImageExt(imageName)||!upload.CheckImageSize(file){
	//		code = e.ERROR_UPLOAD_CHECK_IMAGE_FORMAT
	//	}else {
	//		err := upload.CheckImage(fullPath)
	//		if err != nil {
	//			logging.Warn(err)
	//			code = e.ERROR_UPLOAD_CHECK_IMAGE_FAIL
	//		}else if err := c.SaveUploadedFile(image,src);err != nil{
	//			logging.Warn(err)
	//			code = e.ERROR_UPLOAD_SAVE_IMAGE_FAIL
	//		}else {
	//			//data["image_url"] = upload.GetImageFullUrl(imageName)
	//			data["image_save_url"] = savePath + imageName
	//		}
	//	}
	//}
	//valid := validation.Validation{}
	//valid.Min(tagId, 1, "tag_id").Message("标签ID必须大于0")
	//valid.Required(title, "title").Message("标题不能为空")
	//valid.Required(desc, "desc").Message("简述不能为空")
	//valid.Required(content, "content").Message("内容不能为空")
	//valid.Required(createdBy, "created_by").Message("创建人不能为空")
	//valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	////valid.Required(image,"image").Message("cover_image_url不能为空")
	//valid.Max(len(image),1024,"image").Message("cover_image_url长度不能超过1024")

	//if valid.HasErrors() {
	//	if models.ExistTagByID(tagId) {
	//		data["tag_id"] = tagId
	//		data["title"] = title
	//		data["desc"] = desc
	//		data["content"] = content
	//		data["created_by"] = createdBy
	//		data["state"] = state
	//		data["cover_image_url"] = image
	//		if models.AddArticle(data){
	//			code = e.SUCCESS
	//		}else {
	//			code = e.ERROR_DATABASE_EXCEPTION
	//		}
	//
	//	} else {
	//		code = e.ERROR_NOT_EXIST_TAG
	//	}
	//} else {
	//	for _, err := range valid.Errors {
	//		logging.Info(err.Key+"  "+err.Message)
	//	}
	//}
	//
	//c.JSON(http.StatusOK, gin.H{
	//	"code": code,
	//	"msg":  e.GetMsg(code),
	//	"data": make(map[string]interface{}),
	//})
}

//修改文章

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
// @Summary Update article
// @Produce  json
// @Param id path int true "ID"
// @Param tag_id body string false "TagID"
// @Param title body string false "Title"
// @Param desc body string false "Desc"
// @Param content body string false "Content"
// @Param modified_by body string true "ModifiedBy"
// @Param state body int false "State"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/articles/{id} [put]
func EditArticle(c *gin.Context) {
	appG := app.Gin{c}
	form := EditArticleForm{}
	httpCode,errCode := app.BindAndValid(c,&form)
	if errCode!=e.SUCCESS{
		appG.Response(httpCode,errCode,nil)
		return
	}
	articleService := article_service.Article{
		ID:            form.ID,
		TagID:         form.TagID,
		Title:         form.Title,
		Desc:          form.Desc,
		Content:       form.Content,
		CoverImageUrl: form.CoverImageUrl,
		State:         form.State,
	}
	exists,err:= articleService.ExistByID()
	if err != nil{
		appG.Response(http.StatusInternalServerError,e.ERROR_CHECK_EXIST_ARTICLE_FAIL,nil)
		return
	}
	if !exists{
		appG.Response(http.StatusOK,e.ERROR_NOT_EXIST_ARTICLE,nil)
		return
	}
	tagService := tag_service.Tag{ID: form.TagID}
	exists, err = tagService.ExistByID()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}

	if !exists {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}

	err = articleService.Edit()
	if err != nil{
		appG.Response(http.StatusOK,e.ERROR_EDIT_ARTICLE_FAIL,nil)
		return
	}
	appG.Response(http.StatusOK,e.SUCCESS,nil)
	//id := com.StrTo(c.Param("id")).MustInt()
	//tagId := com.StrTo(c.Query("tag_id")).MustInt()
	//title := c.Query("title")
	//desc := c.Query("desc")
	//content := c.Query("content")
	//modifiedBy := c.Query("modified_by")
	//image := c.Query("image")
	//
	//var state int = -1
	//if arg := c.Query("state"); arg != "" {
	//	state = com.StrTo(arg).MustInt()
	//	valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	//}
	//
	//valid.Min(id, 1, "id").Message("ID必须大于0")
	//valid.MaxSize(title, 100, "title").Message("标题最长为100字符")
	//valid.MaxSize(desc, 255, "desc").Message("简述最长为255字符")
	//valid.MaxSize(content, 65535, "content").Message("内容最长为65535字符")
	//valid.Required(modifiedBy, "modified_by").Message("修改人不能为空")
	//valid.MaxSize(modifiedBy, 100, "modified_by").Message("修改人最长为100字符")
	//valid.Max(len(image),1024,"image").Message("cover_image_url长度不能超过1024")
	//
	//code := e.INVALID_PARAMS
	//if !valid.HasErrors() {
	//	if models.ExistArticleByID(id) {
	//		if models.ExistTagByID(tagId) {
	//			data := make(map[string]interface{})
	//			if tagId > 0 {
	//				data["tag_id"] = tagId
	//			}
	//			if title != "" {
	//				data["title"] = title
	//			}
	//			if desc != "" {
	//				data["desc"] = desc
	//			}
	//			if content != "" {
	//				data["content"] = content
	//			}
	//
	//			data["modified_by"] = modifiedBy
	//			data["cover_image_url"] = image
	//			if models.EditArticle(id, data) {
	//				code = e.SUCCESS
	//			}else {
	//				code = e.ERROR_DATABASE_EXCEPTION
	//			}
	//		} else {
	//			code = e.ERROR_NOT_EXIST_TAG
	//		}
	//	} else {
	//		code = e.ERROR_NOT_EXIST_ARTICLE
	//	}
	//} else {
	//	for _, err := range valid.Errors {
	//		logging.Info(err.Key+"  "+err.Message)
	//	}
	//}
	//
	//c.JSON(http.StatusOK, gin.H{
	//	"code": code,
	//	"msg":  e.GetMsg(code),
	//	"data": make(map[string]string),
	//})
}

//删除文章

// @Summary Delete article
// @Produce  json
// @Param id path int true "ID"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/articles/{id} [delete]
func DeleteArticle(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()
	appG := app.Gin{c}
	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")
	if valid.HasErrors(){
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusBadRequest,e.INVALID_PARAMS,nil)
		return
	}
	articleService := article_service.Article{ID:id}
	exists,err:= articleService.ExistByID()
	if err != nil{
		appG.Response(http.StatusInternalServerError,e.ERROR_CHECK_EXIST_ARTICLE_FAIL,nil)
		return
	}
	if !exists{
		appG.Response(http.StatusOK,e.ERROR_NOT_EXIST_ARTICLE,nil)
		return
	}
	err = articleService.Del()
	if err != nil{
		appG.Response(http.StatusInternalServerError,e.ERROR_DELETE_ARTICLE_FAIL,nil)
		return
	}
	appG.Response(http.StatusOK,e.SUCCESS,nil)
	//code := e.INVALID_PARAMS
	//if !valid.HasErrors() {
	//	if models.ExistArticleByID(id) {
	//		code = e.ERROR_DATABASE_EXCEPTION
	//		if models.DeleteArticle(id) {
	//			code = e.SUCCESS
	//		}
	//	} else {
	//		code = e.ERROR_NOT_EXIST_ARTICLE
	//	}
	//} else {
	//	for _, err := range valid.Errors {
	//		logging.Info(err.Key, err.Message)
	//	}
	//}
	//
	//c.JSON(http.StatusOK, gin.H{
	//	"code": code,
	//	"msg":  e.GetMsg(code),
	//	"data": make(map[string]string),
	//})
}
