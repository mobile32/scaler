package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/mobile32/scaler/src/config"
	"os"
	"path/filepath"

	"github.com/mobile32/scaler/src/utils"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

func sendSNSNotification(session *session.Session) {
	if snsArn := config.Envs.SnsArn; snsArn != "" {
		snsSuccessMessage := config.Envs.SnsSuccessMessage

		svc := sns.New(session)

		result, err := svc.Publish(&sns.PublishInput{
			TopicArn: &snsArn,
			Message:  &snsSuccessMessage,
		})
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		fmt.Println("SNS message was sent (id): " + *result.MessageId)
	}
}

func scaleImages(ctx context.Context, s3Event events.S3Event) {
	for _, record := range s3Event.Records {
		s3 := record.S3
		fmt.Printf("[%s - %s] Bucket = %s, Key = %s \n", record.EventSource, record.EventTime, s3.Bucket.Name, s3.Object.Key)
	}

	session, _ := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1"),
	})

	bucket := utils.FilesManager{
		Session:    session,
		BucketName: config.Envs.BucketName ,
	}

	originalFilesLocations := bucket.GetListOfFilesInBucketPath(config.Envs.SourcePath, true)
	resizedFilesLocations := bucket.GetListOfFilesInBucketPath(config.Envs.TargetPath, true)

	newFilesLocations := utils.CreateBucketFilesDiff(originalFilesLocations, resizedFilesLocations)
	newFilesLocations = utils.RemoveInvalidFilesTypes(newFilesLocations)

	fmt.Println("New files locations", newFilesLocations)

	for _, newFileLocation := range newFilesLocations {
		newFileLocationWithPrefix := filepath.Join(config.Envs.SourcePath, newFileLocation)

		bucket.DownloadFileFromBucket(newFileLocationWithPrefix)
		utils.ScaleImage(newFileLocationWithPrefix)
		bucket.UploadFileToBucket(newFileLocationWithPrefix)
	}

	sendSNSNotification(session)
}

func main() {
	config.Init()
	lambda.Start(scaleImages)
}
