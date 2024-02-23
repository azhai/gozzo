package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	. "github.com/klauspost/cpuid/v2"
)

const MegaByte = 1024 * 1024

var (
	backDirs int    // 回退目录层级
	cfgFile  string // 配置文件位置
	verbose  bool   // 详细输出
)

func init() {
	flag.IntVar(&backDirs, "bd", 0, "回退目录层级") // 默认在bin目录下
	flag.StringVar(&cfgFile, "cf", "settings.hcl", "配置文件位置")
	// 和urfave/cli的version参数冲突，需要在App中设置HideVersion
	flag.BoolVar(&verbose, "vv", false, "详细输出")
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
