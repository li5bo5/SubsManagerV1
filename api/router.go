package api

import (
    "github.com/gin-gonic/gin"
)

// SetupRouter 设置路由
func SetupRouter() *gin.Engine {
    r := gin.Default()

    // 允许跨域
    r.Use(func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        c.Next()
    })

    // API路由组
    api := r.Group("/api")
    {
        // 订阅管理
        api.POST("/subscriptions", ImportSubscription)
        api.GET("/subscriptions", GetSubscriptions)
        api.DELETE("/subscriptions/:id", DeleteSubscription)
        api.POST("/subscriptions/merge", MergeSubscriptions)

        // 节点管理
        api.POST("/nodes/import", ImportNodes)
        api.GET("/nodes/list", GetNodeList)
        api.POST("/nodes/test", TestNodes)
        api.POST("/nodes/filter", FilterNodes)
        api.POST("/nodes/generate", GenerateSubscription)
    }

    // 静态文件服务
    r.Static("/data", "./data")

    return r
} 