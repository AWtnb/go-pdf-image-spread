package convert

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/karmdip-mi/go-fitz"
)

func getFiles(root string) []string {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".pdf" {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
	return files
}

// https://text.baldanders.info/golang/concatenate-images/

func loadImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func max(x int, y int) int {
	if x > y {
		return x
	}
	return y
}

func getFileBasename(name string) string {
	return strings.TrimSuffix(name, filepath.Ext(name))
}

func concImages(lPath string, rPath string, outDir string) error {

	width, height := 0, 0

	limg, err := loadImage(lPath)
	if err != nil {
		return err
	}
	rimg, err := loadImage(rPath)
	if err != nil {
		return err
	}

	width = limg.Bounds().Dx() + rimg.Bounds().Dx()
	height = max(limg.Bounds().Dy(), rimg.Bounds().Dy())

	outImg := image.NewNRGBA(image.Rect(0, 0, width, height))
	offset := 0
	for _, img := range []image.Image{limg, rimg} {
		rect := img.Bounds()
		draw.Draw(outImg, image.Rect(offset, 0, offset+rect.Dx(), rect.Dy()), img, image.Point{0, 0}, draw.Over)
		offset += rect.Dx()
	}

	outName := fmt.Sprintf("%s-%s", getFileBasename(filepath.Base(lPath)), filepath.Base(rPath))
	outPath := filepath.Join(outDir, outName)
	file, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := jpeg.Encode(file, outImg, &jpeg.Options{Quality: 100}); err != nil {
		return err
	}

	return nil
}

func copyFile(path string, destDir string) error {

	src, err := os.Open(path)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(filepath.Join(destDir, filepath.Base(path)))
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}
	return nil
}

func spreadPages(files []string, outDir string, singleTop bool, vertical bool) error {
	pairs := files
	if len(files)%2 != 0 {
		if singleTop {
			top := files[0]
			if err := copyFile(top, outDir); err != nil {
				return err
			}
			pairs = files[1:]
		} else {
			last := files[len(files)-1]
			if err := copyFile(last, outDir); err != nil {
				return err
			}
			pairs = files[:len(files)-1]
		}
	}
	for i := 0; i < len(pairs); i += 2 {
		l, r := pairs[i], pairs[i+1]
		if vertical {
			l, r = pairs[i+1], pairs[i]
		}
		err := concImages(l, r, outDir)
		if err != nil {
			return err
		}
	}
	return nil
}

func Convert(root string, singleTop bool, vertical bool) int {

	files := getFiles(root)
	td, err := os.MkdirTemp("", "temp-for-conv")
	defer os.RemoveAll(td)

	if err != nil {
		fmt.Println(err)
		return 1
	}
	for _, file := range files {
		doc, err := fitz.New(file)
		if err != nil {
			fmt.Println(err)
			return 1
		}
		pages := []string{}
		for i := 0; i < doc.NumPage(); i++ {
			img, err := doc.Image(i)
			if err != nil {
				fmt.Println(err)
				return 1
			}

			ipath := filepath.Join(td, fmt.Sprintf("p%05d.jpg", i+1))
			f, err := os.Create(ipath)
			if err != nil {
				fmt.Println(err)
				return 1
			}

			if err := jpeg.Encode(f, img, &jpeg.Options{Quality: 100}); err != nil {
				fmt.Println(err)
				return 1
			}

			f.Close()
			pages = append(pages, ipath)
		}
		if err := spreadPages(pages, root, singleTop, vertical); err != nil {
			fmt.Println(err)
			return 1
		}
	}
	return 0
}
