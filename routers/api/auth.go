package api

import (
	"My-gin-Project/pkg/app"
	"My-gin-Project/pkg/e"
	"My-gin-Project/pkg/util"
	"My-gin-Project/service/auth_service"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"net/http"
)

type auth struct {
	Username string `valid:"Required;MaxSize(50)"`
	Password string `valid:"Required;MaxSize(50)"`
}

// @Summary Get Auth
// @Produce  json
// @Param username query string true "userName"
// @Param password query string true "password"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /auth [get]
func GetAuth(c *gin.Context) {
	appG := app.Gin{C: c}
	username := c.Query("username")
	password := c.Query("password")

	valid := validation.Validation{}

	a := auth{Username: username, Password: password}
	valid.Valid(&a)
	if valid.HasErrors(){
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	authService := auth_service.Auth{Username: username, Password: password}
	isExist, err := authService.Check()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}

	if !isExist {
		appG.Response(http.StatusUnauthorized, e.ERROR_AUTH, nil)
		return
	}

	token, err := util.GenerateToken(username, password)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_AUTH_TOKEN, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"token": token,
	})
	//data := make(map[string]interface{})
	//code := e.INVALID_PARAMS
	//
	//if ok {
	//	isExist := models.CheckAuth(username, password)
	//
	//	if isExist {
	//		token, err := util.GenerateToken(username, password)
	//		if err != nil {
	//			code = e.ERROR_AUTH_TOKEN
	//		} else {
	//			data["token"] = token
	//
	//			code = e.SUCCESS
	//		}
	//	} else {
	//		for _, err := range valid.Errors {
	//			logging.Info(err.Key+"  "+err.Message)
	//		}
	//	}
	//	// JSON将给定结构作为JSON序列化到响应主体中。
	//	//还将Content-Type设置为“ application / json”。
	//	//type H map[string]interface{}
	//	c.JSON(http.StatusOK, gin.H{
	//		"code": code,
	//		"msg":  e.GetMsg(code),
	//		"data": data,
	//	})
	//}
}
