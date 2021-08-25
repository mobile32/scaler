package utils

import (
	"fmt"
	"github.com/disintegration/imaging"
	"image"
	"log"
	"os"
	"path/filepath"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

func ScaleImage(fileName string, path string) {
	reader, err := os.Open(filepath.Join(path, fileName))
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	srcImage, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}

	newImage := imaging.Resize(srcImage, 128, 128, imaging.Lanczos)

	err = imaging.Save(newImage, filepath.Join(path, fileName))
	if err != nil {
		log.Fatalf("failed to save image: %v", err)
	}
	fmt.Println("Resized", filepath.Join(path, fileName))
}