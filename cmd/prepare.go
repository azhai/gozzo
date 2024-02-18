package cmd

import (
	"fmt"
	"os"
	"runtime"

	. "github.com/klauspost/cpuid/v2"
)

var (
	Version  = "1.3.9"
	MegaByte = 1024 * 1024
)

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
