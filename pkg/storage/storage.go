// Package storage uploads generated assets to Google Cloud Storage.
package storage

import (
	"context"
	"fmt"

	gcs "cloud.google.com/go/storage"
)

// ImageUploader stores PNG images in a publicly readable bucket.
type ImageUploader struct {
	client *gcs.Client
	bucket string
}

func NewImageUploader(ctx context.Context, bucket string) (*ImageUploader, error) {
	client, err := gcs.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("creating storage client: %w", err)
	}
	return &ImageUploader{client: client, bucket: bucket}, nil
}

// UploadPNG stores data as name and returns its public URL.
func (u *ImageUploader) UploadPNG(ctx context.Context, name string, data []byte) (string, error) {
	w := u.client.Bucket(u.bucket).Object(name).NewWriter(ctx)
	w.ContentType = "image/png"
	w.CacheControl = "public, max-age=31536000"
	if _, err := w.Write(data); err != nil {
		w.Close()
		return "", fmt.Errorf("writing %q to bucket %q: %w", name, u.bucket, err)
	}
	if err := w.Close(); err != nil {
		return "", fmt.Errorf("closing upload of %q to bucket %q: %w", name, u.bucket, err)
	}
	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", u.bucket, name), nil
}
