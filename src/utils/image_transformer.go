package utils

import (
	"fmt"
	"github.com/mobile32/scaler/src/config"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
)

func ScaleImage(fileLocation string) {
	reader, err := os.Open(filepath.Join("/tmp", fileLocation))
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	srcImage, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}

	newImage := imaging.Resize(srcImage, config.Envs.ImagesWidth, config.Envs.ImagesHeight, imaging.Lanczos)

	err = imaging.Save(newImage, filepath.Join("/tmp", fileLocation))
	if err != nil {
		log.Fatalf("failed to save image: %v", err)
	}
	fmt.Println("Resized", filepath.Join("/tmp", fileLocation))
}
