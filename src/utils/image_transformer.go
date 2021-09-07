package utils

import (
	"fmt"
	"image"
	"log"
	"os"
	"path/filepath"
	"strconv"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

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

	width, err := strconv.ParseInt(os.Getenv("IMAGES_WIDTH"), 10, 32)
	if err != nil {
		log.Fatal(err)
	}

	height, err := strconv.ParseInt(os.Getenv("IMAGES_HEIGHT"), 10, 16)
	if err != nil {
		log.Fatal(err)
	}

	newImage := imaging.Resize(srcImage, int(width), int(height), imaging.Lanczos)

	err = imaging.Save(newImage, filepath.Join("/tmp", fileLocation))
	if err != nil {
		log.Fatalf("failed to save image: %v", err)
	}
	fmt.Println("Resized", filepath.Join("/tmp", fileLocation))
}
