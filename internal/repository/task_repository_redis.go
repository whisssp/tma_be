package repository

import (
	"fmt"
	"onboarding_test/internal/delivery/http/payload"
	"onboarding_test/internal/entity"
	"onboarding_test/internal/service"
	"strconv"
	"time"
)

const (
	redisKey     = "tasks"
	redisExpTime = 10 * time.Minute
)

type TaskRepoRedis struct {
}

func NewTaskRepoRedis() *TaskRepoRedis {
	return &TaskRepoRedis{}
}

func (taskRepoRedis *TaskRepoRedis) Create(task *entity.Task) error {
	idStr := strconv.FormatInt(task.ID, 10)
	return service.RedisSetHashGenericKey(redisKey, idStr, task, redisExpTime)
}

func (taskRepoRedis *TaskRepoRedis) Update(task *entity.Task) (*entity.Task, error) {
	idStr := strconv.FormatInt(task.ID, 10)
	err := service.RedisSetHashGenericKey(redisKey, idStr, task, redisExpTime)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (taskRepoRedis *TaskRepoRedis) GetTaskById(id int64) *entity.Task {
	var task entity.Task
	service.RedisGetHashGenericKey(redisKey, strconv.FormatInt(id, 10), &task)
	return &task
}

func (taskRepoRedis *TaskRepoRedis) GetAllTasks(pagination *payload.Pagination) []entity.Task {
	mapResult := make(map[string]entity.Task)
	taskIDs := service.GetHashGenericWithPagination(redisKey, pagination)
	fmt.Printf("Tasks: %v", taskIDs)

	return ToTaskArray(mapResult)
}

func (taskRepoRedis *TaskRepoRedis) RemoveTask(task *entity.Task) error {
	return service.RedisRemoveHashGenericKey(redisKey, strconv.FormatInt(task.ID, 10))
}

func (taskRepoRedis *TaskRepoRedis) LoadTasksToRedis(tasks []entity.Task) {
	err := service.RedisSetHashGenericKeySlice(redisKey, tasks, GetTaskID, redisExpTime)
	if err != nil {
		fmt.Printf("\nerror: %v\n", err)
		return
	}
}

func ToTaskArray(mapTasks map[string]entity.Task) []entity.Task {
	arrTasks := make([]entity.Task, 0)
	for _, v := range mapTasks {
		arrTasks = append(arrTasks, v)
	}
	return arrTasks
}

func GetTaskID(task entity.Task) int64 {
	return task.ID
}