package logging

import (
	"My-gin-Project/pkg/setting"
	"fmt"
	"log"
	"os"
	"time"
)

func getLogFilePath() string {
	return fmt.Sprintf("%s", setting.AppSetting.LogSavePath)
}

func getLogFileFullPath() string {
	prefixPath := getLogFilePath()
	suffixPath := fmt.Sprintf("%s%s.%s", setting.AppSetting.LogSaveName, time.Now().Format(setting.AppSetting.TimeFormat), setting.AppSetting.LogFileExt)
	return fmt.Sprintf("%s%s", prefixPath, suffixPath)
}

func openLogFile(file string) *os.File {
	// Stat returns a FileInfo describing the named file.
	_, err := os.Stat(file)
	switch {
	case os.IsNotExist(err):
		mkDir()
	case os.IsPermission(err):
		log.Fatalf("Permission:%v", err)
	}
	// OpenFile is the generalized open call; most users will use Open
	// or Create instead.
	logFile, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	//O_RDONLY int = syscall.O_RDONLY // open the file read-only.
	//O_WRONLY int = syscall.O_WRONLY // open the file write-only.
	//O_RDWR   int = syscall.O_RDWR   // open the file read-write.
	//// The remaining values may be or'ed in to control behavior.
	//O_APPEND int = syscall.O_APPEND // append data to the file when writing.
	//O_CREATE int = syscall.O_CREAT  // create a new file if none exists.
	//O_EXCL   int = syscall.O_EXCL   // used with O_CREATE, file must not exist.
	//O_SYNC   int = syscall.O_SYNC   // open for synchronous I/O.
	//O_TRUNC  int = syscall.O_TRUNC  // truncate regular writable file when o
	if err != nil {
		log.Fatalf("fail to open file :%v", err)
	}
	return logFile
}

// Getwd returns a rooted path name corresponding to the
// current directory.
func mkDir() {
	dir, _ := os.Getwd()
	err := os.MkdirAll(dir+"/"+getLogFilePath(), os.ModePerm)
	if err != nil {
		panic(err)
	}
}
