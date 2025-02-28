package models

import "time"

// SubscriptionFileHistory 订阅文件历史记录
type SubscriptionFileHistory struct {
    FileName    string    `json:"file_name"`    // 文件名
    LocalURL    string    `json:"local_url"`    // 本地访问地址
    GenerateTime time.Time `json:"generate_time"` // 生成时间
    NodeCount   int       `json:"node_count"`   // 节点数量
}

// NodeStatus 节点状态统计
type NodeStatus struct {
    TotalNodes     int `json:"total_nodes"`      // 总节点数
    CurrentNodes   int `json:"current_nodes"`    // 当前筛选后的节点数
    SlowNodes      int `json:"slow_nodes"`       // 慢速节点数
    FaultNodes     int `json:"fault_nodes"`      // 故障节点数（延迟>400ms）
}

// SystemStatus 系统状态
type SystemStatus struct {
    NodeStatus         NodeStatus                 `json:"node_status"`          // 节点状态
    LatestSubFile     string                     `json:"latest_sub_file"`      // 最新的sub文件本地地址
    LatestInputFile   string                     `json:"latest_input_file"`    // 最新的Sub-Input文件本地地址
    SubHistory        []SubscriptionFileHistory  `json:"sub_history"`          // 近三次的订阅历史记录
    LastUpdateTime    time.Time                  `json:"last_update_time"`     // 最后更新时间
    LastTaskTime     time.Time                  `json:"last_task_time"`      // 最后任务执行时间
    LastTaskName     string                     `json:"last_task_name"`      // 最后执行的任务名称
    LastTaskStatus   string                     `json:"last_task_status"`    // 最后任务执行状态
} 