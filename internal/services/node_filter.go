package services

import (
    "fmt"
    "time"
    "subsmanager/internal/models"
    "subsmanager/internal/utils"
    "os"
    "gopkg.in/yaml.v3"
)

// FilterConfig 节点筛选配置
type FilterConfig struct {
    MaxLatency    int     `json:"max_latency"`    // 延迟上限(ms)
    MinSpeed      float64 `json:"min_speed"`      // 速度下限(MB/s)
}

// FilterResult 筛选结果
type FilterResult struct {
    Nodes       []*models.Node    `json:"nodes"`        // 筛选后的节点列表
    TotalNodes  int              `json:"total_nodes"`  // 筛选后的节点总数
    FilterTime  time.Time        `json:"filter_time"`  // 筛选时间
}

// NodeFilter 节点筛选器
type NodeFilter struct {
    config  *FilterConfig
}

// generateFileName 生成订阅文件名
func (nf *NodeFilter) generateFileName() string {
    // 生成格式为 Sub-Input-MM-DD-HH-mm.yaml 的文件名
    return fmt.Sprintf("Sub-Input-%s.yaml", time.Now().Format("01-02-15-04"))
}

// NewNodeFilter 创建节点筛选器
func NewNodeFilter(config *FilterConfig) *NodeFilter {
    return &NodeFilter{
        config: config,
    }
}

// FilterNodes 筛选节点
func (nf *NodeFilter) FilterNodes(nodes []*models.Node) (*FilterResult, error) {
    result := &FilterResult{
        FilterTime: time.Now(),
    }

    utils.LogInfo("开始节点筛选")

    validNodes := make([]*models.Node, 0)
    for _, node := range nodes {
        // 检查节点是否已测试
        if node.LastTestedAt.IsZero() {
            continue
        }

        // 检查延迟
        if node.Latency > nf.config.MaxLatency {
            continue
        }

        // 检查速度
        if node.DownloadSpeed < nf.config.MinSpeed {
            continue
        }

        validNodes = append(validNodes, node)
    }

    result.Nodes = validNodes
    result.TotalNodes = len(validNodes)

    utils.LogNodeFilter(nf.config.MaxLatency, nf.config.MinSpeed, len(validNodes))

    return result, nil
}

// GenerateSubscription 生成订阅
func (nf *NodeFilter) GenerateSubscription(nodes []*models.Node) (*models.Subscription, error) {
    sub := &models.Subscription{
        ID:        fmt.Sprintf("filtered-%s", time.Now().Format("20060102150405")),
        Name:      fmt.Sprintf("优选节点-%s", time.Now().Format("01-02 15:04")),
        Type:      "mixed",
        NodeCount: len(nodes),
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }

    return sub, nil
}

// SaveSubscription 保存订阅到文件
func (nf *NodeFilter) SaveSubscription(sub *models.Subscription, filePath string) error {
    content, err := yaml.Marshal(sub)
    if err != nil {
        utils.LogError("订阅序列化失败: %v", err)
        return fmt.Errorf("订阅序列化失败: %v", err)
    }

    err = os.WriteFile(filePath, content, 0644)
    if err != nil {
        utils.LogError("保存订阅文件失败: %v", err)
        return fmt.Errorf("保存订阅文件失败: %v", err)
    }

    utils.LogInfo("订阅文件保存成功: %s", filePath)
    return nil
} 