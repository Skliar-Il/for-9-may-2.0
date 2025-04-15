package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"log"
	"os"
)

type Config struct {
	AccessKeyID     string `env:"AWS_ACCESS_KEY_ID"`
	SecretAccessKey string `env:"AWS_SECRET_ACCESS_KEY"`
	StorageName     string `env:"AWS_STORAGE_NAME"`
}

func NewS3Client(configAWS *Config) *s3.Client {
	if err := os.Setenv("AWS_ACCESS_KEY_ID", configAWS.AccessKeyID); err != nil {
		log.Fatalf("failed set env AWS_ACCESS_KEY_ID: %v", err)
	}
	if err := os.Setenv("AWS_SECRET_ACCESS_KEY", configAWS.SecretAccessKey); err != nil {
		log.Fatalf("failed set env AWS_SECRET_ACCESS_KEY: %v", err)
	}
	if err := os.Setenv("AWS_REGION", "ru-central1"); err != nil {
		log.Fatalf("faoled set env AWS_REGION: %v", err)
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("ru-central1"),
	)
	if err != nil {
		log.Fatalf("failed load config: %v", err)
	}

	return s3.NewFromConfig(cfg)
}
