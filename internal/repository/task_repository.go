package repository

import (
	"onboarding_test/internal/entity"
)

type TaskRepository interface {
	Create(*entity.Task) error
	Update(*entity.Task) (*entity.Task, error)
	GetTaskById(int64) (*entity.Task, error)
	GetAllTasks() ([]entity.Task, error)
	RemoveTask(*entity.Task) error
}