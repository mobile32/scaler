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