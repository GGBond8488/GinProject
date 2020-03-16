package api

import (
	"My-gin-Project/pkg/e"
	"My-gin-Project/pkg/logging"
	"My-gin-Project/pkg/upload"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UploadImage(c *gin.Context){
	code := e.SUCCESS
	data := make(map[string]string)
	//拿到文件
	//func (r *Request) FormFile(key string) (multipart.File, *multipart.FileHeader, error)
	file,image,err := c.Request.FormFile("image")
	if err != nil {
		logging.Warn(err)
		code = e.ERROR
		c.JSON(http.StatusOK,gin.H{
			"code":code,
			"msg":e.GetMsg(code),
			"data":data,
		})
	}
	if image == nil {
		code = e.INVALID_PARAMS
	}else {
		imageName := upload.GetImageName(image.Filename)
		fullPath := upload.GetImageFullPath()
		savePath := upload.GetImagePath()

		src := fullPath+imageName

		if !upload.CheckImageExt(imageName)||!upload.CheckImageSize(file){
			code = e.ERROR_UPLOAD_CHECK_IMAGE_FORMAT
		}else {
			err := upload.CheckImage(fullPath)
			if err != nil {
				logging.Warn(err)
				code = e.ERROR_UPLOAD_CHECK_IMAGE_FAIL
				//存储文件到指定位置
			}else if err := c.SaveUploadedFile(image,src);err != nil{
				logging.Warn(err)
				code = e.ERROR_UPLOAD_SAVE_IMAGE_FAIL
			}else {
				data["image_url"] = upload.GetImageFullUrl(imageName)
				data["image_save_url"] = savePath + imageName
			}
		}
	}

	c.JSON(http.StatusOK,gin.H{
		"code":code,
		"msg":e.GetMsg(code),
		"data":data,
	})
}