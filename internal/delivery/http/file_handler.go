package http

import (
	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
	"net/http"
	"onboarding_test/internal/delivery/http/payload"
	"onboarding_test/internal/usecase"
	"onboarding_test/utils"
)

type FileHandler struct {
	fileUsecase *usecase.FileUsecase
}

func NewFileHandler(fileUsecase *usecase.FileUsecase) *FileHandler {
	return &FileHandler{
		fileUsecase: fileUsecase,
	}
}

func (h FileHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/files/upload", h.handleUploadFile)
}

// Upload file 			godoc
// @Summary 			Upload a file
// @Description			Upload a file to get link get media
// Tag					File
// @Param				file formData file true "file"
// @Success				200		{object} payload.AppResponse{}
// @Failure      		400  	{object} payload.ErrResponse{}
// @Failure 			500 	{object} payload.ErrResponse{}
// @Router				/files/upload [post]
func (h FileHandler) handleUploadFile(c *gin.Context) {
	file, header, _ := c.Request.FormFile("file")
	fileName := utils.SanitizeFileName(header.Filename)
	mime, err := mimetype.DetectReader(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, payload.ErrDetectFileType(err))
		return
	}
	if _, err := file.Seek(0, 0); err != nil {
		c.JSON(http.StatusInternalServerError, payload.ErrDetectFileType(err))
		return
	}
	defer file.Close()
	contentType := mime.String()
	r, err := h.fileUsecase.UploadImage(fileName, file, contentType)
	if err != nil {
		c.JSON(http.StatusBadRequest, payload.ErrUploadFileFailed(err))
		return
	}
	c.JSON(http.StatusOK, payload.NewSimpleSuccessResponse(r))
}