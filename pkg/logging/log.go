package logging

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

type Level int

var (
	F *os.File

	DefaultPrefix      = ""
	DefaultCallerDepth = 2

	logger     *log.Logger
	logPrefix  = ""
	levelFlags = []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}
)

const (
	DEBUG Level = iota
	INFO
	WARNING
	ERROR
	FATAL
)

func Setup() {
	var err error
	filePath := getLogFilePath()
	fileName := getLogFileName()
	F,err = openLogFile(fileName,filePath)
	if err != nil {
		log.Fatalln(err)
	}
	//log.New：创建一个新的日志记录器。out定义要写入日志数据的IO句柄。
	//prefix定义每个生成的日志行的开头。
	//flag定义了日志记录属性
	//func New(out io.Writer, prefix string, flag int) *Logger {
	//	return &Logger{out: out, prefix: prefix, flag: flag}
	//}
	//log.LstdFlags：日志记录的格式属性之一，其余的选项如下
	//const (
	//	Ldate         = 1 << iota     // the date in the local time zone: 2009/01/23
	//	Ltime                         // the time in the local time zone: 01:23:23
	//	Lmicroseconds                 // microsecond resolution: 01:23:23.123123.  assumes Ltime.
	//	Llongfile                     // full file name and line number: /a/b/c/d.go:23
	//	Lshortfile                    // final file name element and line number: d.go:23. overrides Llongfile
	//	LUTC                          // if Ldate or Ltime is set, use UTC rather than the local time zone
	//	LstdFlags     = Ldate | Ltime // initial values for the standard logger
	//)
	logger = log.New(F, DefaultPrefix, log.LstdFlags)
}

func Debug(v ...interface{}) {
	setPrefix(DEBUG)
	logger.Println(v)
}

func Info(v ...interface{}) {
	setPrefix(INFO)
	logger.Println(v)
}

func Warn(v ...interface{}) {
	setPrefix(WARNING)
	logger.Println(v)
}

func Error(v ...interface{}) {
	setPrefix(ERROR)
	logger.Println(v)
}

func Fatal(v ...interface{}) {
	setPrefix(FATAL)
	logger.Fatalln(v)
}

//func Caller(skip int) (pc uintptr, file string, line int, ok bool) {
//	rpc := make([]uintptr, 1)
//	n := callers(skip+1, rpc[:])
//	if n < 1 {
//		return
//	}
//	frame, _ := CallersFrames(rpc).Next()
//	return frame.PC, frame.File, frame.Line, frame.PC != 0
//}

func setPrefix(level Level) {
	// Caller reports file and line number information about function invocations on
	// the calling goroutine's stack.
	//The argument skip is the number of stack frames
	// to ascend
	_, file, line, ok := runtime.Caller(DefaultCallerDepth)
	if ok {
		logPrefix = fmt.Sprintf("[%s][%s:%d]", levelFlags[level], filepath.Base(file), line)
	} else {
		logPrefix = fmt.Sprintf("[%s]", levelFlags[level])
	}
	logger.SetPrefix(logPrefix)
}
