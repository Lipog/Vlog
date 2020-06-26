package vlog

import (
	"fmt"
	"os"
	"path"
	"time"
)

//该包的作用是往文件里写日志

type FileLogger struct {
	Level LogLevel
	filePath string
	fileName string
	fileObj *os.File
	errFileObj *os.File
	maxFileSize int64
}

//NewFileLogger构造函数
func NewFileLogger(levelStr, fp, fn string, maxSize int64) *FileLogger {
	level, err := paraseLogLevel(levelStr)
	//如果无法解析日志等级字符串，那么就panic
	if err != nil {
		panic(err)
	}
	
	f1 := &FileLogger{
		Level:       level,
		filePath:    fp,
		fileName:    fn,
		maxFileSize: maxSize,
	}

	//按照文件路径和文件名打开
	err = f1.initFile()
	//如果文件都打不开，就要panic了
	if err != nil {
		panic(err)
	}
	return f1
}

func (f *FileLogger) initFile() error {
	fullFileName := path.Join(f.filePath, f.fileName)
	fileObj, err := os.OpenFile(fullFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("open log file failed", err)
		return err
	}

	errFileObj, err := os.OpenFile(fullFileName + ".err", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("open err log file failed", err)
		return err
	}

	//可以走到这一步，说明日志文件都打开了
	f.fileObj = fileObj
	f.errFileObj = errFileObj
	return nil
}

func (f *FileLogger) enable(loglevel LogLevel) bool {
	return f.Level <= loglevel
}

func (f *FileLogger) checkSize(file *os.File) bool {
	info, err := file.Stat()
	if err != nil {
		fmt.Println("get file info err: ", err)
	}
	return info.Size() >= f.maxFileSize
}


func (f *FileLogger)log(level LogLevel, format string, a ...interface{}) {
	if f.enable(level) {
		msg := fmt.Sprintf(format, a...)
		funcName, fileName, line := getInfo(3)
		now := time.Now()
		formatTime := now.Format("2006-01-02 15:04:05")
		//分别输出时间，函数名，函数所在的文件名，以及行号
		level_str := getLogString(level)

		//如果要记录的日志的大小超过了设置的大小，那么就需要切割
		file, err := f.splitFile(f.fileObj)
		if err != nil {
			panic(err)
		}
		f.fileObj = file

		fmt.Fprintf(f.fileObj, "[%v] [%v] [%v:%v:%v] %v \n", formatTime, level_str,
			fileName, funcName, line, msg)

		if level >= ERROR {
			file, err := f.splitFile(f.errFileObj)
			if err != nil {
				panic(err)
			}
			f.errFileObj = file
			//如果要记录的日志大于ERROR级别，还要在err日志文件中记录一遍
			fmt.Fprintf(f.errFileObj, "[%v] [%v] [%v:%v:%v] %v \n", formatTime, level_str,
				fileName, funcName, line, msg)
		}
	}
}

func (f *FileLogger) Debug(format string, a ...interface{}) {
	f.log(DEBUG, format, a...)
}
func (f *FileLogger) Info(format string, a ...interface{}) {
	f.log(INFO, format, a...)
}
func (f *FileLogger) Error(format string, a ...interface{}) {
	f.log(ERROR,format, a...)
}
func (f *FileLogger) Waring(format string, a ...interface{}) {
	f.log(WARNING, format, a...)
}
func (f *FileLogger) Fatal(format string, a ...interface{}) {
	f.log(FATAL, format, a...)
}

//关闭使用的数据流
func (f *FileLogger) Close() {
	f.fileObj.Close()
	f.errFileObj.Close()
}