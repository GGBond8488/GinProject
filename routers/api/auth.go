package api

import (
	"My-gin-Project/models"
	"My-gin-Project/pkg/e"
	"My-gin-Project/pkg/logging"
	"My-gin-Project/pkg/util"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"

	"net/http"
)

type auth struct {
	Username string `valid:"Required;MaxSize(50)"`
	Password string `valid:"Required;MaxSize(50)"`
}

func GetAuth(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	valid := validation.Validation{}

	a := auth{Username: username, Password: password}
	ok, _ := valid.Valid(&a)

	data := make(map[string]interface{})
	code := e.INVALID_PARAMS

	if ok {
		isExist := models.CheckAuth(username, password)

		if isExist {
			token, err := util.GenerateToken(username, password)
			if err != nil {
				code = e.ERROR_AUTH_TOKEN
			} else {
				data["token"] = token

				code = e.SUCCESS
			}
		} else {
			for _, err := range valid.Errors {
				logging.Info(err.Key+"  "+err.Message)
			}
		}
		// JSON将给定结构作为JSON序列化到响应主体中。
		//还将Content-Type设置为“ application / json”。
		//type H map[string]interface{}
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  e.GetMsg(code),
			"data": data,
		})
	}
}
