package config

import (
	"fmt"

	"github.com/azhai/gozzo/logging"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"go.uber.org/zap"
)

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
	Name    string   `hcl:"name,optional" json:"name,omitempty"`
	Version string   `hcl:"version,optional" json:"version,omitempty"`
	Remain  hcl.Body `hcl:",remain"`
}

// LogConfig 日志配置，指定文件夹或URL文件
type LogConfig struct {
	LogLevel string `hcl:"log_level,optional" json:"log_level,omitempty"`
	LogFile  string `hcl:"log_file,optional" json:"log_file,omitempty"`
	LogDir   string `hcl:"log_dir,optional" json:"log_dir,omitempty"`
}

func SetupLog() {
	var logger *zap.SugaredLogger
	cfg := GetLogSettings()
	if cfg.LogFile != "" {
		logger = logging.NewLoggerURL(cfg.LogLevel, cfg.LogFile)
	} else if cfg.LogDir == "" {
		logger = logging.NewLogger(cfg.LogDir)
	} else {
		logger = zap.NewNop().Sugar()
	}
	logging.SetLogger(logger)
}

// LoadConfigFile 读取默认配置文件
func LoadConfigFile(options any) (*RootConfig, error) {
	return ReadConfigFile(cfgFile, verbose, options)
}

// ReadConfigFile 读取配置文件
func ReadConfigFile(cfgFile string, verbose bool, options any) (*RootConfig, error) {
	var err error
	if theSettings == nil {
		theSettings = new(RootConfig)
		if verbose {
			fmt.Println("Config file is", cfgFile)
		}
		err = hclsimple.DecodeFile(cfgFile, nil, theSettings)
	}
	if err == nil && options != nil {
		_ = gohcl.DecodeBody(theSettings.App.Remain, nil, options)
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
