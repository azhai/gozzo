package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	backDirs int    // 回退目录层级
	cfgFile  string // 配置文件位置
	verbose  bool   // 详细输出
)

func init() {
	flag.IntVar(&backDirs, "bcd", 0, "回退目录层级") // 默认在bin目录下
	flag.StringVar(&cfgFile, "hcl", "settings.hcl", "配置文件位置")
	// 和urfave/cli的version参数冲突，需要在App中设置HideVersion
	flag.BoolVar(&verbose, "vvv", false, "详细输出")
}

// SetupEnv 根据不同场景初始化
func SetupEnv(options any) {
	if IsRunTest() {
		_, _ = BackToDir(1) // 从tests退回根目录
	} else {
		_ = BackToAppDir() // 根据backDirs退回APP所在目录，一般不需要
	}
	if cfgFile != "" && !filepath.IsAbs(cfgFile) {
		cfgFile, _ = filepath.Abs(cfgFile) // 配置文件绝对路径
	}
	settings, err := ReadConfigFile(cfgFile, verbose, options)
	if err != nil {
		fmt.Printf("err=%#v\n%#v\n\n", err, settings)
	}
}

// IsRunTest 是否测试模式下
func IsRunTest() bool {
	return strings.HasSuffix(os.Args[0], ".test")
}

// BackToDir 退回上层目录
func BackToDir(back int) (dir string, err error) {
	if back == 0 {
		return
	} else if back < 0 {
		back = 0 - back
	}
	dir = strings.Repeat("../", back)
	if dir, err = filepath.Abs(dir); err == nil {
		err = os.Chdir(dir)
	}
	return
}

// BackToAppDir 如果在子目录下运行，需要先退回上层目录
func BackToAppDir() error {
	dir, err := BackToDir(backDirs)
	if err == nil && dir != "" && verbose {
		fmt.Println("Back to dir", dir)
	}
	return err
}

// Verbose 是否输出详细信息
func Verbose() bool {
	if !flag.Parsed() {
		panic("Verbose called before Parse")
	}
	return verbose
}
