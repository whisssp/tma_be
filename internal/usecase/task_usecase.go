package usecase

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"onboarding_test/internal/delivery/http/payload"
	"onboarding_test/internal/entity"
	"onboarding_test/internal/repository"
)

const (
	entityName   = "tasks"
	actionUpdate = "update"
	actionDelete = "delete"
)

type TaskUsecase struct {
	taskRepoGorm  *repository.TaskRepoGorm
	taskRepoRedis *repository.TaskRepoRedis
}

func NewTaskUsecase(taskRepo *repository.TaskRepoGorm, taskRepoRedis *repository.TaskRepoRedis) *TaskUsecase {
	return &TaskUsecase{
		taskRepoGorm:  taskRepo,
		taskRepoRedis: taskRepoRedis,
	}
}

func (t TaskUsecase) Create(createTaskRequest payload.CreateTaskRequest) (*payload.TaskResponse, error) {
	if err := validateStruct(createTaskRequest); err != nil {
		return nil, payload.ErrInvalidRequest(err)
	}

	if err := validateStatus(createTaskRequest.Status); err != nil {
		return nil, payload.ErrInvalidRequest(err)
	}

	task := createTaskRqToTask(createTaskRequest)

	if err := t.taskRepoGorm.Create(task); err != nil {
		return nil, payload.ErrCannotCreateEntity(entityName, err)
	}
	err := t.taskRepoRedis.Create(task)
	if err != nil {
		fmt.Printf("\n error creating task on redis: %v\n", err)
		go t.LoadTasksToRedis()
	}

	taskResponse := taskToTaskResponse(task)
	return &taskResponse, nil
}

func (t TaskUsecase) UpdateTaskById(updateTaskRequest payload.UpdateTaskRequest) (*payload.TaskResponse, error) {
	if err := validateStruct(updateTaskRequest); err != nil {
		return nil, payload.ErrInvalidRequest(err)
	}

	if err := validateStatus(updateTaskRequest.Status); err != nil {
		return nil, payload.ErrInvalidRequest(err)
	}

	oldTask, err := t.taskRepoGorm.GetTaskById(updateTaskRequest.Id)
	if oldTask == nil {
		return nil, payload.ErrEntityNotFound(entityName, err)
	}

	updateTask(oldTask, updateTaskRqToTask(updateTaskRequest))
	taskUpdated, err := t.taskRepoGorm.Update(oldTask)
	if err != nil {
		return nil, payload.ErrCannotUpdateEntity(entityName, err)
	}
	taskResponse := taskToTaskResponse(taskUpdated)

	_, err = t.taskRepoRedis.Update(taskUpdated)
	if err != nil {
		fmt.Printf("\nerror updating task on Redis: %v\n", err)
		go t.LoadTasksToRedis()
	}
	return &taskResponse, nil
}

func (t TaskUsecase) GetTaskById(id int64) (*payload.TaskResponse, error) {
	task := t.taskRepoRedis.GetTaskById(id)
	var taskResponse payload.TaskResponse
	if task.ID != 0 {
		taskResponse = taskToTaskResponse(task)
		return &taskResponse, nil
	}

	task, err := t.taskRepoGorm.GetTaskById(id)
	if err != nil {
		return nil, payload.ErrEntityNotFound(entityName, err)
	}

	taskResponse = taskToTaskResponse(task)
	go t.LoadTasksToRedis()
	return &taskResponse, nil
}

func (t TaskUsecase) GetAllTasks(limit, page int, sortBy string, filter payload.TaskFilter) (*payload.ListTaskResponse, error) {
	pagination := payload.Pagination{Limit: limit, Page: page, Sort: sortBy}

	//tasks := t.taskRepoRedis.GetAllTasks(&pagination)
	//if len(tasks) != 0 {
	//	tasksResponse := tasksToListTaskResponse(tasks, pagination)
	//	return &tasksResponse, nil
	//
	//}
	tasks, err := t.taskRepoGorm.GetAllTasksWithFilterAndPaginate(&pagination, filter)
	if err != nil {
		return nil, payload.ErrEntityNotFound(entityName, err)
	}
	tasksResponse := tasksToListTaskResponse(tasks, pagination)
	return &tasksResponse, nil
}

func (t TaskUsecase) RemoveTaskById(id int64) error {
	task, err := t.taskRepoGorm.GetTaskById(id)
	if err != nil {
		return payload.ErrEntityNotFound(entityName, err)
	}

	//updateTask(task, nil, actionDelete)
	if err := t.taskRepoGorm.RemoveTask(task); err != nil {
		return payload.ErrCannotDeleteEntity(entityName, err)
	}
	err = t.taskRepoRedis.RemoveTask(task)
	if err != nil {
		fmt.Printf("error remove task on Redis: %v\n", err)
		go t.LoadTasksToRedis()
	}
	return nil
}

func (t TaskUsecase) LoadTasksToRedis() {
	listTasks, err := t.taskRepoGorm.GetAllTasks()
	if err != nil {
		fmt.Printf("error getting tasks from DB: %v\n", err)
		return
	}
	t.taskRepoRedis.LoadTasksToRedis(listTasks)
}

func createTaskRqToTask(reqPayload payload.CreateTaskRequest) *entity.Task {
	return &entity.Task{
		Title:       reqPayload.Title,
		Description: reqPayload.Description,
		Status:      entity.GetListValidStatus()[reqPayload.Status],
		Image:       reqPayload.Image,
	}
}

func updateTaskRqToTask(reqPayload payload.UpdateTaskRequest) *entity.Task {
	return &entity.Task{
		ID:          reqPayload.Id,
		Title:       reqPayload.Title,
		Description: reqPayload.Description,
		Status:      entity.GetListValidStatus()[reqPayload.Status],
		Image:       reqPayload.Image,
	}
}

func taskToTaskResponse(task *entity.Task) payload.TaskResponse {
	return payload.TaskResponse{
		Id:          task.ID,
		Title:       task.Title,
		Image:       task.Image,
		Description: task.Description,
		Status:      task.Status,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}
}

func tasksToListTaskResponse(tasks []entity.Task, pagination payload.Pagination) payload.ListTaskResponse {
	listTasks := make([]payload.TaskResponse, 0)
	for _, t := range tasks {
		taskResponse := taskToTaskResponse(&t)
		listTasks = append(listTasks, taskResponse)
	}
	return payload.ListTaskResponse{
		Tasks:         listTasks,
		Limit:         pagination.Limit,
		Page:          pagination.Page,
		TotalElements: pagination.TotalRows,
	}
}

func updateTask(oldTask, newTask *entity.Task) {
	oldTask.Title = newTask.Title
	oldTask.Description = newTask.Description
	oldTask.Status = newTask.Status
	oldTask.Image = newTask.Image
}

// use map
func validateStatus(status int) error {
	if _, exist := entity.GetListValidStatus()[status]; exist == false {
		return fmt.Errorf("invalid status, status must be in one of these values %v", entity.GetListValidStatus())
	}
	return nil
}

func validateStruct(e interface{}) error {
	return validator.New().Struct(e)
}