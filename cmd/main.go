package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/Kartik-2239/ls-gallery/internal"
)

func main() {
	ImgDir := flag.String("path", "images", "Path to the directory containing images")
	flag.Parse()
	filePath := *ImgDir
	provided := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == "path" {
			provided = true
		}
	})
	if !provided {
		filePath, _ = os.Getwd()
		fmt.Println(filePath)
	}
	files, err := os.ReadDir(filePath)
	if err != nil {
		log.Fatal(err)
	}
	paths := []string{}
	for _, file := range files {
		if !file.IsDir() {
			for _, ext := range []string{".jpg", ".jpeg", ".png", ".gif"} {
				if strings.HasSuffix(strings.ToLower(file.Name()), ext) {
					paths = append(paths, path.Join(filePath, file.Name()))
					break
				}
			}
		}
	}
	fmt.Println(len(paths))
	if len(paths) == 0 {
		log.Fatal("No images found in the specified directory")
	}

	internal.Initialize(paths)
}
