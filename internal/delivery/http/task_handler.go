package http

import (
	"fmt"
	"net/http"
	"onboarding_test/internal/delivery/http/payload"
	"onboarding_test/internal/usecase"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	taskUsecase *usecase.TaskUsecase
}

func NewTaskHandler(taskUsecase *usecase.TaskUsecase) *TaskHandler {
	return &TaskHandler{taskUsecase}
}

func (h TaskHandler) RegisterRoutes(routerGroup *gin.RouterGroup) {
	tasksRoute := routerGroup.Group("/tasks")
	{
		tasksRoute.POST("", h.handleCreateNewTask)
		tasksRoute.GET("", h.handleGetAllTasks)
		tasksRoute.GET("/*id", h.handleGetTaskById)
		tasksRoute.PUT("/*id", h.handleUpdateTaskById)
		tasksRoute.DELETE("/*id", h.handleRemoveTaskById)
	}
}

// CreateNewTask		godoc
// @Summary				Create tasks
// @Description		Save task data in Db.
// @Param				task body payload.CreateTaskRequest true "create task request"
// @Produce				application/json
// @Tags					tasks
// @Success				204
// @Failure      		400  	{object} payload.ErrResponse{}
// @Failure 			404 	{object} payload.ErrResponse{}
// @Router			/tasks 	[post]
func (h TaskHandler) handleCreateNewTask(c *gin.Context) {
	var reqPayload payload.CreateTaskRequest
	if err := c.ShouldBindJSON(&reqPayload); err != nil {
		c.JSON(http.StatusBadRequest, payload.ErrInvalidRequestBody(err))
		return
	}
	taskCreated, err := h.taskUsecase.Create(reqPayload)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, payload.NewSimpleSuccessResponse(taskCreated))
}

// GetAllTasks			godoc
// @Summary				Get all tasks
// @Description		Return list of tasks.
// @Tags					tasks
// @Success				200 {object} payload.AppResponse{}
// @Failure      		400  	{object} payload.ErrResponse{}
// @Failure 			404 	{object} payload.ErrResponse{}
// @Router				/tasks [get]
func (h TaskHandler) handleGetAllTasks(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	sortBy := c.DefaultQuery("sort", "")
	fmt.Printf("Sort Param %v", sortBy)

	var taskFilter payload.TaskFilter
	if err := c.ShouldBindQuery(&taskFilter); err != nil {
		c.JSON(http.StatusBadRequest, payload.ErrInvalidRequest(err))
		return
	}

	tasks, err := h.taskUsecase.GetAllTasks(limit, page, sortBy, taskFilter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, payload.NewSimpleSuccessResponse(tasks))
}

// GetTaskById			godoc
// @Summary				Get a task by id
// @Description		Get a task by id from path variable named id from DB.
// @Param				id path string  true "id"
// @Produce				application/json
// @Tags					tasks
// @Failure      		400  	{object} payload.ErrResponse{}
// @Failure 			404 	{object} payload.ErrResponse{}
// @Router				/tasks/{id} [get]
func (h TaskHandler) handleGetTaskById(c *gin.Context) {
	id, _ := strconv.ParseInt(removeSlashFromParam(c.Param("id")), 0, 64)
	if id == 0 {
		c.JSON(http.StatusBadRequest, payload.ErrMissingParams(fmt.Errorf("missing a %v parameter", "[id]")))
		return
	}

	task, err := h.taskUsecase.GetTaskById(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, err)
			return
		}
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, payload.NewSimpleSuccessResponse(task))
}

// UpdateTaskById		godoc
// @Summary				Update a task by id
// @Description		Update  a task by id.
// @Param				id path string  true "id"
// @Param				task body payload.UpdateTaskRequest true "Update task request"
// @Produce				application/json
// @Tags					tasks
// @Success				200 {object} payload.AppResponse{}
// @Failure      		400  	{object} payload.ErrResponse{}
// @Failure 			404 	{object} payload.ErrResponse{}
// @Router				/tasks/{id} [put]
func (h TaskHandler) handleUpdateTaskById(c *gin.Context) {
	id, _ := strconv.ParseInt(removeSlashFromParam(c.Param("id")), 0, 64)
	if id == 0 {
		c.JSON(http.StatusBadRequest, payload.ErrMissingParams(fmt.Errorf("missing a %v parameter", "[id]")))
		return
	}

	var updateTaskRequest payload.UpdateTaskRequest
	if err := c.ShouldBindJSON(&updateTaskRequest); err != nil {
		c.JSON(http.StatusBadRequest, payload.ErrInvalidRequestBody(err))
		return
	}

	updateTaskRequest.Id = id
	taskUpdated, err := h.taskUsecase.UpdateTaskById(updateTaskRequest)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, err)
			return
		}
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, payload.NewSimpleSuccessResponse(taskUpdated))
}

// RemoveTaskById		godoc
// @Summary				Remove a task by id
// @Description		Remove a task by id in DB.
// @Param				id path string  true "id"
// @Tags					tasks
// @Success				200		{object} payload.AppResponse{}
// @Failure      		400  	{object} payload.ErrResponse{}
// @Failure 			404 	{object} payload.ErrResponse{}
// @Router				/tasks [delete]
func (h TaskHandler) handleRemoveTaskById(c *gin.Context) {
	id, err := strconv.ParseInt(removeSlashFromParam(c.Param("id")), 0, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, payload.ErrInvalidRequest(fmt.Errorf("invalid %v parameter", "[id]")))
		return
	}
	if id == 0 {
		c.JSON(http.StatusBadRequest, payload.ErrMissingParams(fmt.Errorf("invalid %v parameter", "[id]")))
		return
	}

	if err := h.taskUsecase.RemoveTaskById(id); err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, err)
			return
		}
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, payload.NewSimpleSuccessResponse(nil))
}

func removeSlashFromParam(param string) string {
	if strings.Contains(param, "/") {
		param = strings.Replace(param, "/", "", -1)
	}
	if strings.Contains(param, "\\") {
		param = strings.Replace(param, "\\", "", -1)
	}
	return param
}