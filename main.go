package main

import (
	"flag"
	"fmt"
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
	flag.Usage = func() {
		fmt.Println("REPOSITORY: https://github.com/AWtnb/go-pdf-image-spread")
		flag.PrintDefaults()
	}
	flag.Parse()
	os.Exit(run(singleTop, vertical))
}

func run(singleTop bool, vertical bool) int {
	p, _ := os.Executable()
	root := filepath.Dir(p)

	err := convert.Convert(root, singleTop, vertical)
	if err != nil {
		fmt.Println(err)
		return 1
	}
	fmt.Println("\nAll pages have been converted! Check `conc` folder inside output folder.")
	return 0
}
