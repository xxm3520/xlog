# xlog
## 基于zap封装的log组件
### 快速开始
>go get -u github.com/xxm3520/xlog

也可以直接import引用然后执行 go mod tidy 或者 go mod vendor

## 使用示例
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
## 鸣谢
https://github.com/uber-go/zap

https://github.com/gogf/gf