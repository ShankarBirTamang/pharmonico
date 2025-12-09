// Package services provides MinIO object storage operations
package services

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinIOService handles MinIO object storage operations
type MinIOService struct {
	client   *minio.Client
	endpoint string
}

// UploadOptions contains options for file uploads
type UploadOptions struct {
	ContentType string
	Metadata    map[string]string
}

// ObjectInfo represents information about an uploaded object
type ObjectInfo struct {
	Bucket       string
	ObjectName   string
	Size         int64
	ContentType  string
	ETag         string
	LastModified time.Time
}

// NewMinIOService creates a new MinIO service client
func NewMinIOService(endpoint, accessKey, secretKey string, useSSL bool) (*MinIOService, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}

	return &MinIOService{
		client:   client,
		endpoint: endpoint,
	}, nil
}

// Upload uploads a file to MinIO
// bucket: the bucket name (e.g., "insurance-cards", "shipping-labels", "ncpdp-raw")
// objectName: the object name/path in the bucket
// reader: the data to upload
// size: size of the data (-1 for unknown size)
// opts: optional upload options
func (s *MinIOService) Upload(ctx context.Context, bucket, objectName string, reader io.Reader, size int64, opts *UploadOptions) (*ObjectInfo, error) {
	// Ensure bucket exists
	exists, err := s.client.BucketExists(ctx, bucket)
	if err != nil {
		return nil, fmt.Errorf("failed to check bucket existence: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("bucket %s does not exist", bucket)
	}

	// Prepare upload options
	putOptions := minio.PutObjectOptions{}
	if opts != nil {
		if opts.ContentType != "" {
			putOptions.ContentType = opts.ContentType
		}
		if opts.Metadata != nil {
			putOptions.UserMetadata = opts.Metadata
		}
	}

	// Upload the object
	_, err = s.client.PutObject(ctx, bucket, objectName, reader, size, putOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to upload object: %w", err)
	}

	// Get object info to retrieve full metadata
	objInfo, err := s.client.StatObject(ctx, bucket, objectName, minio.StatObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get object info after upload: %w", err)
	}

	return &ObjectInfo{
		Bucket:       bucket,
		ObjectName:   objInfo.Key,
		Size:         objInfo.Size,
		ContentType:  objInfo.ContentType,
		ETag:         objInfo.ETag,
		LastModified: objInfo.LastModified,
	}, nil
}

// GenerateSignedURL generates a presigned URL for accessing an object
// bucket: the bucket name
// objectName: the object name/path
// expiry: duration until the URL expires (default: 1 hour)
func (s *MinIOService) GenerateSignedURL(ctx context.Context, bucket, objectName string, expiry time.Duration) (string, error) {
	if expiry == 0 {
		expiry = 1 * time.Hour // Default 1 hour
	}

	url, err := s.client.PresignedGetObject(ctx, bucket, objectName, expiry, nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate signed URL: %w", err)
	}

	return url.String(), nil
}

// GenerateSignedPutURL generates a presigned URL for uploading an object
// bucket: the bucket name
// objectName: the object name/path
// expiry: duration until the URL expires (default: 1 hour)
func (s *MinIOService) GenerateSignedPutURL(ctx context.Context, bucket, objectName string, expiry time.Duration) (string, error) {
	if expiry == 0 {
		expiry = 1 * time.Hour // Default 1 hour
	}

	url, err := s.client.PresignedPutObject(ctx, bucket, objectName, expiry)
	if err != nil {
		return "", fmt.Errorf("failed to generate signed PUT URL: %w", err)
	}

	return url.String(), nil
}

// Delete deletes an object from MinIO
// bucket: the bucket name
// objectName: the object name/path
func (s *MinIOService) Delete(ctx context.Context, bucket, objectName string) error {
	err := s.client.RemoveObject(ctx, bucket, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete object: %w", err)
	}
	return nil
}

// DeleteMultiple deletes multiple objects from MinIO
// bucket: the bucket name
// objectNames: slice of object names/paths to delete
func (s *MinIOService) DeleteMultiple(ctx context.Context, bucket string, objectNames []string) error {
	objectsCh := make(chan minio.ObjectInfo)
	go func() {
		defer close(objectsCh)
		for _, objectName := range objectNames {
			objectsCh <- minio.ObjectInfo{
				Key: objectName,
			}
		}
	}()

	errorCh := s.client.RemoveObjects(ctx, bucket, objectsCh, minio.RemoveObjectsOptions{})
	for err := range errorCh {
		if err.Err != nil {
			return fmt.Errorf("failed to delete object %s: %w", err.ObjectName, err.Err)
		}
	}

	return nil
}

// List lists objects in a bucket
// bucket: the bucket name
// prefix: optional prefix to filter objects (e.g., "insurance-cards/2024/")
// recursive: whether to list recursively
func (s *MinIOService) List(ctx context.Context, bucket, prefix string, recursive bool) ([]ObjectInfo, error) {
	opts := minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: recursive,
	}

	var objects []ObjectInfo
	for object := range s.client.ListObjects(ctx, bucket, opts) {
		if object.Err != nil {
			return nil, fmt.Errorf("failed to list objects: %w", object.Err)
		}

		objects = append(objects, ObjectInfo{
			Bucket:       bucket,
			ObjectName:   object.Key,
			Size:         object.Size,
			ContentType:  object.ContentType,
			ETag:         object.ETag,
			LastModified: object.LastModified,
		})
	}

	return objects, nil
}

// GetObject retrieves an object from MinIO
// bucket: the bucket name
// objectName: the object name/path
func (s *MinIOService) GetObject(ctx context.Context, bucket, objectName string) (io.ReadCloser, error) {
	obj, err := s.client.GetObject(ctx, bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get object: %w", err)
	}
	return obj, nil
}

// GetObjectInfo retrieves metadata about an object without downloading it
// bucket: the bucket name
// objectName: the object name/path
func (s *MinIOService) GetObjectInfo(ctx context.Context, bucket, objectName string) (*ObjectInfo, error) {
	info, err := s.client.StatObject(ctx, bucket, objectName, minio.StatObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get object info: %w", err)
	}

	return &ObjectInfo{
		Bucket:       bucket,
		ObjectName:   info.Key,
		Size:         info.Size,
		ContentType:  info.ContentType,
		ETag:         info.ETag,
		LastModified: info.LastModified,
	}, nil
}

// BucketExists checks if a bucket exists
func (s *MinIOService) BucketExists(ctx context.Context, bucket string) (bool, error) {
	exists, err := s.client.BucketExists(ctx, bucket)
	if err != nil {
		return false, fmt.Errorf("failed to check bucket existence: %w", err)
	}
	return exists, nil
}
