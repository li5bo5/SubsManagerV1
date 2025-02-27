package api

import (
    "net/http"
)

// InitRouter 初始化路由
func InitRouter() {
    // 订阅相关接口
    http.HandleFunc("/api/subscription/import", ImportSubscription)
    http.HandleFunc("/api/subscription/list", ListSubscriptions)
    http.HandleFunc("/api/subscription/delete", DeleteSubscription)
    http.HandleFunc("/api/subscription/merge", MergeSubscriptions)
    
    // 节点相关接口
    http.HandleFunc("/api/nodes/speedtest", SpeedTest)
    http.HandleFunc("/api/nodes/filter", FilterNodes)
    
    // 系统日志接口
    http.HandleFunc("/api/logs", GetLogs)
    
    // 静态文件服务
    http.Handle("/", http.FileServer(http.Dir("web/dist")))
} 