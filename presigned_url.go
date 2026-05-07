package main

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func generatePresignedURL(s3Client *s3.Client, bucket, key string, expireTime time.Duration) (string, error) {
	// Create a presign client
	presignClient := s3.NewPresignClient(s3Client)

	// Create presigned URL with expiration time
	presignRequest, err := presignClient.PresignGetObject(
		context.Background(),
		&s3.GetObjectInput{
			Bucket: &bucket,
			Key:    &key,
		},
		s3.WithPresignExpires(expireTime),
	)
	if err != nil {
		return "", fmt.Errorf("couldn't create presigned URL: %w", err)
	}

	// Return the presigned URL
	return presignRequest.URL, nil
}
