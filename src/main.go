package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/mobile32/scaler/src/utils"
)

func main() {
	ORIGINAL_BUCKET_NAME := "owl-original-photos-bucket"
	RESIZED_BUCKET_NAME := "owl-resized-photos-bucket"
	TMP_PATH := "tmp_path"

	svc, _ := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1"),
	})

	originalImagesBucket := utils.FilesManager{
		Session:    svc,
		BucketName: ORIGINAL_BUCKET_NAME,
		TmpPath: TMP_PATH,
	}

	resizedImagesBucket := utils.FilesManager{
		Session:    svc,
		BucketName: RESIZED_BUCKET_NAME,
		TmpPath: TMP_PATH,
	}

	originalFilesNames := originalImagesBucket.GetListOfFilesInBucket()
	resizedFilesNames := resizedImagesBucket.GetListOfFilesInBucket()

	newFilesNames := make([]string, 0)
	OUTER: for _, originalFileName := range originalFilesNames {
		for _, resizedFilesName := range resizedFilesNames {
			if originalFileName == resizedFilesName {
				continue OUTER
			}
		}

		newFilesNames = append(newFilesNames, originalFileName)
	}

	fmt.Println(newFilesNames)

	for _, newFileName := range newFilesNames {
		originalImagesBucket.DownladFileFromBucket(newFileName)
		utils.ScaleImage(newFileName, TMP_PATH)
		//resizedImagesBucket.UploadFileToBucket(newFileName)
	}
}
