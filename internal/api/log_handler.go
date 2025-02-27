package api

import (
    "net/http"
    "subsmanager/internal/utils"
)

// GetLogs 获取系统日志
func GetLogs(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    logs := utils.GetRecentLogs()
    ResponseJSON(w, http.StatusOK, logs)
} 