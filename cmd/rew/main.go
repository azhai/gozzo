package main

import (
	"flag"
	"fmt"

	"github.com/azhai/gozzo/config"
	fs "github.com/azhai/gozzo/filesystem"
	"github.com/azhai/gozzo/rewrite"
)

const Version = "1.4.0"

var verbose bool

func init() {
	config.PrepareEnv(20)

	flag.BoolVar(&verbose, "v", false, "display more information")
	flag.Usage = usage
	flag.Parse()
}

func main() {
	fmt.Println("NOTE: The below files with an * at the beginning have been modified.")
	if flag.NArg() == 0 {
		prettifyDir(".")
		return
	}
	for _, dir := range flag.Args() {
		prettifyDir(dir)
	}
}

// usage 使用帮助
func usage() {
	out := flag.CommandLine.Output()
	desc := `Version: v%s
Usage: rew [flags] [dir ...]
`
	_, _ = fmt.Fprintf(out, desc, Version)
	flag.PrintDefaults()
}

// prettifyDir 美化目录下的go代码文件
func prettifyDir(dir string) {
	files, err := fs.FindFiles(dir, ".go", "vendor/", ".git/")
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
