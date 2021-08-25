package utils

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"

	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type FilesManager struct {
	Session    *session.Session
	BucketName string
	TmpPath    string
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

func (filesManager FilesManager) DownladFileFromBucket(fileName string) {
	file, err := filesManager.createFileInPath(fileName)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	svc := s3manager.NewDownloader(filesManager.Session)

	numBytes, err := svc.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(filesManager.BucketName),
			Key:    aws.String(fileName),
		})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Downloaded", file.Name(), numBytes, "bytes")
}

func (filesManager FilesManager) UploadFileToBucket(fileName string) {
	file, err := os.Open(filepath.Join(filesManager.TmpPath, fileName))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	size := fileInfo.Size()
	buffer := make([]byte, size)
	file.Read(buffer)

	_, err = s3.New(filesManager.Session).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(filesManager.BucketName),
		Key:                  aws.String(fileName),
		ACL:                  aws.String("private"),
		Body:                 bytes.NewReader(buffer),
		ContentLength:        aws.Int64(size),
		ContentType:          aws.String(http.DetectContentType(buffer)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Uploaded", file.Name())
}

func (filesManager FilesManager) createFileInPath(fileName string) (*os.File, error) {
	if filesManager.TmpPath != "" {
		tmpPath := filepath.Join(".", filesManager.TmpPath)

		if _, existingErr := os.Stat(tmpPath); os.IsNotExist(existingErr) {
			err := os.Mkdir(tmpPath, 0755)
			if err != nil {
				panic(err)
			}
		}
	}

	return os.Create(filepath.Join(filesManager.TmpPath, fileName))
}