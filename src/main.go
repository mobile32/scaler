package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/mobile32/scaler/src/utils"
)

type Params struct {
	OriginalBucketName string `json:"originalBucketName"`
	ResizedBucketName  string `json:"resizedBucketName"`
}

func ScaleImages(ctx context.Context, params Params) (string, error) {
	svc, _ := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1"),
	})

	originalImagesBucket := utils.FilesManager{
		Session:    svc,
		BucketName: params.OriginalBucketName,
	}

	resizedImagesBucket := utils.FilesManager{
		Session:    svc,
		BucketName: params.ResizedBucketName,
	}

	originalFilesNames := originalImagesBucket.GetListOfFilesInBucket()
	resizedFilesNames := resizedImagesBucket.GetListOfFilesInBucket()

	newFilesNames := make([]string, 0)

OUTER:
	for _, originalFileName := range originalFilesNames {
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
		utils.ScaleImage(newFileName)
		resizedImagesBucket.UploadFileToBucket(newFileName)
	}

	return fmt.Sprintf("Process finished"), nil
}

func main() {
	lambda.Start(ScaleImages)
}
