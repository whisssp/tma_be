package payload

import "net/http"

type AppResponse struct {
	StatusCode int         `json:"code"`
	Message    *string     `json:"message"`
	Data       interface{} `json:"data"`
}

func NewSimpleSuccessResponse(data interface{}) *AppResponse {
	return &AppResponse{
		Data:       data,
		StatusCode: http.StatusOK,
		Message:    nil,
	}
}