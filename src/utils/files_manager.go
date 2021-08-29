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
}

func (filesManager FilesManager) GetListOfFilesInBucketPath(bucketPath string, removePrefix bool) []string {
	svc := s3.New(filesManager.Session)

	params := &s3.ListObjectsInput{
		Bucket: aws.String(filesManager.BucketName),
		Prefix: aws.String(filepath.Clean(bucketPath)),
	}

	resp, _ := svc.ListObjects(params)
	fmt.Println("Listing in bucket", bucketPath)

	fileNames := make([]string, len(resp.Contents))
	for i, key := range resp.Contents {
		if removePrefix {
			fileNames[i] = removePathPrefix(*key.Key, bucketPath)
		} else {
			fileNames[i] = *key.Key
		}
	}
	fmt.Println("Files listed in bucket", fileNames)

	return fileNames
}

func (filesManager FilesManager) DownloadFileFromBucket(fileLocation string) {
	createNewDirectory(filepath.Join("/tmp", filepath.Dir(fileLocation)))

	file, err := os.Create(filepath.Join("/tmp", fileLocation))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fmt.Println("Empty file created", file.Name())

	svc := s3manager.NewDownloader(filesManager.Session)

	numBytes, err := svc.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(filesManager.BucketName),
			Key:    aws.String(fileLocation),
		})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Downloaded", file.Name(), numBytes, "bytes")
}

func (filesManager FilesManager) UploadFileToBucket(fileLocation string) {
	file, err := os.Open(filepath.Join("/tmp", fileLocation))
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
		Key:                  aws.String(createTargetPath(fileLocation)),
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

func RemoveInvalidFilesTypes(filesLocations []string) []string {
	validFilesLocations := make([]string, 0)

	for _, fileLocation := range filesLocations {
		ext := filepath.Ext(fileLocation)
		if isValidType(ext) {
			validFilesLocations = append(validFilesLocations, fileLocation)
		}
	}

	return validFilesLocations
}

func CreateBucketFilesDiff(sourceFilesLocations []string, targetFilesLocations []string) []string {
	diffFilesLocations := make([]string, 0)

OUTER:
	for _, targetFileLocation := range sourceFilesLocations {
		for _, resizedFilesLocation := range targetFilesLocations {
			if targetFileLocation == resizedFilesLocation {
				continue OUTER
			}
		}

		diffFilesLocations = append(diffFilesLocations, targetFileLocation)
	}

	fmt.Println("Diff locations", diffFilesLocations)
	return diffFilesLocations
}

func createNewDirectory(newDirectoryPath string) {
	if _, err := os.Stat(newDirectoryPath); os.IsNotExist(err) {
		err := os.MkdirAll(newDirectoryPath, 0777)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Directory created", newDirectoryPath)
	}
}

func removePathPrefix(fileLocation string, pathPrefix string) string {
	sourcePathPrefixLen := len(filepath.Clean(pathPrefix))
	return fileLocation[sourcePathPrefixLen:]
}

func createTargetPath(fileLocation string) string {
	originalPath := removePathPrefix(fileLocation, os.Getenv("SOURCE_PATH"))
	return filepath.Join(os.Getenv("TARGET_PATH"), originalPath)
}

func isValidType(originalType string) bool {
	validTypes := []string{".jpg", ".jpeg", ".png", ".gif"}

	for _, validType := range validTypes {
		if validType == originalType {
			return true
		}
	}

	return false
}
