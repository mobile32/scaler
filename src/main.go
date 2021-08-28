package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/mobile32/scaler/src/utils"
)

func ScaleImages() {
	svc, _ := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1"),
	})

	bucket := utils.FilesManager{
		Session:    svc,
		BucketName: os.Getenv("BUCKET_NAME"),
	}

	originalFilesLocations := bucket.GetListOfFilesInBucketPath(os.Getenv("SOURCE_PATH"))
	resizedFilesLocations := bucket.GetListOfFilesInBucketPath(os.Getenv("TARGET_PATH"))

	newFilesLocations := make([]string, 0)

OUTER:
	for _, originalFileLocation := range originalFilesLocations {
		for _, resizedFilesLocation := range resizedFilesLocations {
			if originalFileLocation == resizedFilesLocation {
				continue OUTER
			}
		}

		newFilesLocations = append(newFilesLocations, originalFileLocation)
	}

	fmt.Println(newFilesLocations)

	for _, newFileName := range newFilesLocations {
		bucket.DownloadFileFromBucket(newFileName)
		utils.ScaleImage(newFileName)
		bucket.UploadFileToBucket(newFileName)
	}
}

func main() {
	lambda.Start(ScaleImages)
}
