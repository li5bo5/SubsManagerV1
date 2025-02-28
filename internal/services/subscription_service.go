package services

import (
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"
    "subsmanager/config"
    "subsmanager/internal/models"
    "subsmanager/internal/utils"
    "time"
    "gopkg.in/yaml.v3"
    "math/rand"
)

type SubscriptionService struct {
    subscriptions map[string]*models.Subscription
    nodes        map[string]*models.Node
    history      map[string]*models.SubscriptionHistory // 订阅历史记录
}

var DefaultSubscriptionService = &SubscriptionService{
    subscriptions: make(map[string]*models.Subscription),
    nodes:        make(map[string]*models.Node),
    history:      make(map[string]*models.SubscriptionHistory),
}

// ImportSubscription 导入订阅
func (s *SubscriptionService) ImportSubscription(name, url string) (*models.Subscription, error) {
    // 解析订阅
    result, err := utils.ParseSubscription(url)
    if err != nil {
        return nil, fmt.Errorf("parse subscription failed: %v", err)
    }

    // 创建订阅记录
    sub := &models.Subscription{
        ID:        fmt.Sprintf("sub_%d", time.Now().Unix()),
        Name:      name,
        Type:      string(result.Type),
        URL:       url,
        NodeCount: result.NodeCount,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
    
    // 保存订阅
    s.subscriptions[sub.ID] = sub

    // 保存节点
    for _, node := range result.Nodes {
        node.ID = fmt.Sprintf("node_%d", time.Now().UnixNano())
        node.SubscriptionID = sub.ID
        s.nodes[node.ID] = node
    }

    // 记录解析统计信息
    utils.LogInfo("Subscription imported: %s (ID: %s), Stats: Total=%d, Success=%d, Failed=%v",
        name, sub.ID, result.Stats.Total, result.Stats.Success, result.Stats.Failed)

    // 保存到文件
    if err := s.SaveToFile(); err != nil {
        return nil, fmt.Errorf("save to file failed: %v", err)
    }

    return sub, nil
}

// DeleteSubscription 删除订阅
func (s *SubscriptionService) DeleteSubscription(id string) error {
    if _, exists := s.subscriptions[id]; !exists {
        return fmt.Errorf("subscription not found: %s", id)
    }
    delete(s.subscriptions, id)
    return nil
}

// GetSubscriptions 获取所有订阅
func (s *SubscriptionService) GetSubscriptions() []*models.Subscription {
    subs := make([]*models.Subscription, 0, len(s.subscriptions))
    for _, sub := range s.subscriptions {
        subs = append(subs, sub)
    }
    return subs
}

// MergeSubscriptions 合并订阅
func (s *SubscriptionService) MergeSubscriptions(ids []string) (*models.MergeResult, error) {
    utils.LogInfo("开始合并订阅，订阅数量：%d", len(ids))

    // 验证所有订阅ID是否存在
    nodes := make([]*models.Node, 0)
    for _, id := range ids {
        if _, exists := s.subscriptions[id]; !exists {
            return nil, fmt.Errorf("subscription not found: %s", id)
        }

        // 收集该订阅下的所有节点
        for _, node := range s.nodes {
            if node.SubscriptionID == id {
                nodes = append(nodes, node)
            }
        }
    }

    // 创建合并结果
    result := &models.MergeResult{
        Timestamp: time.Now(),
        NodeCount: len(nodes),
    }

    // 生成合并后的文件名
    fileName := fmt.Sprintf("merged_%s.yaml", time.Now().Format("20060102150405"))
    filePath := filepath.Join(config.GlobalConfig.Storage.Path, fileName)
    
    // 保存合并后的节点到文件
    if err := s.saveNodesToFile(nodes, filePath); err != nil {
        return nil, fmt.Errorf("save merged nodes failed: %v", err)
    }

    // 设置文件访问URL
    result.FileURL = fmt.Sprintf("/subscriptions/%s", fileName)

    // 记录合并日志
    utils.LogSubscriptionMerge(len(ids))

    return result, nil
}

// saveNodesToFile 保存节点到文件
func (s *SubscriptionService) saveNodesToFile(nodes []*models.Node, filePath string) error {
    // 创建YAML格式的节点列表
    nodeList := struct {
        Proxies []map[string]interface{} `yaml:"proxies"`
    }{
        Proxies: make([]map[string]interface{}, 0, len(nodes)),
    }

    // 转换节点格式
    for _, node := range nodes {
        proxy := map[string]interface{}{
            "name":   node.Alias,
            "type":   node.Type,
            "server": node.Address,
            "port":   node.Port,
        }
        nodeList.Proxies = append(nodeList.Proxies, proxy)
    }

    // 序列化为YAML
    data, err := yaml.Marshal(nodeList)
    if err != nil {
        return fmt.Errorf("marshal nodes failed: %v", err)
    }

    // 写入文件
    return os.WriteFile(filePath, data, 0644)
}

// TestNodes 测试节点
func (s *SubscriptionService) TestNodes(config models.SpeedTestConfig) (*models.SpeedTestResult, error) {
    // 初始化测试结果
    result := &models.SpeedTestResult{
        TotalCount:  len(s.nodes),
        TestedNodes: make([]*models.Node, 0),
    }

    // 创建工作池
    type workItem struct {
        node *models.Node
        err  error
    }
    jobs := make(chan *models.Node, result.TotalCount)
    results := make(chan workItem, result.TotalCount)

    // 启动工作协程
    for i := 0; i < config.Concurrent; i++ {
        go func() {
            for node := range jobs {
                // 测试延迟
                latency, err := s.testNodeLatency(node, config.Timeout)
                if err != nil {
                    results <- workItem{node: node, err: err}
                    continue
                }

                // 更新节点延迟
                node.Latency = latency
                result.LatencyTested++

                // 如果延迟超过阈值，跳过下载测速
                if latency > config.MaxLatency {
                    result.LatencyDropped++
                    results <- workItem{node: node, err: nil}
                    continue
                }

                // 测试下载速度
                speed, err := s.testNodeSpeed(node, config.TestURL, config.Timeout)
                if err != nil {
                    results <- workItem{node: node, err: err}
                    continue
                }

                // 更新节点下载速度
                node.DownloadSpeed = speed
                result.SpeedTested++
                
                // 更新最后测试时间
                node.LastTestedAt = time.Now()

                results <- workItem{node: node, err: nil}
            }
        }()
    }

    // 发送任务
    go func() {
        for _, node := range s.nodes {
            jobs <- node
        }
        close(jobs)
    }()

    // 收集结果
    for i := 0; i < result.TotalCount; i++ {
        work := <-results
        if work.err != nil {
            utils.LogError("Test node failed: %v", work.err)
            continue
        }
        result.TestedNodes = append(result.TestedNodes, work.node)
        result.Progress = float64(i+1) / float64(result.TotalCount) * 100
    }

    // 记录测速日志
    utils.LogSpeedTest(
        result.TotalCount,
        result.LatencyTested,
        result.LatencyDropped,
        result.SpeedTested,
    )

    // 保存更新后的节点信息
    if err := s.SaveToFile(); err != nil {
        return nil, fmt.Errorf("save test results failed: %v", err)
    }

    return result, nil
}

// testNodeLatency 测试节点延迟
func (s *SubscriptionService) testNodeLatency(node *models.Node, timeout int) (int, error) {
    // TODO: 实现实际的延迟测试逻辑
    // 这里先使用模拟数据
    time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
    return rand.Intn(1000), nil
}

// testNodeSpeed 测试节点下载速度
func (s *SubscriptionService) testNodeSpeed(node *models.Node, testURL string, timeout int) (float64, error) {
    // TODO: 实现实际的下载速度测试逻辑
    // 这里先使用模拟数据
    time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
    return rand.Float64() * 10, nil
}

// FilterNodes 筛选节点
func (s *SubscriptionService) FilterNodes(condition models.FilterCondition) ([]*models.Node, error) {
    // TODO: 实现节点筛选逻辑
    return nil, nil
}

// GenerateSubscription 生成订阅文件
func (s *SubscriptionService) GenerateSubscription(nodes []*models.Node) (string, error) {
    // TODO: 实现生成订阅文件逻辑
    return "", nil
}

// AddSubscriptionHistory 添加订阅历史记录
func (s *SubscriptionService) AddSubscriptionHistory(subscriptionID string, action string, nodeCount int, details string) error {
    history := &models.SubscriptionHistory{
        ID:             fmt.Sprintf("hist_%d", time.Now().UnixNano()),
        SubscriptionID: subscriptionID,
        Action:         action,
        NodeCount:      nodeCount,
        CreatedAt:      time.Now(),
        Details:        details,
    }
    
    s.history[history.ID] = history

    // 记录日志
    utils.LogInfo("Added subscription history: Action=%s, SubscriptionID=%s, NodeCount=%d",
        action, subscriptionID, nodeCount)

    // 保存到文件
    return s.SaveToFile()
}

// SaveToFile 保存数据到文件
func (s *SubscriptionService) SaveToFile() error {
    data := struct {
        Subscriptions map[string]*models.Subscription        `json:"subscriptions"`
        Nodes        map[string]*models.Node                `json:"nodes"`
        History      map[string]*models.SubscriptionHistory `json:"history"`
    }{
        Subscriptions: s.subscriptions,
        Nodes:        s.nodes,
        History:      s.history,
    }
    
    jsonData, err := json.MarshalIndent(data, "", "    ")
    if err != nil {
        return err
    }
    
    dataPath := filepath.Join(config.GlobalConfig.Storage.Path, "data.json")
    return os.WriteFile(dataPath, jsonData, 0644)
}

// LoadFromFile 从文件加载数据
func (s *SubscriptionService) LoadFromFile() error {
    dataPath := filepath.Join(config.GlobalConfig.Storage.Path, "data.json")
    data, err := os.ReadFile(dataPath)
    if err != nil {
        if os.IsNotExist(err) {
            return nil
        }
        return err
    }
    
    var stored struct {
        Subscriptions map[string]*models.Subscription        `json:"subscriptions"`
        Nodes        map[string]*models.Node                `json:"nodes"`
        History      map[string]*models.SubscriptionHistory `json:"history"`
    }
    
    if err := json.Unmarshal(data, &stored); err != nil {
        return err
    }
    
    s.subscriptions = stored.Subscriptions
    s.nodes = stored.Nodes
    s.history = stored.History
    return nil
}

// ImportNodesFromFile 从YAML文件导入节点
func (s *SubscriptionService) ImportNodesFromFile(filePath string) (*models.ImportResult, error) {
    // 读取YAML文件
    data, err := os.ReadFile(filePath)
    if err != nil {
        return nil, fmt.Errorf("read file failed: %v", err)
    }

    // 解析YAML
    var config struct {
        Proxies []map[string]interface{} `yaml:"proxies"`
    }
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("parse yaml failed: %v", err)
    }

    result := &models.ImportResult{
        TotalCount: len(config.Proxies),
        Nodes:      make([]*models.Node, 0),
    }

    // 节点去重map
    nodeMap := make(map[string]bool)

    // 解析节点
    for _, proxy := range config.Proxies {
        nodeType, _ := proxy["type"].(string)
        server, _ := proxy["server"].(string)
        port, _ := proxy["port"].(int)
        name, _ := proxy["name"].(string)

        // 生成唯一标识
        key := fmt.Sprintf("%s-%s-%d", nodeType, server, port)
        if nodeMap[key] {
            result.DuplicateCount++
            continue
        }

        // 创建节点
        node := &models.Node{
            ID:        fmt.Sprintf("node_%d", time.Now().UnixNano()),
            Type:      nodeType,
            Alias:     name,
            Address:   server,
            Port:      port,
            Protocol:  s.getNodeProtocol(proxy),
            Group:     "imported",
        }

        nodeMap[key] = true
        result.Nodes = append(result.Nodes, node)
        result.ImportedCount++

        // 保存到内存
        s.nodes[node.ID] = node
    }

    // 保存到文件
    if err := s.SaveToFile(); err != nil {
        return nil, fmt.Errorf("save to file failed: %v", err)
    }

    return result, nil
}

// GetNodeList 获取节点列表
func (s *SubscriptionService) GetNodeList(query models.NodeListQuery) (*models.NodeList, error) {
    // 过滤节点
    var filteredNodes []*models.Node
    for _, node := range s.nodes {
        if query.Type != "" && node.Type != query.Type {
            continue
        }
        filteredNodes = append(filteredNodes, node)
    }

    // 计算分页
    total := len(filteredNodes)
    start := (query.Page - 1) * query.PageSize
    end := start + query.PageSize
    if end > total {
        end = total
    }

    // 返回分页结果
    return &models.NodeList{
        Total:    total,
        Page:     query.Page,
        PageSize: query.PageSize,
        Nodes:    filteredNodes[start:end],
    }, nil
}

// getNodeProtocol 获取节点传输协议
func (s *SubscriptionService) getNodeProtocol(proxy map[string]interface{}) string {
    nodeType, _ := proxy["type"].(string)
    switch nodeType {
    case "vmess":
        if network, ok := proxy["network"].(string); ok {
            return network
        }
        return "tcp"
    case "ss":
        return "shadowsocks"
    case "hysteria2":
        return "hysteria2"
    case "trojan":
        return "trojan"
    default:
        return "unknown"
    }
}

// UpdateAllSubscriptions 更新所有订阅
func (s *SubscriptionService) UpdateAllSubscriptions() error {
    utils.LogInfo("开始更新所有订阅")
    
    for _, sub := range s.subscriptions {
        // 更新订阅
        result, err := utils.ParseSubscription(sub.URL)
        if err != nil {
            utils.LogError("更新订阅失败: %v", err)
            continue
        }

        // 更新节点
        for _, node := range result.Nodes {
            s.nodes[node.ID] = node
        }

        // 更新订阅信息
        sub.NodeCount = len(result.Nodes)
        sub.UpdatedAt = time.Now()
    }

    utils.LogInfo("所有订阅更新完成")
    return nil
}

// TestAllNodes 测试所有节点
func (s *SubscriptionService) TestAllNodes() error {
    utils.LogInfo("开始测试所有节点")

    config := models.SpeedTestConfig{
        MaxLatency: config.GlobalConfig.Filter.MaxLatency,
        TestURL:    "http://speedtest.net",  // 这里应该从配置文件中读取
        Timeout:    10,                      // 这里应该从配置文件中读取
        Concurrent: config.GlobalConfig.Subscription.MaxConcurrent,
    }

    result, err := s.TestNodes(config)
    if err != nil {
        utils.LogError("节点测试失败: %v", err)
        return err
    }

    utils.LogSpeedTest(
        result.TotalCount,
        result.LatencyTested,
        result.LatencyDropped,
        result.SpeedTested,
    )

    return nil
} 