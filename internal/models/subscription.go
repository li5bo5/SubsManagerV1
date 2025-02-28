package models

import "time"

// Subscription 订阅信息
type Subscription struct {
    ID        string    `json:"id"`
    Name      string    `json:"name"`
    Type      string    `json:"type"`
    URL       string    `json:"url"`
    NodeCount int       `json:"node_count"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// Node 节点信息
type Node struct {
    ID              string  `json:"id"`
    Type            string  `json:"type"`
    Alias           string  `json:"alias"`
    Address         string  `json:"address"`
    Port            int     `json:"port"`
    Protocol        string  `json:"protocol"`
    SubscriptionID  string  `json:"subscription_id"`
    Group           string  `json:"group"`
    Latency         int     `json:"latency"`         // 延迟(ms)
    DownloadSpeed   float64 `json:"download_speed"`  // 下载速度(MB/s)
    LastTestedAt    time.Time `json:"last_tested_at"`
}

// TestResult 节点测试结果
type TestResult struct {
    NodeID        string    `json:"node_id"`
    Latency      int       `json:"latency"`       // 延迟(ms)
    DownloadSpeed float64  `json:"download_speed"` // 下载速度(MB/s)
    TestedAt     time.Time `json:"tested_at"`
}

// FilterCondition 节点筛选条件
type FilterCondition struct {
    MaxLatency       int     `json:"max_latency"`        // 最大延迟(ms)
    MinDownloadSpeed float64 `json:"min_download_speed"` // 最小下载速度(MB/s)
}

// MergeResult 订阅整合结果
type MergeResult struct {
    Timestamp  time.Time `json:"timestamp"`  // 整合时间
    FileURL    string    `json:"file_url"`   // 订阅文件URL
    NodeCount  int       `json:"node_count"` // 节点总数
}

// ImportResult 节点导入结果
type ImportResult struct {
    TotalCount     int       `json:"total_count"`      // 总节点数
    ImportedCount  int       `json:"imported_count"`   // 成功导入数
    DuplicateCount int       `json:"duplicate_count"`  // 重复节点数
    Nodes          []*Node   `json:"nodes"`            // 导入的节点列表
}

// NodeList 节点列表（支持分页）
type NodeList struct {
    Total    int     `json:"total"`     // 总节点数
    Page     int     `json:"page"`      // 当前页码
    PageSize int     `json:"page_size"` // 每页大小
    Nodes    []*Node `json:"nodes"`     // 节点列表
}

// NodeListQuery 节点列表查询参数
type NodeListQuery struct {
    Page     int    `form:"page" binding:"required,min=1"`
    PageSize int    `form:"page_size" binding:"required,min=10,max=100"`
    Type     string `form:"type"`       // 节点类型筛选
}

// SpeedTestResult 节点测速结果
type SpeedTestResult struct {
    TotalCount      int       `json:"total_count"`       // 总节点数
    LatencyTested   int       `json:"latency_tested"`    // 延迟测速节点数
    LatencyDropped  int       `json:"latency_dropped"`   // 延迟测速丢弃数
    SpeedTested     int       `json:"speed_tested"`      // 下载测速节点数
    Progress        float64   `json:"progress"`          // 测速进度(0-100)
    TestedNodes     []*Node   `json:"tested_nodes"`      // 已测速节点
}

// SpeedTestConfig 测速配置
type SpeedTestConfig struct {
    MaxLatency      int     `json:"max_latency"`       // 延迟阈值(ms)
    TestURL         string  `json:"test_url"`          // 下载测试URL
    Timeout         int     `json:"timeout"`           // 超时时间(秒)
    Concurrent      int     `json:"concurrent"`        // 并发数
}

// SpeedTestStats 测速统计
type SpeedTestStats struct {
    StartTime       time.Time `json:"start_time"`      // 开始时间
    EndTime         time.Time `json:"end_time"`        // 结束时间
    TotalNodes      int       `json:"total_nodes"`     // 总节点数
    TestedNodes     int       `json:"tested_nodes"`    // 已测试节点数
    SuccessNodes    int       `json:"success_nodes"`   // 测试成功节点数
    FailedNodes     int       `json:"failed_nodes"`    // 测试失败节点数
}

// SubscriptionHistory 订阅历史记录
type SubscriptionHistory struct {
    ID             string    `json:"id"`              // 历史记录ID
    SubscriptionID string    `json:"subscription_id"` // 订阅ID
    Action         string    `json:"action"`          // 操作类型
    NodeCount      int       `json:"node_count"`      // 节点数量
    CreatedAt      time.Time `json:"created_at"`      // 创建时间
    Details        string    `json:"details"`         // 详细信息
}

// SubscriptionAction 订阅操作类型
const (
    ActionImport   = "import"   // 导入订阅
    ActionUpdate   = "update"   // 更新订阅
    ActionDelete   = "delete"   // 删除订阅
    ActionGenerate = "generate" // 生成订阅
) 