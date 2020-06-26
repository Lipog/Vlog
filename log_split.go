package vlog

import (
	"fmt"
	"os"
	"path"
	"time"
)

func (f *FileLogger) splitFile(file *os.File) (*os.File, error) {
	//0.如果不满足分割的条件，那么就返回原来的file
	if !f.checkSize(file) {
		return file, nil
	}
	//1.既然选择切割了，那么就要把原来的file关闭掉
	//2.然后就是起一个新的名字，首先获取到当前的时间戳
	nowStr := time.Now().Format("20060102150405000")
	//2.1然后获取到旧的文件名
	info, err := file.Stat()
	//如果获取信息错误，就返回错误
	if err != nil {
		fmt.Println("get file info err:", err)
		return nil, err
	}
	//2.2然后获取当前log日志文件的名称,并拿到当前日志文件的完整路径
	oldFileName := path.Join(f.filePath, info.Name())
	newLogName := fmt.Sprintf("%s.backup.%s", oldFileName, nowStr)
	file.Close()
	//所以一定要在重命名之前，进行file.Close()
	//!!!!!!!!!!!!!!!!!要注意重命名文件的时候，要关闭文件的使用!!!!!!!!!!!!!!!!!!!!!!!
	err = os.Rename(oldFileName, newLogName)
	if err != nil {
		fmt.Println("重命名错误:", err)
		return nil, err
	}
	//然后打开一个新的file流，并赋值给f.file
	openFile, err := os.OpenFile(oldFileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("open newfile failed, err:", err)
		return nil, err
	}
	return openFile, nil
}