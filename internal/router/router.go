package router

import (
	"github.com/gorilla/mux"
	"subsmanager/internal/handlers"
)

func NewRouter(handler *handlers.Handler) *mux.Router {
	router := mux.NewRouter()

	// ... existing routes ...

	// 节点筛选相关路由
	router.HandleFunc("/api/nodes/filter", handler.FilterNodesHandler).Methods("POST")
	router.HandleFunc("/subscriptions/{filename}", handler.GetSubscriptionHandler).Methods("GET")

	// 状态监控路由
	router.HandleFunc("/api/status", handler.GetSystemStatus).Methods("GET")

	return router
}

// SetupRouter 设置路由
func SetupRouter(
	subscriptionHandler *handlers.SubscriptionHandler,
	statusHandler *handlers.StatusHandler,
	taskHandler *handlers.TaskHandler,
) *mux.Router {
	router := mux.NewRouter()

	// 订阅相关路由
	router.POST("/api/subscriptions", subscriptionHandler.ImportSubscription)
	router.GET("/api/subscriptions", subscriptionHandler.ListSubscriptions)
	router.GET("/api/subscriptions/:id", subscriptionHandler.GetSubscription)
	router.DELETE("/api/subscriptions/:id", subscriptionHandler.DeleteSubscription)
	router.POST("/api/subscriptions/merge", subscriptionHandler.MergeSubscriptions)
	router.POST("/api/nodes/test", subscriptionHandler.TestNodes)

	// 状态相关路由
	router.GET("/api/status", statusHandler.GetStatus)
	router.GET("/api/status/history", statusHandler.GetHistory)

	// 任务相关路由
	router.POST("/api/tasks", taskHandler.CreateTask)
	router.GET("/api/tasks", taskHandler.ListTasks)
	router.GET("/api/tasks/:id", taskHandler.GetTask)
	router.PUT("/api/tasks/:id", taskHandler.UpdateTask)
	router.DELETE("/api/tasks/:id", taskHandler.DeleteTask)

	return router
} 