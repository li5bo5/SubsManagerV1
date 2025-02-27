package services

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

// LogLevel 定义日志级别
type LogLevel string

const (
	LogLevelInfo    LogLevel = "INFO"
	LogLevelError   LogLevel = "ERROR"
	LogLevelSuccess LogLevel = "SUCCESS"
	LogLevelWarning LogLevel = "WARNING"
)

// LogEntry 定义日志条目
type LogEntry struct {
	Timestamp time.Time         `json:"timestamp"`
	Level     LogLevel         `json:"level"`
	Message   string           `json:"message"`
	Details   string           `json:"details,omitempty"`
	CodeInfo  string           `json:"code_info,omitempty"`
	Data      json.RawMessage  `json:"data,omitempty"`
}

// LogService 日志服务
type LogService struct {
	entries []LogEntry
	mu      sync.RWMutex
	maxSize int
	logFile *os.File
}

// NewLogService 创建新的日志服务
func NewLogService(logPath string, maxSize int) (*LogService, error) {
	// 确保日志目录存在
	if err := os.MkdirAll(filepath.Dir(logPath), 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %v", err)
	}

	// 打开日志文件
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %v", err)
	}

	return &LogService{
		entries: make([]LogEntry, 0),
		maxSize: maxSize,
		logFile: logFile,
	}, nil
}

// getCodeInfo 获取代码位置信息
func getCodeInfo() string {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		return "unknown"
	}
	return fmt.Sprintf("%s:%d", filepath.Base(file), line)
}

// addEntry 添加日志条目
func (s *LogService) addEntry(level LogLevel, message string, details string, data interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 序列化额外数据
	var rawData json.RawMessage
	if data != nil {
		if jsonData, err := json.Marshal(data); err == nil {
			rawData = jsonData
		}
	}

	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   message,
		Details:   details,
		CodeInfo:  getCodeInfo(),
		Data:      rawData,
	}

	// 添加到内存
	s.entries = append(s.entries, entry)
	if len(s.entries) > s.maxSize {
		s.entries = s.entries[1:]
	}

	// 写入文件
	jsonEntry, _ := json.Marshal(entry)
	s.logFile.Write(append(jsonEntry, '\n'))
}

// GetLogs 获取所有日志
func (s *LogService) GetLogs() []LogEntry {
	s.mu.RLock()
	defer s.mu.RUnlock()

	logs := make([]LogEntry, len(s.entries))
	copy(logs, s.entries)
	return logs
}

// Close 关闭日志服务
func (s *LogService) Close() error {
	return s.logFile.Close()
}

// 日志记录方法
func (s *LogService) Info(message string, data ...interface{}) {
	s.addEntry(LogLevelInfo, message, "", data)
}

func (s *LogService) Error(message string, data ...interface{}) {
	s.addEntry(LogLevelError, message, "", data)
}

func (s *LogService) Success(message string, data ...interface{}) {
	s.addEntry(LogLevelSuccess, message, "", data)
}

func (s *LogService) Warning(message string, data ...interface{}) {
	s.addEntry(LogLevelWarning, message, "", data)
}

// 特定业务日志方法
func (s *LogService) LogSubscriptionImport(name string, nodeCount int) {
	s.addEntry(LogLevelSuccess, "订阅导入完成",
		fmt.Sprintf("已导入订阅：%s，节点数：%d个", name, nodeCount),
		map[string]interface{}{
			"subscription_name": name,
			"node_count":       nodeCount,
		})
}

func (s *LogService) LogSubscriptionMerge(count int) {
	s.addEntry(LogLevelSuccess, "订阅整合完成",
		fmt.Sprintf("已整合%d条订阅", count),
		map[string]interface{}{
			"merged_count": count,
		})
}

func (s *LogService) LogSubscriptionDelete(name string) {
	s.addEntry(LogLevelSuccess, "订阅删除完成",
		fmt.Sprintf("已删除订阅：%s", name),
		map[string]interface{}{
			"subscription_name": name,
		})
}

func (s *LogService) LogSpeedTest(total, latencyTested, latencyDropped, speedTested int) {
	s.addEntry(LogLevelSuccess, "节点测速完成",
		fmt.Sprintf("共计%d个节点，延迟测速%d个节点，延迟测速丢弃%d个节点，下载测速%d个节点",
			total, latencyTested, latencyDropped, speedTested),
		map[string]interface{}{
			"total_nodes":      total,
			"latency_tested":   latencyTested,
			"latency_dropped":  latencyDropped,
			"speed_tested":     speedTested,
		})
}

func (s *LogService) LogNodeFilter(latencyLimit int, speedLimit float64, validCount int) {
	s.addEntry(LogLevelSuccess, "节点优选完成",
		fmt.Sprintf("筛选条件：延迟上限：%dms，速度下限：%.1fM/s，筛选完成，符合条件节点共计:%d个",
			latencyLimit, speedLimit, validCount),
		map[string]interface{}{
			"latency_limit": latencyLimit,
			"speed_limit":   speedLimit,
			"valid_count":   validCount,
		})
}

func (s *LogService) LogSubscriptionGenerate(url string) {
	s.addEntry(LogLevelSuccess, "优选订阅生成完成",
		fmt.Sprintf("订阅地址：%s", url),
		map[string]interface{}{
			"subscription_url": url,
		})
}

func (s *LogService) LogParseError(subscriptionName, nodeType string, err error) {
	s.addEntry(LogLevelError, "订阅解析错误",
		fmt.Sprintf("订阅：%s，节点类型：%s，错误：%v", subscriptionName, nodeType, err),
		map[string]interface{}{
			"subscription_name": subscriptionName,
			"node_type":        nodeType,
			"error":            err.Error(),
		})
}

func (s *LogService) LogTaskExecution(taskID, taskName, taskType string, status string, duration time.Duration, err error) {
	data := map[string]interface{}{
		"task_id":   taskID,
		"task_name": taskName,
		"task_type": taskType,
		"duration":  duration.String(),
	}

	var level LogLevel
	var message string
	var details string

	switch status {
	case "success":
		level = LogLevelSuccess
		message = fmt.Sprintf("任务执行成功：%s", taskName)
		details = fmt.Sprintf("执行时长：%s", duration)
	case "failed":
		level = LogLevelError
		message = fmt.Sprintf("任务执行失败：%s", taskName)
		details = fmt.Sprintf("执行时长：%s，错误：%v", duration, err)
		data["error"] = err.Error()
	case "timeout":
		level = LogLevelWarning
		message = fmt.Sprintf("任务执行超时：%s", taskName)
		details = fmt.Sprintf("执行时长：%s", duration)
	}

	s.addEntry(level, message, details, data)
} 