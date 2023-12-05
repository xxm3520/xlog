# xlog
## 基于zap封装的log组件
### 快速开始
>go get -u github.com/xxm3520/xlog

也可以直接import引用然后执行 go mod tidy 或者 go mod vendor

## 使用示例

#### 正常使用

```go
package main

import (
	"errors"
	"github.com/xxm3520/xlog"
)

func main() {
	//	初始化配置
	xlog.InitConfig("项目名", "./log")
	xlog.New().SetAdditionalInfo("test", "这里可以写任意内容").Info("这是一条Info信息")
	err := errors.New("新建一个错误信息")
	xlog.New().SetAdditionalInfo("test", "这里可以写任意内容").Error("这是一条err信息", err)
	xlog.New().SetAdditionalInfo("test", "这里可以写任意内容").Warn("这是一条warn信息")
	xlog.New().SetAdditionalInfo("test", "这里可以写任意内容").Debug("这是一条debug信息")
	xlog.New().SetAdditionalInfo("test", "这里可以写任意内容").Info("输出日志并打印到控制台").Print()
	xlog.New().SetAdditionalInfo("test", "这里可以写任意内容").Info("输出日志并打印到控制台，换行打印").Println()
}

```

####  自定义个钩子

使用场景主要就是在输出日志到文件的同时，需要根据日志内容做一些其他的操作，比如发送邮件，发送短信，或者在控制台上输出彩色的日志级别，更加有利于CI中的调试等等

```go
func PrintColorfulLog(c *xlog.LogCore) *xlog.LogCore {

	LOG_BUFFER := ""
	if c.level == "error" {
		LOG_BUFFER = fmt.Sprintf("Date: %s | Level: %s | Message: %s | AddtionalInfo: %+v | Error: %s",
			time.Now().Format("2006-01-02 15:04:05.000"),
			c.level,
			c.Message,
			c.AdditionalInfo,
			c.Err)
	} else {
		LOG_BUFFER = fmt.Sprintf("Date: %s | Level: %s | Message: %s | AddtionalInfo: %+v",
			time.Now().Format("2006-01-02 15:04:05.000"),
			c.level,
			c.Message,
			c.AdditionalInfo)
	}

	if c.level == "error" {
		// 把LOG BUFFER 标记为红色 在控制台上输出
		fmt.Printf("\033[31m%s\033[0m\n", LOG_BUFFER)
	}

	if c.level == "warn" {
		// 把LOG BUFFER 标记为黄色 在控制台上输出
		fmt.Printf("\033[33m%s\033[0m\n", LOG_BUFFER)
	}

	if c.level == "info" {
		// 把LOG BUFFER 标记为绿色 在控制台上输出
		fmt.Printf("\033[32m%s\033[0m\n", LOG_BUFFER)
	}

	if c.level == "debug" {
		// 把LOG BUFFER 标记为蓝色 在控制台上输出
		fmt.Printf("\033[34m%s\033[0m\n", LOG_BUFFER)
	}
}




//	初始化配置
xlog.InitConfig("项目名", "./log")
xlog.New().SetHookFunc(PrintColorfulLog).SetAdditionalInfo("test", "这里可以写任意内容").Info("这是一条Info信息")
xlog.New().SetHookFunc(PrintColorfulLog).SetAdditionalInfo("test", "这里可以写任意内容").Debug("这是一条Debug信息")
xlog.New().SetHookFunc(PrintColorfulLog).SetAdditionalInfo("test", "这里可以写任意内容").Error("这是一条Error信息")
```



## 鸣谢
https://github.com/uber-go/zap

https://github.com/gogf/gf