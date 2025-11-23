package storage

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	appCfg "oma-library/internal/config"
)

type R2Client struct {
	client *s3.Client
	bucket string
	ctx   context.Context
}

func NewR2Client(ctx context.Context, r2cfg appCfg.R2Config) (*R2Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(r2cfg.AccessKey, r2cfg.SecretKey, ""),
		),
		config.WithRegion("auto"),
		config.WithBaseEndpoint(r2cfg.Endpoint),
	)
	if err != nil {
		return nil, err
	}

	return &R2Client{
		client: s3.NewFromConfig(cfg),
		bucket: r2cfg.Bucket,
		ctx:   ctx,
	}, nil
}

func (r2 *R2Client) UploadFileToR2(ctx context.Context, key string, file io.Reader)  error {
	_, err := r2.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &r2.bucket,
		Key:    &key,
		Body:   file,
	})
	if err != nil {
		return err
	}
	return  nil
}

func (r2 *R2Client) DownloadFileFromR2(ctx context.Context, key string) (*s3.GetObjectOutput, error) {
	
	result, err := r2.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &r2.bucket,
		Key:    &key,
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}
