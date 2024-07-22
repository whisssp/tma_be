package payload

import "time"

type TaskResponse struct {
	Id          int64     `json:"id"`
	Title       string    `json:"title"`
	Image       string    `json:"imageUrl"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type ListTaskResponse struct {
	Tasks         []TaskResponse `json:"tasks"`
	Limit         int            `json:"pageSize"`
	Page          int            `json:"pageNumber"`
	TotalElements int64          `json:"totalElements"`
}