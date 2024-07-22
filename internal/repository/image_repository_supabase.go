package repository

import (
	storage_go "github.com/supabase-community/storage-go"
	"mime/multipart"
	"onboarding_test/internal/service"
)

const (
	ImageBucket = "images"
)

type ImgRepoSupabase struct {
}

func NewImgRepoSupabase() *ImgRepoSupabase {
	//driver := supaProvider.GetDriver()
	//_, err := driver.GetBucket("images")
	//if err != nil {
	//}
	return &ImgRepoSupabase{}
}

func (imgRepoSupa ImgRepoSupabase) UploadImage(fileName string, file multipart.File, contentType string) (string, error) {
	Upsert := true
	//contentType := "image/jpeg"
	_, e := service.UploadFileToBucket(ImageBucket, fileName, file, storage_go.FileOptions{
		ContentType: &contentType,
		Upsert:      &Upsert,
	})
	if e != nil {
		return "", e
	}
	urlResponse := service.GetUrl(ImageBucket, fileName)
	return urlResponse, nil
}