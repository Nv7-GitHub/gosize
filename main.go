package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

var module = *flag.String("module", "", "Go module to compile. Defaults to current directory.")
var port = *flag.Int("port", 8080, "Port for webserver to listen on.")
var outdir string

func init() {
	cachedir, err := os.UserCacheDir()
	handle(err)
	flag.StringVar(&outdir, "output", filepath.Join(cachedir, "gosize"), "Output JS directory. Defaults to gosize folder in OS cache directory.")
}

func handle(err error) {
	if err != nil {
		fmt.Println("error:", err)
		flag.Usage()
		os.Exit(1)
		return
	}
}

func main() {
	flag.Parse()
	err := os.MkdirAll(outdir, os.ModePerm)
	handle(err)

	sizes := build(module)
	mapped := tab2dic(sizes)

	fmt.Println(mapped)
}
