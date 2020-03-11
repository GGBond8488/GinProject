package file

import (
	"io/ioutil"
	"mime/multipart"
	"os"
	"path"
)

//GetSize：获取文件大小
//GetExt：获取文件后缀
//CheckExist：检查文件是否存在
//CheckPermission：检查文件权限
//IsNotExistMkDir：如果不存在则新建文件夹
//MkDir：新建文件夹
//Open：打开文件

func GetSize(f multipart.File)(int,error){
	content,err := ioutil.ReadAll(f)

	return len(content),err
}

func GetExt(filename string)string  {
	return path.Ext(filename)
}
// Stat returns a FileInfo describing the named file.
// If there is an error, it will be of type *PathError.
func CheckNotExist(src string)bool  {
	_,err := os.Stat(src)
	// IsNotExist returns a boolean indicating whether the error is known to
	// report that a file or directory does not exist. It is satisfied by
	// ErrNotExist as well as some syscall errors.
	return os.IsNotExist(err)
}

func CheckPermission(src string)bool  {
	_,err := os.Stat(src)

	return os.IsPermission(err)
}

func IsNotExistMkDir(src string) error{
	if Exist := CheckNotExist(src); Exist{
		if err := MkDir(src);err != nil{
			return err
		}
	}
	return nil
}

func MkDir (src string) error{
	err := os.MkdirAll(src,os.ModePerm)
	if err!=nil{
		return err
	}
	return nil
}

func Open(name string,flag int, perm os.FileMode)(*os.File,error)  {
	f,err := os.OpenFile(name,flag,perm)
	if err!=nil{
		return nil,err
	}
	return f,nil
}