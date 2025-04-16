package storage

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Config struct {
	AccessKeyID     string `env:"AWS_ACCESS_KEY_ID"`
	SecretAccessKey string `env:"AWS_SECRET_ACCESS_KEY"`
	StorageName     string `env:"AWS_STORAGE_NAME"`
}

func NewS3Client(configAWS Config) (*s3.Client, error) {

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if service == s3.ServiceID {
			return aws.Endpoint{
				PartitionID: "yc",
				URL:         "https://storage.yandexcloud.net",
			}, nil
		}
		return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested")
	})

	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithCredentialsProvider(
			&credentials.StaticCredentialsProvider{
				Value: aws.Credentials{
					AccessKeyID:     configAWS.AccessKeyID,
					SecretAccessKey: configAWS.SecretAccessKey,
				},
			},
		),
		config.WithRegion("auto"),
		config.WithClientLogMode(aws.LogRequest|aws.LogResponse|aws.LogRetries), // Включаем логирование
	)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})
	return client, nil
}
