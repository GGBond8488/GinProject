package jwt

import (
	"My-gin-Project/pkg/e"
	"My-gin-Project/pkg/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

/*
权限验证的中间件
*/

func JWT() gin.HandlerFunc{
	return func (c *gin.Context){
		var code int
		var data interface{}

		code = e.SUCCESS
		token:=c.Query("token")
		if token == ""{
			code = e.INVALID_PARAMS
		}else {
			claims,err := util.ParseToken(token)
			if err != nil{
				code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
			}else if time.Now().Unix()>claims.ExpiresAt{
				code = e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
			}
		}

		if code != e.SUCCESS{
			c.JSON(http.StatusUnauthorized,gin.H{
				"code" : code,
				"msg" : e.GetMsg(code),
				"data" : data,
			})
			//Abort阻止挂起的处理程序被调用。不会停止当前的处理程序。
			//假设您有一个授权中间件，用于验证当前请求是否得到授权。
			//如果授权失败（例如：密码不匹配），请调用中止以确保不调用此请求的其余处理程序。
			c.Abort()
			return
		}
		//Next只能在中间件内部使用。
		//它在调用处理程序内的链中执行挂起的处理程序。
		//这里如果鉴权成功则继续业务逻辑
		c.Next()
	}
}