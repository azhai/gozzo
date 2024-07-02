package logging

import (
	"fmt"
	"io"
	"net/url"
	"path/filepath"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Output 输出配置
type Output struct {
	Start, Stop string
	OutPaths    []string
}

// LogConfig 日志配置
type LogConfig struct {
	zap.Config
	MinLevel   string
	LevelCase  string
	TimeFormat string
	Outputs    []Output
}

// NewLogger 指定日志目录，普通和错误日志分文件存放
func NewLogger(dir string) *zap.SugaredLogger {
	return NewLoggerCustom(DefaultConfig(), dir)
}

// NewLoggerURL 单路径日志，指定日志级别和路径
func NewLoggerURL(level, url string) *zap.SugaredLogger {
	if strings.Contains(url, "$FILE") {
		return NewLoggerCustom(SingleFileConfig(level, ""), url)
	}
	return NewLoggerCustom(SingleFileConfig(level, url), "")
}

// NewLoggerCustom 根据配置产生记录器
func NewLoggerCustom(cfg *LogConfig, dir string) *zap.SugaredLogger {
	zl, err := cfg.BuildLogger(dir)
	if err == nil {
		return zl.Sugar()
	}
	panic(err)
}

// DefaultConfig 默认配置，使用两个文件分别记录警告和错误
func DefaultConfig() *LogConfig {
	return &LogConfig{
		Config: zap.Config{
			Encoding:         "console",
			OutputPaths:      []string{},
			ErrorOutputPaths: []string{"stderr"}, // zap内部错误输出
		},
		MinLevel:   "debug",
		LevelCase:  "cap",
		TimeFormat: "2006-01-02 15:04:05",
		Outputs: []Output{
			{Start: "debug", Stop: "info", OutPaths: []string{"access.log"}},
			{Start: "warn", Stop: "fatal", OutPaths: []string{"error.log"}},
		},
	}
}

// SingleFileConfig 使用单个文件的记录器
func SingleFileConfig(level, file string) *LogConfig {
	cfg := DefaultConfig()
	cfg.MinLevel = level
	cfg.Outputs = []Output{
		{Start: level, Stop: "fatal", OutPaths: []string{file}},
	}
	return cfg
}

// BuildLogger 生成日志记录器
func (c *LogConfig) BuildLogger(dir string, opts ...zap.Option) (*zap.Logger, error) {
	if c.IsNop() {
		return zap.NewNop(), nil
	}
	if c.BuildLevel().Enabled(zapcore.InfoLevel) {
		c.Development = true
		c.Sampling = nil
	}
	dir = strings.TrimSpace(dir)
	if cores := c.BuildCores(dir); len(cores) > 1 {
		opts = append(opts, ReplaceCores(cores))
	}
	return c.Config.Build(opts...)
}

// IsNop 是否空日志
func (c *LogConfig) IsNop() bool {
	return len(c.Outputs) == 0 && len(c.OutputPaths) == 0
}

// BuildLevel 当前日志的最低级别
func (c *LogConfig) BuildLevel() zap.AtomicLevel {
	var level zapcore.Level
	c.MinLevel, level = GetZapLevel(c.MinLevel)
	c.Level = zap.NewAtomicLevelAt(level)
	return c.Level
}

// BuildCores 产生记录器内核
func (c *LogConfig) BuildCores(dir string) []zapcore.Core {
	var (
		cores []zapcore.Core
		ws    zapcore.WriteSyncer
		err   error
	)
	enc := c.GetEncoder()
	for _, out := range c.Outputs {
		enabler := GetLevelEnabler(out.Start, out.Stop, c.MinLevel)
		if enabler == nil || len(out.OutPaths) == 0 {
			continue
		}
		c.OutputPaths = GetLogPath(dir, out.OutPaths)
		if len(c.OutputPaths) == 0 || c.OutputPaths[0] == "/dev/null" {
			ws = zapcore.AddSync(io.Discard)
		} else if ws, _, err = zap.Open(c.OutputPaths...); err != nil {
			continue
		}
		cores = append(cores, zapcore.NewCore(enc, ws, enabler))
	}
	return cores
}

// GetEncoder 根据编码配置设置日志格式
func (c *LogConfig) GetEncoder() zapcore.Encoder {
	c.Config.EncoderConfig = NewEncoderConfig(c.TimeFormat, c.LevelCase)
	if strings.ToLower(c.Encoding) == "json" {
		return zapcore.NewJSONEncoder(c.Config.EncoderConfig)
	}
	return zapcore.NewConsoleEncoder(c.Config.EncoderConfig)
}

// ReplaceCores 替换为多种输出的Core
func ReplaceCores(cores []zapcore.Core) zap.Option {
	return zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return zapcore.NewTee(cores...)
	})
}

// GetLogPath 使用绝对路径
func GetLogPath(dir string, files []string) []string {
	if dir = strings.TrimSpace(dir); dir == "/dev/null" {
		return nil
	}
	var err error
	for i, file := range files {
		if dir == "" && strings.HasPrefix(file, "std") {
			files[i] = file
			continue
		}
		if strings.Contains(dir, "$FILE") {
			file = strings.Replace(dir, "$FILE", file, 1)
		} else if dir != "" {
			file = filepath.Join(dir, file)
		}
		if file, err = GetAbsPath(file, false); err == nil {
			files[i] = file
		}
	}
	return files
}

// GetAbsPath 使用真实的绝对路径
func GetAbsPath(file string, onlyFile bool) (path string, err error) {
	var u *url.URL
	if u, err = url.Parse(file); err != nil {
		return
	}
	var scheme string
	if scheme = u.Scheme; scheme == "" {
		scheme = "file"
	}
	if onlyFile && scheme == "file" {
		path = file
		return // 只能处理文件类型
	}
	// 去掉scheme，重新解析，以便接着处理相对路径问题
	if u.Scheme != "" {
		file = file[len(u.Scheme+"://"):]
	}
	u, err = url.Parse(file)
	path, _ = filepath.Abs(u.Path)
	path = ignoreWinDisk(path)
	if scheme != "file" { // 重新拼接
		path = fmt.Sprintf("%s://%s?%s", scheme, path, u.RawQuery)
	}
	return
}
