package routers

import (
	_ "My-gin-Project/docs"
	"My-gin-Project/middleware/jwt"
	"My-gin-Project/pkg/setting"
	"My-gin-Project/routers/api"
	"My-gin-Project/routers/api/version1"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func InitRouter() *gin.Engine {
	/*
	func Default() *Engine {
		debugPrintWARNINGDefault()
		engine := New()
		engine.Use(Logger(), Recovery())
		return engine
	}
	*/
	// New returns a new blank Engine instance without any middleware attached.
	r := gin.New()
	// Use将全局中间件附加到路由器。 即通过Use（）连接的中间件将是
	//包含在每个单个请求的处理程序链中。 甚至404、405，静态文件...
	//例如，这是记录器或错误管理中间件的正确位置。
	r.Use(gin.Logger())
	//Recover返回从panic中回复的中间件，并写入500
	r.Use(gin.Recovery())

	gin.SetMode(setting.ServerSetting.RunMode)

	r.GET("/auth", api.GetAuth)
	r.GET("/swagger/*any",ginSwagger.WrapHandler(swaggerFiles.Handler))

	apiv1 := r.Group("/api/v1")
	apiv1.Use(jwt.JWT())
	{	// url /api/v1/tags
		//获取标签列表
		apiv1.GET("/tags", version1.GetTags)
		//新建标签
		apiv1.POST("/tags", version1.AddTag)
		//更新指定标签
		apiv1.PUT("/tags/:id", version1.EditTag)
		//删除指定标签
		apiv1.DELETE("/tags/:id", version1.DeleteTag)


		//获取文章列表
		apiv1.GET("/articles", version1.GetArticles)
		//获取指定文章
		apiv1.GET("/articles/:id", version1.GetArticle)
		//新建文章
		apiv1.POST("/articles", version1.AddArticle)
		//更新指定文章
		apiv1.PUT("/articles/:id", version1.EditArticle)
		//删除指定文章
		apiv1.DELETE("/articles/:id", version1.DeleteArticle)

	}

	return r
}