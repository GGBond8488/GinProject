package upload

import (
	"My-gin-Project/pkg/file"
	"My-gin-Project/pkg/logging"
	"My-gin-Project/pkg/setting"
	"My-gin-Project/pkg/util"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"path"
	"strings"
)


//GetImageFullUrl：获取图片完整访问URL
//GetImageName：获取图片名称
//GetImagePath：获取图片路径
//GetImageFullPath：获取图片完整路径
//CheckImageExt：检查图片后缀
//CheckImageSize：检查图片大小
//CheckImage：检查图片

func GetImageFullUrl(name string) string {
	return setting.AppSetting.ImagePrefixUrl + "/" + GetImagePath() + name
}

func GetImageName(name string)string{
	ext := path.Ext(name)
	fileName := strings.TrimSuffix(name,ext)
	fileName = util.EncodeMD5(fileName)

	return fileName+ext
}

func GetImagePath()string  {
	return setting.AppSetting.ImageSavePath
}

func GetImageFullPath()string  {
	return setting.AppSetting.RuntimeRootPath+setting.AppSetting.ImageSavePath
}

func CheckImageExt(filename string)bool  {
	ext := path.Ext(filename)
	for _,TrueExt := range setting.AppSetting.ImageAllowExts{
		if strings.ToUpper(TrueExt) == strings.ToUpper(ext){
			return true
		}
	}

	return false
}

func CheckImageSize(f multipart.File)bool{
	size,err := file.GetSize(f)
	if err!=nil{
		log.Println(err)
		logging.Warn(err)
		return false
	}

	return size<=setting.AppSetting.ImageMaxSize
}

func CheckImage(src string) error {
	dir ,err := os.Getwd()
	if err!=nil{
		return fmt.Errorf("os.Getwd err:%v",err)
	}

	err = file.IsNotExistMkDir(dir + "/"+src)
	if err != nil {
		return fmt.Errorf("file.IsNotExistMkDir err:%v",err)
	}

	perm := file.CheckPermission(src)
	if perm{
		return fmt.Errorf("file.CheckPermission Permmission denied src %s",src)
	}
	return nil
}