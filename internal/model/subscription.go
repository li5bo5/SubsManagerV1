package model

import (
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"
    "subsmanager/config"
    "subsmanager/internal/utils"
    "sync"
    "time"
)

// Subscription 订阅信息结构
type Subscription struct {
    ID          string    `json:"id"`           // 订阅唯一标识
    Name        string    `json:"name"`         // 订阅名称
    URL         string    `json:"url"`          // 订阅链接
    NodeCount   int       `json:"nodeCount"`    // 节点数量
    Type        string    `json:"type"`         // 订阅类型（base64/yaml/json）
    UpdateTime  time.Time `json:"updateTime"`   // 更新时间
    Nodes       []Node    `json:"nodes"`        // 节点列表
}

// Node 节点信息结构
type Node struct {
    Name       string                 `json:"name"`       // 节点名称
    Type       string                 `json:"type"`       // 节点类型（vmess/ss/hy2/trojan）
    Server     string                 `json:"server"`     // 服务器地址
    Port       int                    `json:"port"`       // 端口
    Group      string                 `json:"group"`      // 分组信息
    Config     map[string]interface{} `json:"config"`     // 节点配置
}

var (
    subscriptions     []Subscription
    subscriptionMutex sync.RWMutex
    subscriptionFile  string
)

func init() {
    // 初始化订阅文件路径
    subscriptionFile = filepath.Join(config.GlobalConfig.Storage.Path, "subscriptions.json")
    // 加载已有订阅
    loadSubscriptions()
}

// loadSubscriptions 从文件加载订阅信息
func loadSubscriptions() {
    subscriptionMutex.Lock()
    defer subscriptionMutex.Unlock()

    // 如果文件不存在，创建空的订阅列表
    if _, err := os.Stat(subscriptionFile); os.IsNotExist(err) {
        subscriptions = make([]Subscription, 0)
        return
    }

    // 读取文件内容
    data, err := os.ReadFile(subscriptionFile)
    if err != nil {
        utils.LogError("Failed to read subscriptions file: %v", err)
        subscriptions = make([]Subscription, 0)
        return
    }

    // 解析JSON
    if err := json.Unmarshal(data, &subscriptions); err != nil {
        utils.LogError("Failed to parse subscriptions file: %v", err)
        subscriptions = make([]Subscription, 0)
        return
    }
}

// saveSubscriptions 保存订阅信息到文件
func saveSubscriptions() error {
    subscriptionMutex.RLock()
    defer subscriptionMutex.RUnlock()

    data, err := json.MarshalIndent(subscriptions, "", "  ")
    if err != nil {
        return err
    }

    return os.WriteFile(subscriptionFile, data, 0644)
}

// GetSubscriptions 获取所有订阅
func GetSubscriptions() []Subscription {
    subscriptionMutex.RLock()
    defer subscriptionMutex.RUnlock()

    result := make([]Subscription, len(subscriptions))
    copy(result, subscriptions)
    return result
}

// AddSubscription 添加新订阅
func AddSubscription(sub Subscription) error {
    subscriptionMutex.Lock()
    defer subscriptionMutex.Unlock()

    subscriptions = append(subscriptions, sub)
    return saveSubscriptions()
}

// DeleteSubscription 删除订阅
func DeleteSubscription(id string) error {
    subscriptionMutex.Lock()
    defer subscriptionMutex.Unlock()

    for i, sub := range subscriptions {
        if sub.ID == id {
            // 从切片中删除元素
            subscriptions = append(subscriptions[:i], subscriptions[i+1:]...)
            return saveSubscriptions()
        }
    }
    return nil
}

// MergeSubscriptions 整合所有订阅节点
func MergeSubscriptions() ([]Node, error) {
    subscriptionMutex.RLock()
    defer subscriptionMutex.RUnlock()

    var mergedNodes []Node
    nodeMap := make(map[string]bool) // 用于去重

    for _, sub := range subscriptions {
        for _, node := range sub.Nodes {
            // 使用节点类型、服务器和端口作为唯一标识
            key := fmt.Sprintf("%s-%s-%d", node.Type, node.Server, node.Port)
            if !nodeMap[key] {
                mergedNodes = append(mergedNodes, node)
                nodeMap[key] = true
            }
        }
    }

    return mergedNodes, nil
} 