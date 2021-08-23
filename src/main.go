package main

import (
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/mobile32/scaler/src/utils"
)

func main() {
	BUCKET_NAME := "owl-original-photos-bucket"

	session, _ := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1"),
	})

	filesManager := utils.FilesManager{
		Session:    session,
		BucketName: BUCKET_NAME,
	}

	filesNames := filesManager.GetListOfFilesInBucket()
	for _, fileName := range filesNames {
		filesManager.DownladFileFromBucket(fileName)
	}
}
