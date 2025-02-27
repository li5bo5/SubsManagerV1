package utils

import (
    "encoding/base64"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "strings"
    "subsmanager/internal/models"
    "time"

    "gopkg.in/yaml.v3"
)

// SubscriptionType 订阅类型
type SubscriptionType string

const (
    TypeUnknown SubscriptionType = "unknown"
    TypeBase64  SubscriptionType = "base64"
    TypeYAML    SubscriptionType = "yaml"
    TypeJSON    SubscriptionType = "json"
)

// NodeType 节点类型
const (
    NodeTypeVmess      = "vmess"
    NodeTypeSS         = "ss"
    NodeTypeHysteria2  = "hysteria2"
    NodeTypeTrojan     = "trojan"
)

// ParseStats 解析统计
type ParseStats struct {
    Total     int
    Success   int
    Failed    map[string]int // key: node type, value: failed count
}

// SubscriptionParseResult 订阅解析结果
type SubscriptionParseResult struct {
    Type      SubscriptionType
    NodeCount int
    Nodes     []*models.Node
    Stats     ParseStats
}

// OpenClashConfig OpenClash配置结构
type OpenClashConfig struct {
    Proxies []map[string]interface{} `yaml:"proxies"`
}

// ParseSubscription 解析订阅链接
func ParseSubscription(url string) (*SubscriptionParseResult, error) {
    // 获取订阅内容
    content, err := fetchSubscriptionContent(url)
    if err != nil {
        return nil, fmt.Errorf("fetch subscription content failed: %v", err)
    }

    // 识别订阅类型
    subType := detectSubscriptionType(content)

    // 根据类型解析节点
    var nodes []*models.Node
    switch subType {
    case TypeBase64:
        nodes, err = parseBase64Subscription(content)
    case TypeYAML:
        nodes, err = parseYAMLSubscription(content)
    case TypeJSON:
        nodes, err = parseJSONSubscription(content)
    default:
        return nil, fmt.Errorf("unsupported subscription type")
    }

    if err != nil {
        return nil, fmt.Errorf("parse subscription failed: %v", err)
    }

    return &SubscriptionParseResult{
        Type:      subType,
        NodeCount: len(nodes),
        Nodes:     nodes,
    }, nil
}

// fetchSubscriptionContent 获取订阅内容
func fetchSubscriptionContent(url string) (string, error) {
    resp, err := http.Get(url)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }

    return string(body), nil
}

// detectSubscriptionType 识别订阅类型
func detectSubscriptionType(content string) SubscriptionType {
    // 尝试base64解码
    decoded, err := base64.StdEncoding.DecodeString(content)
    if err == nil && isValidBase64NodeList(string(decoded)) {
        return TypeBase64
    }

    // 检查是否是YAML格式
    if strings.Contains(content, "proxies:") {
        return TypeYAML
    }

    // 检查是否是JSON格式
    if strings.HasPrefix(strings.TrimSpace(content), "{") {
        return TypeJSON
    }

    return TypeUnknown
}

// isValidBase64NodeList 检查是否是有效的Base64节点列表
func isValidBase64NodeList(content string) bool {
    lines := strings.Split(content, "\n")
    for _, line := range lines {
        line = strings.TrimSpace(line)
        if line == "" {
            continue
        }
        prefix := strings.ToLower(line)
        if !strings.HasPrefix(prefix, "vmess://") && 
           !strings.HasPrefix(prefix, "ss://") && 
           !strings.HasPrefix(prefix, "hysteria2://") && 
           !strings.HasPrefix(prefix, "trojan://") {
            return false
        }
    }
    return true
}

// parseBase64Subscription 解析Base64编码的订阅
func parseBase64Subscription(content string) ([]*models.Node, error) {
    // Base64解码
    decoded, err := base64.StdEncoding.DecodeString(content)
    if err != nil {
        return nil, fmt.Errorf("base64 decode failed: %v", err)
    }

    // 分割成单个节点
    lines := strings.Split(string(decoded), "\n")
    nodes := make([]*models.Node, 0)
    stats := ParseStats{
        Total:   0,
        Success: 0,
        Failed:  make(map[string]int),
    }

    for _, line := range lines {
        line = strings.TrimSpace(line)
        if line == "" {
            continue
        }

        stats.Total++
        var node *models.Node
        var err error

        switch {
        case strings.HasPrefix(strings.ToLower(line), "vmess://"):
            node, err = parseVmessNode(line)
            if err != nil {
                stats.Failed[NodeTypeVmess]++
                LogParseError("base64", NodeTypeVmess, err)
            }
        case strings.HasPrefix(strings.ToLower(line), "ss://"):
            node, err = parseSSNode(line)
            if err != nil {
                stats.Failed[NodeTypeSS]++
                LogParseError("base64", NodeTypeSS, err)
            }
        case strings.HasPrefix(strings.ToLower(line), "hysteria2://"):
            node, err = parseHysteria2Node(line)
            if err != nil {
                stats.Failed[NodeTypeHysteria2]++
                LogParseError("base64", NodeTypeHysteria2, err)
            }
        case strings.HasPrefix(strings.ToLower(line), "trojan://"):
            node, err = parseTrojanNode(line)
            if err != nil {
                stats.Failed[NodeTypeTrojan]++
                LogParseError("base64", NodeTypeTrojan, err)
            }
        }

        if err != nil {
            continue
        }

        if node != nil {
            stats.Success++
            nodes = append(nodes, node)
        }
    }

    LogInfo("Base64 subscription parse stats: Total=%d, Success=%d, Failed=%v", 
        stats.Total, stats.Success, stats.Failed)

    return nodes, nil
}

// parseVmessNode 解析Vmess节点
func parseVmessNode(line string) (*models.Node, error) {
    // 移除前缀
    encoded := strings.TrimPrefix(line, "vmess://")
    
    // Base64解码
    decoded, err := base64.StdEncoding.DecodeString(encoded)
    if err != nil {
        return nil, err
    }

    // 解析JSON
    var vmessInfo struct {
        Add  string `json:"add"`
        Port int    `json:"port"`
        ID   string `json:"id"`
        Net  string `json:"net"`
        Type string `json:"type"`
        PS   string `json:"ps"`
    }

    if err := json.Unmarshal(decoded, &vmessInfo); err != nil {
        return nil, err
    }

    return &models.Node{
        Type:         "vmess",
        Alias:        vmessInfo.PS,
        Address:      vmessInfo.Add,
        Port:         vmessInfo.Port,
        Protocol:     vmessInfo.Net,
        LastTestedAt: time.Time{},
    }, nil
}

// parseSSNode 解析Shadowsocks节点
func parseSSNode(line string) (*models.Node, error) {
    // 移除前缀
    encoded := strings.TrimPrefix(line, "ss://")
    
    // 分离标签
    parts := strings.Split(encoded, "#")
    if len(parts) != 2 {
        return nil, fmt.Errorf("invalid ss link format")
    }

    // URL解码标签
    alias, err := url.QueryUnescape(parts[1])
    if err != nil {
        alias = parts[1]
    }

    // Base64解码配置部分
    decoded, err := base64.StdEncoding.DecodeString(parts[0])
    if err != nil {
        return nil, err
    }

    // 解析配置
    config := strings.Split(string(decoded), "@")
    if len(config) != 2 {
        return nil, fmt.Errorf("invalid ss config format")
    }

    // 解析地址和端口
    addrParts := strings.Split(config[1], ":")
    if len(addrParts) != 2 {
        return nil, fmt.Errorf("invalid ss address format")
    }

    port := 0
    fmt.Sscanf(addrParts[1], "%d", &port)

    return &models.Node{
        Type:         "ss",
        Alias:        alias,
        Address:      addrParts[0],
        Port:         port,
        Protocol:     "shadowsocks",
        LastTestedAt: time.Time{},
    }, nil
}

// parseHysteria2Node 解析Hysteria2节点
func parseHysteria2Node(line string) (*models.Node, error) {
    // 移除前缀
    uri := strings.TrimPrefix(line, "hysteria2://")
    
    // 解析URL
    u, err := url.Parse("hysteria2://" + uri)
    if err != nil {
        return nil, err
    }

    // 解析端口
    port := u.Port()
    if port == "" {
        port = "443" // 默认端口
    }
    portNum := 0
    fmt.Sscanf(port, "%d", &portNum)

    return &models.Node{
        Type:         NodeTypeHysteria2,
        Alias:        u.Fragment,
        Address:      u.Hostname(),
        Port:         portNum,
        Protocol:     "hysteria2",
        LastTestedAt: time.Time{},
    }, nil
}

// parseTrojanNode 解析Trojan节点
func parseTrojanNode(line string) (*models.Node, error) {
    // 移除前缀
    uri := strings.TrimPrefix(line, "trojan://")
    
    // 解析URL
    u, err := url.Parse("trojan://" + uri)
    if err != nil {
        return nil, err
    }

    // 解析端口
    port := u.Port()
    if port == "" {
        port = "443" // 默认端口
    }
    portNum := 0
    fmt.Sscanf(port, "%d", &portNum)

    return &models.Node{
        Type:         NodeTypeTrojan,
        Alias:        u.Fragment,
        Address:      u.Hostname(),
        Port:         portNum,
        Protocol:     "trojan",
        LastTestedAt: time.Time{},
    }, nil
}

// parseYAMLSubscription 解析YAML格式的订阅
func parseYAMLSubscription(content string) ([]*models.Node, error) {
    var config OpenClashConfig
    if err := yaml.Unmarshal([]byte(content), &config); err != nil {
        return nil, fmt.Errorf("yaml unmarshal failed: %v", err)
    }

    nodes := make([]*models.Node, 0)
    stats := ParseStats{
        Total:   len(config.Proxies),
        Success: 0,
        Failed:  make(map[string]int),
    }

    for _, proxy := range config.Proxies {
        nodeType, ok := proxy["type"].(string)
        if !ok {
            continue
        }

        var node *models.Node
        var err error

        switch strings.ToLower(nodeType) {
        case NodeTypeVmess:
            node, err = parseYAMLVmessNode(proxy)
        case NodeTypeSS:
            node, err = parseYAMLSSNode(proxy)
        case NodeTypeHysteria2:
            node, err = parseYAMLHysteria2Node(proxy)
        case NodeTypeTrojan:
            node, err = parseYAMLTrojanNode(proxy)
        default:
            continue
        }

        if err != nil {
            stats.Failed[nodeType]++
            LogParseError("yaml", nodeType, err)
            continue
        }

        if node != nil {
            stats.Success++
            nodes = append(nodes, node)
        }
    }

    LogInfo("YAML subscription parse stats: Total=%d, Success=%d, Failed=%v",
        stats.Total, stats.Success, stats.Failed)

    return nodes, nil
}

// parseYAMLVmessNode 解析YAML格式的Vmess节点
func parseYAMLVmessNode(proxy map[string]interface{}) (*models.Node, error) {
    name, _ := proxy["name"].(string)
    server, _ := proxy["server"].(string)
    port, _ := proxy["port"].(int)
    network, _ := proxy["network"].(string)

    return &models.Node{
        Type:         NodeTypeVmess,
        Alias:        name,
        Address:      server,
        Port:         port,
        Protocol:     network,
        LastTestedAt: time.Time{},
    }, nil
}

// parseYAMLSSNode 解析YAML格式的Shadowsocks节点
func parseYAMLSSNode(proxy map[string]interface{}) (*models.Node, error) {
    name, _ := proxy["name"].(string)
    server, _ := proxy["server"].(string)
    port, _ := proxy["port"].(int)

    return &models.Node{
        Type:         NodeTypeSS,
        Alias:        name,
        Address:      server,
        Port:         port,
        Protocol:     "shadowsocks",
        LastTestedAt: time.Time{},
    }, nil
}

// parseYAMLHysteria2Node 解析YAML格式的Hysteria2节点
func parseYAMLHysteria2Node(proxy map[string]interface{}) (*models.Node, error) {
    name, _ := proxy["name"].(string)
    server, _ := proxy["server"].(string)
    port, _ := proxy["port"].(int)

    return &models.Node{
        Type:         NodeTypeHysteria2,
        Alias:        name,
        Address:      server,
        Port:         port,
        Protocol:     "hysteria2",
        LastTestedAt: time.Time{},
    }, nil
}

// parseYAMLTrojanNode 解析YAML格式的Trojan节点
func parseYAMLTrojanNode(proxy map[string]interface{}) (*models.Node, error) {
    name, _ := proxy["name"].(string)
    server, _ := proxy["server"].(string)
    port, _ := proxy["port"].(int)

    return &models.Node{
        Type:         NodeTypeTrojan,
        Alias:        name,
        Address:      server,
        Port:         port,
        Protocol:     "trojan",
        LastTestedAt: time.Time{},
    }, nil
}

// parseJSONSubscription 解析JSON格式的订阅
func parseJSONSubscription(content string) ([]*models.Node, error) {
    // 尝试解析为OpenClash格式的JSON
    var config OpenClashConfig
    if err := json.Unmarshal([]byte(content), &config); err == nil {
        return parseYAMLSubscription(content)
    }

    return nil, fmt.Errorf("unsupported json format")
} 