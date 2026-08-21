package gcs

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"time"

	"cloud.google.com/go/storage"
	"github.com/riazahmedshah/vfs/internal/config"
	"google.golang.org/api/option"
)

type Disposition string

const (
	DispositionInline     Disposition = "inline"
	DispositionAttachment Disposition = "attachment"
)

type GCSClient struct {
	cfg    *config.Config
	client *storage.Client
}

func NewGCSClient(cfg *config.Config) (*GCSClient, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithAuthCredentialsFile(option.ServiceAccount, cfg.GCS.GoogleCredentialPath))
	if err != nil {
		return nil, fmt.Errorf("failed to create storage client: %w", err)
	}
	return &GCSClient{
		cfg:    cfg,
		client: client,
	}, nil
}

func (g *GCSClient) UploadFile(ctx context.Context, objectName string, reader io.Reader) error {
	wc := g.client.Bucket(g.cfg.GCS.GCSBucketName).Object(objectName).NewWriter(ctx)

	if _, err := io.Copy(wc, reader); err != nil {
		return fmt.Errorf("failed to copy data: %w", err)
	}

	if err := wc.Close(); err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}

	return nil
}

func (g *GCSClient) GenerateSignedURL(objectName string, disposition Disposition, expiry time.Duration) (string, error) {
	opts := &storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV4,
		Method:  "GET",
		Expires: time.Now().Add(expiry),
		QueryParameters: url.Values{
			"response-content-disposition": []string{string(disposition)},
		},
	}

	url, err := g.client.Bucket(g.cfg.GCS.GCSBucketName).SignedURL(objectName, opts)
	if err != nil {
		return "", fmt.Errorf("failed to generate signed URL: %w", err)
	}
	return url, nil
}

func (g *GCSClient) DeleteFile(ctx context.Context, objectName string) error {
	if err := g.client.Bucket(g.cfg.GCS.GCSBucketName).Object(objectName).Delete(ctx); err != nil {
		return fmt.Errorf("failed to delete object: %w", err)
	}
	return nil
}
