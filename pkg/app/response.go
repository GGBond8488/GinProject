package app

import (
	"My-gin-Project/pkg/e"
	"github.com/gin-gonic/gin"
)

type Gin struct {
	C *gin.Context
}
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func (g *Gin)Response(httpCode,errCode int,data interface{}) {
	g.C.JSON(httpCode,gin.H{
		"code" : errCode,
		"msg":e.GetMsg(errCode),
		"data":data,
	})

	return
}