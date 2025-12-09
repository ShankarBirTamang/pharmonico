// Package services provides MinIO service tests
package services

import (
	"bytes"
	"context"
	"os"
	"strings"
	"testing"
	"time"
)

// getTestMinIOService returns a MinIO service client for testing
func getTestMinIOService(t *testing.T) *MinIOService {
	endpoint := os.Getenv("MINIO_ENDPOINT")
	if endpoint == "" {
		endpoint = "localhost:9000"
	}

	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	if accessKey == "" {
		accessKey = "minioadmin"
	}

	secretKey := os.Getenv("MINIO_SECRET_KEY")
	if secretKey == "" {
		secretKey = "minioadmin"
	}

	// Determine if we should use SSL (typically false for local development)
	useSSL := os.Getenv("MINIO_USE_SSL") == "true"

	service, err := NewMinIOService(endpoint, accessKey, secretKey, useSSL)
	if err != nil {
		t.Fatalf("Failed to create MinIO service: %v", err)
	}

	return service
}

// TestMinIOService_Upload tests the upload functionality
func TestMinIOService_Upload(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	service := getTestMinIOService(t)
	ctx := context.Background()

	// Test bucket
	bucket := "insurance-cards"
	testObjectName := "test/upload-test.txt"
	testContent := "This is a test file for MinIO upload"

	// Ensure bucket exists
	exists, err := service.BucketExists(ctx, bucket)
	if err != nil {
		t.Fatalf("Failed to check bucket existence: %v", err)
	}
	if !exists {
		t.Skipf("Bucket %s does not exist. Please ensure MinIO is running and buckets are initialized.", bucket)
	}

	// Clean up any existing test object
	_ = service.Delete(ctx, bucket, testObjectName)

	// Upload test
	reader := bytes.NewReader([]byte(testContent))
	opts := &UploadOptions{
		ContentType: "text/plain",
		Metadata: map[string]string{
			"test": "true",
		},
	}

	info, err := service.Upload(ctx, bucket, testObjectName, reader, int64(len(testContent)), opts)
	if err != nil {
		t.Fatalf("Upload failed: %v", err)
	}

	if info == nil {
		t.Fatal("Upload returned nil info")
	}

	if info.Bucket != bucket {
		t.Errorf("Expected bucket %s, got %s", bucket, info.Bucket)
	}

	if info.ObjectName != testObjectName {
		t.Errorf("Expected object name %s, got %s", testObjectName, info.ObjectName)
	}

	if info.Size != int64(len(testContent)) {
		t.Errorf("Expected size %d, got %d", len(testContent), info.Size)
	}

	// Clean up
	defer func() {
		if err := service.Delete(ctx, bucket, testObjectName); err != nil {
			t.Logf("Failed to clean up test object: %v", err)
		}
	}()

	t.Logf("Upload test passed: %+v", info)
}

// TestMinIOService_GenerateSignedURL tests signed URL generation
func TestMinIOService_GenerateSignedURL(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	service := getTestMinIOService(t)
	ctx := context.Background()

	bucket := "insurance-cards"
	testObjectName := "test/signed-url-test.txt"
	testContent := "Test content for signed URL"

	// Ensure bucket exists
	exists, err := service.BucketExists(ctx, bucket)
	if err != nil {
		t.Fatalf("Failed to check bucket existence: %v", err)
	}
	if !exists {
		t.Skipf("Bucket %s does not exist. Please ensure MinIO is running and buckets are initialized.", bucket)
	}

	// Upload a test object first
	reader := bytes.NewReader([]byte(testContent))
	_, err = service.Upload(ctx, bucket, testObjectName, reader, int64(len(testContent)), nil)
	if err != nil {
		t.Fatalf("Failed to upload test object: %v", err)
	}

	// Clean up
	defer func() {
		if err := service.Delete(ctx, bucket, testObjectName); err != nil {
			t.Logf("Failed to clean up test object: %v", err)
		}
	}()

	// Generate signed URL
	expiry := 15 * time.Minute
	url, err := service.GenerateSignedURL(ctx, bucket, testObjectName, expiry)
	if err != nil {
		t.Fatalf("Failed to generate signed URL: %v", err)
	}

	if url == "" {
		t.Fatal("Generated signed URL is empty")
	}

	// URL should contain the bucket and object name
	if !strings.Contains(url, bucket) || !strings.Contains(url, testObjectName) {
		t.Errorf("Signed URL does not contain expected bucket or object name: %s", url)
	}

	t.Logf("Generated signed URL: %s", url)
}

// TestMinIOService_Delete tests the delete functionality
func TestMinIOService_Delete(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	service := getTestMinIOService(t)
	ctx := context.Background()

	bucket := "insurance-cards"
	testObjectName := "test/delete-test.txt"
	testContent := "Test content for deletion"

	// Ensure bucket exists
	exists, err := service.BucketExists(ctx, bucket)
	if err != nil {
		t.Fatalf("Failed to check bucket existence: %v", err)
	}
	if !exists {
		t.Skipf("Bucket %s does not exist. Please ensure MinIO is running and buckets are initialized.", bucket)
	}

	// Upload a test object first
	reader := bytes.NewReader([]byte(testContent))
	_, err = service.Upload(ctx, bucket, testObjectName, reader, int64(len(testContent)), nil)
	if err != nil {
		t.Fatalf("Failed to upload test object: %v", err)
	}

	// Verify object exists
	info, err := service.GetObjectInfo(ctx, bucket, testObjectName)
	if err != nil {
		t.Fatalf("Failed to get object info before deletion: %v", err)
	}
	if info == nil {
		t.Fatal("Object info is nil before deletion")
	}

	// Delete the object
	err = service.Delete(ctx, bucket, testObjectName)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Verify object is deleted
	_, err = service.GetObjectInfo(ctx, bucket, testObjectName)
	if err == nil {
		t.Fatal("Object still exists after deletion")
	}

	t.Log("Delete test passed")
}

// TestMinIOService_List tests the list functionality
func TestMinIOService_List(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	service := getTestMinIOService(t)
	ctx := context.Background()

	bucket := "insurance-cards"
	testPrefix := "test/list-test-"

	// Ensure bucket exists
	exists, err := service.BucketExists(ctx, bucket)
	if err != nil {
		t.Fatalf("Failed to check bucket existence: %v", err)
	}
	if !exists {
		t.Skipf("Bucket %s does not exist. Please ensure MinIO is running and buckets are initialized.", bucket)
	}

	// Upload multiple test objects
	testObjects := []string{
		testPrefix + "1.txt",
		testPrefix + "2.txt",
		testPrefix + "3.txt",
	}

	for _, objName := range testObjects {
		content := "Test content for " + objName
		reader := bytes.NewReader([]byte(content))
		_, err := service.Upload(ctx, bucket, objName, reader, int64(len(content)), nil)
		if err != nil {
			t.Fatalf("Failed to upload test object %s: %v", objName, err)
		}
	}

	// Clean up
	defer func() {
		for _, objName := range testObjects {
			if err := service.Delete(ctx, bucket, objName); err != nil {
				t.Logf("Failed to clean up test object %s: %v", objName, err)
			}
		}
	}()

	// List objects with prefix
	objects, err := service.List(ctx, bucket, testPrefix, false)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(objects) < len(testObjects) {
		t.Errorf("Expected at least %d objects, got %d", len(testObjects), len(objects))
	}

	// Verify all test objects are in the list
	found := make(map[string]bool)
	for _, obj := range objects {
		found[obj.ObjectName] = true
	}

	for _, objName := range testObjects {
		if !found[objName] {
			t.Errorf("Expected object %s not found in list", objName)
		}
	}

	t.Logf("List test passed: found %d objects", len(objects))
}

// TestMinIOService_GetObject tests retrieving an object
func TestMinIOService_GetObject(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	service := getTestMinIOService(t)
	ctx := context.Background()

	bucket := "insurance-cards"
	testObjectName := "test/get-test.txt"
	testContent := "Test content for retrieval"

	// Ensure bucket exists
	exists, err := service.BucketExists(ctx, bucket)
	if err != nil {
		t.Fatalf("Failed to check bucket existence: %v", err)
	}
	if !exists {
		t.Skipf("Bucket %s does not exist. Please ensure MinIO is running and buckets are initialized.", bucket)
	}

	// Upload a test object first
	reader := bytes.NewReader([]byte(testContent))
	_, err = service.Upload(ctx, bucket, testObjectName, reader, int64(len(testContent)), nil)
	if err != nil {
		t.Fatalf("Failed to upload test object: %v", err)
	}

	// Clean up
	defer func() {
		if err := service.Delete(ctx, bucket, testObjectName); err != nil {
			t.Logf("Failed to clean up test object: %v", err)
		}
	}()

	// Get the object
	objReader, err := service.GetObject(ctx, bucket, testObjectName)
	if err != nil {
		t.Fatalf("GetObject failed: %v", err)
	}
	defer objReader.Close()

	// Read the content
	buf := make([]byte, len(testContent))
	n, err := objReader.Read(buf)
	if err != nil && err.Error() != "EOF" {
		t.Fatalf("Failed to read object content: %v", err)
	}

	retrievedContent := string(buf[:n])
	if retrievedContent != testContent {
		t.Errorf("Expected content %s, got %s", testContent, retrievedContent)
	}

	t.Log("GetObject test passed")
}

// TestMinIOService_GetObjectInfo tests retrieving object metadata
func TestMinIOService_GetObjectInfo(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	service := getTestMinIOService(t)
	ctx := context.Background()

	bucket := "insurance-cards"
	testObjectName := "test/info-test.txt"
	testContent := "Test content for info"

	// Ensure bucket exists
	exists, err := service.BucketExists(ctx, bucket)
	if err != nil {
		t.Fatalf("Failed to check bucket existence: %v", err)
	}
	if !exists {
		t.Skipf("Bucket %s does not exist. Please ensure MinIO is running and buckets are initialized.", bucket)
	}

	// Upload a test object first
	reader := bytes.NewReader([]byte(testContent))
	_, err = service.Upload(ctx, bucket, testObjectName, reader, int64(len(testContent)), nil)
	if err != nil {
		t.Fatalf("Failed to upload test object: %v", err)
	}

	// Clean up
	defer func() {
		if err := service.Delete(ctx, bucket, testObjectName); err != nil {
			t.Logf("Failed to clean up test object: %v", err)
		}
	}()

	// Get object info
	info, err := service.GetObjectInfo(ctx, bucket, testObjectName)
	if err != nil {
		t.Fatalf("GetObjectInfo failed: %v", err)
	}

	if info == nil {
		t.Fatal("GetObjectInfo returned nil")
	}

	if info.Bucket != bucket {
		t.Errorf("Expected bucket %s, got %s", bucket, info.Bucket)
	}

	if info.ObjectName != testObjectName {
		t.Errorf("Expected object name %s, got %s", testObjectName, info.ObjectName)
	}

	if info.Size != int64(len(testContent)) {
		t.Errorf("Expected size %d, got %d", len(testContent), info.Size)
	}

	t.Logf("GetObjectInfo test passed: %+v", info)
}

// TestMinIOService_DeleteMultiple tests deleting multiple objects
func TestMinIOService_DeleteMultiple(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	service := getTestMinIOService(t)
	ctx := context.Background()

	bucket := "insurance-cards"
	testPrefix := "test/delete-multiple-"

	// Ensure bucket exists
	exists, err := service.BucketExists(ctx, bucket)
	if err != nil {
		t.Fatalf("Failed to check bucket existence: %v", err)
	}
	if !exists {
		t.Skipf("Bucket %s does not exist. Please ensure MinIO is running and buckets are initialized.", bucket)
	}

	// Upload multiple test objects
	testObjects := []string{
		testPrefix + "1.txt",
		testPrefix + "2.txt",
		testPrefix + "3.txt",
	}

	for _, objName := range testObjects {
		content := "Test content for " + objName
		reader := bytes.NewReader([]byte(content))
		_, err := service.Upload(ctx, bucket, objName, reader, int64(len(content)), nil)
		if err != nil {
			t.Fatalf("Failed to upload test object %s: %v", objName, err)
		}
	}

	// Delete multiple objects
	err = service.DeleteMultiple(ctx, bucket, testObjects)
	if err != nil {
		t.Fatalf("DeleteMultiple failed: %v", err)
	}

	// Verify all objects are deleted
	for _, objName := range testObjects {
		_, err := service.GetObjectInfo(ctx, bucket, objName)
		if err == nil {
			t.Errorf("Object %s still exists after deletion", objName)
		}
	}

	t.Log("DeleteMultiple test passed")
}
