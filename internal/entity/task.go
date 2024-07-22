package entity

import (
	"gorm.io/gorm"
)

const (
	StatusPending   = "pending"
	StatusCompleted = "completed"
	StatusDeleted   = "deleted"
)

var taskStatus = map[int]string{
	1: StatusPending,
	2: StatusCompleted,
}

//var taskStatus = []string{StatusPending, StatusCompleted}

type Task struct {
	ID          int64  `gorm:"type:bigserial;primaryKey" json:"id"`
	Title       string `gorm:"type:varchar(255)" validate:"required" json:"title"`
	Image       string `gorm:"type:text" json:"imageUrl"`
	Description string `gorm:"type:varchar(255)" json:"description"`
	Status      string `gorm:"type:varchar(50)" validate:"required" json:"status"`
	gorm.Model
}

func GetListValidStatus() map[int]string {
	return taskStatus
}