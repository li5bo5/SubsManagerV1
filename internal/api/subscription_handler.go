package api

import (
    "encoding/json"
    "net/http"
    "subsmanager/internal/model"
    "subsmanager/internal/utils"
)

// ListSubscriptions 获取订阅列表
func ListSubscriptions(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    subs := model.GetSubscriptions()
    ResponseJSON(w, http.StatusOK, subs)
}

// DeleteSubscription 删除订阅
func DeleteSubscription(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodDelete {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    id := r.URL.Query().Get("id")
    if id == "" {
        http.Error(w, "Missing subscription ID", http.StatusBadRequest)
        return
    }

    if err := model.DeleteSubscription(id); err != nil {
        http.Error(w, "Failed to delete subscription", http.StatusInternalServerError)
        return
    }

    // 记录删除日志
    utils.LogSubscriptionDelete(id)
    ResponseJSON(w, http.StatusOK, map[string]string{"message": "Subscription deleted successfully"})
}

// MergeSubscriptions 整合订阅
func MergeSubscriptions(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // 获取所有节点并去重
    mergedNodes, err := model.MergeSubscriptions()
    if err != nil {
        http.Error(w, "Failed to merge subscriptions", http.StatusInternalServerError)
        return
    }

    // 记录整合日志
    utils.LogSubscriptionMerge(len(mergedNodes))

    // 返回整合后的节点列表
    ResponseJSON(w, http.StatusOK, map[string]interface{}{
        "message": "Subscriptions merged successfully",
        "nodes": mergedNodes,
    })
} 