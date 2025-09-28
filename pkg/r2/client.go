package r2client

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func NewR2Client(ctx context.Context) (*s3.Client, error) {
	accountId := os.Getenv("R2_ACCOUNT_ID")
	accessKeyId := os.Getenv("R2_ACCESS_KEY_ID")
	accessKeySecret := os.Getenv("R2_ACCESS_KEY_SECRET")
	r2Region := os.Getenv("R2_REGION") 

	if accountId == "" || accessKeyId == "" || accessKeySecret == "" {
		log.Println("R2 configuration missing in environment variables.")
		return nil, fmt.Errorf("R2_ACCOUNT_ID, R2_ACCESS_KEY_ID, or R2_ACCESS_KEY_SECRET is missing")
	}

	if r2Region == "" {
        r2Region = "auto" 
    }

	// 1. Muat konfigurasi dasar dengan Static Credentials
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyId, accessKeySecret, "")),
		config.WithRegion(r2Region),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// 2. Buat S3 Client dan Terapkan Endpoint Override
    // Pola endpoint R2: https://<accountid>.r2.cloudflarestorage.com
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountId))
	})

	return client, nil
}
