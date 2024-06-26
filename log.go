package xlog

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gogf/gf/v2/os/glog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

type log struct {
	ContextName string `json:"context_name"` // 项目名
	logger      *zap.Logger
}
type LogCore struct {
	AdditionalInfo map[string]interface{} `json:"additional_info"` // 附加信息
	Message        string                 `json:"message"`
	Err            string                 `json:"error"`
	Context        context.Context        `json:"context"` //上下文信息
	Stack          string                 `json:"stack"`   // 日志堆栈
	HookFunc       func(c *LogCore)       `json:"-"`       // 钩子函数
	Level          string                 `json:"-"`       // 日志级别
}

var name string
var path string
var zapLogger *zap.Logger

const (
	INFO_LEVEL  = "info"
	DEBUG_LEVEL = "debug"
	WARN_LEVEL  = "warn"
	ERROR_LEVEL = "error"
	FATAL_LEVEL = "fatal"
)

var levelMap map[string]zapcore.Level = map[string]zapcore.Level{
	"info":  zap.InfoLevel,
	"debug": zap.DebugLevel,
	"warn":  zap.WarnLevel,
	"error": zap.ErrorLevel,
	"fatal": zap.FatalLevel,
}

func InitConfig(projectName string, logPath string, level string) error {
	name = projectName
	path = logPath

	if _, ok := levelMap[strings.ToLower(strings.TrimSpace(level))]; ok {
		zapLogger = initLogger(logPath, levelMap[level])
	} else {
		zapLogger = nil
		return fmt.Errorf("未知的日志级别: %s", level)
	}

	return nil
}

func initLogger(path string, level zapcore.Level) *zap.Logger {
	logWriter, _ := rotatelogs.New(
		path+"/"+name+"-%Y%m%d.log",
		rotatelogs.WithLinkName(path+".log"),
		rotatelogs.WithMaxAge(24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)

	// Create a file encoder
	fileEncoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
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
	})

	// Create a zap logCore
	logCore := zapcore.NewCore(
		fileEncoder,
		zapcore.AddSync(logWriter),
		level,
	)

	teeCore := zapcore.NewTee(
		logCore,
	)

	// Create a logger
	return zap.New(teeCore, zap.AddCaller())
}

func logInit(level string) *log {
	if name == "" {
		panic(errors.New("未设置项目名称"))
	}
	if path == "" {
		panic(errors.New("未设置日志路径"))
	}

	_ = level

	logs := new(log)
	logs.ContextName = name
	logs.logger = zapLogger
	return logs
}
func New() *LogCore {
	core := new(LogCore)
	core.Context = context.TODO()
	core.HookFunc = nil
	return core
}

func (c *LogCore) SetHookFunc(f func(c *LogCore)) *LogCore {
	c.HookFunc = f
	return c
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
	logger := logs.logger
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
		c.Level = "info"
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
	logger := logs.logger
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
		c.Level = "error"
		c.HookFunc(c)
	}

	logger.Sync()
	return c

}
func (c *LogCore) Warn(msg string) *LogCore {
	c.Message = msg
	logs := logInit("warn")
	logger := logs.logger

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
		c.Level = "warn"
		c.HookFunc(c)
	}

	logger.Sync()
	return c
}
func (c *LogCore) Debug(msg string) *LogCore {
	c.Message = msg
	logs := logInit("debug")
	logger := logs.logger
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
		c.Level = "debug"
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
