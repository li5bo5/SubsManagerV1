package services

import (
    "fmt"
    "time"
)

// FilterConfig 节点筛选配置
type FilterConfig struct {
    MaxLatency    int     `json:"max_latency"`    // 延迟上限(ms)
    MinSpeed      float64 `json:"min_speed"`      // 速度下限(MB/s)
}

// FilterResult 筛选结果
type FilterResult struct {
    Nodes       []*Node    `json:"nodes"`        // 筛选后的节点列表
    TotalNodes  int        `json:"total_nodes"`  // 筛选后的节点总数
    FilterTime  time.Time  `json:"filter_time"`  // 筛选时间
}

// NodeFilter 节点筛选器
type NodeFilter struct {
    config  *FilterConfig
    logger  *Logger
}

// generateFileName 生成订阅文件名
func (nf *NodeFilter) generateFileName() string {
    // 生成格式为 Sub-Input-MM-DD-HH-mm.yaml 的文件名
    return fmt.Sprintf("Sub-Input-%s.yaml", time.Now().Format("01-02-15-04"))
}

// NewNodeFilter 创建节点筛选器
func NewNodeFilter(config *FilterConfig, logger *Logger) *NodeFilter {
    return &NodeFilter{
        config: config,
        logger: logger,
    }
}

// FilterNodes 筛选节点
func (nf *NodeFilter) FilterNodes(nodes []*Node) (*FilterResult, error) {
    result := &FilterResult{
        FilterTime: time.Now(),
    }

    // 记录开始筛选
    nf.logger.Info("开始节点筛选", map[string]interface{}{
        "max_latency": nf.config.MaxLatency,
        "min_speed":   nf.config.MinSpeed,
        "total_nodes": len(nodes),
    })

    // 筛选符合条件的节点
    for _, node := range nodes {
        // 跳过未测速的节点
        if node.TestTime.IsZero() {
            continue
        }

        // 检查延迟和速度是否满足条件
        if node.Latency <= nf.config.MaxLatency && node.Speed >= nf.config.MinSpeed {
            result.Nodes = append(result.Nodes, node)
        }
    }

    result.TotalNodes = len(result.Nodes)

    // 记录筛选结果
    nf.logger.Info("节点筛选完成", map[string]interface{}{
        "filtered_nodes": result.TotalNodes,
        "total_nodes":   len(nodes),
    })

    return result, nil
}

// GenerateSubscription 生成订阅文件
func (nf *NodeFilter) GenerateSubscription(nodes []*Node) (*Subscription, error) {
    if len(nodes) == 0 {
        return nil, fmt.Errorf("没有可用的节点")
    }

    // 创建新的订阅
    sub := &Subscription{
        ID:        fmt.Sprintf("sub_%d", time.Now().Unix()),
        Name:      "优选节点订阅",
        Type:      "yaml", // 默认使用yaml格式
        UpdatedAt: time.Now(),
        Nodes:     nodes,
    }

    // 记录订阅生成
    nf.logger.Info("生成订阅文件", map[string]interface{}{
        "node_count": len(nodes),
        "sub_type":   sub.Type,
    })

    // 添加订阅历史记录
    if err := DefaultSubscriptionService.AddSubscriptionHistory(
        sub.ID,
        models.ActionGenerate,
        len(nodes),
        fmt.Sprintf("生成优选节点订阅，共%d个节点", len(nodes)),
    ); err != nil {
        nf.logger.Error("添加订阅历史记录失败", map[string]interface{}{
            "error": err.Error(),
        })
        // 不中断流程，继续返回订阅
    }

    return sub, nil
}

// SaveSubscription 保存订阅到文件
func (nf *NodeFilter) SaveSubscription(sub *Subscription, filePath string) error {
    // 将订阅内容序列化为YAML格式
    content, err := sub.ToYAML()
    if err != nil {
        nf.logger.Error("订阅序列化失败", map[string]interface{}{
            "error": err.Error(),
        })
        return fmt.Errorf("订阅序列化失败: %v", err)
    }

    // 写入文件
    err = WriteFile(filePath, content)
    if err != nil {
        nf.logger.Error("保存订阅文件失败", map[string]interface{}{
            "error": err.Error(),
            "path":  filePath,
        })
        return fmt.Errorf("保存订阅文件失败: %v", err)
    }

    nf.logger.Info("订阅文件保存成功", map[string]interface{}{
        "path": filePath,
    })

    return nil
} 