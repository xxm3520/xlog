package main

import (
	"errors"

	"github.com/xxm3520/xlog"
)

func main() {
	//	初始化配置
	xlog.InitConfig("项目名", "./log", xlog.ERROR_LEVEL)
	for {
		xlog.New().SetAdditionalInfo("test", "这里可以写任意内容").Info("这是一条Info信息")
		err := errors.New("新建一个错误信息")
		xlog.New().SetAdditionalInfo("test", "这里可以写任意内容").Error("这是一条err信息", err)
		xlog.New().SetAdditionalInfo("test", "这里可以写任意内容").Warn("这是一条warn信息")
		xlog.New().SetAdditionalInfo("test", "这里可以写任意内容").Debug("这是一条debug信息")
		xlog.New().SetAdditionalInfo("test", "这里可以写任意内容").Info("输出日志并打印到控制台").Print()
		xlog.New().SetAdditionalInfo("test", "这里可以写任意内容").Info("输出日志并打印到控制台，换行打印").Println()
	}

}
