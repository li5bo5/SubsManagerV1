package api

import (
    "net/http"
    "subsmanager/internal/models"
    "subsmanager/internal/services"

    "github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}

// ImportSubscriptionRequest 导入订阅请求
type ImportSubscriptionRequest struct {
    Name string `json:"name" binding:"required"`
    URL  string `json:"url" binding:"required,url"`
}

// ImportSubscription 导入订阅
func ImportSubscription(c *gin.Context) {
    var req ImportSubscriptionRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, Response{
            Code:    400,
            Message: "Invalid request parameters",
        })
        return
    }

    sub, err := services.DefaultSubscriptionService.ImportSubscription(req.Name, req.URL)
    if err != nil {
        c.JSON(http.StatusInternalServerError, Response{
            Code:    500,
            Message: err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, Response{
        Code:    200,
        Message: "Success",
        Data:    sub,
    })
}

// GetSubscriptions 获取所有订阅
func GetSubscriptions(c *gin.Context) {
    subs := services.DefaultSubscriptionService.GetSubscriptions()
    c.JSON(http.StatusOK, Response{
        Code:    200,
        Message: "Success",
        Data:    subs,
    })
}

// DeleteSubscription 删除订阅
func DeleteSubscription(c *gin.Context) {
    id := c.Param("id")
    if err := services.DefaultSubscriptionService.DeleteSubscription(id); err != nil {
        c.JSON(http.StatusInternalServerError, Response{
            Code:    500,
            Message: err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, Response{
        Code:    200,
        Message: "Success",
    })
}

// MergeSubscriptionsRequest 合并订阅请求
type MergeSubscriptionsRequest struct {
    IDs []string `json:"ids" binding:"required,min=1"`
}

// MergeSubscriptions 合并订阅
func MergeSubscriptions(c *gin.Context) {
    var req MergeSubscriptionsRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, Response{
            Code:    400,
            Message: "Invalid request parameters",
        })
        return
    }

    result, err := services.DefaultSubscriptionService.MergeSubscriptions(req.IDs)
    if err != nil {
        c.JSON(http.StatusInternalServerError, Response{
            Code:    500,
            Message: err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, Response{
        Code:    200,
        Message: "Success",
        Data:    result,
    })
}

// TestNodesRequest 节点测试请求
type TestNodesRequest struct {
    MaxLatency int     `json:"max_latency" binding:"required"`
    TestURL    string  `json:"test_url" binding:"required"`
    Timeout    int     `json:"timeout" binding:"required"`
    Concurrent int     `json:"concurrent" binding:"required"`
}

// TestNodes 测试节点速度
func TestNodes(c *gin.Context) {
    var req TestNodesRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // 验证参数
    if req.MaxLatency <= 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "max_latency must be greater than 0"})
        return
    }
    if req.Timeout <= 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "timeout must be greater than 0"})
        return
    }
    if req.Concurrent <= 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "concurrent must be greater than 0"})
        return
    }
    if req.TestURL == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "test_url is required"})
        return
    }

    // 创建测试配置
    config := models.SpeedTestConfig{
        MaxLatency: req.MaxLatency,
        TestURL:    req.TestURL,
        Timeout:    req.Timeout,
        Concurrent: req.Concurrent,
    }

    // 执行节点测试
    result, err := services.GetSubscriptionService().TestNodes(config)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, result)
}

// FilterNodesRequest 筛选节点请求
type FilterNodesRequest struct {
    MaxLatency       int     `json:"max_latency" binding:"required,min=0"`
    MinDownloadSpeed float64 `json:"min_download_speed" binding:"required,min=0"`
}

// FilterNodes 筛选节点
func FilterNodes(c *gin.Context) {
    var req FilterNodesRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, Response{
            Code:    400,
            Message: "Invalid request parameters",
        })
        return
    }

    condition := models.FilterCondition{
        MaxLatency:       req.MaxLatency,
        MinDownloadSpeed: req.MinDownloadSpeed,
    }

    nodes, err := services.DefaultSubscriptionService.FilterNodes(condition)
    if err != nil {
        c.JSON(http.StatusInternalServerError, Response{
            Code:    500,
            Message: err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, Response{
        Code:    200,
        Message: "Success",
        Data:    nodes,
    })
}

// GenerateSubscription 生成订阅
func GenerateSubscription(c *gin.Context) {
    // TODO: 实现生成订阅的处理逻辑
    c.JSON(http.StatusOK, Response{
        Code:    200,
        Message: "Success",
    })
}

// ImportNodesRequest 导入节点请求
type ImportNodesRequest struct {
    FilePath string `json:"file_path" binding:"required"`
}

// ImportNodes 导入节点
func ImportNodes(c *gin.Context) {
    var req ImportNodesRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, Response{
            Code:    400,
            Message: "Invalid request parameters",
        })
        return
    }

    result, err := services.DefaultSubscriptionService.ImportNodesFromFile(req.FilePath)
    if err != nil {
        c.JSON(http.StatusInternalServerError, Response{
            Code:    500,
            Message: err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, Response{
        Code:    200,
        Message: "Success",
        Data:    result,
    })
}

// GetNodeList 获取节点列表
func GetNodeList(c *gin.Context) {
    var query models.NodeListQuery
    if err := c.ShouldBindQuery(&query); err != nil {
        c.JSON(http.StatusBadRequest, Response{
            Code:    400,
            Message: "Invalid query parameters",
        })
        return
    }

    result, err := services.DefaultSubscriptionService.GetNodeList(query)
    if err != nil {
        c.JSON(http.StatusInternalServerError, Response{
            Code:    500,
            Message: err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, Response{
        Code:    200,
        Message: "Success",
        Data:    result,
    })
} 