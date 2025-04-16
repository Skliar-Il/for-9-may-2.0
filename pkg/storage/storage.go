package storage

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

type InterfaceStorage interface {
	LoadJPG(image []byte) (string, error)
}

type Storage struct {
	Config Config
	Client *s3.Client
}

func NewStorage(cfg Config) (*Storage, error) {
	client, err := NewS3Client(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create S3 client: %w", err)
	}

	return &Storage{
		Config: cfg,
		Client: client,
	}, nil
}

func (s Storage) LoadJPG(image []byte) (string, error) {
	imageID := uuid.New()
	key := fmt.Sprintf("photos/%s.jpg", imageID.String())

	_, err := s.Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(s.Config.StorageName),
		Key:         aws.String(key),
		Body:        bytes.NewReader(image),
		ContentType: aws.String("image/jpeg"),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload image to storage: %w", err)
	}

	imageURL := fmt.Sprintf("https://%s.storage.yandexcloud.net/%s", s.Config.StorageName, key)
	return imageURL, nil
}
