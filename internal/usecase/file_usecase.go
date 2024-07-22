package usecase

import (
	"fmt"
	"mime/multipart"
	"onboarding_test/internal/repository"
	"strconv"
	"time"
)

type FileUsecase struct {
	imgRepoSupabase *repository.ImgRepoSupabase
}

func NewFileUsecase(imgRepoSupabase *repository.ImgRepoSupabase) *FileUsecase {
	return &FileUsecase{
		imgRepoSupabase: imgRepoSupabase,
	}
}

func (fuc FileUsecase) UploadImage(fileName string, f multipart.File, contentType string) (string, error) {
	instant := strconv.FormatInt(time.Now().Unix(), 10)
	result, err := fuc.imgRepoSupabase.UploadImage(fmt.Sprintf("%v%v", instant, fileName), f, contentType)
	if err != nil {
		fmt.Printf("Upload Failed: %v", err)
		return "", err
	}

	return result, nil
}