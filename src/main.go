package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mobile32/scaler/src/utils"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

func scaleImages() {
	svc, _ := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1"),
	})

	bucket := utils.FilesManager{
		Session:    svc,
		BucketName: os.Getenv("BUCKET_NAME"),
	}

	originalFilesLocations := bucket.GetListOfFilesInBucketPath(os.Getenv("SOURCE_PATH"), true)
	resizedFilesLocations := bucket.GetListOfFilesInBucketPath(os.Getenv("TARGET_PATH"), true)

	newFilesLocations := utils.CreateBucketFilesDiff(originalFilesLocations, resizedFilesLocations)
	newFilesLocations = utils.RemoveInvalidFilesTypes(newFilesLocations)

	fmt.Println("New files locations", newFilesLocations)

	for _, newFileLocation := range newFilesLocations {
		newFileLocationWithPrefix := filepath.Join(os.Getenv("SOURCE_PATH"), newFileLocation)

		bucket.DownloadFileFromBucket(newFileLocationWithPrefix)
		utils.ScaleImage(newFileLocationWithPrefix)
		bucket.UploadFileToBucket(newFileLocationWithPrefix)
	}
}

func main() {
	lambda.Start(scaleImages)
}
