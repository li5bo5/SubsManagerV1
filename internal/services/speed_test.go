package services

import (
    "context"
    "fmt"
    "io"
    "net"
    "net/http"
    "time"
)

// 测试相关配置常量
const (
    // 延迟测试配置
    DefaultLatencyTimeout = 5 * time.Second    // 默认延迟测试超时时间
    DefaultMaxConcurrent = 10                  // 默认最大并发数
    
    // 下载测试配置
    DefaultSpeedTimeout        = 30 * time.Second  // 默认下载测试超时时间
    DefaultDialTimeout        = 30 * time.Second  // 默认连接超时时间
    DefaultKeepAlive         = 30 * time.Second  // 默认连接保持时间
    DefaultTLSTimeout        = 10 * time.Second  // 默认TLS握手超时时间
    DefaultIdleTimeout       = 90 * time.Second  // 默认空闲连接超时时间
    DefaultExpectTimeout     = 1 * time.Second   // 默认 Expect: 100-continue 超时时间
    DefaultBufferSize        = 8192              // 默认缓冲区大小
    
    // 其他配置
    DefaultMaxIdleConns      = 100               // 默认最大空闲连接数
)

// TestConfig 测试配置
type TestConfig struct {
    LatencyTimeout  time.Duration // 延迟测试超时时间
    SpeedTimeout    time.Duration // 下载测试超时时间
    MaxConcurrent   int          // 最大并发数
    BufferSize      int          // 缓冲区大小
}

// NewDefaultTestConfig 创建默认测试配置
func NewDefaultTestConfig() *TestConfig {
    return &TestConfig{
        LatencyTimeout: DefaultLatencyTimeout,
        SpeedTimeout:   DefaultSpeedTimeout,
        MaxConcurrent:  DefaultMaxConcurrent,
        BufferSize:     DefaultBufferSize,
    }
}

// LatencyTestResult 存储延迟测试结果
type LatencyTestResult struct {
    NodeID      string
    Latency     int       // 延迟(ms)
    TestTime    time.Time
    Error       string
}

// SpeedTestResult 存储速度测试结果
type SpeedTestResult struct {
    NodeID      string
    Speed       float64   // 下载速度(MB/s)
    TestTime    time.Time
    Error       string
}

// TestRecord 测试记录
type TestRecord struct {
    ID          string    `json:"id"`
    NodeID      string    `json:"node_id"`
    Latency     int       `json:"latency"`      // ms
    Speed       float64   `json:"speed"`        // MB/s
    TestTime    time.Time `json:"test_time"`
    Error       string    `json:"error"`
}

// TestProgress 测试进度
type TestProgress struct {
    TotalNodes      int     `json:"total_nodes"`
    CompletedNodes  int     `json:"completed_nodes"`
    CurrentProgress float64 `json:"current_progress"` // 0-100
    Stage          string  `json:"stage"`            // "latency" or "speed"
}

// 测速服务器配置
var speedTestServers = []string{
    "http://cachefly.cachefly.net/100mb.test",
    "http://speedtest.tele2.net/100MB.zip",
    // 可以添加更多备用测速服务器
}

// TestManager 管理测试过程
type TestManager struct {
    config        *TestConfig
    Results       chan *LatencyTestResult
    logger        *Logger
    db            *sql.DB
}

// NewTestManager 创建测试管理器
func NewTestManager(config *TestConfig, logger *Logger, db *sql.DB) *TestManager {
    if config == nil {
        config = NewDefaultTestConfig()
    }
    
    return &TestManager{
        config:        config,
        Results:       make(chan *LatencyTestResult, config.MaxConcurrent),
        logger:        logger,
        db:            db,
    }
}

// testNodeLatency 测试节点延迟
func (tm *TestManager) testNodeLatency(node *Node) (*LatencyTestResult, error) {
    result := &LatencyTestResult{
        NodeID:   node.ID,
        TestTime: time.Now(),
    }

    ctx, cancel := context.WithTimeout(context.Background(), tm.config.LatencyTimeout)
    defer cancel()

    start := time.Now()
    var d net.Dialer
    conn, err := d.DialContext(ctx, "tcp", fmt.Sprintf("%s:%d", node.Address, node.Port))
    if err != nil {
        result.Error = fmt.Sprintf("连接失败: %v", err)
        return result, err
    }
    defer conn.Close()

    result.Latency = int(time.Since(start).Milliseconds())
    return result, nil
}

// testNodeSpeed 测试节点下载速度
func (tm *TestManager) testNodeSpeed(node *Node) (*SpeedTestResult, error) {
    result := &SpeedTestResult{
        NodeID:   node.ID,
        TestTime: time.Now(),
    }

    testURL := speedTestServers[0]

    client := &http.Client{
        Timeout: tm.config.SpeedTimeout,
        Transport: &http.Transport{
            Proxy: http.ProxyFromEnvironment,
            DialContext: (&net.Dialer{
                Timeout:   DefaultDialTimeout,
                KeepAlive: DefaultKeepAlive,
            }).DialContext,
            MaxIdleConns:          DefaultMaxIdleConns,
            IdleConnTimeout:       DefaultIdleTimeout,
            TLSHandshakeTimeout:   DefaultTLSTimeout,
            ExpectContinueTimeout: DefaultExpectTimeout,
        },
    }

    start := time.Now()
    resp, err := client.Get(testURL)
    if err != nil {
        result.Error = fmt.Sprintf("请求失败: %v", err)
        return result, err
    }
    defer resp.Body.Close()

    buf := make([]byte, tm.config.BufferSize)
    var totalBytes int64
    for {
        n, err := resp.Body.Read(buf)
        if n > 0 {
            totalBytes += int64(n)
        }
        if err == io.EOF {
            break
        }
        if err != nil {
            result.Error = fmt.Sprintf("读取失败: %v", err)
            return result, err
        }
    }

    duration := time.Since(start).Seconds()
    result.Speed = float64(totalBytes) / 1024 / 1024 / duration
    return result, nil
}

// StartLatencyTest 开始批量延迟测试
func (tm *TestManager) StartLatencyTest(nodes []*Node) {
    sem := make(chan struct{}, tm.config.MaxConcurrent)
    
    for _, node := range nodes {
        sem <- struct{}{} // 获取信号量
        
        go func(n *Node) {
            defer func() { <-sem }() // 释放信号量
            
            result, err := tm.testNodeLatency(n)
            if err != nil {
                tm.logger.Error("节点延迟测试失败", map[string]interface{}{
                    "nodeID": n.ID,
                    "error":  err.Error(),
                })
            }
            
            tm.Results <- result
        }(node)
    }
}

// StartSpeedTest 开始批量速度测试
func (tm *TestManager) StartSpeedTest(nodes []*Node) {
    sem := make(chan struct{}, tm.config.MaxConcurrent)
    
    for _, node := range nodes {
        // 只测试延迟<=400ms的节点
        if node.Latency > 400 {
            continue
        }

        sem <- struct{}{} // 获取信号量
        
        go func(n *Node) {
            defer func() { <-sem }() // 释放信号量
            
            result, err := tm.testNodeSpeed(n)
            if err != nil {
                tm.logger.Error("节点速度测试失败", map[string]interface{}{
                    "nodeID": n.ID,
                    "error":  err.Error(),
                })
            }
            
            // 更新节点速度信息
            n.Speed = result.Speed
            
            // 可以在这里实现结果持久化
        }(node)
    }
}

// SaveTestResult 保存测试结果
func (tm *TestManager) SaveTestResult(record *TestRecord) error {
    // 创建SQL语句
    query := `
        INSERT INTO node_test_records (
            id, node_id, latency, speed, test_time, error
        ) VALUES (?, ?, ?, ?, ?, ?)
    `
    
    // 执行SQL
    _, err := tm.db.Exec(query,
        record.ID,
        record.NodeID,
        record.Latency,
        record.Speed,
        record.TestTime,
        record.Error,
    )
    
    if err != nil {
        tm.logger.Error("保存测试记录失败", map[string]interface{}{
            "error": err.Error(),
            "record": record,
        })
        return err
    }
    
    return nil
}

// GetTestHistory 获取测试历史记录
func (tm *TestManager) GetTestHistory(nodeID string, limit int) ([]*TestRecord, error) {
    query := `
        SELECT id, node_id, latency, speed, test_time, error
        FROM node_test_records
        WHERE node_id = ?
        ORDER BY test_time DESC
        LIMIT ?
    `
    
    rows, err := tm.db.Query(query, nodeID, limit)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var records []*TestRecord
    for rows.Next() {
        record := &TestRecord{}
        err := rows.Scan(
            &record.ID,
            &record.NodeID,
            &record.Latency,
            &record.Speed,
            &record.TestTime,
            &record.Error,
        )
        if err != nil {
            return nil, err
        }
        records = append(records, record)
    }
    
    return records, nil
}

// UpdateTestProgress 更新测试进度
func (tm *TestManager) UpdateTestProgress(progress *TestProgress) {
    // 这里可以实现进度更新的逻辑
    // 比如通过WebSocket推送到前端
    // 或者存储到Redis等
} 