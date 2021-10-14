package config

import (
	"github.com/caarlos0/env/v6"
	"log"
)

type EnvsConfig struct {
	BucketName        string `env:"BUCKET_NAME,required"`
	SourcePath        string `env:"SOURCE_PATH,required"`
	TargetPath        string `env:"TARGET_PATH,required"`
	ImagesWidth       int `env:"IMAGES_WIDTH,required"`
	ImagesHeight      int `env:"IMAGES_HEIGHT,required"`
	SnsArn            string `env:"SNS_ARN"`
	SnsSuccessMessage string `env:"SNS_SUCCESS_MESSAGE"`
}

var Envs EnvsConfig

func Init()  {
	if err := env.Parse(Envs); err != nil {
		log.Fatal(err)
	}
}