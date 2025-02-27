package utils

import (
    "encoding/json"
    "fmt"
    "log"
    "os"
    "path/filepath"
    "subsmanager/config"
    "sync"
    "time"
)

// LogLevel 定义日志级别
type LogLevel string

const (
    INFO    LogLevel = "INFO"
    ERROR   LogLevel = "ERROR"
    SUCCESS LogLevel = "SUCCESS"
    WARNING LogLevel = "WARNING"
)

// LogEntry 定义日志条目结构
type LogEntry struct {
    Timestamp time.Time `json:"timestamp"`
    Level     LogLevel  `json:"level"`
    Message   string    `json:"message"`
    Details   string    `json:"details,omitempty"`
}

var (
    logger *log.Logger
    // 内存中保存最近的日志记录
    logEntries     []LogEntry
    logEntriesMux  sync.RWMutex
    maxLogEntries  = 1000 // 最多保存1000条日志
)

func init() {
    // 创建日志文件
    logPath := filepath.Join(config.GlobalConfig.Storage.Path, "subsmanager.log")
    logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
        log.Fatal(err)
    }

    // 设置日志格式
    logger = log.New(logFile, "", log.Ldate|log.Ltime)
    logEntries = make([]LogEntry, 0)
}

// addLogEntry 添加日志条目到内存
func addLogEntry(level LogLevel, message string, details string) {
    entry := LogEntry{
        Timestamp: time.Now(),
        Level:     level,
        Message:   message,
        Details:   details,
    }

    logEntriesMux.Lock()
    defer logEntriesMux.Unlock()

    // 如果超过最大数量，移除最旧的日志
    if len(logEntries) >= maxLogEntries {
        logEntries = logEntries[1:]
    }
    logEntries = append(logEntries, entry)

    // 写入文件
    logStr, _ := json.Marshal(entry)
    logger.Printf("[%s] %s", level, string(logStr))
}

// GetRecentLogs 获取最近的日志记录
func GetRecentLogs() []LogEntry {
    logEntriesMux.RLock()
    defer logEntriesMux.RUnlock()
    
    result := make([]LogEntry, len(logEntries))
    copy(result, logEntries)
    return result
}

// LogSubscriptionImport 记录订阅导入日志
func LogSubscriptionImport(subscriptionName string, nodeCount int) {
    msg := fmt.Sprintf("订阅导入完成")
    details := fmt.Sprintf("已导入订阅：%s，节点数：%d个", subscriptionName, nodeCount)
    addLogEntry(SUCCESS, msg, details)
}

// LogSubscriptionMerge 记录订阅整合日志
func LogSubscriptionMerge(count int) {
    msg := fmt.Sprintf("订阅整合完成")
    details := fmt.Sprintf("已整合%d条订阅", count)
    addLogEntry(SUCCESS, msg, details)
}

// LogSubscriptionDelete 记录订阅删除日志
func LogSubscriptionDelete(subscriptionName string) {
    msg := fmt.Sprintf("订阅删除完成")
    details := fmt.Sprintf("已删除订阅：%s", subscriptionName)
    addLogEntry(SUCCESS, msg, details)
}

// LogSpeedTest 记录节点测速日志
func LogSpeedTest(total, latencyTested, latencyDropped, speedTested int) {
    msg := fmt.Sprintf("测速完成")
    details := fmt.Sprintf("共计%d个节点，延迟测速%d个节点，延迟测速丢弃%d个节点，下载测速%d个节点",
        total, latencyTested, latencyDropped, speedTested)
    addLogEntry(SUCCESS, msg, details)
}

// LogNodeFilter 记录节点优选日志
func LogNodeFilter(latencyLimit int, speedLimit float64, validCount int) {
    msg := fmt.Sprintf("节点优选完成")
    details := fmt.Sprintf("筛选条件：延迟上限：%dms，速度下限：%.1fM/s，筛选完成，符合条件节点共计:%d个",
        latencyLimit, speedLimit, validCount)
    addLogEntry(SUCCESS, msg, details)
}

// LogSubscriptionGenerate 记录优选订阅生成日志
func LogSubscriptionGenerate(url string) {
    msg := fmt.Sprintf("优选订阅生成完成")
    details := fmt.Sprintf("订阅地址：%s", url)
    addLogEntry(SUCCESS, msg, details)
}

// LogError 记录错误日志
func LogError(format string, v ...interface{}) {
    msg := fmt.Sprintf(format, v...)
    addLogEntry(ERROR, msg, "")
}

// LogInfo 记录信息日志
func LogInfo(format string, v ...interface{}) {
    msg := fmt.Sprintf(format, v...)
    addLogEntry(INFO, msg, "")
}

// LogParseError 记录解析错误
func LogParseError(subscriptionName string, nodeType string, err error) {
    msg := fmt.Sprintf("订阅解析错误")
    details := fmt.Sprintf("订阅：%s，节点类型：%s，错误：%v", subscriptionName, nodeType, err)
    addLogEntry(ERROR, msg, details)
} 