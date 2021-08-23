package utils

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type FilesManager struct {
	Session    *session.Session
	BucketName string
}

func (filesManager FilesManager) GetListOfFilesInBucket() []string {
	svc := s3.New(filesManager.Session)

	params := &s3.ListObjectsInput{
		Bucket: aws.String(filesManager.BucketName),
	}

	resp, _ := svc.ListObjects(params)

	fileNames := make([]string, len(resp.Contents))
	for i, key := range resp.Contents {
		fileNames[i] = *key.Key
	}

	return fileNames
}

func (filesManager FilesManager) DownladFileFromBucket(item string) {
	file, err := os.Create(item)
	if err != nil {
		fmt.Println(err)
	}

	defer file.Close()

	svc := s3manager.NewDownloader(filesManager.Session)

	numBytes, err := svc.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(filesManager.BucketName),
			Key:    aws.String(item),
		})
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Downloaded", file.Name(), numBytes, "bytes")
}
