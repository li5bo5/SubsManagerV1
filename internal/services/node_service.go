package services

import (
    "subsmanager/internal/models"
    "sync"
)

// NodeService 节点服务
type NodeService struct {
    nodes     map[string]*models.Node
    nodeMutex sync.RWMutex
}

// NewNodeService 创建节点服务
func NewNodeService() *NodeService {
    return &NodeService{
        nodes: make(map[string]*models.Node),
    }
}

// GetAllNodes 获取所有节点
func (s *NodeService) GetAllNodes() []*models.Node {
    s.nodeMutex.RLock()
    defer s.nodeMutex.RUnlock()

    nodes := make([]*models.Node, 0, len(s.nodes))
    for _, node := range s.nodes {
        nodes = append(nodes, node)
    }
    return nodes
}

// GetAllTestedNodes 获取所有已测试的节点
func (s *NodeService) GetAllTestedNodes() []*models.Node {
    s.nodeMutex.RLock()
    defer s.nodeMutex.RUnlock()

    nodes := make([]*models.Node, 0)
    for _, node := range s.nodes {
        if !node.LastTestedAt.IsZero() {
            nodes = append(nodes, node)
        }
    }
    return nodes
}

// GetFilteredNodes 获取筛选后的节点
func (s *NodeService) GetFilteredNodes() ([]*models.Node, error) {
    s.nodeMutex.RLock()
    defer s.nodeMutex.RUnlock()

    nodes := make([]*models.Node, 0)
    for _, node := range s.nodes {
        if !node.LastTestedAt.IsZero() && node.Latency > 0 && node.DownloadSpeed > 0 {
            nodes = append(nodes, node)
        }
    }
    return nodes, nil
}

// AddNode 添加节点
func (s *NodeService) AddNode(node *models.Node) {
    s.nodeMutex.Lock()
    defer s.nodeMutex.Unlock()
    s.nodes[node.ID] = node
}

// UpdateNode 更新节点
func (s *NodeService) UpdateNode(node *models.Node) {
    s.nodeMutex.Lock()
    defer s.nodeMutex.Unlock()
    s.nodes[node.ID] = node
}

// DeleteNode 删除节点
func (s *NodeService) DeleteNode(nodeID string) {
    s.nodeMutex.Lock()
    defer s.nodeMutex.Unlock()
    delete(s.nodes, nodeID)
}

// GetNode 获取节点
func (s *NodeService) GetNode(nodeID string) (*models.Node, bool) {
    s.nodeMutex.RLock()
    defer s.nodeMutex.RUnlock()
    node, exists := s.nodes[nodeID]
    return node, exists
} 