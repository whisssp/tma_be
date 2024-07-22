package service

import (
	"fmt"
	storage_go "github.com/supabase-community/storage-go"
	"io"
	"onboarding_test/pkg/supabase"
)

var storageClient *storage_go.Client

func NewSupaStorageService(supaStorageClient *supabase.SupaStorageClient) {
	storageClient = supaStorageClient.GetDriver()
}

func CreateBucket(bucketId string, bucketOption storage_go.BucketOptions) error {
	_, err := storageClient.CreateBucket(bucketId, bucketOption)
	return fmt.Errorf("error creating bucket - %v", err)
}

func GetBucket(bucketId string) (*storage_go.Bucket, error) {
	result, err := storageClient.GetBucket(bucketId)
	if err != nil {
		return nil, fmt.Errorf("error getting bucket id: %v - %v", bucketId, err)
	}

	return &result, nil
}

func RemoveAllElementsOfBucket(bucketId string) (string, error) {
	result, err := storageClient.EmptyBucket(bucketId)
	if err != nil {
		return "", fmt.Errorf("error removing elements of bucket id: %v - %v", bucketId, err)
	}
	return result.Message, nil
}

func GetAllElementsFromBucket() ([]storage_go.Bucket, error) {
	result, err := storageClient.ListBuckets()
	if err != nil {
		return nil, fmt.Errorf("error getting all elements of bucket - %v", err)
	}
	return result, nil
}

func UploadFileToBucket(bucketId, filePath string, file io.Reader, fileOptions storage_go.FileOptions) (*storage_go.FileUploadResponse, error) {
	result, err := storageClient.UploadFile(bucketId, filePath, file, fileOptions)
	if err != nil {
		return nil, fmt.Errorf("error uploading file to bucket - %v", err)
	}
	return &result, nil
}

func GetUrl(bucketId, filePath string) string {
	url := storageClient.GetPublicUrl(bucketId, filePath, storage_go.UrlOptions{
		Transform: nil,
		Download:  false,
	})
	return url.SignedURL
}