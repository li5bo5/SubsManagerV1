package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"SubsManagerV1/internal/models"
	"SubsManagerV1/internal/services"
)

// TaskHandler 处理任务相关的请求
type TaskHandler struct {
	schedulerService *services.SchedulerService
}

// NewTaskHandler 创建新的任务处理器
func NewTaskHandler(schedulerService *services.SchedulerService) *TaskHandler {
	return &TaskHandler{
		schedulerService: schedulerService,
	}
}

// CreateTask 创建新任务
func (h *TaskHandler) CreateTask(c *gin.Context) {
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// 生成任务ID
	task.ID = uuid.New().String()

	// 添加任务
	if err := h.schedulerService.AddTask(&task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, task)
}

// UpdateTask 更新任务
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	taskID := c.Param("id")
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	task.ID = taskID
	if err := h.schedulerService.UpdateTask(&task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, task)
}

// DeleteTask 删除任务
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	taskID := c.Param("id")
	if err := h.schedulerService.RemoveTask(taskID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetTask 获取任务信息
func (h *TaskHandler) GetTask(c *gin.Context) {
	taskID := c.Param("id")
	task, err := h.schedulerService.GetTask(taskID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, task)
}

// ListTasks 列出所有任务
func (h *TaskHandler) ListTasks(c *gin.Context) {
	tasks := h.schedulerService.ListTasks()
	c.JSON(http.StatusOK, tasks)
} 