package repository

import (
	"fmt"
	"gorm.io/gorm"
	"math"
	"onboarding_test/internal/delivery/http/payload"
	"onboarding_test/internal/entity"
)

type TaskRepoGorm struct {
	db *gorm.DB
}

func NewTaskRepoGorm(db *gorm.DB) *TaskRepoGorm {
	return &TaskRepoGorm{
		db: db,
	}
}

func (r TaskRepoGorm) Create(task *entity.Task) error {
	tx := r.db.Begin()
	if err := tx.Error; err != nil {
		tx.Rollback()
		return payload.ErrDB(err)
	}
	if err := tx.Create(&task).Error; err != nil {
		tx.Rollback()
		return payload.ErrDB(err)
	}
	if err := tx.Commit().Error; err != nil {
		return payload.ErrDB(err)
	}
	return nil
}

func (r TaskRepoGorm) Update(task *entity.Task) (*entity.Task, error) {
	tx := r.db.Begin()
	if err := tx.Error; err != nil {
		tx.Rollback()
		return nil, payload.ErrDB(err)
	}

	if err := tx.Updates(task).Error; err != nil {
		tx.Rollback()
		return nil, payload.ErrDB(err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, payload.ErrDB(err)
	}

	return task, nil
}

func (r TaskRepoGorm) GetTaskById(id int64) (*entity.Task, error) {
	var task entity.Task
	if err := r.db.Where("id = ?", id).First(&task).Error; err != nil {
		return nil, payload.ErrDB(err)
	}
	return &task, nil
}

func (r TaskRepoGorm) GetAllTasksWithFilterAndPaginate(pagination *payload.Pagination, filter payload.TaskFilter) ([]entity.Task, error) {
	var totalRows int64
	tasks := make([]entity.Task, 0)
	db := r.db.Scopes(filterTasks(r.db, filter))
	if pagination == nil {
		if err := db.Find(&tasks).Count(&totalRows).Error; err != nil {
			return nil, payload.ErrDB(err)
		}
		pagination.TotalRows = totalRows
		pagination.TotalPages = int(math.Ceil(float64(totalRows) / float64(pagination.Limit)))
		return tasks, nil
	}
	if err := db.Scopes(paginate(tasks, pagination, r.db)).Find(&tasks).Count(&totalRows).Error; err != nil {
		return nil, payload.ErrDB(err)
	}

	pagination.TotalRows = totalRows
	pagination.TotalPages = int(math.Ceil(float64(totalRows) / float64(pagination.Limit)))

	return tasks, nil
}

func (r TaskRepoGorm) GetAllTasks() ([]entity.Task, error) {
	tasks := make([]entity.Task, 0)
	if err := r.db.Find(&tasks).Error; err != nil {
		return nil, payload.ErrDB(err)
	}
	return tasks, nil
}

func (r TaskRepoGorm) RemoveTask(task *entity.Task) error {
	tx := r.db.Begin()
	if err := tx.Error; err != nil {
		tx.Rollback()
		return payload.ErrDB(err)
	}

	if task.DeletedAt.Valid {
		tx.Rollback()
		return payload.ErrDB(fmt.Errorf("this is task already deleted"))
	}

	if err := tx.Delete(task).Error; err != nil {
		return payload.ErrDB(err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return payload.ErrDB(err)
	}
	return nil
}

func paginate(value interface{}, pagination *payload.Pagination, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	//var totalRows int64
	//db.Model(value).Count(&totalRows)
	//pagination.TotalRows = totalRows
	//totalPages := int(math.Ceil(float64(totalRows) / float64(pagination.Limit)))
	//pagination.TotalPages = totalPages
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit()).Order(pagination.GetSort())
	}
}

func filterTasks(db *gorm.DB, filter payload.TaskFilter) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if filter.Keyword != "" {
			keyword := "%" + filter.Keyword + "%"
			db = db.Where("title LIKE ? OR description LIKE ?", keyword, keyword)
		}
		if filter.ID != 0 {
			db = db.Where("id = ?", filter.ID)
		}
		if filter.Title != "" {
			db = db.Where("title LIKE ?", "%"+filter.Title+"%")
		}
		if filter.Description != "" {
			db = db.Where("description LIKE ?", "%"+filter.Description+"%")
		}
		if filter.Status != 0 {
			db = db.Where("status = ?", entity.GetListValidStatus()[filter.Status])
		}
		if filter.CreatedAtFrom != nil {
			db = db.Where("created_at >= ?", *filter.CreatedAtFrom)
		}
		if filter.CreatedAtTo != nil {
			db = db.Where("created_at <= ?", *filter.CreatedAtTo)
		}
		if filter.UpdatedAtFrom != nil {
			db = db.Where("updated_at >= ?", *filter.UpdatedAtFrom)
		}
		if filter.UpdatedAtTo != nil {
			db = db.Where("updated_at <= ?", *filter.UpdatedAtTo)
		}
		return db
	}
}