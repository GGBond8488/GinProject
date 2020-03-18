package version1

import (
	"My-gin-Project/pkg/app"
	"My-gin-Project/pkg/e"
	"My-gin-Project/pkg/logging"
	"My-gin-Project/pkg/setting"
	"My-gin-Project/pkg/util"
	"My-gin-Project/service/tag_service"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"net/http"
)

//获取多个文章标签
//c *gin.Context是Gin很重要的组成部分，可以理解为上下文，
//它允许我们在中间件之间传递变量、管理流、验证请求的JSON和呈现JSON响应


// @Summary Get multiple article tags
// @Produce  json
// @Param name query string false "Name"
// @Param state query int false "State"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/tags [get]
func GetTags(c *gin.Context) {
	//c.Query可用于获取?name=test&state=1这类URL参数，
	//而c.DefaultQuery则支持设置一个默认值
	appG := app.Gin{C: c}
	//valid := validation.Validation{}

	name := c.Query("name")
	state := -1
	if args := c.Query("state"); args != "" {
		state = com.StrTo(args).MustInt()
	}
	logging.Info(state)
	//valid.Required(name,"name")
	//valid.Range(state,0,1,"state")

	//if valid.HasErrors(){
	//	app.MarkErrors(valid.Errors)
	//	appG.Response(http.StatusBadRequest,e.INVALID_PARAMS,nil)
	//	return
	//}
	tagService := tag_service.Tag{
		Name:       name,
		State:      state,
		PageNum:    util.GetPage(c),
		PageSize:   setting.AppSetting.PageSize,
	}
	total,err := tagService.Count()
	if err != nil{
		appG.Response(http.StatusInternalServerError,e.ERROR_COUNT_TAG_FAIL,nil)
		return
	}
	tags,err:= tagService.GetAll()
	if err != nil{
		appG.Response(http.StatusInternalServerError,e.ERROR_GET_TAGS_FAIL,nil)
		return
	}
	data := make(map[string]interface{})
	data["total"] = total
	data["tags"] = tags

	appG.Response(http.StatusOK,e.SUCCESS,data)
	//var ok1,ok2 bool
 	//maps := make(map[string]interface{})
 	//data := make(map[string]interface{})
	//
 	//maps["delete_on"] = 0
	//
 	//if name!=""{
 	//	maps["name"] = name
	//}
	//var state int = 1
	//if args := c.Query("state");args!=""{
	//	state = com.StrTo(args).MustInt()
	//	maps["state"] = state
	//}
	//
	//code := e.ERROR_DATABASE_EXCEPTION
	//data["lists"],ok1 = models.GetTags(util.GetPage(c),setting.AppSetting.PageSize,maps)
	//data["total"],ok2 = models.GetTagTotal(maps)
	//if ok1&&ok2{
	//	code = e.SUCCESS
	//}
	//c.JSON(http.StatusOK,gin.H{
	//	"code" 	:code,
	//	"msg"	:e.GetMsg(code),
	//	"data"	:data,
	//})
}

// @Summary Add article tag
// @Produce  json
// @Param name body string true "Name"
// @Param state body int false "State"
// @Param created_by body int false "CreatedBy"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/tags [post]
type AddTagForm struct {
	Name      string `form:"name" valid:"Required;MaxSize(100)"`
	CreatedBy string `form:"created_by" valid:"Required;MaxSize(100)"`
	State     int    `form:"state" valid:"Range(0,1)"`
}
func AddTag(c *gin.Context) {
	appG := app.Gin{c}
	form := AddTagForm{}
	httpCode,errCode := app.BindAndValid(c,&form)
	if errCode!=e.SUCCESS{
		appG.Response(httpCode,errCode,nil)
	}
	tagService := tag_service.Tag{
		Name:       form.Name,
		CreatedBy:  form.CreatedBy,
		State:      form.State,
	}

	exists,err := tagService.ExistByName()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}
	if exists {
		appG.Response(http.StatusOK, e.ERROR_EXIST_TAG, nil)
		return
	}
	err = tagService.Add()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_TAG_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
	//name := c.Query("name")
	//state := com.StrTo(c.DefaultQuery("state", "0")).MustInt()
	//createdBy := c.Query("created_by")
	//
	//valid := validation.Validation{}
	//valid.Required(name, "name").Message("名称不能为空")
	//valid.MaxSize(name, 100, "name").Message("名称最长为100字符")
	//valid.Required(createdBy, "created_by").Message("创建人不能为空")
	//valid.MaxSize(createdBy, 100, "created_by").Message("创建人最长为100字符")
	//valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	//
	//code := e.INVALID_PARAMS
	//
	//if ! valid.HasErrors() {
	//	if ! models.ExistTagByName(name)&&models.AddTag(name, state, createdBy) {
	//		code = e.SUCCESS
	//	} else {
	//		code = e.ERROR_EXIST_TAG
	//	}
	//}
	//c.JSON(http.StatusOK, gin.H{
	//	"code" : code,
	//	"msg" : e.GetMsg(code),
	//	"data" : make(map[string]string),
	//})

}

//修改文章标签
type EditTagFrom struct {
	ID 				int 	`form:"id" valid:"Required;Min(1)"`
	Name			string	`form:"name" valid:"Required;MaxSize(100)"`
	ModifiedBy 		string	`form:"modified_by" valid:"Required;MaxSize(100)"`
	State 			int		`form:"state" valid:"Range(0,1)"`
}

// @Summary Update article tag
// @Produce  json
// @Param id path int true "ID"
// @Param name body string true "Name"
// @Param state body int false "State"
// @Param modified_by body string true "ModifiedBy"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/tags/{id} [put]
func EditTag(c *gin.Context) {
	appG := app.Gin{c}
	form := EditTagFrom{}
	httpCode,errCode := app.BindAndValid(c,&form)
	if errCode!=e.SUCCESS{
		appG.Response(httpCode,errCode,nil)
	}
	tagService := tag_service.Tag{
		ID:         form.ID,
		Name:       form.Name,
		ModifiedBy: form.ModifiedBy,
		State:      form.State,
	}
	exists,err := tagService.ExistByID()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}

	if !exists {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}

	err = tagService.Edit()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EDIT_TAG_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
	//id := com.StrTo(c.Param("id")).MustInt()
	//name := c.Query("name")
	//modifiedBy := c.Query("modified_by")
	//
	//valid := validation.Validation{}
	//
	//var state int = -1
	//if arg := c.Query("state"); arg != "" {
	//	state = com.StrTo(arg).MustInt()
	//	valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	//}
	//
	//valid.Required(id, "id").Message("ID不能为空")
	//valid.Required(modifiedBy, "modified_by").Message("修改人不能为空")
	//valid.MaxSize(modifiedBy, 100, "modified_by").Message("修改人最长为100字符")
	//valid.MaxSize(name, 100, "name").Message("名称最长为100字符")
	//
	//code := e.INVALID_PARAMS
	//if ! valid.HasErrors() {
	//	code = e.SUCCESS
	//	if models.ExistTagByID(id) {
	//		data := make(map[string]interface{})
	//		data["modified_by"] = modifiedBy
	//		if name != "" {
	//			data["name"] = name
	//		}
	//		if state != -1 {
	//			data["state"] = state
	//		}
	//		if !models.EditTag(id, data){
	//			code = e.ERROR_NOT_EXIST_TAG
	//		}
	//	} else {
	//		code = e.ERROR_NOT_EXIST_TAG
	//	}
	//}
	//
	//c.JSON(http.StatusOK, gin.H{
	//	"code" : code,
	//	"msg" : e.GetMsg(code),
	//	"data" : make(map[string]string),
	//})
}

//删除文章标签

// @Summary Delete article tag
// @Produce  json
// @Param id path int true "ID"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/tags/{id} [delete]
func DeleteTag(c *gin.Context) {
	appG := app.Gin{c}
	id := com.StrTo(c.Param("id")).MustInt()
	//
	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	if valid.HasErrors(){
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusBadRequest,e.INVALID_PARAMS,nil)
		return
	}
	tagService := tag_service.Tag{
		ID:         id,
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

	if err := tagService.Delete(); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_DELETE_TAG_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
	//
	//code := e.INVALID_PARAMS
	//if ! valid.HasErrors() {
	//	code = e.SUCCESS
	//	if models.ExistTagByID(id)&&models.DeleteTag(id) {
	//
	//	} else {
	//		code = e.ERROR_NOT_EXIST_TAG
	//	}
	//}
	//
	//c.JSON(http.StatusOK, gin.H{
	//	"code" : code,
	//	"msg" : e.GetMsg(code),
	//	"data" : make(map[string]string),
	//})
}


