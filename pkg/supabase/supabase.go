package supabase

import (
	"fmt"
	supa "github.com/nedpals/supabase-go"
	storage_go "github.com/supabase-community/storage-go"
	"onboarding_test/internal/config"
)

type SupabaseClient struct {
	supabaseClient *supa.Client
}

type SupaStorageClient struct {
	storageClient *storage_go.Client
}

func NewSupabaseClient() *SupabaseClient {
	supaClient := supa.CreateClient(config.Envs.SupUrl, config.Envs.SupKey)
	if supaClient == nil {
		fmt.Println("Failed to connect to Supabase")
	}
	fmt.Println("Connected to Supabase successfully")
	return &SupabaseClient{
		supabaseClient: supaClient,
	}
}

func NewSupaStorageClient() *SupaStorageClient {
	supaClient := storage_go.NewClient(config.Envs.SupStorageRawUrl, config.Envs.SupKey, nil)

	if supaClient != nil {
		fmt.Println("Connected to Supabase Storage successfully")
		return &SupaStorageClient{
			storageClient: supaClient,
		}
	}
	fmt.Println("Failed to connect to Supabase Storage")
	return nil
}

func (storageClient *SupaStorageClient) GetDriver() *storage_go.Client {
	return storageClient.storageClient
}