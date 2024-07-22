package payload

import (
	"errors"
	"fmt"
	"net/http"
)

type ErrResponse struct {
	StatusCode int         `json:"code"`
	Message    string      `json:"message"`
	ErrKey     string      `json:"-"`
	RootErr    error       `json:"-"`
	Data       interface{} `json:"data"`
}

func NewErrResponse(code int, message, key string, RootErr error) *ErrResponse {
	return &ErrResponse{
		StatusCode: code,
		Message:    fmt.Sprintf("%s: %v", key, message),
		ErrKey:     key,
		RootErr:    RootErr,
		Data:       nil,
	}
}

func (e *ErrResponse) RootError() error {
	if err, ok := e.RootErr.(*ErrResponse); ok {
		return err.RootError()
	}
	return e.RootErr
}

func (e *ErrResponse) Error() string {
	return e.RootError().Error()
}

func NewCustomError(code int, root error, msg, key string) *ErrResponse {
	if root != nil {
		return NewErrResponse(code, msg, key, root)
	}
	return NewErrResponse(code, root.Error(), key, errors.New(msg))
}

func ErrMissingParams(err error) *ErrResponse {
	return NewErrResponse(http.StatusBadRequest, err.Error(), "ErrMissingParam", err)
}

func ErrDB(err error) *ErrResponse {
	return NewErrResponse(http.StatusInternalServerError, err.Error(), "ERR_DB", err)
}

func ErrInvalidRequest(err error) *ErrResponse {
	return NewCustomError(http.StatusBadRequest, err, err.Error(), "ErrInvalidRequest")
}

func ErrEntityNotFound(entity string, err error) *ErrResponse {
	return NewCustomError(http.StatusNotFound, err, err.Error(), fmt.Sprintf("Err%sNotFound", entity))
}

func ErrCannotDeleteEntity(entity string, err error) *ErrResponse {
	return NewCustomError(http.StatusBadRequest, err, err.Error(), fmt.Sprintf("ErrCannotDelete%s", entity))
}

func ErrCannotCreateEntity(entity string, err error) *ErrResponse {
	if err == nil {
		return NewErrResponse(http.StatusBadRequest, err.Error(), fmt.Sprintf("ErrCannotCreate%s", entity), nil)
	}
	return NewErrResponse(http.StatusBadRequest, err.Error(), fmt.Sprintf("ErrCannotCreate%s", entity), err)
}

func ErrCannotUpdateEntity(entity string, err error) *ErrResponse {
	return NewCustomError(http.StatusBadRequest, err, err.Error(), fmt.Sprintf("ErrCannotUpdate%s", entity))
}

func ErrInvalidRequestBody(err error) *ErrResponse {
	return NewCustomError(http.StatusBadRequest, err, err.Error(), "InvalidRequestBody")
}

func ErrConvertQueryParamFailed(err error) *ErrResponse {
	return NewCustomError(http.StatusBadRequest, err, err.Error(), "ConvertQueryParamFailed")
}

// File Error

func ErrUploadFileFailed(err error) *ErrResponse {
	return NewCustomError(http.StatusInternalServerError, err, err.Error(), "UploadFileFailed")
}

func ErrDetectFileType(err error) *ErrResponse {
	return NewCustomError(http.StatusInternalServerError, err, err.Error(), "ErrDetectFileType")
}

func ErrResetFilPointer(err error) *ErrResponse {
	return NewCustomError(http.StatusInternalServerError, err, err.Error(), "ErrResetFilePointer")
}