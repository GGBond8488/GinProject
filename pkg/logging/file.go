package logging

import (
	"My-gin-Project/pkg/file"
	"My-gin-Project/pkg/setting"
	"fmt"
	"os"
	"time"
)

func getLogFilePath() string {
	return fmt.Sprintf("%s%s",setting.AppSetting.RuntimeRootPath,setting.AppSetting.LogSavePath)
}

func getLogFileName() string {
	return fmt.Sprintf("%s%s.%s",
		setting.AppSetting.LogSaveName,
		time.Now().Format(setting.AppSetting.TimeFormat),
		setting.AppSetting.LogFileExt)
}

func openLogFile(fileName,filePath string) (*os.File,error) {
	dir ,err := os.Getwd()
	if err!=nil{
		return nil,fmt.Errorf("os.Getwd err :%v",err)
	}

	src := dir + "/" +filePath
	perm := file.CheckPermission(src)

	if perm{
		return nil,fmt.Errorf("file.CheckPermission Permission denied src:%s",src)
	}
	err = file.IsNotExistMkDir(src)
	if err!=nil{
		return nil,fmt.Errorf("file.IsNotExistMkDir src: %s, err: %v", src, err)
	}

	f,err := file.Open(src+fileName,os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("Fail to OpenFile :%v", err)
	}
	return f,nil
	// Stat returns a FileInfo describing the named file.
	//_, err := os.Stat(file)
	//switch {
	//case os.IsNotExist(err):
	//	mkDir()
	//case os.IsPermission(err):
	//	log.Fatalf("Permission:%v", err)
	//}
	// OpenFile is the generalized open call; most users will use Open
	// or Create instead.
	//logFile, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	//O_RDONLY int = syscall.O_RDONLY // open the file read-only.
	//O_WRONLY int = syscall.O_WRONLY // open the file write-only.
	//O_RDWR   int = syscall.O_RDWR   // open the file read-write.
	//// The remaining values may be or'ed in to control behavior.
	//O_APPEND int = syscall.O_APPEND // append data to the file when writing.
	//O_CREATE int = syscall.O_CREAT  // create a new file if none exists.
	//O_EXCL   int = syscall.O_EXCL   // used with O_CREATE, file must not exist.
	//O_SYNC   int = syscall.O_SYNC   // open for synchronous I/O.
	//O_TRUNC  int = syscall.O_TRUNC  // truncate regular writable file when o
	//if err != nil {
	//	log.Fatalf("fail to open file :%v", err)
	//}
	//return logFile
}

// Getwd returns a rooted path name corresponding to the
// current directory.
//func mkDir() {
//	dir, _ := os.Getwd()
//	err := os.MkdirAll(dir+"/"+getLogFilePath(), os.ModePerm)
//	if err != nil {
//		panic(err)
//	}
//}
