package xlog

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/os/glog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type log struct {
	ContextName      string `json:"context_name"` // 项目名
	Core             zapcore.Core
	lumberjackLogger *lumberjack.Logger
}
type LogCore struct {
	AdditionalInfo map[string]interface{} `json:"additional_info"` // 附加信息
	Message        string                 `json:"message"`
	Err            string                 `json:"error"`
	Context        context.Context        `json:"context"` //上下文信息
	Stack          string                 `json:"stack"`   // 日志堆栈
	HookFunc       func(c *LogCore)       `json:"-"`       // 钩子函数
}

var name string
var path string

func InitConfig(projectName string, logPath string) {
	name = projectName
	path = logPath
}
func logInit(level string) *log {
	if name == "" {
		panic(errors.New("未设置项目名称"))
	}
	if path == "" {
		panic(errors.New("未设置日志路径"))
	}
	var logLevel zapcore.LevelEnabler
	var fileName string
	switch level {
	case "info":
		logLevel = zap.InfoLevel
		fileName = fmt.Sprintf("%s/%s_info.log", path, time.Now().Format("2006-01-02"))
	case "error":
		logLevel = zap.ErrorLevel
		fileName = fmt.Sprintf("%s/%s_err.log", path, time.Now().Format("2006-01-02"))
	case "warn":
		logLevel = zap.WarnLevel
		fileName = fmt.Sprintf("%s/%s_warn.log", path, time.Now().Format("2006-01-02"))
	default:
		logLevel = zap.DebugLevel
		fileName = fmt.Sprintf("%s/%s_debug.log", path, time.Now().Format("2006-01-02"))

	}
	// 创建Lumberjack实例，用于日志文件的分割
	lumberjackLogger := &lumberjack.Logger{
		Filename:   fileName, // 日志文件路径
		MaxSize:    10,       // 单个日志文件最大大小（MB）
		MaxBackups: 5,        // 最多保留的旧日志文件数量
		MaxAge:     30,       // 保留的旧日志文件的最大天数
		Compress:   true,     // 是否压缩旧日志文件
	}
	// 创建编码器配置，用于将日志格式化为JSON
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(lumberjackLogger),
		logLevel,
	)
	logs := new(log)
	logs.ContextName = name
	logs.Core = core
	logs.lumberjackLogger = lumberjackLogger
	return logs
}
func New() *LogCore {
	core := new(LogCore)
	core.Context = context.TODO()
	core.HookFunc = nil
	return core
}

func (c *LogCore) SetHookFunc(f func(c *LogCore)) {
	c.HookFunc = f
}

func (c *LogCore) getStack() {
	c.Stack = glog.GetStack()
}
func (c *LogCore) SetAdditionalInfo(key string, value interface{}) *LogCore {
	if c.AdditionalInfo == nil {
		c.AdditionalInfo = make(map[string]interface{})
	}
	c.AdditionalInfo[key] = value
	return c
}
func (c *LogCore) Info(msg string) *LogCore {
	c.Message = msg
	logs := logInit("info")
	defer logs.lumberjackLogger.Close()
	logger := zap.New(logs.Core)
	if c.AdditionalInfo == nil {
		logger.Info(
			msg,
			zap.Any("project_name", name),
			zap.Any("log_path", path))
	} else {
		jsons, _ := json.Marshal(c.AdditionalInfo)
		logger.Info(
			msg,
			zap.Any("additional_info", string(jsons)),
			zap.Any("project_name", name),
			zap.Any("log_path", path))
	}

	if c.HookFunc != nil {
		c.HookFunc(c)
	}

	logger.Sync()
	return c
}
func (c *LogCore) Error(msg string, err error) *LogCore {
	c.Message = msg
	if err != nil {
		c.Err = err.Error()
		c.getStack()
	} else {
		c.Stack = ""
		c.Err = ""
	}
	logs := logInit("error")
	defer logs.lumberjackLogger.Close()
	logger := zap.New(logs.Core)
	if c.AdditionalInfo == nil {
		logger.Error(
			msg,
			zap.Any("project_name", name),
			zap.Any("log_path", path))
	} else {
		jsons, _ := json.Marshal(c.AdditionalInfo)
		logger.Error(
			msg,
			zap.Any("additional_info", string(jsons)),
			zap.Any("stack", c.Stack),
			zap.Any("err", c.Err),
			zap.Any("project_name", name),
			zap.Any("log_path", path))
	}

	if c.HookFunc != nil {
		c.HookFunc(c)
	}

	logger.Sync()
	return c

}
func (c *LogCore) Warn(msg string) *LogCore {
	c.Message = msg
	logs := logInit("warn")
	defer logs.lumberjackLogger.Close()
	logger := zap.New(logs.Core)

	if c.AdditionalInfo == nil {
		logger.Warn(
			msg,
			zap.Any("project_name", name),
			zap.Any("log_path", path))
	} else {
		jsons, _ := json.Marshal(c.AdditionalInfo)
		logger.Warn(
			msg,
			zap.Any("additional_info", string(jsons)),
			zap.Any("project_name", name),
			zap.Any("log_path", path))
	}

	if c.HookFunc != nil {
		c.HookFunc(c)
	}

	logger.Sync()
	return c
}
func (c *LogCore) Debug(msg string) *LogCore {
	c.Message = msg
	logs := logInit("debug")
	defer logs.lumberjackLogger.Close()
	logger := zap.New(logs.Core)
	if c.AdditionalInfo == nil {
		logger.Debug(
			msg,
			zap.Any("project_name", name),
			zap.Any("log_path", path))
	} else {
		jsons, _ := json.Marshal(c.AdditionalInfo)
		logger.Debug(
			msg,
			zap.Any("additional_info", string(jsons)),
			zap.Any("project_name", name),
			zap.Any("log_path", path))
	}

	if c.HookFunc != nil {
		c.HookFunc(c)
	}

	logger.Sync()
	return c
}

func (c *LogCore) Print() {
	content, _ := json.Marshal(c)
	fmt.Print(string(content))
}
func (c *LogCore) Println() {
	content, _ := json.Marshal(c)
	fmt.Println(string(content))
}
