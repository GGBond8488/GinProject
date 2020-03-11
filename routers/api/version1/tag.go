package version1

import (
	"My-gin-Project/models"
	"My-gin-Project/pkg/e"
	"My-gin-Project/pkg/setting"
	"My-gin-Project/pkg/util"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"net/http"
)

//获取多个文章标签
//c *gin.Context是Gin很重要的组成部分，可以理解为上下文，
//它允许我们在中间件之间传递变量、管理流、验证请求的JSON和呈现JSON响应



func GetTags(c *gin.Context)  {
	//c.Query可用于获取?name=test&state=1这类URL参数，
	//而c.DefaultQuery则支持设置一个默认值
 	name := c.Query("name")
	var ok1,ok2 bool
 	maps := make(map[string]interface{})
 	data := make(map[string]interface{})

 	maps["delete_on"] = 0

 	if name!=""{
 		maps["name"] = name
	}
	var state int = 1
	if args := c.Query("state");args!=""{
		state = com.StrTo(args).MustInt()
		maps["state"] = state
	}

	code := e.ERROR_DATABASE_EXCEPTION
	data["lists"],ok1 = models.GetTags(util.GetPage(c),setting.AppSetting.PageSize,maps)
	data["total"],ok2 = models.GetTagTotal(maps)
	if ok1&&ok2{
		code = e.SUCCESS
	}
	c.JSON(http.StatusOK,gin.H{
		"code" 	:code,
		"msg"	:e.GetMsg(code),
		"data"	:data,
	})
}

//新增文章标签


func AddTag(c *gin.Context) {
	name := c.Query("name")
	state := com.StrTo(c.DefaultQuery("state", "0")).MustInt()
	createdBy := c.Query("created_by")

	valid := validation.Validation{}
	valid.Required(name, "name").Message("名称不能为空")
	valid.MaxSize(name, 100, "name").Message("名称最长为100字符")
	valid.Required(createdBy, "created_by").Message("创建人不能为空")
	valid.MaxSize(createdBy, 100, "created_by").Message("创建人最长为100字符")
	valid.Range(state, 0, 1, "state").Message("状态只允许0或1")

	code := e.INVALID_PARAMS

	if ! valid.HasErrors() {
		if ! models.ExistTagByName(name)&&models.AddTag(name, state, createdBy) {
			code = e.SUCCESS
		} else {
			code = e.ERROR_EXIST_TAG
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code" : code,
		"msg" : e.GetMsg(code),
		"data" : make(map[string]string),
	})

}

//修改文章标签

func EditTag(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()
	name := c.Query("name")
	modifiedBy := c.Query("modified_by")

	valid := validation.Validation{}

	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	}

	valid.Required(id, "id").Message("ID不能为空")
	valid.Required(modifiedBy, "modified_by").Message("修改人不能为空")
	valid.MaxSize(modifiedBy, 100, "modified_by").Message("修改人最长为100字符")
	valid.MaxSize(name, 100, "name").Message("名称最长为100字符")

	code := e.INVALID_PARAMS
	if ! valid.HasErrors() {
		code = e.SUCCESS
		if models.ExistTagByID(id) {
			data := make(map[string]interface{})
			data["modified_by"] = modifiedBy
			if name != "" {
				data["name"] = name
			}
			if state != -1 {
				data["state"] = state
			}
			if !models.EditTag(id, data){
				code = e.ERROR_NOT_EXIST_TAG
			}
		} else {
			code = e.ERROR_NOT_EXIST_TAG
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code" : code,
		"msg" : e.GetMsg(code),
		"data" : make(map[string]string),
	})
}

//删除文章标签


func DeleteTag(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	code := e.INVALID_PARAMS
	if ! valid.HasErrors() {
		code = e.SUCCESS
		if models.ExistTagByID(id)&&models.DeleteTag(id) {

		} else {
			code = e.ERROR_NOT_EXIST_TAG
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code" : code,
		"msg" : e.GetMsg(code),
		"data" : make(map[string]string),
	})
}


