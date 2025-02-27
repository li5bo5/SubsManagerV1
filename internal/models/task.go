package models

import "time"

// TaskType 定义任务类型
type TaskType string

const (
	TaskTypeSubscriptionUpdate TaskType = "subscription_update" // 订阅更新任务
	TaskTypeNodeTest          TaskType = "node_test"           // 节点测速任务
)

// TaskStatus 定义任务状态
type TaskStatus string

const (
	TaskStatusEnabled  TaskStatus = "enabled"  // 任务启用
	TaskStatusDisabled TaskStatus = "disabled" // 任务禁用
)

// Task 定义自动化任务
type Task struct {
	ID          string     `json:"id"`           // 任务ID
	Type        TaskType   `json:"type"`         // 任务类型
	Name        string     `json:"name"`         // 任务名称
	Status      TaskStatus `json:"status"`       // 任务状态
	Cron        string     `json:"cron"`         // Cron表达式
	LastRunTime time.Time  `json:"last_run_time"` // 上次运行时间
	CreateTime  time.Time  `json:"create_time"`   // 创建时间
	UpdateTime  time.Time  `json:"update_time"`   // 更新时间
}

// TaskResult 定义任务执行结果
type TaskResult struct {
	TaskID      string    `json:"task_id"`       // 任务ID
	StartTime   time.Time `json:"start_time"`    // 开始时间
	EndTime     time.Time `json:"end_time"`      // 结束时间
	Status      string    `json:"status"`        // 执行状态
	Message     string    `json:"message"`       // 执行信息
	Error       string    `json:"error"`         // 错误信息
} 