package main

import (
    "fmt"
    "log"
    "os"
    "subsmanager/api"
    "subsmanager/config"
    "subsmanager/internal/services"
)

func init() {
    // 初始化配置
    if err := config.Init(); err != nil {
        log.Fatalf("Failed to initialize config: %v", err)
    }

    // 创建数据目录
    if err := os.MkdirAll(config.GlobalConfig.Storage.Path, 0755); err != nil {
        log.Fatalf("Failed to create data directory: %v", err)
    }

    // 加载数据
    if err := services.DefaultSubscriptionService.LoadFromFile(); err != nil {
        log.Printf("Failed to load data from file: %v", err)
    }
}

func main() {
    // 设置路由
    r := api.SetupRouter()

    // 启动服务器
    addr := fmt.Sprintf("%s:%d", config.GlobalConfig.Server.Host, config.GlobalConfig.Server.Port)
    log.Printf("Server starting on %s", addr)
    if err := r.Run(addr); err != nil {
        log.Fatalf("Server failed to start: %v", err)
    }
} 