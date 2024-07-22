package payload

type CreateTaskRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
	Status      int    `json:"status"`
	Image       string `json:"imageUrl"`
}

type UpdateTaskRequest struct {
	Id          int64  `json:"id"`
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
	Status      int    `json:"status" validate:"required"`
	Image       string `json:"imageUrl"`
}