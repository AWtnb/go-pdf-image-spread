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

	"github.com/cheggaaa/pb/v3"
	"github.com/gen2brain/go-fitz"
)

func getFiles(root string, ext string) ([]string, error) {
	var files []string
	wd, err := os.Open(root)
	if err != nil {
		return files, err
	}
	defer wd.Close()
	names, err := wd.Readdirnames(0)
	if err != nil {
		return files, err
	}
	for _, file := range names {
		if filepath.Ext(file) == ext {
			files = append(files, filepath.Join(root, file))
		}
	}
	return files, err
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

func trimExt(name string) string {
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

	outName := fmt.Sprintf("%s-%s", trimExt(filepath.Base(lPath)), filepath.Base(rPath))
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

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}
	return nil
}

func allocate(files []string, outDir string, singleTop bool, vertical bool) error {
	pairs := files
	if singleTop {
		top := pairs[0]
		if err := copyFile(top, outDir); err != nil {
			return err
		}
		pairs = pairs[1:]
	}
	if len(pairs)%2 != 0 {
		last := pairs[len(pairs)-1]
		if err := copyFile(last, outDir); err != nil {
			return err
		}
		pairs = pairs[:len(pairs)-1]
	}

	fmt.Println("\n- concatenating images...")
	bar := pb.StartNew(len(pairs) / 2)

	for i := 0; i < len(pairs); i += 2 {
		l, r := pairs[i], pairs[i+1]
		if vertical {
			l, r = pairs[i+1], pairs[i]
		}
		err := concImages(l, r, outDir)
		if err != nil {
			return err
		}
		bar.Increment()
	}
	bar.Finish()
	return nil
}

func toImage(file string, outDir string) ([]string, error) {
	pages := []string{}
	doc, err := fitz.New(file)
	if err != nil {
		return pages, err
	}

	fmt.Println("\n- converting to image...")
	bar := pb.StartNew(doc.NumPage())

	for i := 0; i < doc.NumPage(); i++ {
		img, err := doc.Image(i)
		if err != nil {
			return pages, err
		}

		ipath := filepath.Join(outDir, fmt.Sprintf("p%05d.jpg", i+1))
		f, err := os.Create(ipath)
		if err != nil {
			return pages, err
		}
		if err := jpeg.Encode(f, img, &jpeg.Options{Quality: 100}); err != nil {
			return pages, err
		}
		defer f.Close()

		pages = append(pages, ipath)
		bar.Increment()
	}
	bar.Finish()
	return pages, nil
}

func Convert(root string, singleTop bool, vertical bool) error {
	files, err := getFiles(root, ".pdf")
	if err != nil {
		return err
	}
	for _, file := range files {

		fmt.Printf("\nProcessing: %s\n", filepath.Base(file))

		outDir := trimExt(file)
		if singleTop {
			outDir = outDir + "-singletop"
		}
		if vertical {
			outDir = outDir + "-vertical"
		}
		if err := os.Mkdir(outDir, os.ModePerm); err != nil {
			return err
		}
		pages, err := toImage(file, outDir)
		if err != nil {
			return err
		}
		concDir := filepath.Join(outDir, "conc")
		if err := os.Mkdir(concDir, os.ModePerm); err != nil {
			return err
		}
		if err := allocate(pages, concDir, singleTop, vertical); err != nil {
			return err
		}
	}
	return nil
}
