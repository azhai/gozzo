package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/azhai/gozzo/logging"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsimple"
	. "github.com/klauspost/cpuid/v2"
	"go.uber.org/zap"
)

const MegaByte = 1024 * 1024

var theSettings *RootConfig

// RootConfig 顶层配置，包含其他配置块
type RootConfig struct {
	Debug  bool       `hcl:"debug" json:"debug"`
	App    *AppConfig `hcl:"app,block" json:"app"`
	Log    *LogConfig `hcl:"log,block" json:"log,omitempty"`
	Remain hcl.Body   `hcl:",remain"`
}

// AppConfig App配置，包括App名称和自定义配置
type AppConfig struct {
	Name    string `hcl:"name,optional" json:"name,omitempty"`
	Version string `hcl:"version,optional" json:"version,omitempty"`
}

// LogConfig 日志配置，指定文件夹或URL文件
type LogConfig struct {
	LogLevel string `hcl:"log_level,optional" json:"log_level,omitempty"`
	LogFile  string `hcl:"log_file,optional" json:"log_file,omitempty"`
	LogDir   string `hcl:"log_dir,optional" json:"log_dir,omitempty"`
}

// PrepareEnv 初始化环境
func PrepareEnv(size int) {
	if size > 0 { // 压舱石，阻止频繁GC
		ballast := make([]byte, size*MegaByte)
		runtime.KeepAlive(ballast)
	}

	if level := os.Getenv("GOAMD64"); level == "" {
		level = fmt.Sprintf("v%d", CPU.X64Level())
		fmt.Printf("请设置环境变量 export GOAMD64=%s\n\n", level)
	}
	if IsRunTest() {
		BackToDir(1) // 从tests退回根目录
	}
}

// IsRunTest 是否测试模式下
func IsRunTest() bool {
	return strings.HasSuffix(os.Args[0], ".test")
}

// BackToDir 退回上层目录
func BackToDir(back int) (err error) {
	if back == 0 {
		return
	} else if back < 0 {
		back = 0 - back
	}
	dir := strings.Repeat("../", back)
	if dir, err = filepath.Abs(dir); err == nil {
		err = os.Chdir(dir)
	}
	return
}

// SetupLog 根据配置初始化日志单例
func SetupLog() {
	var logger *zap.SugaredLogger
	cfg := GetLogSettings()
	if cfg.LogFile != "" {
		logger = logging.NewLoggerURL(cfg.LogLevel, cfg.LogFile)
	} else if cfg.LogDir != "" {
		logger = logging.NewLogger(cfg.LogDir)
	} else {
		logger = zap.NewNop().Sugar()
	}
	logging.SetLogger(logger)
}

// ReadConfigFile 读取配置文件
func ReadConfigFile(cfgFile string, verbose bool, others any) (*RootConfig, error) {
	var err error
	if theSettings == nil {
		theSettings = new(RootConfig)
		if verbose {
			fmt.Println("Config file is", cfgFile)
		}
		err = hclsimple.DecodeFile(cfgFile, nil, theSettings)
	}
	if err == nil && others != nil {
		_ = gohcl.DecodeBody(theSettings.Remain, nil, others)
	}
	return theSettings, err
}

// GetTheSettings 返回主配置单例
func GetTheSettings() *RootConfig {
	return theSettings
}

// GetAppSettings 返回App配置单例
func GetAppSettings() *AppConfig {
	if theSettings != nil {
		return theSettings.App
	}
	return new(AppConfig)
}

// GetLogSettings 返回日志配置单例
func GetLogSettings() *LogConfig {
	if theSettings != nil {
		return theSettings.Log
	}
	return new(LogConfig)
}
