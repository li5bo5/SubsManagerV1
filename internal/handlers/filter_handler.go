package handlers

import (
    "encoding/json"
    "net/http"
    "path/filepath"
)

// FilterNodesRequest 节点筛选请求
type FilterNodesRequest struct {
    MaxLatency int     `json:"max_latency"` // 延迟上限(ms)
    MinSpeed   float64 `json:"min_speed"`   // 速度下限(MB/s)
}

// FilterNodesResponse 节点筛选响应
type FilterNodesResponse struct {
    Nodes      []*Node `json:"nodes"`       // 筛选后的节点列表
    TotalNodes int     `json:"total_nodes"` // 节点总数
    SubURL     string  `json:"sub_url"`     // 订阅文件URL
}

// FilterNodesHandler 处理节点筛选请求
func (h *Handler) FilterNodesHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        h.responseError(w, http.StatusMethodNotAllowed, "方法不允许")
        return
    }

    // 解析请求
    var req FilterNodesRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        h.responseError(w, http.StatusBadRequest, "请求格式错误")
        return
    }

    // 创建筛选配置
    config := &services.FilterConfig{
        MaxLatency: req.MaxLatency,
        MinSpeed:   req.MinSpeed,
    }

    // 创建节点筛选器
    filter := services.NewNodeFilter(config, h.logger)

    // 获取所有已测速的节点
    nodes, err := h.nodeService.GetAllTestedNodes()
    if err != nil {
        h.responseError(w, http.StatusInternalServerError, "获取节点失败")
        return
    }

    // 筛选节点
    result, err := filter.FilterNodes(nodes)
    if err != nil {
        h.responseError(w, http.StatusInternalServerError, "节点筛选失败")
        return
    }

    // 生成订阅文件
    sub, err := filter.GenerateSubscription(result.Nodes)
    if err != nil {
        h.responseError(w, http.StatusInternalServerError, "生成订阅失败")
        return
    }

    // 保存订阅文件
    fileName := filter.generateFileName()
    filePath := filepath.Join(h.config.DataDir, fileName)
    err = filter.SaveSubscription(sub, filePath)
    if err != nil {
        h.responseError(w, http.StatusInternalServerError, "保存订阅失败")
        return
    }

    // 构建订阅URL
    subURL := fmt.Sprintf("http://%s/subscriptions/%s", r.Host, fileName)

    // 返回响应
    resp := FilterNodesResponse{
        Nodes:      result.Nodes,
        TotalNodes: result.TotalNodes,
        SubURL:     subURL,
    }

    h.responseJSON(w, http.StatusOK, resp)
}

// GetSubscriptionHandler 获取订阅文件
func (h *Handler) GetSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        h.responseError(w, http.StatusMethodNotAllowed, "方法不允许")
        return
    }

    // 获取文件名
    filename := filepath.Base(r.URL.Path)
    filePath := filepath.Join(h.config.DataDir, filename)

    // 检查文件是否存在
    if !FileExists(filePath) {
        h.responseError(w, http.StatusNotFound, "订阅文件不存在")
        return
    }

    // 设置响应头
    w.Header().Set("Content-Type", "application/x-yaml")
    w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

    // 发送文件
    http.ServeFile(w, r, filePath)
} 