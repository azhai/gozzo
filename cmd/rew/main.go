package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/azhai/gozzo/filesystem"
	"github.com/azhai/gozzo/rewrite"
	. "github.com/klauspost/cpuid/v2"
)

var (
	MegaByte = 1024 * 1024
	verbose  bool
)

func init() {
	// 压舱石，阻止频繁GC
	ballast := make([]byte, 20*MegaByte)
	runtime.KeepAlive(ballast)

	if level := os.Getenv("GOAMD64"); level == "" {
		level = fmt.Sprintf("v%d", CPU.X64Level())
		fmt.Printf("请设置环境变量 export GOAMD64=%s\n\n", level)
	}

	flag.BoolVar(&verbose, "v", false, "是否输出详情")
	flag.Parse()
}

func main() {
	if flag.NArg() == 0 {
		prettifyDir(".")
		return
	}
	for _, dir := range flag.Args() {
		prettifyDir(dir)
	}
}

// prettifyDir 美化目录下的go代码文件
func prettifyDir(dir string) {
	files, err := filesystem.FindFiles(dir, ".go", "vendor/", ".git/")
	if err != nil {
		panic(err)
	}
	var chg bool
	for filename := range files {
		chg, err = rewrite.PrettifyGolangFile(filename, true)
		if verbose {
			if chg {
				fmt.Println("*", filename)
			} else {
				fmt.Println("|", filename)
			}
		}
		if err != nil {
			fmt.Println(err)
		}
	}
}
