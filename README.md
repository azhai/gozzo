# gozzo 尜舟

The utils in golang.

## 浮点数 decimal

```go
package main
import (
    "fmt"
	
    "github.com/azhai/gozzo/decimal"
)

// 浮点数
func main() {
    x := 123.45678
    a := decimal.NewDecimal(decimal.RoundN(x, 2), 2)
    fmt.Println(a.String()) // 123.45
    b := decimal.ParseDecimal(a.String(), 2)
    fmt.Println(b.String()) // 123.45
}
```

## 文件操作 filesystem

```go
package main

import (
	"fmt"

	"github.com/azhai/gozzo/filesystem"
)

// 文件计行
func main() {
	filename := "README.md"
	count := filesystem.LineCount(filename)

	// 逐行返回，适用于大文件
	var lines []string
	r := filesystem.NewLineReader(filename)
	for r.Reading() {
		lines = append(lines, r.Text())
	}
	if len(lines) == count {
		fmt.Printf("%s have %d lines\n", filename, count)
	} else {
		fmt.Println("Error !")
	}
}
```

## 文件日志 logging

```go
package main

import (
	"math"
	"time"

	"github.com/azhai/gozzo/logging"
)

// CalcAge 计算年龄
func CalcAge(birthday string) int {
	birth, err := time.Parse("2006-01-02", birthday)
	if err != nil {
		return -1
	}
	hours := time.Since(birth).Hours()
	return int(math.Round(hours / 24 / 365.25))
}

func main() {
	birthday := "1996-02-29"
	age := CalcAge(birthday)
	logger := logging.NewLoggerURL("debug", "stdout") // 输出到屏幕
	logger.Infof("I was born on %s, I am %d years old.", birthday, age)
}
```

## go代码美化 rewrite

```
➜ make && ./bin/rew -h

#/usr/local/go/bin/go clean
rm -f ./bin/*
Clean all.
Compile rew ...
GOOS=darwin GOARCH=amd64 GOAMD64=v3 CGO_ENABLED=1 /usr/local/go/bin/go build -ldflags="-s -w" -o ./bin/rew ./cmd/rew
Build success.

Version: v1.3.7
Usage: rew [flags] [dir ...]
  -v	display more information
```
