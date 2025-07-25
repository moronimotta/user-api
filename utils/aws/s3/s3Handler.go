package awsS3

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	bucketName = os.Getenv("BUCKET_NAME")
	region     = os.Getenv("AWS_REGION")
	s3Client   *s3.Client
)

func GeneratePresignedUploadURL(fileName, contentType string) (string, string, error) {
	cfg, _ := config.LoadDefaultConfig(context.TODO())
	client := s3.NewFromConfig(cfg)
	presignClient := s3.NewPresignClient(client)

	key := fmt.Sprintf("uploads/%s", fileName)

	params := &s3.PutObjectInput{
		Bucket:      &bucketName,
		Key:         &key,
		ContentType: &contentType,
	}

	presignedReq, err := presignClient.PresignPutObject(context.TODO(), params, func(o *s3.PresignOptions) {
		o.Expires = 15 * time.Minute
	})
	if err != nil {
		return "", "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}
	return presignedReq.URL, key, nil
}
