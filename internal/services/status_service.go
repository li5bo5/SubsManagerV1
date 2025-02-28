package services

import (
    "fmt"
    "path/filepath"
    "sort"
    "subsmanager/config"
    "subsmanager/internal/models"
    "subsmanager/internal/utils"
    "sync"
    "time"
)

// StatusService 状态监控服务
type StatusService struct {
    subscriptionService *SubscriptionService
    nodeService        *NodeService
    config            *config.Config
    status            *models.SystemStatus
    historyMutex      sync.RWMutex
    history           []models.SubscriptionFileHistory
}

// NewStatusService 创建状态监控服务
func NewStatusService(subService *SubscriptionService, nodeService *NodeService, config *config.Config) *StatusService {
    return &StatusService{
        subscriptionService: subService,
        nodeService:        nodeService,
        config:            config,
        status:            &models.SystemStatus{},
        history:           make([]models.SubscriptionFileHistory, 0),
    }
}

// UpdateNodeStatus 更新节点状态
func (s *StatusService) UpdateNodeStatus() {
    nodes := s.nodeService.GetAllNodes()
    testedNodes := s.nodeService.GetAllTestedNodes()
    
    status := models.NodeStatus{
        TotalNodes:   len(nodes),
        CurrentNodes: 0,
        SlowNodes:    0,
        FaultNodes:   0,
    }

    // 统计节点状态
    for _, node := range testedNodes {
        if node.Latency > 400 { // 延迟大于400ms的节点视为故障节点
            status.FaultNodes++
        }
        if node.DownloadSpeed < s.config.Filter.MinSpeed {
            status.SlowNodes++
        }
    }

    // 获取当前筛选后的节点数
    filteredNodes, _ := s.nodeService.GetFilteredNodes()
    status.CurrentNodes = len(filteredNodes)

    s.status.NodeStatus = status
    s.status.LastUpdateTime = time.Now()
}

// AddSubscriptionHistory 添加订阅历史记录
func (s *StatusService) AddSubscriptionHistory(fileName string, nodeCount int) {
    s.historyMutex.Lock()
    defer s.historyMutex.Unlock()

    // 创建新的历史记录
    history := models.SubscriptionFileHistory{
        FileName:     fileName,
        LocalURL:     fmt.Sprintf("http://localhost:%d/subscriptions/%s", s.config.Server.Port, fileName),
        GenerateTime: time.Now(),
        NodeCount:    nodeCount,
    }

    // 添加到历史记录列表
    s.history = append(s.history, history)

    // 按生成时间排序
    sort.Slice(s.history, func(i, j int) bool {
        return s.history[i].GenerateTime.After(s.history[j].GenerateTime)
    })

    // 只保留最近三条记录
    if len(s.history) > 3 {
        s.history = s.history[:3]
    }

    // 更新系统状态中的历史记录
    s.status.SubHistory = s.history
}

// UpdateLatestFiles 更新最新文件信息
func (s *StatusService) UpdateLatestFiles() {
    s.status.LatestSubFile = filepath.Join(s.config.Storage.Path, "sub.yaml")
    s.status.LatestInputFile = filepath.Join(s.config.Storage.Path, s.getLatestInputFile())
}

// getLatestInputFile 获取最新的Sub-Input文件
func (s *StatusService) getLatestInputFile() string {
    if len(s.history) > 0 {
        return s.history[0].FileName
    }
    return ""
}

// AddTaskHistory 添加任务历史记录
func (s *StatusService) AddTaskHistory(result *models.TaskResult) {
    s.historyMutex.Lock()
    defer s.historyMutex.Unlock()

    // 更新系统状态
    s.status.LastTaskTime = result.EndTime
    s.status.LastTaskName = result.TaskID
    s.status.LastTaskStatus = result.Status

    // 记录日志
    if result.Error != "" {
        utils.LogError("任务执行失败: %s - %v", result.TaskID, result.Error)
    } else {
        utils.LogInfo("任务执行成功: %s", result.TaskID)
    }
}

// GetSystemStatus 获取系统状态
func (s *StatusService) GetSystemStatus() *models.SystemStatus {
    s.UpdateNodeStatus()
    s.UpdateLatestFiles()
    return s.status
} 