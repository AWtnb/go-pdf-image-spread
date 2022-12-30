package main

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/AWtnb/go-pdf-image-spread/convert"
)

func main() {
	var (
		singleTop bool
		vertical  bool
	)
	flag.BoolVar(&singleTop, "singleTop", false, "switch to start with non-spread page")
	flag.BoolVar(&vertical, "vertical", false, "switch to start allocate pages from right to left")
	flag.Parse()
	os.Exit(run(singleTop, vertical))
}

func run(singleTop bool, vertical bool) int {

	p, _ := os.Executable()
	root := filepath.Dir(p)

	return convert.Convert(root, singleTop, vertical)
}
