package services

import (
	"fmt"
	"sync"
	"time"
	"context"

	"github.com/robfig/cron/v3"
	"subsmanager/internal/models"
)

// SchedulerService 定义调度器服务
type SchedulerService struct {
	cron          *cron.Cron
	tasks         map[string]*models.Task
	taskEntries   map[string]cron.EntryID
	subService    *SubscriptionService
	statusService *StatusService
	logService    *LogService
	mu            sync.RWMutex
	timeout       time.Duration // 任务超时时间
}

// NewSchedulerService 创建新的调度器服务
func NewSchedulerService(subService *SubscriptionService, statusService *StatusService, logService *LogService) *SchedulerService {
	return &SchedulerService{
		cron:          cron.New(cron.WithSeconds()),
		tasks:         make(map[string]*models.Task),
		taskEntries:   make(map[string]cron.EntryID),
		subService:    subService,
		statusService: statusService,
		logService:    logService,
		timeout:       10 * time.Minute, // 默认超时时间10分钟
	}
}

// Start 启动调度器
func (s *SchedulerService) Start() {
	s.cron.Start()
}

// Stop 停止调度器
func (s *SchedulerService) Stop() {
	s.cron.Stop()
}

// AddTask 添加新任务
func (s *SchedulerService) AddTask(task *models.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 检查任务是否已存在
	if _, exists := s.tasks[task.ID]; exists {
		return fmt.Errorf("task with ID %s already exists", task.ID)
	}

	// 根据任务类型创建对应的执行函数
	var jobFunc func()
	switch task.Type {
	case models.TaskTypeSubscriptionUpdate:
		jobFunc = s.createSubscriptionUpdateJob(task)
	case models.TaskTypeNodeTest:
		jobFunc = s.createNodeTestJob(task)
	default:
		return fmt.Errorf("unsupported task type: %s", task.Type)
	}

	// 添加定时任务
	entryID, err := s.cron.AddFunc(task.Cron, jobFunc)
	if err != nil {
		return fmt.Errorf("failed to add cron job: %v", err)
	}

	// 保存任务信息
	s.tasks[task.ID] = task
	s.taskEntries[task.ID] = entryID

	return nil
}

// RemoveTask 移除任务
func (s *SchedulerService) RemoveTask(taskID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entryID, exists := s.taskEntries[taskID]
	if !exists {
		return fmt.Errorf("task with ID %s not found", taskID)
	}

	s.cron.Remove(entryID)
	delete(s.tasks, taskID)
	delete(s.taskEntries, taskID)

	return nil
}

// UpdateTask 更新任务
func (s *SchedulerService) UpdateTask(task *models.Task) error {
	if err := s.RemoveTask(task.ID); err != nil {
		return err
	}
	return s.AddTask(task)
}

// GetTask 获取任务信息
func (s *SchedulerService) GetTask(taskID string) (*models.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, exists := s.tasks[taskID]
	if !exists {
		return nil, fmt.Errorf("task with ID %s not found", taskID)
	}
	return task, nil
}

// ListTasks 列出所有任务
func (s *SchedulerService) ListTasks() []*models.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tasks := make([]*models.Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}
	return tasks
}

// executeWithTimeout 带超时控制的任务执行
func (s *SchedulerService) executeWithTimeout(task *models.Task, operation func() error) *models.TaskResult {
	result := &models.TaskResult{
		TaskID:    task.ID,
		StartTime: time.Now(),
		Status:    "running",
	}

	// 创建带超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	// 创建错误通道
	done := make(chan error)

	// 在goroutine中执行任务
	go func() {
		done <- operation()
	}()

	// 等待任务完成或超时
	select {
	case err := <-done:
		result.EndTime = time.Now()
		if err != nil {
			result.Status = "failed"
			result.Error = err.Error()
			// 记录失败日志
			s.logService.Error("Task execution failed",
				"taskId", task.ID,
				"taskName", task.Name,
				"taskType", task.Type,
				"error", err.Error(),
				"duration", result.EndTime.Sub(result.StartTime).String())
		} else {
			result.Status = "success"
			result.Message = fmt.Sprintf("Successfully executed task: %s", task.Name)
			// 记录成功日志
			s.logService.Info("Task execution succeeded",
				"taskId", task.ID,
				"taskName", task.Name,
				"taskType", task.Type,
				"duration", result.EndTime.Sub(result.StartTime).String())
		}
	case <-ctx.Done():
		result.EndTime = time.Now()
		result.Status = "timeout"
		result.Error = "task execution timed out"
		// 记录超时日志
		s.logService.Warning("Task execution timed out",
			"taskId", task.ID,
			"taskName", task.Name,
			"taskType", task.Type,
			"timeout", s.timeout.String(),
			"duration", result.EndTime.Sub(result.StartTime).String())
	}

	return result
}

// createSubscriptionUpdateJob 创建订阅更新任务
func (s *SchedulerService) createSubscriptionUpdateJob(task *models.Task) func() {
	return func() {
		// 记录任务开始日志
		s.logService.Info("Starting subscription update task",
			"taskId", task.ID,
			"taskName", task.Name)

		// 执行任务并记录结果
		result := s.executeWithTimeout(task, s.subService.UpdateAllSubscriptions)

		// 更新任务状态
		task.LastRunTime = result.StartTime
		task.UpdateTime = result.EndTime

		// 记录执行结果
		s.statusService.AddTaskHistory(result)
	}
}

// createNodeTestJob 创建节点测试任务
func (s *SchedulerService) createNodeTestJob(task *models.Task) func() {
	return func() {
		// 记录任务开始日志
		s.logService.Info("Starting node test task",
			"taskId", task.ID,
			"taskName", task.Name)

		// 执行任务并记录结果
		result := s.executeWithTimeout(task, s.subService.TestAllNodes)

		// 更新任务状态
		task.LastRunTime = result.StartTime
		task.UpdateTime = result.EndTime

		// 记录执行结果
		s.statusService.AddTaskHistory(result)
	}
}

// SetTimeout 设置任务超时时间
func (s *SchedulerService) SetTimeout(timeout time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.timeout = timeout
} 