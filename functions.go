//package utils
//
//import {
//	"fmt"
//	"os"
//
//	_ "image/gif"
//	_ "image/jpeg"
//	_ "image/png"
//
//	"github.com/aws/aws-sdk-go/aws"
//	"github.com/aws/aws-sdk-go/aws/session"
//	"github.com/aws/aws-sdk-go/service/s3"
//	"github.com/aws/aws-sdk-go/service/s3/s3manager"
//}
//
//func ScaleImage() {
//	reader, err := os.Open("./test_images/owl_2.jpeg")
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer reader.Close()
//
//	srcImage, _, err := image.Decode(reader)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	newImage := imaging.Resize(srcImage, 128, 128, imaging.Lanczos)
//
//	// Save the resulting image as JPEG.
//	err = imaging.Save(newImage, "owl_2.jpeg")
//	if err != nil {
//		log.Fatalf("failed to save image: %v", err)
//	}
//}
//
//type MyEvent struct {
//	Name string `json:"name"`
//}
//
//func HandleRequest(ctx context.Context, name MyEvent) (string, error) {
//	return fmt.Sprintf("Hello %s!", name.Name), nil
//}
//
//func UseLambda() {
//	lambda.Start(HandleRequest)
//}